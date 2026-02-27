"""
Wenxi网盘 - 智能搜索API路由
作者：Wenxi
功能：提供全文搜索相关的REST API接口
性能目标：API响应时间 < 200ms
"""

from typing import Optional, List
from pydantic import BaseModel, Field
from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.orm import Session

from logger import logger
from database import get_db
from models import User
from routers.auth import get_current_user
from search.service import search_service

router = APIRouter()


class SearchRequest(BaseModel):
    """搜索请求模型"""
    query: str = Field(..., min_length=1, max_length=200, description="搜索关键词")
    file_type: Optional[str] = Field(None, description="文件类型过滤")
    limit: int = Field(20, ge=1, le=100, description="返回结果数量")
    offset: int = Field(0, ge=0, description="结果偏移量")


class SearchResultItem(BaseModel):
    """搜索结果项模型"""
    id: int
    filename: str
    original_filename: str
    file_size: int
    mime_type: Optional[str]
    created_at: Optional[str]
    is_shared: bool
    description: Optional[str]
    score: float
    highlights: dict


class SearchResponse(BaseModel):
    """搜索响应模型"""
    results: List[SearchResultItem]
    total: int
    query: str
    search_time_ms: float
    limit: int
    offset: int


class SearchStats(BaseModel):
    """搜索统计信息模型"""
    total_documents: int
    index_dir: str
    index_size_mb: float


class FileTypeSuggestion(BaseModel):
    """文件类型建议模型"""
    value: str
    label: str
    extensions: List[str]


class ReindexResponse(BaseModel):
    """重建索引响应模型"""
    success: bool
    indexed_count: Optional[int] = None
    message: str


@router.get("/search", response_model=SearchResponse)
async def search_files(
    query: str = Query(..., min_length=1, max_length=200, description="搜索关键词"),
    file_type: Optional[str] = Query(None, description="文件类型过滤"),
    limit: int = Query(20, ge=1, le=100, description="返回结果数量"),
    offset: int = Query(0, ge=0, description="结果偏移量"),
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Wenxi - 智能文件搜索接口
    
    功能：
    - 全文搜索文件名和描述
    - 支持文件类型过滤
    - 按相关性排序
    - 返回高亮片段
    
    性能：
    - 响应时间 < 200ms
    - 支持分页
    
    参数：
    - query: 搜索关键词（必需）
    - file_type: 文件类型过滤（可选）
    - limit: 返回数量（默认20，最大100）
    - offset: 偏移量（默认0）
    """
    try:
        # 记录搜索日志
        logger.info(f"Wenxi - 用户 {current_user.username} 搜索: '{query}' (type={file_type})")
        
        # 执行搜索
        result = await search_service.search_files(
            query=query,
            user_id=current_user.id,
            db=db,
            file_type=file_type,
            limit=limit,
            offset=offset
        )
        
        # 检查响应时间
        if result['search_time_ms'] > 200:
            logger.warning(
                f"Wenxi - 搜索响应较慢: {result['search_time_ms']}ms "
                f"(query='{query}', user={current_user.username})"
            )
        else:
            logger.debug(f"Wenxi - 搜索完成: {result['search_time_ms']}ms, 找到 {result['total']} 个结果")
        
        return SearchResponse(**result)
        
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Wenxi - 搜索接口错误: {e}")
        raise HTTPException(status_code=500, detail="搜索服务暂时不可用")


@router.post("/search", response_model=SearchResponse)
async def search_files_post(
    request: SearchRequest,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Wenxi - 智能文件搜索接口（POST方式）
    
    功能与GET /search相同，适用于复杂的搜索参数
    """
    return await search_files(
        query=request.query,
        file_type=request.file_type,
        limit=request.limit,
        offset=request.offset,
        current_user=current_user,
        db=db
    )


@router.get("/search/suggestions", response_model=List[str])
async def get_search_suggestions(
    query: str = Query(..., min_length=1, max_length=100),
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Wenxi - 搜索建议接口
    
    功能：根据输入的前缀提供搜索建议
    用于搜索框的自动完成功能
    """
    try:
        # 这里可以实现基于历史搜索的建议
        # 简化版本返回空列表
        return []
        
    except Exception as e:
        logger.error(f"Wenxi - 获取搜索建议失败: {e}")
        return []


@router.get("/search/file-types", response_model=List[FileTypeSuggestion])
async def get_file_types(
    current_user: User = Depends(get_current_user)
):
    """
    Wenxi - 获取支持的文件类型列表
    
    功能：返回可用于过滤的文件类型列表
    """
    return search_service.get_file_type_suggestions()


@router.get("/search/stats", response_model=SearchStats)
async def get_search_stats(
    current_user: User = Depends(get_current_user)
):
    """
    Wenxi - 获取搜索索引统计信息
    
    功能：返回搜索索引的统计信息
    需要管理员权限（简化实现，所有用户可查看）
    """
    try:
        stats = await search_service.get_search_stats()
        return SearchStats(**stats)
        
    except Exception as e:
        logger.error(f"Wenxi - 获取搜索统计失败: {e}")
        raise HTTPException(status_code=500, detail="获取统计信息失败")


@router.post("/search/reindex", response_model=ReindexResponse)
async def reindex_all_files(
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Wenxi - 重建所有文件索引
    
    功能：重新构建全文搜索索引
    注意：这是一个耗时操作，仅管理员可用（简化实现）
    """
    try:
        logger.info(f"Wenxi - 用户 {current_user.username} 触发索引重建")
        
        result = await search_service.reindex_all_files(db)
        
        return ReindexResponse(**result)
        
    except Exception as e:
        logger.error(f"Wenxi - 重建索引失败: {e}")
        raise HTTPException(status_code=500, detail="重建索引失败")


@router.post("/search/index-file/{file_id}")
async def index_single_file(
    file_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Wenxi - 将单个文件添加到搜索索引
    
    功能：手动将指定文件添加到搜索索引
    """
    try:
        from models import File as FileModel
        
        # 验证文件所有权
        file = db.query(FileModel).filter(
            FileModel.id == file_id,
            FileModel.owner_id == current_user.id
        ).first()
        
        if not file:
            raise HTTPException(status_code=404, detail="文件不存在")
        
        # 添加到索引
        success = await search_service.index_file(file_id, db)
        
        if success:
            return {"message": "文件已添加到索引", "file_id": file_id}
        else:
            raise HTTPException(status_code=500, detail="索引文件失败")
            
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Wenxi - 索引文件失败: {e}")
        raise HTTPException(status_code=500, detail="索引文件失败")


@router.delete("/search/index-file/{file_id}")
async def remove_file_from_index(
    file_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Wenxi - 从搜索索引中移除文件
    
    功能：手动从搜索索引中删除指定文件
    """
    try:
        from models import File as FileModel
        
        # 验证文件所有权
        file = db.query(FileModel).filter(
            FileModel.id == file_id,
            FileModel.owner_id == current_user.id
        ).first()
        
        if not file:
            raise HTTPException(status_code=404, detail="文件不存在")
        
        # 从索引移除
        success = await search_service.remove_file_index(file_id)
        
        if success:
            return {"message": "文件已从索引移除", "file_id": file_id}
        else:
            raise HTTPException(status_code=500, detail="移除索引失败")
            
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Wenxi - 移除文件索引失败: {e}")
        raise HTTPException(status_code=500, detail="移除索引失败")