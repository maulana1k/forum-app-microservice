from starlette.middleware.base import BaseHTTPMiddleware
from starlette.requests import Request
import logging
import time

# Setup clean logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s:\t %(levelname)s - %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)
logger = logging.getLogger()

class LoggingMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        # Skip docs and health checks
        if request.url.path in ['/docs', '/redoc', '/openapi.json', '/health']:
            return await call_next(request)
            
        start_time = time.time()
        response = await call_next(request)
        process_time = (time.time() - start_time) * 1000
        
        logger.info(
            f"{request.method} {request.url.path} "
            f"Status: {response.status_code} "
            f"Duration: {process_time:.2f} ms"
        )
        return response