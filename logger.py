import logging
import sys

# 配置日志格式
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(sys.stdout),
        logging.FileHandler('app.log', encoding='utf-8')
    ]
)

# 创建全局logger实例
logger = logging.getLogger("wenxi-disk")
logger.setLevel(logging.INFO)

# 导出logger供其他模块使用
__all__ = ['logger']