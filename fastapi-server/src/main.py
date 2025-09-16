import logging
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager
import asyncio

from src.config.settings import settings
from src.config.database import engine
from src.routers.router import api_router
from src.routers.predict import predict_router
from src.logging import LoggingMiddleware
from src.core.exceptions import setup_exception_handlers
from src.core.grpc import grpc_server
from src.services.sentiment.manager import SentimentManager

sentiment_manager = SentimentManager()


@asynccontextmanager
async def lifespan(app: FastAPI):

    # --- Startup ---
    try:
        await grpc_server.start()
        logging.info("gRPC server started")

        await sentiment_manager.load_model()
        consumer = sentiment_manager.get_consumer()
        await consumer.start()
        logging.info("Sentiment analysis consumer started")

        yield  # FastAPI runs the app here

    # --- SHUTDOWN ---
    finally:
        # Stop consumer
        try:
            await consumer.stop()
            logging.info("Sentiment analysis consumer stopped")
        except asyncio.CancelledError:
            logging.warning("Consumer shutdown cancelled, proceeding anyway")
        except Exception as e:
            logging.error(f"Error stopping consumer: {e}")

        # Stop gRPC server with grace period
        try:
            await grpc_server.stop(grace=5)
            logging.info("gRPC server stopped")
        except asyncio.CancelledError:
            logging.warning("gRPC shutdown cancelled, proceeding anyway")
        except Exception as e:
            logging.error(f"Error stopping gRPC server: {e}")


def create_application() -> FastAPI:
    app = FastAPI(
        title=settings.PROJECT_NAME,
        description=settings.PROJECT_DESCRIPTION,
        version=settings.VERSION,
        openapi_url=f"{settings.API_V1_STR}/openapi.json",
        docs_url="/docs",
        redoc_url="/redoc",
        lifespan=lifespan,
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
    app.include_router(predict_router, prefix=settings.API_V1_STR)

    return app


app = create_application()

if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "src.main:app", host="0.0.0.0", port=8000, reload=True, log_level="info"
    )
