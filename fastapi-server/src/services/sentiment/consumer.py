import asyncio
import json
import logging

from src.core.rabbitmq import RabbitMQManager
from .model import SentimentModel
from src.config.settings import settings

logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)


class SentimentConsumer:
    def __init__(self, model: SentimentModel):
        self.model = model
        self.queue = None
        self.consumer_tag = None
        self.consume_queue = "post-create"
        self.sentiment_post_queue = "post-sentiment"
        self.sentiment_threshold = 0.8

    async def start(self):
        async def handler(message):
            async with message.process(requeue=False):
                try:
                    payload = json.loads(message.body.decode())
                    post_id = payload.get("post_id")
                    content = payload.get("content", "")
                    if not post_id:
                        logger.warning("Message missing 'id', skipping.")
                        return

                    score = await asyncio.to_thread(
                        self.model.predict_sentiment, content
                    )
                    flagged = score > self.sentiment_threshold

                    logger.info(f"Sentiment score for post {post_id}: {score}")
                    if flagged:
                        await RabbitMQManager.publish(
                            {
                                "post_id": post_id,
                                "content": content,
                                "score": score,
                            },
                            routing_key=self.sentiment_post_queue,
                        )
                except Exception as e:
                    logger.exception(f"Failed to process message: {e}")

        self.queue, self.consumer_tag = await RabbitMQManager.consume(
            self.consume_queue, handler, no_ack=False
        )

    async def stop(self):
        if self.queue and self.consumer_tag:
            await self.queue.cancel(self.consumer_tag)
        await RabbitMQManager.close()
