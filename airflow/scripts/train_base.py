import os
import mlflow
from minio import Minio
from config.settings import settings

OUTPUT_DIR = "/tmp/model"
os.makedirs(OUTPUT_DIR, exist_ok=True)


def preprocess_dataset(df, model_name: str):
    """Placeholder: preprocess dataset for each model."""
    # TODO: implement model-specific preprocessing
    return df


def get_model_and_tokenizer(model_name: str):
    """Initialize model and tokenizer for a given model name."""
    model_config = settings.MODELS.get(model_name)
    if not model_config:
        raise ValueError(f"Model config not found for {model_name}")

    from transformers import AutoModelForSequenceClassification, AutoTokenizer

    tokenizer = AutoTokenizer.from_pretrained(model_config["base_model"])
    model = AutoModelForSequenceClassification.from_pretrained(
        model_config["base_model"], num_labels=2
    )
    return model, tokenizer


def train_model(model, tokenizer, dataset, output_dir: str):
    """Placeholder for training logic."""
    # TODO: implement HuggingFace Trainer or PyTorch Lightning training
    pass


def export_onnx(model, tokenizer, output_dir: str, model_name: str):
    """Export ONNX model."""
    import torch

    dummy_input = tokenizer("contoh teks", return_tensors="pt")
    torch.onnx.export(
        model,
        (dummy_input["input_ids"], dummy_input["attention_mask"]),
        os.path.join(output_dir, f"{model_name}.onnx"),
        input_names=["input_ids", "attention_mask"],
        output_names=["logits"],
        dynamic_axes={"input_ids": {0: "batch"}, "attention_mask": {0: "batch"}},
    )


def upload_to_minio(output_dir: str, model_name: str):
    """Upload model artifacts to MinIO."""
    client = Minio(
        settings.MINIO_ENDPOINT,
        access_key=settings.MINIO_ACCESS_KEY,
        secret_key=settings.MINIO_SECRET_KEY,
        secure=False,
    )
    if not client.bucket_exists(settings.MINIO_BUCKET):
        client.make_bucket(settings.MINIO_BUCKET)
    for f in os.listdir(output_dir):
        client.fput_object(
            settings.MINIO_BUCKET, f"{model_name}/{f}", os.path.join(output_dir, f)
        )


def log_mlflow(output_dir: str, model_name: str, params: dict = None):
    """Log artifacts and params to MLflow."""
    mlflow.set_tracking_uri(settings.MLFLOW_URI)
    mlflow.set_experiment("batch_ml_training")
    with mlflow.start_run(run_name=model_name):
        if params:
            mlflow.log_params(params)
        mlflow.log_artifacts(output_dir)
