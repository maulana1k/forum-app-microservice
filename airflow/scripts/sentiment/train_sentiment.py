import pandas as pd
from train_base import (
    preprocess_dataset,
    get_model_and_tokenizer,
    train_model,
    export_onnx,
    upload_to_minio,
    log_mlflow,
)

MODEL_NAME = "sentiment"

# --------------------------
# Load dataset (batch example)
# --------------------------
df = pd.DataFrame(
    {"text": ["contoh teks positif", "contoh teks negatif"], "label": [0, 1]}
)

# --------------------------
# Preprocess dataset
# --------------------------
dataset = preprocess_dataset(df, MODEL_NAME)

# --------------------------
# Initialize model/tokenizer
# --------------------------
model, tokenizer = get_model_and_tokenizer(MODEL_NAME)

# --------------------------
# Train model
# --------------------------
train_model(model, tokenizer, dataset, "/tmp/model")

# --------------------------
# Export ONNX
# --------------------------
export_onnx(model, tokenizer, "/tmp/model", MODEL_NAME)

# --------------------------
# Upload to MinIO
# --------------------------
upload_to_minio("/tmp/model", MODEL_NAME)

# --------------------------
# Log to MLflow
# --------------------------
log_mlflow("/tmp/model", MODEL_NAME, params={"model_name": MODEL_NAME})
