import os
import torch
import logging
from dataclasses import dataclass
from typing import Optional

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

@dataclass
class TrainingConfig:
    batch_size: int = 16
    epochs: int = 10
    learning_rate: float = 2e-5
    eps: float = 1e-8
    warmup_ratio: float = 0.1
    model_name: str = "sergeyzh/rubert-tiny-turbo"
    device: torch.device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
    output_dir: str = "./output"
    grad_clip: float = 1.0

    def create_output_dir(self) -> None:
        """Create the output directory if it does not exist."""
        if not os.path.exists(self.output_dir):
            os.makedirs(self.output_dir)
            logger.info(f"Created output directory at {self.output_dir}")