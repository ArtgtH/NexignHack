import datetime
import json
import logging
import os

from onnx_model import ONNXClassifier


logger = logging.getLogger(os.getenv("WORKER_ID")+".task_service")


class TaskService:
    def __init__(self, chanel_in, task_in, redis):
        self._channel_in = chanel_in
        self._task_queue = task_in
        self._redis = redis
        self._predictor = ONNXClassifier(
            onnx_path="model.onnx", tokenizer_name="sergeyzh/rubert-tiny-turbo"
        )

        self._channel_in.queue_declare(queue=self._task_queue, passive=True)

    def start(self):
        self._channel_in.basic_consume(
            queue=self._task_queue, on_message_callback=self._callback, auto_ack=True
        )
        self._channel_in.start_consuming()

    def _callback(self, ch, method, properties, body):
        logger.info(f"Получили таску {datetime.datetime.now()}")
        string = body.decode("utf-8")
        data = json.loads(string)
        texts = [x["messageText"] for x in data["messages"]]

        logger.info(f"Старт МЛ {datetime.datetime.now()}")
        results = self._predictor.predict(texts)
        logger.info(f"Финиш МЛ {datetime.datetime.now()}")

        for idx, res in enumerate(results):
            data["messages"][idx]["result"] = res

        self._publish(data)

    def _publish(self, message):
        self._redis.set(message["id"], json.dumps(message))
        logger.info(f"Отправили таску {datetime.datetime.now()}")
