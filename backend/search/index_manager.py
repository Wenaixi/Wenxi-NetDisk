"""
Wenxi网盘 - 全文搜索引擎模块
作者：Wenxi
功能：基于Whoosh实现文件全文索引和搜索
性能优化：异步索引更新、分词优化、缓存机制
"""

import os
import re
import asyncio
from datetime import datetime
from typing import List, Dict, Optional, Tuple
from concurrent.futures import ThreadPoolExecutor

from whoosh import index
from whoosh.fields import Schema, TEXT, ID, DATETIME, NUMERIC, KEYWORD
from whoosh.qparser import MultifieldParser, QueryParser
from whoosh.query import Term, And, Or
from whoosh.analysis import StandardAnalyzer, RegexTokenizer, LowercaseFilter, StopFilter
from whoosh.sorting import FieldFacet, ScoreFacet

from logger import logger

# 全局配置
INDEX_DIR = os.path.join(os.path.dirname(__file__), "..", "search_index")
executor = ThreadPoolExecutor(max_workers=4)

# 自定义中文分析器
class ChineseAnalyzer:
    """支持中英文混合分词的分析器"""
    
    def __init__(self):
        # 使用正则分词，支持中英文
        self.tokenizer = RegexTokenizer(r"\w+|[^\w\s]")
        
    def __call__(self, value):
        """分词处理"""
        if not value:
            return
        
        # 转为小写
        value = value.lower()
        
        # 中文字符分割（每2-3个字符作为一个词）
        chinese_chars = re.findall(r'[\u4e00-\u9fff]', value)
        
        # 英文单词分割
        english_words = re.findall(r'[a-zA-Z0-9_]+', value)
        
        # 数字分割
        numbers = re.findall(r'\d+', value)
        
        # 生成位置信息
        pos = 0
        all_tokens = []
        
        # 处理英文单词
        for word in english_words:
            all_tokens.append((word, pos))
            pos += 1
            
        # 处理中文（2-gram）
        chinese_text = ''.join(chinese_chars)
        for i in range(len(chinese_text) - 1):
            bigram = chinese_text[i:i+2]
            all_tokens.append((bigram, pos))
            pos += 1
            
        # 处理数字
        for num in numbers:
            all_tokens.append((num, pos))
            pos += 1
        
        # 去重并排序
        seen = set()
        for token, p in all_tokens:
            if token not in seen:
                seen.add(token)
                yield token


# 定义Whoosh索引Schema
FILE_SCHEMA = Schema(
    file_id=ID(stored=True, unique=True),  # 文件ID（唯一）
    filename=TEXT(stored=True, analyzer=StandardAnalyzer()),  # 文件名
    original_filename=TEXT(stored=True, analyzer=StandardAnalyzer()),  # 原始文件名
    content=TEXT(stored=True, analyzer=StandardAnalyzer()),  # 文件内容（文本文件）
    description=TEXT(stored=True, analyzer=StandardAnalyzer()),  # 文件描述
    file_type=KEYWORD(stored=True, lowercase=True),  # 文件类型
    owner_id=ID(stored=True),  # 所有者ID
    file_size=NUMERIC(stored=True),  # 文件大小
    created_at=DATETIME(stored=True),  # 创建时间
    updated_at=DATETIME(stored=True),  # 更新时间
    mime_type=ID(stored=True),  # MIME类型
)


class SearchIndexManager:
    """
    Wenxi - 全文索引管理器
    功能：管理文件索引的创建、更新、删除和搜索
    """
    
    _instance = None
    _lock = asyncio.Lock()
    
    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance._initialized = False
        return cls._instance
    
    def __init__(self):
        if self._initialized:
            return
            
        self._initialized = True
        self.index_dir = os.path.abspath(INDEX_DIR)
        self._ix = None
        self._ensure_index_exists()
    
    def _ensure_index_exists(self):
        """确保索引目录和索引存在"""
        try:
            if not os.path.exists(self.index_dir):
                os.makedirs(self.index_dir)
                logger.info(f"Wenxi - 创建搜索索引目录: {self.index_dir}")
            
            if not index.exists_in(self.index_dir):
                self._ix = index.create_in(self.index_dir, FILE_SCHEMA)
                logger.info("Wenxi - 创建新的全文搜索索引")
            else:
                self._ix = index.open_dir(self.index_dir)
                logger.info("Wenxi - 打开现有全文搜索索引")
                
        except Exception as e:
            logger.error(f"Wenxi - 搜索索引初始化失败: {e}")
            raise
    
    @property
    def index(self):
        """获取索引对象"""
        if self._ix is None:
            self._ensure_index_exists()
        return self._ix
    
    async def add_file(self, file_data: Dict) -> bool:
        """
        异步添加文件到索引
        
        Args:
            file_data: 文件信息字典，包含:
                - file_id: 文件ID
                - filename: 文件名
                - original_filename: 原始文件名
                - content: 文件内容（可选）
                - description: 描述（可选）
                - file_type: 文件类型
                - owner_id: 所有者ID
                - file_size: 文件大小
                - created_at: 创建时间
                - updated_at: 更新时间
                - mime_type: MIME类型
        """
        try:
            loop = asyncio.get_event_loop()
            return await loop.run_in_executor(executor, self._add_file_sync, file_data)
        except Exception as e:
            logger.error(f"Wenxi - 添加文件到索引失败: {e}")
            return False
    
    def _add_file_sync(self, file_data: Dict) -> bool:
        """同步添加文件到索引（内部方法）"""
        try:
            writer = self.index.writer()
            
            # 提取文件类型
            file_type = self._extract_file_type(file_data.get('original_filename', ''))
            
            writer.add_document(
                file_id=str(file_data['file_id']),
                filename=file_data.get('filename', ''),
                original_filename=file_data.get('original_filename', ''),
                content=file_data.get('content', ''),
                description=file_data.get('description', ''),
                file_type=file_type,
                owner_id=str(file_data['owner_id']),
                file_size=file_data.get('file_size', 0),
                created_at=file_data.get('created_at', datetime.utcnow()),
                updated_at=file_data.get('updated_at', datetime.utcnow()),
                mime_type=file_data.get('mime_type', 'application/octet-stream'),
            )
            
            writer.commit()
            logger.debug(f"Wenxi - 文件添加到索引: {file_data.get('original_filename')}")
            return True
            
        except Exception as e:
            logger.error(f"Wenxi - 同步添加文件到索引失败: {e}")
            return False
    
    async def update_file(self, file_data: Dict) -> bool:
        """异步更新文件索引"""
        try:
            # 先删除旧索引，再添加新索引
            await self.delete_file(str(file_data['file_id']))
            return await self.add_file(file_data)
        except Exception as e:
            logger.error(f"Wenxi - 更新文件索引失败: {e}")
            return False
    
    async def delete_file(self, file_id: str) -> bool:
        """异步从索引中删除文件"""
        try:
            loop = asyncio.get_event_loop()
            return await loop.run_in_executor(executor, self._delete_file_sync, file_id)
        except Exception as e:
            logger.error(f"Wenxi - 从索引删除文件失败: {e}")
            return False
    
    def _delete_file_sync(self, file_id: str) -> bool:
        """同步删除文件索引（内部方法）"""
        try:
            writer = self.index.writer()
            writer.delete_by_term('file_id', file_id)
            writer.commit()
            logger.debug(f"Wenxi - 从索引删除文件: {file_id}")
            return True
        except Exception as e:
            logger.error(f"Wenxi - 同步删除文件索引失败: {e}")
            return False
    
    async def search(
        self,
        query_text: str,
        owner_id: Optional[int] = None,
        file_type: Optional[str] = None,
        limit: int = 20,
        offset: int = 0
    ) -> Tuple[List[Dict], int]:
        """
        执行全文搜索
        
        Args:
            query_text: 搜索关键词
            owner_id: 所有者ID过滤（可选）
            file_type: 文件类型过滤（可选）
            limit: 返回结果数量限制
            offset: 结果偏移量
            
        Returns:
            (搜索结果列表, 总匹配数)
        """
        try:
            loop = asyncio.get_event_loop()
            return await loop.run_in_executor(
                executor, 
                self._search_sync, 
                query_text, 
                owner_id, 
                file_type, 
                limit, 
                offset
            )
        except Exception as e:
            logger.error(f"Wenxi - 搜索执行失败: {e}")
            return [], 0
    
    def _search_sync(
        self,
        query_text: str,
        owner_id: Optional[int],
        file_type: Optional[str],
        limit: int,
        offset: int
    ) -> Tuple[List[Dict], int]:
        """同步搜索（内部方法）"""
        try:
            with self.index.searcher() as searcher:
                # 构建查询
                # 在多个字段中搜索
                fields = ['original_filename', 'filename', 'content', 'description']
                parser = MultifieldParser(fields, schema=self.index.schema)
                
                # 解析查询文本
                query = parser.parse(query_text)
                
                # 添加过滤器
                filters = []
                if owner_id:
                    filters.append(Term('owner_id', str(owner_id)))
                if file_type:
                    filters.append(Term('file_type', file_type.lower()))
                
                if filters:
                    query = And([query] + filters)
                
                # 执行搜索，按相关性排序
                results = searcher.search(
                    query,
                    limit=limit + offset,
                    sortedby=ScoreFacet()
                )
                
                total = len(results)
                
                # 转换结果为字典列表
                hits = []
                for hit in results[offset:offset + limit]:
                    hits.append({
                        'file_id': hit['file_id'],
                        'filename': hit['filename'],
                        'original_filename': hit['original_filename'],
                        'file_type': hit['file_type'],
                        'file_size': hit['file_size'],
                        'created_at': hit['created_at'].isoformat() if hit['created_at'] else None,
                        'mime_type': hit['mime_type'],
                        'score': hit.score,
                        'highlights': self._get_highlights(hit, query_text)
                    })
                
                return hits, total
                
        except Exception as e:
            logger.error(f"Wenxi - 同步搜索失败: {e}")
            return [], 0
    
    def _get_highlights(self, hit, query_text: str) -> Dict[str, str]:
        """获取搜索高亮片段"""
        highlights = {}
        
        try:
            # 高亮文件名
            if 'original_filename' in hit:
                highlights['filename'] = self._highlight_text(
                    hit['original_filename'], 
                    query_text
                )
            
            # 高亮内容片段
            if hit.get('content'):
                highlights['content'] = self._extract_snippet(
                    hit['content'],
                    query_text
                )
                
        except Exception as e:
            logger.debug(f"生成高亮失败: {e}")
            
        return highlights
    
    def _highlight_text(self, text: str, query: str) -> str:
        """高亮匹配文本"""
        if not text or not query:
            return text
            
        # 简单的文本高亮
        pattern = re.compile(re.escape(query), re.IGNORECASE)
        return pattern.sub(lambda m: f"**{m.group()}**", text)
    
    def _extract_snippet(self, content: str, query: str, max_length: int = 200) -> str:
        """提取包含关键词的文本片段"""
        if not content:
            return ""
            
        # 查找关键词位置
        pattern = re.compile(re.escape(query), re.IGNORECASE)
        match = pattern.search(content)
        
        if match:
            start = max(0, match.start() - 50)
            end = min(len(content), match.end() + 50)
            snippet = content[start:end]
            
            if start > 0:
                snippet = "..." + snippet
            if end < len(content):
                snippet = snippet + "..."
                
            return self._highlight_text(snippet, query)
        
        # 如果没有匹配，返回前max_length个字符
        return content[:max_length] + ("..." if len(content) > max_length else "")
    
    def _extract_file_type(self, filename: str) -> str:
        """从文件名提取文件类型"""
        if not filename:
            return "unknown"
            
        ext = os.path.splitext(filename)[1].lower()
        
        # 文件类型映射
        type_mapping = {
            '.txt': 'text',
            '.md': 'text',
            '.doc': 'document',
            '.docx': 'document',
            '.pdf': 'document',
            '.xls': 'spreadsheet',
            '.xlsx': 'spreadsheet',
            '.ppt': 'presentation',
            '.pptx': 'presentation',
            '.jpg': 'image',
            '.jpeg': 'image',
            '.png': 'image',
            '.gif': 'image',
            '.mp4': 'video',
            '.avi': 'video',
            '.mov': 'video',
            '.mp3': 'audio',
            '.wav': 'audio',
            '.zip': 'archive',
            '.rar': 'archive',
            '.tar': 'archive',
            '.gz': 'archive',
            '.py': 'code',
            '.js': 'code',
            '.html': 'code',
            '.css': 'code',
            '.java': 'code',
            '.cpp': 'code',
            '.c': 'code',
        }
        
        return type_mapping.get(ext, 'other')
    
    async def reindex_all(self, files_data: List[Dict]) -> bool:
        """
        重建所有索引
        
        Args:
            files_data: 所有文件数据的列表
        """
        try:
            loop = asyncio.get_event_loop()
            return await loop.run_in_executor(executor, self._reindex_all_sync, files_data)
        except Exception as e:
            logger.error(f"Wenxi - 重建索引失败: {e}")
            return False
    
    def _reindex_all_sync(self, files_data: List[Dict]) -> bool:
        """同步重建索引（内部方法）"""
        try:
            # 清除现有索引
            import shutil
            if os.path.exists(self.index_dir):
                shutil.rmtree(self.index_dir)
            
            # 重新创建索引
            self._ix = None
            self._ensure_index_exists()
            
            # 批量添加文档
            writer = self.index.writer()
            
            for file_data in files_data:
                file_type = self._extract_file_type(file_data.get('original_filename', ''))
                
                writer.add_document(
                    file_id=str(file_data['file_id']),
                    filename=file_data.get('filename', ''),
                    original_filename=file_data.get('original_filename', ''),
                    content=file_data.get('content', ''),
                    description=file_data.get('description', ''),
                    file_type=file_type,
                    owner_id=str(file_data['owner_id']),
                    file_size=file_data.get('file_size', 0),
                    created_at=file_data.get('created_at', datetime.utcnow()),
                    updated_at=file_data.get('updated_at', datetime.utcnow()),
                    mime_type=file_data.get('mime_type', 'application/octet-stream'),
                )
            
            writer.commit()
            logger.info(f"Wenxi - 索引重建完成，共 {len(files_data)} 个文件")
            return True
            
        except Exception as e:
            logger.error(f"Wenxi - 同步重建索引失败: {e}")
            return False
    
    async def get_stats(self) -> Dict:
        """获取索引统计信息"""
        try:
            with self.index.searcher() as searcher:
                return {
                    'total_documents': searcher.doc_count(),
                    'index_dir': self.index_dir,
                    'index_size_mb': self._get_index_size()
                }
        except Exception as e:
            logger.error(f"Wenxi - 获取索引统计失败: {e}")
            return {'total_documents': 0, 'error': str(e)}
    
    def _get_index_size(self) -> float:
        """获取索引目录大小（MB）"""
        try:
            total_size = 0
            for dirpath, dirnames, filenames in os.walk(self.index_dir):
                for f in filenames:
                    fp = os.path.join(dirpath, f)
                    total_size += os.path.getsize(fp)
            return round(total_size / (1024 * 1024), 2)
        except:
            return 0.0


# 全局搜索管理器实例
search_manager = SearchIndexManager()


async def extract_text_content(file_path: str, mime_type: str) -> str:
    """
    从文件中提取文本内容（用于索引）
    
    Args:
        file_path: 文件路径
        mime_type: 文件MIME类型
        
    Returns:
        提取的文本内容
    """
    try:
        # 文本文件
        if mime_type and mime_type.startswith('text/'):
            async with aiofiles.open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
                content = await f.read()
                # 限制索引内容长度
                return content[:100000]  # 最多索引10万字符
        
        # 其他类型文件暂不提取内容
        return ""
        
    except Exception as e:
        logger.debug(f"提取文件内容失败: {e}")
        return ""


# 为了消除循环导入，延迟导入
import aiofiles