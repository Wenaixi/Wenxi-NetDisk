"""
Wenxi网盘 - 用户认证模块
作者：Wenxi
功能：处理用户注册、登录、JWT认证
"""

from datetime import datetime, timedelta, timezone
from typing import Optional

from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from pydantic import BaseModel, EmailStr
from sqlalchemy.orm import Session
from passlib.context import CryptContext
import jwt
from jwt.exceptions import InvalidTokenError

from logger import logger
from database import get_db
from models import User, File as FileModel
import os


# JWT配置
SECRET_KEY = "wenxi-super-secret-key-change-in-production"
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 30

# 密码加密
pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="api/auth/login")

router = APIRouter()


class UserCreate(BaseModel):
    """用户注册模型"""
    username: str
    email: EmailStr
    password: str


class UserLogin(BaseModel):
    """用户登录模型"""
    username: str
    password: str


class Token(BaseModel):
    """JWT令牌模型"""
    access_token: str
    token_type: str


class UserResponse(BaseModel):
    """用户响应模型"""
    id: int
    username: str
    email: str
    created_at: datetime


def verify_password(plain_password: str, hashed_password: str) -> bool:
    """验证密码"""
    return pwd_context.verify(plain_password, hashed_password)


def get_password_hash(password: str) -> str:
    """加密密码"""
    return pwd_context.hash(password)


def create_access_token(data: dict, expires_delta: Optional[timedelta] = None):
    """创建JWT令牌"""
    to_encode = data.copy()
    if expires_delta:
        expire = datetime.now(timezone.utc) + expires_delta
    else:
        expire = datetime.now(timezone.utc) + timedelta(minutes=15)
    
    to_encode.update({"exp": expire})
    encoded_jwt = jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)
    return encoded_jwt


async def get_current_user(token: str = Depends(oauth2_scheme), db: Session = Depends(get_db)):
    """获取当前用户"""
    credentials_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="无法验证凭据",
        headers={"WWW-Authenticate": "Bearer"},
    )
    
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        username: str = payload.get("sub")
        if username is None:
            raise credentials_exception
    except InvalidTokenError:
        raise credentials_exception
    
    user = db.query(User).filter(User.username == username).first()
    if user is None:
        raise credentials_exception
    return user


async def get_current_user_from_token(token: str, db: Session = Depends(get_db)):
    """通过token字符串获取当前用户（用于iframe下载）"""
    credentials_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="无效的认证令牌",
    )
    
    try:
        # 移除Bearer前缀（如果有）
        if token.startswith('Bearer '):
            token = token[7:]
            
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        username: str = payload.get("sub")
        if username is None:
            raise credentials_exception
    except InvalidTokenError:
        raise credentials_exception
    
    user = db.query(User).filter(User.username == username).first()
    if user is None:
        raise credentials_exception
    return user


@router.post("/register", response_model=UserResponse)
async def register(user: UserCreate, db: Session = Depends(get_db)):
    """用户注册"""
    logger.info(f"注册新用户: {user.username}")
    
    # 检查用户名是否已存在
    db_user = db.query(User).filter(User.username == user.username).first()
    if db_user:
        logger.warning(f"用户名已存在: {user.username}")
        raise HTTPException(
            status_code=400,
            detail="用户名已存在"
        )
    
    # 检查邮箱是否已存在
    db_user = db.query(User).filter(User.email == user.email).first()
    if db_user:
        logger.warning(f"邮箱已存在: {user.email}")
        raise HTTPException(
            status_code=400,
            detail="邮箱已存在"
        )
    
    # 创建新用户
    hashed_password = get_password_hash(user.password)
    db_user = User(
        username=user.username,
        email=user.email,
        hashed_password=hashed_password
    )
    
    db.add(db_user)
    db.commit()
    db.refresh(db_user)
    
    logger.info(f"用户注册成功: {user.username} (ID: {db_user.id})")
    return db_user


@router.post("/login", response_model=Token)
async def login(form_data: OAuth2PasswordRequestForm = Depends(), db: Session = Depends(get_db)):
    """用户登录"""
    logger.info(f"用户登录: {form_data.username}")
    
    # 查找用户
    user = db.query(User).filter(User.username == form_data.username).first()
    if not user or not verify_password(form_data.password, user.hashed_password):
        logger.warning(f"登录失败: {form_data.username}")
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="用户名或密码错误",
            headers={"WWW-Authenticate": "Bearer"},
        )
    
    # 创建访问令牌
    access_token_expires = timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    access_token = create_access_token(
        data={"sub": user.username}, expires_delta=access_token_expires
    )
    
    logger.info(f"用户登录成功: {user.username}")
    return {"access_token": access_token, "token_type": "bearer"}


@router.get("/me", response_model=UserResponse)
async def read_users_me(current_user: User = Depends(get_current_user)):
    """获取当前用户信息"""
    return current_user


class UpdateUsernameRequest(BaseModel):
    """修改用户名请求模型"""
    username: str


class UpdateEmailRequest(BaseModel):
    """修改邮箱请求模型"""
    email: EmailStr
    password: str


class UpdatePasswordRequest(BaseModel):
    """修改密码请求模型"""
    old_password: str
    new_password: str


class DeleteAccountRequest(BaseModel):
    """删除账户请求模型"""
    username: str
    password: str
    email: str


@router.put("/update-username")
async def update_username(
    request: UpdateUsernameRequest,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """修改用户名"""
    logger.info(f"用户 {current_user.username} 尝试修改用户名为: {request.username}")
    
    # 检查新用户名是否已存在
    existing_user = db.query(User).filter(
        User.username == request.username,
        User.id != current_user.id
    ).first()
    
    if existing_user:
        logger.warning(f"用户名已存在: {request.username}")
        raise HTTPException(
            status_code=400,
            detail="用户名已存在"
        )
    
    # 更新用户名
    current_user.username = request.username
    db.commit()
    
    logger.info(f"用户名修改成功: {current_user.id} -> {request.username}")
    return {"message": "用户名修改成功"}


@router.put("/update-email")
async def update_email(
    request: UpdateEmailRequest,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """修改邮箱"""
    logger.info(f"用户 {current_user.username} 尝试修改邮箱为: {request.email}")
    
    # 验证原密码
    if not verify_password(request.password, current_user.hashed_password):
        logger.warning(f"密码验证失败: {current_user.username}")
        raise HTTPException(
            status_code=400,
            detail="密码错误"
        )
    
    # 检查新邮箱是否已存在
    existing_user = db.query(User).filter(
        User.email == request.email,
        User.id != current_user.id
    ).first()
    
    if existing_user:
        logger.warning(f"邮箱已存在: {request.email}")
        raise HTTPException(
            status_code=400,
            detail="邮箱已存在"
        )
    
    # 更新邮箱
    current_user.email = request.email
    db.commit()
    
    logger.info(f"邮箱修改成功: {current_user.id} -> {request.email}")
    return {"message": "邮箱修改成功"}


@router.put("/update-password")
async def update_password(
    request: UpdatePasswordRequest,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """修改密码"""
    logger.info(f"用户 {current_user.username} 尝试修改密码")
    
    # 验证原密码
    if not verify_password(request.old_password, current_user.hashed_password):
        logger.warning(f"原密码验证失败: {current_user.username}")
        raise HTTPException(
            status_code=400,
            detail="原密码错误"
        )
    
    # 更新密码
    current_user.hashed_password = get_password_hash(request.new_password)
    db.commit()
    
    logger.info(f"密码修改成功: {current_user.username}")
    return {"message": "密码修改成功"}


@router.delete("/delete-account")
async def delete_account(
    request: DeleteAccountRequest,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """删除账户 - 永久删除用户及其所有数据"""
    logger.info(f"用户 {current_user.username} 尝试删除账户")
    
    # 验证用户名
    if request.username != current_user.username:
        logger.warning(f"用户名验证失败: {request.username} != {current_user.username}")
        raise HTTPException(
            status_code=400,
            detail="用户名不正确"
        )
    
    # 验证邮箱
    if request.email != current_user.email:
        logger.warning(f"邮箱验证失败: {request.email} != {current_user.email}")
        raise HTTPException(
            status_code=400,
            detail="邮箱不正确"
        )
    
    # 验证密码
    if not verify_password(request.password, current_user.hashed_password):
        logger.warning(f"密码验证失败: {current_user.username}")
        raise HTTPException(
            status_code=400,
            detail="密码不正确"
        )
    
    try:
        # 获取用户的所有文件
        user_files = db.query(FileModel).filter(FileModel.owner_id == current_user.id).all()
        
        # 删除物理文件
        for file_record in user_files:
            try:
                file_path = os.path.join(os.path.dirname(__file__), "..", file_record.file_path)
                if os.path.exists(file_path):
                    os.remove(file_path)
                    logger.info(f"已删除文件: {file_record.original_filename}")
            except Exception as e:
                logger.error(f"删除文件失败: {file_record.file_path}, 错误: {e}")
        
        # 删除数据库中的文件记录
        db.query(FileModel).filter(FileModel.owner_id == current_user.id).delete()
        
        # 删除用户记录
        db.delete(current_user)
        db.commit()
        
        logger.info(f"账户删除成功: {current_user.username} (ID: {current_user.id})")
        
        return {"message": "账户已成功删除"}
        
    except Exception as e:
        db.rollback()
        logger.error(f"删除账户失败: {current_user.username}, 错误: {e}")
        raise HTTPException(
            status_code=500,
            detail="删除账户时发生错误"
        )


@router.get("/health")
async def health_check():
    """健康检查端点"""
    return {"status": "healthy", "service": "auth", "timestamp": datetime.now(timezone.utc).isoformat()}