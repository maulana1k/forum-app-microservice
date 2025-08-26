import os
from pydantic_settings import BaseSettings
from typing import List, Optional
import secrets


class Settings(BaseSettings):
    # Project
    PROJECT_NAME: str = "FastAPI ML Server"
    PROJECT_DESCRIPTION: str = "FastAPI server with ML inference capabilities"
    VERSION: str = "1.0.0"
    API_V1_STR: str = "/api/v1"
    
    # Security
    SECRET_KEY: str = secrets.token_urlsafe(32)
    ACCESS_TOKEN_EXPIRE_MINUTES: int = 60 * 24 * 8  # 8 days
    ALGORITHM: str = "HS256"
    
    # CORS
    ALLOWED_HOSTS: List[str] = ["*"]
    
    # Database
    POSTGRES_HOST: str = "localhost"
    POSTGRES_USER: str = "dev"
    POSTGRES_PASSWORD: str = "dev"
    POSTGRES_DB: str = "fastapi_db"
    POSTGRES_PORT: str = "5432"
    
    DOCKER_ENV: str = "false"
    
    @property
    def DATABASE_URL(self) -> str:
        host = "postgres" if self.DOCKER_ENV else self.POSTGRES_HOST
        return f"postgresql://{self.POSTGRES_USER}:{self.POSTGRES_PASSWORD}@{host}:{self.POSTGRES_PORT}/{self.POSTGRES_DB}"
    
    # Redis (optional for caching)
    REDIS_URL: Optional[str] = "redis://localhost:6379"
    
    # ML Models
    MODEL_PATH: str = "./dataset/models/"
    MAX_MODEL_SIZE: int = 100 * 1024 * 1024  # 100MB
    
    # Logging
    LOG_LEVEL: str = "INFO"
    
    class Config:
        env_file = ".env"


settings = Settings()