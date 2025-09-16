import os
import secrets
from typing import List, Dict
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    # ------------------------
    # Project
    # ------------------------
    PROJECT_NAME: str = "ML Scheduler Service"
    VERSION: str = "1.0.0"

    # ------------------------
    # MLflow
    # ------------------------
    MLFLOW_URI: str = os.getenv("MLFLOW_URI", "http://mlflow:5000")

    # ------------------------
    # MinIO / S3
    # ------------------------
    MINIO_ENDPOINT: str = os.getenv("MINIO_ENDPOINT", "minio:9000")
    MINIO_ACCESS_KEY: str = os.getenv("MINIO_ACCESS_KEY", "minioadmin")
    MINIO_SECRET_KEY: str = os.getenv("MINIO_SECRET_KEY", "minioadmin")
    MINIO_BUCKET: str = os.getenv("MINIO_BUCKET", "models")

    # ------------------------
    # Docker / Environment
    # ------------------------
    DOCKER_ENV: bool = os.getenv("DOCKER_ENV", "false") == "true"

    # ------------------------
    # Training Defaults
    # ------------------------
    BATCH_SIZE: int = 32
    MAX_LEN: int = 256
    NUM_EPOCHS: int = 3
    USE_CUDA: bool = False

    # ------------------------
    # Models Configuration
    # ------------------------
    # You can add more models here as you scale
    MODELS: Dict[str, Dict] = {
        "sentiment": {
            "base_model": "w11wo/indonesian-roberta-base-sentiment-classifier",
            "tokenizer_path": "models/sentiment/tokenizer",
            "onnx_path": "models/sentiment/model.onnx",
        },
        # "placeholder": {
        #     "base_model": "bert-base-uncased",
        #     "tokenizer_path": "models/placeholder/tokenizer",
        #     "onnx_path": "models/placeholder/model.onnx",
        # },
    }

    # ------------------------
    # ETL / Scheduler
    # ------------------------
    TRAIN_INTERVAL_MINUTES: int = 240

    # ------------------------
    # Logging
    # ------------------------
    LOG_LEVEL: str = "INFO"

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"


# Singleton instance
settings = Settings()
