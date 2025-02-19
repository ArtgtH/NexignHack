from training_config import TrainingConfig

import os
import torch
import pandas as pd
from transformers import AdamW, get_linear_schedule_with_warmup, AutoModelForSequenceClassification
from tqdm import tqdm
from sklearn.metrics import accuracy_score, f1_score
from sklearn.utils.class_weight import compute_class_weight
from torch.utils.data import Dataset
import matplotlib.pyplot as plt
from typing import List, Dict, Any, Optional, Tuple

import logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class SentimentDataset(Dataset):
    """
    A PyTorch Dataset class for sentiment analysis tasks.

    Args:
        encodings (Dict[str, torch.Tensor]): Tokenized input encodings (e.g., input_ids, attention_mask).
        labels (List[int]): List of labels corresponding to the input encodings.
    """

    def __init__(self, encodings: Dict[str, torch.Tensor], labels: List[int]):
        self.encodings = encodings
        self.labels = torch.tensor(labels, dtype=torch.long)

    def __getitem__(self, idx: int) -> Dict[str, torch.Tensor]:
        return {
            "input_ids": self.encodings["input_ids"][idx],
            "attention_mask": self.encodings["attention_mask"][idx],
            "labels": self.labels[idx]
        }

    def __len__(self) -> int:
        return len(self.labels)


class ModelTrainer:
    """
    A class to train and evaluate a sentiment analysis model.

    Args:
        model (AutoModelForSequenceClassification): Pre-trained model for sequence classification.
        train_loader (torch.utils.data.DataLoader): DataLoader for the training dataset.
        val_loader (torch.utils.data.DataLoader): DataLoader for the validation dataset.
        config (TrainingConfig): Configuration object containing hyperparameters and settings.
    """

    def __init__(self, model: AutoModelForSequenceClassification, train_loader: torch.utils.data.DataLoader,
                 val_loader: torch.utils.data.DataLoader, config: TrainingConfig = TrainingConfig()):
        self.model = model.to(config.device)
        self.train_loader = train_loader
        self.val_loader = val_loader
        self.config = config
        self.config.create_output_dir()
        self.optimizer = AdamW(self.model.parameters(), lr=config.learning_rate, eps=config.eps)
        self.scheduler = self._create_scheduler()
        self.best_val_loss = float('inf')
        self.metrics: List[Dict[str, Any]] = []
        self.class_weights = self._calculate_class_weights()
        if self.class_weights is not None:
            self.criterion = torch.nn.CrossEntropyLoss(weight=self.class_weights.to(config.device))
        else:
            self.criterion = None

    def _calculate_class_weights(self) -> Optional[torch.Tensor]:
        """
        Calculate class weights for imbalanced datasets.

        Returns:
            Optional[torch.Tensor]: Tensor of class weights if labels are provided, otherwise None.
        """
        labels = []
        for batch in self.train_loader:
            labels.extend(batch["labels"].tolist())
        if labels:
            class_weights = compute_class_weight("balanced", classes=torch.unique(torch.tensor(labels)).tolist(), y=labels)
            return torch.tensor(class_weights, dtype=torch.float32)
        return None

    def _create_scheduler(self):
        total_steps = len(self.train_loader) * self.config.epochs
        warmup_steps = int(total_steps * self.config.warmup_ratio)
        return get_linear_schedule_with_warmup(
            self.optimizer,
            num_warmup_steps=warmup_steps,
            num_training_steps=total_steps
        )

    def train_epoch(self, epoch: int) -> float:
        """
        Train the model for one epoch.

        Args:
            epoch (int): Current epoch number.

        Returns:
            float: Average training loss for the epoch.
        """
        self.model.train()
        epoch_loss = 0.0
        progress_bar = tqdm(self.train_loader, desc=f"Epoch {epoch + 1}/{self.config.epochs} [Train]")

        for batch in progress_bar:
            self.optimizer.zero_grad()
            inputs = self._prepare_inputs(batch)
            outputs = self.model(**inputs)
            logits = outputs.logits
            labels = inputs["labels"]
            if self.criterion is not None:
                loss = self.criterion(logits, labels)
            else:
                loss = outputs.loss
            loss.backward()
            torch.nn.utils.clip_grad_norm_(self.model.parameters(), self.config.grad_clip)
            self.optimizer.step()
            self.scheduler.step()

            epoch_loss += loss.item()
            progress_bar.set_postfix({"loss": loss.item()})

        return epoch_loss / len(self.train_loader)

    def validate(self) -> Tuple[float, float, float]:
        """
        Validate the model on the validation dataset.

        Returns:
            Tuple[float, float, float]: Average validation loss, accuracy, and F1 score.
        """
        self.model.eval()
        epoch_loss = 0.0
        all_preds = []
        all_labels = []

        with torch.no_grad():
            for batch in tqdm(self.val_loader, desc="Validating"):
                inputs = self._prepare_inputs(batch)
                outputs = self.model(**inputs)
                logits = outputs.logits
                labels = inputs["labels"]
                if self.criterion is not None:
                    loss = self.criterion(logits, labels)
                else:
                    loss = outputs.loss
                epoch_loss += loss.item()

                preds = torch.argmax(logits, dim=1).cpu().numpy()
                all_preds.extend(preds)
                all_labels.extend(labels.cpu().numpy())

        avg_loss = epoch_loss / len(self.val_loader)
        accuracy = accuracy_score(all_labels, all_preds)
        f1 = f1_score(all_labels, all_preds, average="weighted")

        return avg_loss, accuracy, f1

    def _prepare_inputs(self, batch: Dict[str, torch.Tensor]) -> Dict[str, torch.Tensor]:
        return {
            "input_ids": batch["input_ids"].to(self.config.device),
            "attention_mask": batch["attention_mask"].to(self.config.device),
            "labels": batch["labels"].to(self.config.device)
        }

    def train(self) -> List[Dict[str, Any]]:
        """
        Train the model for the specified number of epochs.

        Returns:
            List[Dict[str, Any]]: List of metrics for each epoch.
        """
        for epoch in range(self.config.epochs):
            train_loss = self.train_epoch(epoch)
            val_loss, val_acc, val_f1 = self.validate()

            self.metrics.append({
                "epoch": epoch + 1,
                "train_loss": train_loss,
                "val_loss": val_loss,
                "val_accuracy": val_acc,
                "val_f1": val_f1
            })

            logger.info(f"\nEpoch {epoch + 1}/{self.config.epochs}")
            logger.info(f"Train Loss: {train_loss:.4f} | Val Loss: {val_loss:.4f}")
            logger.info(f"Val Accuracy: {val_acc:.4f} | Val F1: {val_f1:.4f}\n")

            if val_loss < self.best_val_loss:
                self.best_val_loss = val_loss
                self.save_model()

        self.plot_metrics()
        self.save_metrics_to_excel()
        return self.metrics

    def save_model(self):
        torch.save(self.model.state_dict(), f"{self.config.output_dir}/best_model.pt")
        logger.info(f"Model saved to {self.config.output_dir}/best_model.pt")

    def plot_metrics(self):
        """
        Plot training and validation metrics.
        """
        try:
            plt.style.use('ggplot')
        except ImportError:
            logger.warning("Matplotlib is required to plot metrics. Please install it.")
            return

        epochs = [m['epoch'] for m in self.metrics]
        train_loss = [m['train_loss'] for m in self.metrics]
        val_loss = [m['val_loss'] for m in self.metrics]
        val_f1 = [m['val_f1'] for m in self.metrics]

        fig, ax1 = plt.subplots(figsize=(10, 6))

        color = 'tab:blue'
        ax1.set_xlabel('Epoch')
        ax1.set_ylabel('Loss', color=color)
        ax1.plot(epochs, train_loss, label='Train Loss', color=color, linestyle='-', linewidth=2)
        ax1.plot(epochs, val_loss, label='Val Loss', color='tab:orange', linestyle='-', linewidth=2)
        ax1.tick_params(axis='y', labelcolor=color)
        ax1.legend(loc='upper left')

        ax2 = ax1.twinx()
        color = 'tab:green'
        ax2.set_ylabel('F1 Score', color=color)
        ax2.plot(epochs, val_f1, label='Val F1', color=color, linestyle='--', linewidth=2)
        ax2.tick_params(axis='y', labelcolor=color)
        ax2.legend(loc='upper right')

        plt.title('Training Loss, Validation Loss, and Validation F1 Score')
        fig.tight_layout()
        plt.savefig(os.path.join(self.config.output_dir, 'loss_f1_plot.png'), bbox_inches='tight', dpi=300)
        plt.close()

    def save_metrics_to_excel(self):
        """
        Save training/validation metrics to an Excel file.
        """
        df = pd.DataFrame(self.metrics)
        df.to_excel(os.path.join(self.config.output_dir, "training_metrics.xlsx"), index=False)
        logger.info(f"Metrics saved to {os.path.join(self.config.output_dir, 'training_metrics.xlsx')}")