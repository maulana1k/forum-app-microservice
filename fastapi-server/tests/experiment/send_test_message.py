import asyncio
import json
from aio_pika import connect_robust, Message, DeliveryMode


async def send_test_message():
    connection = await connect_robust("amqp://guest:guest@localhost:5672/")
    channel = await connection.channel()
    queue_name = "posts"  # match your settings

    message = {"id": "test-1", "content": "This is a toxic test message"}
    await channel.default_exchange.publish(
        Message(json.dumps(message).encode(), delivery_mode=DeliveryMode.PERSISTENT),
        routing_key=queue_name,
    )
    print("Test message sent")
    await connection.close()


asyncio.run(send_test_message())
