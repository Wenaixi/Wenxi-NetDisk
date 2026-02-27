"""
Wenxi网盘 - 智能搜索模块
作者：Wenxi
功能：提供全文搜索功能
"""

from search.index_manager import search_manager, extract_text_content
from search.service import search_service

__all__ = ['search_manager', 'search_service', 'extract_text_content']