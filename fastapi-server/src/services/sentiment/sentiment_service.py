# app/services/toxicity_service.py
from transformers import AutoTokenizer, AutoModelForSequenceClassification, pipeline
import threading
import os
import json
from typing import List, Dict
from src.core.config import MODEL_BASE_PATH

class ModelManager:
    def __init__(self):
        self.model = None
        self.tokenizer = None
        self.pipeline = None
        self.version = None
        self.lock = threading.Lock()
        self.load_latest_model()

    def load_latest_model(self):
        meta_path = os.path.join(MODEL_BASE_PATH, "latest.json")
        if not os.path.exists(meta_path):
            raise FileNotFoundError("No model metadata found")
        with open(meta_path) as f:
            meta = json.load(f)
        self.version = meta["version"]
        model_path = os.path.join(MODEL_BASE_PATH, self.version)
        self.tokenizer = AutoTokenizer.from_pretrained(model_path)
        self.model = AutoModelForSequenceClassification.from_pretrained(model_path)
        self.pipeline = pipeline("text-classification", model=self.model, tokenizer=self.tokenizer)
        print(f"Loaded model version {self.version}")

    def update_model(self, version: str):
        model_path = os.path.join(MODEL_BASE_PATH, version)
        tokenizer = AutoTokenizer.from_pretrained(model_path)
        model = AutoModelForSequenceClassification.from_pretrained(model_path)
        pipeline_obj = pipeline("text-classification", model=model, tokenizer=tokenizer)
        with self.lock:
            self.model = model
            self.tokenizer = tokenizer
            self.pipeline = pipeline_obj
            self.version = version
        print(f"Model swapped to version {version} successfully")

    # def predict(self, texts: List[str]) -> List[Dict]:
    #     with self.lock:
    #         return self.pipeline(texts)
