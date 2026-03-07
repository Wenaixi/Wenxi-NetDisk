"""
Wenxi网盘 - 搜索服务模块
作者：Wenxi
功能：提供搜索相关的业务逻辑和API接口
"""

from typing import List, Optional, Dict
from datetime import datetime
from fastapi import HTTPException
from sqlalchemy.orm import Session

from logger import logger
from models import File as FileModel
from search.index_manager import search_manager, extract_text_content


class SearchService:
    """
    Wenxi - 搜索服务类
    功能：处理搜索相关的业务逻辑
    """
    
    @staticmethod
    async def index_file(file_id: int, db: Session) -> bool:
        """
        将文件添加到搜索索引
        
        Args:
            file_id: 文件ID
            db: 数据库会话
            
        Returns:
            是否成功
        """
        try:
            # 获取文件信息
            file = db.query(FileModel).filter(FileModel.id == file_id).first()
            if not file:
                logger.warning(f"文件不存在，无法索引: {file_id}")
                return False
            
            # 构建文件数据
            file_data = {
                'file_id': file.id,
                'filename': file.filename,
                'original_filename': file.original_filename,
                'description': file.description or '',
                'owner_id': file.owner_id,
                'file_size': file.file_size,
                'created_at': file.created_at,
                'updated_at': file.updated_at,
                'mime_type': file.mime_type or 'application/octet-stream',
                'content': ''  # 默认空内容
            }
            
            # 尝试提取文本文件内容
            # 注意：实际生产环境中可能需要解密文件后再提取
            # 这里简化处理，仅提取文件名的关键词信息
            
            # 添加到索引
            success = await search_manager.add_file(file_data)
            
            if success:
                logger.info(f"Wenxi - 文件已索引: {file.original_filename}")
            
            return success
            
        except Exception as e:
            logger.error(f"Wenxi - 索引文件失败: {e}")
            return False
    
    @staticmethod
    async def update_file_index(file_id: int, db: Session) -> bool:
        """
        更新文件索引
        
        Args:
            file_id: 文件ID
            db: 数据库会话
            
        Returns:
            是否成功
        """
        try:
            file = db.query(FileModel).filter(FileModel.id == file_id).first()
            if not file:
                logger.warning(f"文件不存在，无法更新索引: {file_id}")
                return False
            
            file_data = {
                'file_id': file.id,
                'filename': file.filename,
                'original_filename': file.original_filename,
                'description': file.description or '',
                'owner_id': file.owner_id,
                'file_size': file.file_size,
                'created_at': file.created_at,
                'updated_at': file.updated_at,
                'mime_type': file.mime_type or 'application/octet-stream',
                'content': ''
            }
            
            success = await search_manager.update_file(file_data)
            
            if success:
                logger.info(f"Wenxi - 文件索引已更新: {file.original_filename}")
            
            return success
            
        except Exception as e:
            logger.error(f"Wenxi - 更新文件索引失败: {e}")
            return False
    
    @staticmethod
    async def remove_file_index(file_id: int) -> bool:
        """
        从索引中移除文件
        
        Args:
            file_id: 文件ID
            
        Returns:
            是否成功
        """
        try:
            success = await search_manager.delete_file(str(file_id))
            
            if success:
                logger.info(f"Wenxi - 文件索引已删除: {file_id}")
            
            return success
            
        except Exception as e:
            logger.error(f"Wenxi - 删除文件索引失败: {e}")
            return False
    
    @staticmethod
    async def search_files(
        query: str,
        user_id: int,
        db: Session,
        file_type: Optional[str] = None,
        limit: int = 20,
        offset: int = 0
    ) -> Dict:
        """
        搜索文件
        
        Args:
            query: 搜索关键词
            user_id: 用户ID
            db: 数据库会话
            file_type: 文件类型过滤（可选）
            limit: 结果数量限制
            offset: 结果偏移量
            
        Returns:
            搜索结果字典
        """
        try:
            import time
            start_time = time.time()
            
            # 执行全文搜索
            hits, total = await search_manager.search(
                query_text=query,
                owner_id=user_id,
                file_type=file_type,
                limit=limit,
                offset=offset
            )
            
            # 获取完整的文件信息
            results = []
            for hit in hits:
                file_id = int(hit['file_id'])
                file = db.query(FileModel).filter(
                    FileModel.id == file_id,
                    FileModel.owner_id == user_id
                ).first()
                
                if file:
                    results.append({
                        'id': file.id,
                        'filename': file.filename,
                        'original_filename': file.original_filename,
                        'file_size': file.file_size,
                        'mime_type': file.mime_type,
                        'created_at': file.created_at.isoformat() if file.created_at else None,
                        'is_shared': file.is_shared,
                        'description': file.description,
                        'score': hit.get('score', 0),
                        'highlights': hit.get('highlights', {})
                    })
            
            search_time = (time.time() - start_time) * 1000  # 转换为毫秒
            
            return {
                'results': results,
                'total': total,
                'query': query,
                'search_time_ms': round(search_time, 2),
                'limit': limit,
                'offset': offset
            }
            
        except Exception as e:
            logger.error(f"Wenxi - 搜索文件失败: {e}")
            raise HTTPException(status_code=500, detail=f"搜索失败: {str(e)}")
    
    @staticmethod
    async def reindex_all_files(db: Session) -> Dict:
        """
        重建所有文件的索引
        
        Args:
            db: 数据库会话
            
        Returns:
            重建结果
        """
        try:
            # 获取所有文件
            files = db.query(FileModel).all()
            
            # 构建文件数据列表
            files_data = []
            for file in files:
                files_data.append({
                    'file_id': file.id,
                    'filename': file.filename,
                    'original_filename': file.original_filename,
                    'description': file.description or '',
                    'owner_id': file.owner_id,
                    'file_size': file.file_size,
                    'created_at': file.created_at,
                    'updated_at': file.updated_at,
                    'mime_type': file.mime_type or 'application/octet-stream',
                    'content': ''
                })
            
            # 重建索引
            success = await search_manager.reindex_all(files_data)
            
            if success:
                logger.info(f"Wenxi - 索引重建完成: {len(files_data)} 个文件")
                return {
                    'success': True,
                    'indexed_count': len(files_data),
                    'message': f'成功索引 {len(files_data)} 个文件'
                }
            else:
                return {
                    'success': False,
                    'message': '索引重建失败'
                }
                
        except Exception as e:
            logger.error(f"Wenxi - 重建索引失败: {e}")
            return {
                'success': False,
                'message': f'重建索引失败: {str(e)}'
            }
    
    @staticmethod
    async def get_search_stats() -> Dict:
        """
        获取搜索索引统计信息
        
        Returns:
            统计信息字典
        """
        return await search_manager.get_stats()
    
    @staticmethod
    def get_file_type_suggestions() -> List[Dict]:
        """
        获取文件类型建议列表
        
        Returns:
            文件类型列表
        """
        return [
            {'value': 'text', 'label': '文本文档', 'extensions': ['.txt', '.md']},
            {'value': 'document', 'label': '办公文档', 'extensions': ['.doc', '.docx', '.pdf']},
            {'value': 'spreadsheet', 'label': '电子表格', 'extensions': ['.xls', '.xlsx']},
            {'value': 'presentation', 'label': '演示文稿', 'extensions': ['.ppt', '.pptx']},
            {'value': 'image', 'label': '图片', 'extensions': ['.jpg', '.jpeg', '.png', '.gif']},
            {'value': 'video', 'label': '视频', 'extensions': ['.mp4', '.avi', '.mov']},
            {'value': 'audio', 'label': '音频', 'extensions': ['.mp3', '.wav']},
            {'value': 'archive', 'label': '压缩文件', 'extensions': ['.zip', '.rar', '.tar', '.gz']},
            {'value': 'code', 'label': '代码文件', 'extensions': ['.py', '.js', '.html', '.css', '.java']},
        ]


# 全局搜索服务实例
search_service = SearchService()