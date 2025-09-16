from pathlib import Path
import secrets
from typing import List, Optional
from pydantic import field_validator
from pydantic_settings import BaseSettings

BASE_DIR = Path(__file__).resolve().parent.parent.parent


class Settings(BaseSettings):
    # ------------------------
    # Project
    # ------------------------
    PROJECT_NAME: str = "FastAPI ML Server"
    PROJECT_DESCRIPTION: str = "FastAPI server with ML inference capabilities"
    VERSION: str = "1.0.0"
    API_V1_STR: str = "/api/v1"

    # ------------------------
    # Security
    # ------------------------
    SECRET_KEY: str = secrets.token_urlsafe(32)
    ACCESS_TOKEN_EXPIRE_MINUTES: int = 60 * 24 * 8  # 8 days
    ALGORITHM: str = "HS256"

    # ------------------------
    # CORS
    # ------------------------
    ALLOWED_HOSTS: List[str] = ["*"]

    # ------------------------
    # Database
    # ------------------------
    POSTGRES_HOST: str = "localhost"
    POSTGRES_USER: str = "dev"
    POSTGRES_PASSWORD: str = "dev"
    POSTGRES_DB: str = "fastapi_db"
    POSTGRES_PORT: str = "5432"

    # ------------------------
    # Docker detection
    # ------------------------
    DOCKER_ENV: str = "false"  # can be overridden by env

    @property
    def DATABASE_URL(self) -> str:
        host = "postgres" if self.DOCKER_ENV == "true" else self.POSTGRES_HOST
        return f"postgresql://{self.POSTGRES_USER}:{self.POSTGRES_PASSWORD}@{host}:{self.POSTGRES_PORT}/{self.POSTGRES_DB}"

    @property
    def MLFLOW_URI(self) -> str:
        if self.DOCKER_ENV == "true":
            return "http://mlflow:5000"
        return "http://127.0.0.1:5000"

    # ------------------------
    # Redis (optional)
    # ------------------------
    REDIS_URL: Optional[str] = "redis://localhost:6379"

    # ------------------------
    # ML Models Sentiment
    # ------------------------
    # MODEL_NAME: str = "w11wo/indonesian-roberta-base-sentiment-classifier"
    MODEL_SENTIMENT_NAME: str = "indonesian_roberta_sentiment.onnx"

    @property
    def MODEL_SENTIMENT_PATH(self) -> str:
        return str(BASE_DIR / "models" / "sentiment" / self.MODEL_SENTIMENT_NAME)

    MODEL_SENTIMENT_TOKENIZER_PATH: str = str(BASE_DIR / "models" / "sentiment")
    MAX_MODEL_SIZE: int = 100 * 1024 * 1024  # 100MB
    BATCH_SIZE: int = 100
    TRAIN_INTERVAL_MINUTES: int = 240
    USE_CUDA: bool = False

    # ------------------------
    # RabbitMQ
    # ------------------------
    RABBITMQ_URL: str = "amqp://guest:guest@localhost:5672/"
    RABBITMQ_CONSUME_QUEUE: str = "posts"
    RABBITMQ_TOXIC_QUEUE: str = "toxic_posts"
    TOXICITY_THRESHOLD: float = 0.7

    @field_validator("RABBITMQ_URL", mode="before")
    @classmethod
    def set_rabbitmq_url(cls, v, info):
        # If env provides a value, use it
        if v:
            return v
        # Otherwise, auto switch based on DOCKER_ENV
        docker_env = info.data.get("DOCKER_ENV", "false") == "true"
        return (
            "amqp://guest:guest@rabbitmq:5672/"
            if docker_env
            else "amqp://guest:guest@localhost:5672/"
        )

    # ------------------------
    # Logging
    # ------------------------
    LOG_LEVEL: str = "INFO"

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"
        extra = "allow"  # prevents extra_forbidden errors


# Singleton instance
settings = Settings()
