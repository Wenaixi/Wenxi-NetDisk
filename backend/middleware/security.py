"""
Wenxi网盘 - 安全中间件和性能监控模块
作者：Wenxi
功能：提供API限流、安全头、性能监控
"""

import time
import functools
from typing import Callable, Optional
from fastapi import Request, Response, HTTPException
from fastapi.middleware.base import BaseHTTPMiddleware
from starlette.middleware.base import RequestResponseEndpoint
import redis.asyncio as redis
import logging

logger = logging.getLogger(__name__)


class RateLimiter:
    """
    基于Redis的API限流器
    使用滑动窗口算法
    """
    def __init__(
        self,
        redis_client: redis.Redis,
        requests_per_minute: int = 60,
        burst_size: int = 10
    ):
        self.redis = redis_client
        self.requests_per_minute = requests_per_minute
        self.burst_size = burst_size
        self.window_size = 60  # 60秒窗口
    
    async def is_allowed(self, key: str) -> bool:
        """
        检查请求是否被允许
        
        Args:
            key: 限流键（通常是IP+端点）
        
        Returns:
            True if allowed, False if rate limited
        """
        now = time.time()
        window_start = now - self.window_size
        
        # 使用Redis sorted set存储请求时间戳
        pipe = self.redis.pipeline()
        
        # 移除窗口外的旧记录
        pipe.zremrangebyscore(key, 0, window_start)
        
        # 获取当前窗口内的请求数
        pipe.zcard(key)
        
        # 添加当前请求
        pipe.zadd(key, {str(now): now})
        
        # 设置过期时间
        pipe.expire(key, self.window_size)
        
        results = await pipe.execute()
        current_requests = results[1]
        
        # 检查是否超过限制
        return current_requests < self.requests_per_minute
    
    async def get_remaining(self, key: str) -> dict:
        """获取限流状态信息"""
        now = time.time()
        window_start = now - self.window_size
        
        pipe = self.redis.pipeline()
        pipe.zremrangebyscore(key, 0, window_start)
        pipe.zcard(key)
        pipe.ttl(key)
        
        results = await pipe.execute()
        current = results[1]
        ttl = results[2]
        
        return {
            "limit": self.requests_per_minute,
            "remaining": max(0, self.requests_per_minute - current),
            "reset": int(now + ttl) if ttl > 0 else int(now + self.window_size),
            "window": "60"
        }


class SecurityHeadersMiddleware(BaseHTTPMiddleware):
    """
    安全响应头中间件
    添加OWASP推荐的安全头
    """
    
    async def dispatch(
        self,
        request: Request,
        call_next: RequestResponseEndpoint
    ) -> Response:
        response = await call_next(request)
        
        # 防止点击劫持
        response.headers["X-Frame-Options"] = "DENY"
        
        # XSS保护
        response.headers["X-Content-Type-Options"] = "nosniff"
        
        # HSTS (HTTPS强制)
        response.headers["Strict-Transport-Security"] = "max-age=31536000; includeSubDomains"
        
        # CSP内容安全策略
        response.headers["Content-Security-Policy"] = (
            "default-src 'self'; "
            "script-src 'self' 'unsafe-inline' 'unsafe-eval'; "
            "style-src 'self' 'unsafe-inline'; "
            "img-src 'self' data: https:; "
            "font-src 'self'; "
            "connect-src 'self';"
        )
        
        # 禁用浏览器特性
        response.headers["Permissions-Policy"] = (
            "camera=(), microphone=(), geolocation=(), "
            "payment=(), usb=(), magnetometer=(), gyroscope=()"
        )
        
        # Referrer策略
        response.headers["Referrer-Policy"] = "strict-origin-when-cross-origin"
        
        return response


class RateLimitMiddleware(BaseHTTPMiddleware):
    """
    限流中间件
    基于IP地址进行限流
    """
    
    def __init__(
        self,
        app,
        redis_client: Optional[redis.Redis] = None,
        requests_per_minute: int = 60
    ):
        super().__init__(app)
        self.limiter = RateLimiter(redis_client, requests_per_minute) if redis_client else None
    
    async def dispatch(
        self,
        request: Request,
        call_next: RequestResponseEndpoint
    ) -> Response:
        if not self.limiter:
            return await call_next(request)
        
        # 获取客户端IP
        client_ip = self._get_client_ip(request)
        endpoint = request.url.path
        key = f"rate_limit:{client_ip}:{endpoint}"
        
        # 检查限流
        is_allowed = await self.limiter.is_allowed(key)
        
        if not is_allowed:
            raise HTTPException(
                status_code=429,
                detail="Too many requests. Please try again later."
            )
        
        # 获取限流信息
        rate_info = await self.limiter.get_remaining(key)
        
        # 处理请求
        response = await call_next(request)
        
        # 添加限流响应头
        response.headers["X-RateLimit-Limit"] = str(rate_info["limit"])
        response.headers["X-RateLimit-Remaining"] = str(rate_info["remaining"])
        response.headers["X-RateLimit-Reset"] = str(rate_info["reset"])
        
        return response
    
    def _get_client_ip(self, request: Request) -> str:
        """获取客户端真实IP"""
        # 优先从X-Forwarded-For获取（如果有代理）
        forwarded = request.headers.get("X-Forwarded-For")
        if forwarded:
            return forwarded.split(",")[0].strip()
        
        # 从X-Real-IP获取
        real_ip = request.headers.get("X-Real-IP")
        if real_ip:
            return real_ip
        
        # 直接连接IP
        if request.client:
            return request.client.host
        
        return "unknown"


class PerformanceMonitor:
    """
    性能监控器
    记录请求处理时间和系统指标
    """
    
    def __init__(self):
        self.request_times: list = []
        self.error_count: int = 0
        self.request_count: int = 0
        self.max_history = 1000
    
    def record_request(self, duration: float, status_code: int):
        """记录请求指标"""
        self.request_count += 1
        
        if status_code >= 400:
            self.error_count += 1
        
        self.request_times.append({
            "duration": duration,
            "status_code": status_code,
            "timestamp": time.time()
        })
        
        # 限制历史记录大小
        if len(self.request_times) > self.max_history:
            self.request_times = self.request_times[-self.max_history:]
    
    def get_stats(self) -> dict:
        """获取性能统计"""
        if not self.request_times:
            return {
                "total_requests": 0,
                "error_rate": 0.0,
                "avg_response_time": 0.0,
                "p95_response_time": 0.0,
                "p99_response_time": 0.0
            }
        
        durations = [r["duration"] for r in self.request_times]
        durations.sort()
        
        total = len(durations)
        p95_idx = int(total * 0.95)
        p99_idx = int(total * 0.99)
        
        return {
            "total_requests": self.request_count,
            "recent_requests": total,
            "error_rate": self.error_count / max(self.request_count, 1),
            "avg_response_time": sum(durations) / len(durations),
            "p95_response_time": durations[min(p95_idx, total - 1)],
            "p99_response_time": durations[min(p99_idx, total - 1)]
        }


class PerformanceMiddleware(BaseHTTPMiddleware):
    """
    性能监控中间件
    记录每个请求的处理时间
    """
    
    def __init__(self, app):
        super().__init__(app)
        self.monitor = PerformanceMonitor()
    
    async def dispatch(
        self,
        request: Request,
        call_next: RequestResponseEndpoint
    ) -> Response:
        start_time = time.time()
        
        try:
            response = await call_next(request)
            status_code = response.status_code
        except Exception as e:
            status_code = 500
            raise
        finally:
            duration = time.time() - start_time
            self.monitor.record_request(duration, status_code)
            
            # 记录慢请求
            if duration > 5.0:  # 超过5秒认为是慢请求
                logger.warning(
                    f"Slow request: {request.method} {request.url.path} "
                    f"took {duration:.2f}s"
                )
        
        # 添加性能响应头
        response.headers["X-Response-Time"] = f"{duration:.3f}s"
        
        return response
    
    def get_stats(self) -> dict:
        """获取性能统计"""
        return self.monitor.get_stats()


def require_auth(roles: Optional[list] = None):
    """
    权限验证装饰器
    
    Args:
        roles: 允许访问的角色列表
    """
    def decorator(func: Callable):
        @functools.wraps(func)
        async def wrapper(*args, **kwargs):
            # 获取request对象
            request = kwargs.get('request')
            if not request and args:
                for arg in args:
                    if isinstance(arg, Request):
                        request = arg
                        break
            
            if not request:
                raise HTTPException(status_code=500, detail="Request object not found")
            
            # 检查认证（这里需要与现有的auth系统集成）
            # 简化版本，实际使用时需要导入get_current_user
            
            return await func(*args, **kwargs)
        return wrapper
    return decorator


# 便捷函数
def setup_security_middleware(app, redis_client: Optional[redis.Redis] = None):
    """
    设置所有安全中间件
    
    Args:
        app: FastAPI应用实例
        redis_client: Redis客户端（用于限流）
    """
    # 安全头中间件
    app.add_middleware(SecurityHeadersMiddleware)
    
    # 性能监控中间件
    app.add_middleware(PerformanceMiddleware)
    
    # 限流中间件（如果有Redis）
    if redis_client:
        app.add_middleware(
            RateLimitMiddleware,
            redis_client=redis_client,
            requests_per_minute=60
        )
    
    logger.info("安全中间件已配置")


# 全局性能监控器实例
performance_monitor = PerformanceMonitor()
