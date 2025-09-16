import asyncio
import json
from aio_pika import connect_robust, Message, DeliveryMode


async def send_bulk_messages(num_messages: int = 100):
    connection = await connect_robust("amqp://guest:guest@localhost:5672/")
    channel = await connection.channel()
    queue_name = "posts"  # make sure it matches your queue

    for i in range(1, num_messages + 1):
        message = {"id": f"{i}", "content": f"This is test message #{i}"}
        await channel.default_exchange.publish(
            Message(
                json.dumps(message).encode(), delivery_mode=DeliveryMode.PERSISTENT
            ),
            routing_key=queue_name,
        )
        if i % 100 == 0:
            print(f"{i} messages sent")

    print(f"All {num_messages} messages sent")
    await connection.close()


if __name__ == "__main__":
    asyncio.run(send_bulk_messages())  # adjust number as needed
