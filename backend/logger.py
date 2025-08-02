"""
Wenxi网盘 - 高级日志系统
作者：Wenxi
功能：提供分级日志系统，支持DEBUG/INFO/WARNING/ERROR级别
环境变量：WENXI_LOG_LEVEL 控制日志级别（默认INFO）
"""

import os
import logging
import sys
from datetime import datetime
from typing import Optional


class WenxiLogger:
    """Wenxi网盘专用日志管理器"""
    
    def __init__(self, name: str = "wenxi-netdisk"):
        self.logger = logging.getLogger(name)
        self._setup_logger()
    
    def _setup_logger(self):
        """配置日志格式和级别"""
        # 从环境变量获取日志级别
        log_level = os.getenv("WENXI_LOG_LEVEL", "INFO").upper()
        level_map = {
            "DEBUG": logging.DEBUG,
            "INFO": logging.INFO,
            "WARNING": logging.WARNING,
            "ERROR": logging.ERROR
        }
        
        self.logger.setLevel(level_map.get(log_level, logging.INFO))
        
        # 避免重复添加处理器
        if not self.logger.handlers:
            # 控制台处理器
            console_handler = logging.StreamHandler(sys.stdout)
            console_handler.setLevel(level_map.get(log_level, logging.INFO))
            
            # 详细格式
            formatter = logging.Formatter(
                '[%(asctime)s] [%(levelname)s] [%(name)s:%(lineno)d] %(message)s',
                datefmt='%Y-%m-%d %H:%M:%S'
            )
            console_handler.setFormatter(formatter)
            
            self.logger.addHandler(console_handler)
    
    def debug(self, message: str, *args, **kwargs):
        """调试信息"""
        self.logger.debug(message, *args, **kwargs)
    
    def info(self, message: str, *args, **kwargs):
        """一般信息"""
        self.logger.info(message, *args, **kwargs)
    
    def warning(self, message: str, *args, **kwargs):
        """警告信息"""
        self.logger.warning(message, *args, **kwargs)
    
    def error(self, message: str, *args, **kwargs):
        """错误信息"""
        self.logger.error(message, *args, **kwargs)
    
    def exception(self, message: str, *args, **kwargs):
        """异常信息"""
        self.logger.exception(message, *args, **kwargs)


# 全局日志实例
logger = WenxiLogger()


def get_logger(name: str) -> WenxiLogger:
    """获取指定名称的日志器"""
    return WenxiLogger(name)


if __name__ == "__main__":
    # 测试日志系统
    test_logger = get_logger("test")
    test_logger.debug("这是调试信息")
    test_logger.info("这是普通信息")
    test_logger.warning("这是警告信息")
    test_logger.error("这是错误信息")
    try:
        1/0
    except Exception:
        test_logger.exception("发生异常")