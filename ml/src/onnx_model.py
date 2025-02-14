import onnxruntime
import numpy as np
from typing import List
from transformers import AutoTokenizer


class ONNXClassifier:
    def __init__(self, onnx_path: str, tokenizer_name: str, max_length: int = 128):
        self.tokenizer = AutoTokenizer.from_pretrained(tokenizer_name)
        self.max_length = max_length
        self.session = onnxruntime.InferenceSession(onnx_path, providers=['CPUExecutionProvider'])

    def predict(self, texts: List[str]) -> List[int]:
        inputs = self.tokenizer(
            texts,
            padding=True,
            truncation=True,
            max_length=self.max_length,
            return_tensors="np"
        )
        input_ids = inputs['input_ids'].astype(np.int64)
        attention_mask = inputs['attention_mask'].astype(np.int64)
        logits = self.session.run(
            None,
            {
                'input_ids': input_ids,
                'attention_mask': attention_mask
            }
        )[0]
        class_ids = np.argmax(logits, axis=1).tolist()
        return class_ids