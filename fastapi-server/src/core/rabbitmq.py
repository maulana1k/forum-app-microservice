import asyncio
import json
import logging
from typing import Callable, Optional
import aio_pika
from aio_pika import Message, DeliveryMode
from src.config.settings import settings

logger = logging.getLogger(__name__)


class RabbitMQManager:
    _connection: Optional[aio_pika.abc.AbstractRobustConnection] = None
    _channel: Optional[aio_pika.abc.AbstractChannel] = None

    @classmethod
    async def connect(cls):
        if cls._connection is None or cls._connection.is_closed:
            logger.info("Connecting to RabbitMQ...")
            cls._connection = await aio_pika.connect_robust(settings.RABBITMQ_URL)
            cls._channel = await cls._connection.channel()
            logger.info("RabbitMQ connected and channel created.")
        return cls._connection, cls._channel

    @classmethod
    async def close(cls):
        if cls._connection and not cls._connection.is_closed:
            await cls._connection.close()
            cls._connection = None
            cls._channel = None
            logger.info("RabbitMQ connection closed.")

    @classmethod
    async def declare_queue(cls, queue_name: str, durable=True):
        await cls.connect()
        return await cls._channel.declare_queue(queue_name, durable=durable)  # type: ignore

    @classmethod
    async def declare_exchange(cls, exchange_name: str, type="direct", durable=True):
        await cls.connect()
        return await cls._channel.declare_exchange(exchange_name, type=type, durable=durable)  # type: ignore

    @classmethod
    async def publish(cls, body: dict, routing_key: str):
        await cls.connect()
        queue = await cls.declare_queue(routing_key)
        exchange = await cls.declare_exchange(routing_key + "-exchange")  # type: ignore
        await queue.bind(exchange=exchange, routing_key=routing_key)
        message = Message(
            body=json.dumps(body).encode(), delivery_mode=DeliveryMode.PERSISTENT
        )
        await exchange.publish(message, routing_key=routing_key)
        logger.info(f"Published message to {routing_key}: {body}")

    @classmethod
    async def consume(cls, queue_name: str, handler: Callable, no_ack=False):
        queue = await cls.declare_queue(queue_name)
        consumer_tag = await queue.consume(handler, no_ack=no_ack)
        logger.info(f"Consumer started on queue {queue_name}")
        return queue, consumer_tag
