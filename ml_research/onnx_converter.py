import torch
import logging
from typing import Dict, Any
from transformers import PreTrainedTokenizer
from training_config import TrainingConfig

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class ONNXExporter:
    """
    A class to export a PyTorch model to ONNX format.

    Args:
        model (torch.nn.Module): The PyTorch model to export.
        tokenizer (PreTrainedTokenizer): Tokenizer associated with the model.
        config (TrainingConfig): Configuration object containing output directory and other settings.
    """

    def __init__(self, model: torch.nn.Module, tokenizer: PreTrainedTokenizer, config: TrainingConfig = TrainingConfig()):
        self.model = model
        self.tokenizer = tokenizer
        self.config = config

    def export(self, dummy_text: str = "Пример текста") -> None:
        """
        Export the model to ONNX format.

        Args:
            dummy_text (str): Example text to trace the model's forward pass. Defaults to "Пример текста".
        """
        class WrappedModel(torch.nn.Module):
            """
            A wrapper class to adapt the model's forward pass for ONNX export.
            """
            def __init__(self, model: torch.nn.Module):
                super().__init__()
                self.model = model

            def forward(self, input_ids: torch.Tensor, attention_mask: torch.Tensor) -> torch.Tensor:
                """
                Forward pass for ONNX export.

                Args:
                    input_ids (torch.Tensor): Tokenized input IDs.
                    attention_mask (torch.Tensor): Attention mask.

                Returns:
                    torch.Tensor: Model logits.
                """
                return self.model(input_ids, attention_mask).logits

        try:
            # Wrap the model for ONNX export
            wrapped_model = WrappedModel(self.model).eval()

            # Tokenize dummy text
            inputs = self.tokenizer(dummy_text, return_tensors="pt")

            # Export the model to ONNX
            torch.onnx.export(
                wrapped_model,
                (inputs["input_ids"], inputs["attention_mask"]),
                f"{self.config.output_dir}/model.onnx",
                input_names=['input_ids', 'attention_mask'],
                output_names=['logits'],
                dynamic_axes={
                    'input_ids': {0: 'batch_size', 1: 'sequence_length'},
                    'attention_mask': {0: 'batch_size', 1: 'sequence_length'},
                    'logits': {0: 'batch_size'}
                },
                opset_version=14
            )

            logger.info(f"Model successfully exported to {self.config.output_dir}/model.onnx")

        except Exception as e:
            logger.error(f"Failed to export model to ONNX: {e}")
            raise