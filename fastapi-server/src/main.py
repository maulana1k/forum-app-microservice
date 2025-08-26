from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager

from src.config.settings import settings
from src.config.database import engine
from src.api.v1.api import api_router
from src.api.middleware.logging import LoggingMiddleware
from src.core.exceptions import setup_exception_handlers
from src.ml.models.predictor import MLModelManager


# @asynccontextmanager
# async def lifespan(app: FastAPI):
#     # Startup
#     await MLModelManager.load_model("models/model.pkl")
#     yield
#     # Shutdown
#     await MLModelManager.cleanup()


def create_application() -> FastAPI:
    app = FastAPI(
        title=settings.PROJECT_NAME,
        description=settings.PROJECT_DESCRIPTION,
        version=settings.VERSION,
        openapi_url=f"{settings.API_V1_STR}/openapi.json",
        docs_url="/docs",                     # Swagger UI
        redoc_url="/redoc"  
        # lifespan=lifespan,
    )

    # Set up CORS
    app.add_middleware(
        CORSMiddleware,
        allow_origins=settings.ALLOWED_HOSTS,
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    # Custom middleware
    app.add_middleware(LoggingMiddleware)

    # Exception handlers
    setup_exception_handlers(app)

    # Include routers
    app.include_router(api_router, prefix=settings.API_V1_STR)

    return app


app = create_application()


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "src.main:app",
        host="0.0.0.0",
        port=8000,
        reload=True,
        log_level="info"
    )