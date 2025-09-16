import logging
import asyncio
from typing import Optional
from src.config.settings import settings
from .model import SentimentModel
from .consumer import SentimentConsumer

logger = logging.getLogger(__name__)


class SentimentManager:
    def __init__(self):
        self._model: Optional[SentimentModel] = None
        self._consumer: Optional[SentimentConsumer] = None
        self._lock = asyncio.Lock()

    async def load_model(
        self,
        onnx_model_path: Optional[str] = None,
        tokenizer_path: Optional[str] = None,
    ):
        """
        Load or reload the sentiment analysis ONNX model.
        """
        async with self._lock:
            onnx_model_path = onnx_model_path or settings.MODEL_SENTIMENT_PATH
            tokenizer_path = tokenizer_path or settings.MODEL_SENTIMENT_TOKENIZER_PATH
            logger.info(f"Loading sentiment model from ONNX: {onnx_model_path}")

            # Offload CPU-bound model loading to a thread
            self._model = await asyncio.to_thread(
                SentimentModel, onnx_model_path, tokenizer_path
            )

            # Stop existing consumer if any
            if self._consumer:
                logger.info("Stopping existing consumer for reload")
                await self._consumer.stop()
                self._consumer = None

    def get_model(self) -> SentimentModel:
        if not self._model:
            raise RuntimeError(
                "Sentiment model is not loaded. Call load_model() first."
            )
        return self._model

    def get_consumer(self) -> SentimentConsumer:
        if not self._model:
            raise RuntimeError("Cannot start consumer: sentiment model not loaded.")
        if not self._consumer:
            self._consumer = SentimentConsumer(self._model)
        return self._consumer
