import onnxruntime as ort
from transformers import AutoTokenizer
import numpy as np


class SentimentModel:
    def __init__(self, onnx_model_path: str, tokenizer_path: str):
        """
        Load ONNX model and tokenizer for inference
        """
        # Initialize ONNX runtime session
        self.session = ort.InferenceSession(
            onnx_model_path, providers=["CPUExecutionProvider"]
        )

        # Load tokenizer
        self.tokenizer = AutoTokenizer.from_pretrained(tokenizer_path)

    def predict_sentiment(self, text: str) -> float:
        """
        Predict toxicity (or sentiment) using ONNX model
        """
        # Tokenize input
        inputs = self.tokenizer(
            text,
            return_tensors="np",  # return numpy arrays for ONNX
            truncation=True,
            padding="max_length",
            max_length=256,
        )

        # Prepare inputs for ONNXRuntime
        ort_inputs = {k: inputs[k] for k in ["input_ids", "attention_mask"]}

        # Run ONNX inference
        logits = self.session.run(None, ort_inputs)[
            0
        ]  # shape: (batch_size, num_labels)

        # Convert logits to probabilities
        # probs = np.exp(logits) / np.sum(np.exp(logits), axis=1, keepdims=True)  # type: ignore

        # shifted_logits = logits - np.max(logits, axis=1, keepdims=True)
        shifted_logits = logits - np.max(logits, axis=1, keepdims=True)  # type: ignore
        probs = np.exp(shifted_logits) / np.sum(
            np.exp(shifted_logits), axis=1, keepdims=True
        )

        return float(probs[0][2])
