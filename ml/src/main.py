import logging
import os

import pika
import redis

from task_service import TaskService


logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    handlers=[logging.StreamHandler()],
)


logger = logging.getLogger(os.getenv("WORKER_ID"))


def main():
    logger.info("Starting RabbitMQ Task Service")

    rabbit = os.getenv("RABBITMQ_URL")
    rabbit_q = os.getenv("RABBITMQ_TASK_QUEUE")
    rdb = os.getenv("REDIS_URL")

    connection = pika.BlockingConnection(pika.URLParameters(rabbit))
    channel_in = connection.channel()
    r = redis.Redis.from_url(rdb)

    task_svc = TaskService(channel_in, rabbit_q, r)
    task_svc.start()


if __name__ == "__main__":
    main()
