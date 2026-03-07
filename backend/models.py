"""
Wenxi网盘 - 数据库模型
作者：Wenxi
功能：定义用户和文件的数据库模型
"""

from datetime import datetime, timezone
from sqlalchemy import Column, Integer, String, DateTime, Boolean, ForeignKey, Text
from sqlalchemy.orm import declarative_base
from sqlalchemy.orm import relationship

Base = declarative_base()


class User(Base):
    """用户模型 - Wenxi权限增强版"""

    __tablename__ = "users"

    id = Column(Integer, primary_key=True, index=True)
    username = Column(String(50), unique=True, index=True, nullable=False)
    email = Column(String(100), unique=True, index=True, nullable=False)
    hashed_password = Column(String(100), nullable=False)
    is_active = Column(Boolean, default=True)

    # 权限系统字段
    role = Column(Integer, default=2)  # 0:超级管理员, 1:管理员, 2:普通用户, 3:访客
    permissions = Column(Integer, default=15)  # 权限位掩码，默认普通用户权限(0x0F)

    created_at = Column(DateTime, default=lambda: datetime.now(timezone.utc))
    updated_at = Column(
        DateTime,
        default=lambda: datetime.now(timezone.utc),
        onupdate=lambda: datetime.now(timezone.utc),
    )

    # 关联文件
    files = relationship("File", back_populates="owner")

    def has_permission(self, permission: int) -> bool:
        """检查用户是否拥有指定权限"""
        return bool(self.permissions & permission)

    def add_permission(self, permission: int):
        """添加权限"""
        self.permissions |= permission

    def remove_permission(self, permission: int):
        """移除权限"""
        self.permissions &= ~permission

    def is_admin(self) -> bool:
        """检查是否为管理员"""
        return self.role <= 1  # 0: 超级管理员, 1: 管理员

    def is_super_admin(self) -> bool:
        """检查是否为超级管理员"""
        return self.role == 0


class File(Base):
    """文件模型"""

    __tablename__ = "files"

    id = Column(Integer, primary_key=True, index=True)
    filename = Column(String(255), nullable=False)
    original_filename = Column(String(255), nullable=False)
    file_path = Column(String(500), nullable=False)  # 存储相对路径
    file_size = Column(Integer, nullable=False)  # 字节
    mime_type = Column(String(100))

    # 关联用户
    owner_id = Column(Integer, ForeignKey("users.id"), nullable=False)
    owner = relationship("User", back_populates="files")

    # 文件状态
    is_shared = Column(Boolean, default=False)
    share_token = Column(String(32), unique=True, index=True, nullable=True)

    # 时间戳
    created_at = Column(DateTime, default=lambda: datetime.now(timezone.utc))
    updated_at = Column(
        DateTime,
        default=lambda: datetime.now(timezone.utc),
        onupdate=lambda: datetime.now(timezone.utc),
    )

    # 额外信息
    description = Column(Text, nullable=True)
    checksum = Column(String(64), nullable=True)  # 文件校验和

    # 版本控制
    versions = relationship(
        "FileVersion",
        back_populates="file",
        order_by="FileVersion.version_number.desc()",
    )

    # 垃圾桶功能字段
    is_deleted = Column(Boolean, default=False)  # 软删除标记
    deleted_at = Column(DateTime, nullable=True)  # 删除时间
    original_path = Column(String(500), nullable=True)  # 恢复时使用的原始路径
