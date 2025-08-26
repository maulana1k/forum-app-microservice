# src/ml/models/predictor.py
import asyncio
import joblib
import os

class MLModelManager:
    _model = None
    _model_path = "models/model.pkl"

    @classmethod
    async def load_model(cls, model_path: str = _model_path):
        """Async load ML model"""
        if not os.path.exists(cls._model_path):
            raise FileNotFoundError(f"Model not found at {cls._model_path}")

        cls._model = joblib.load(cls._model_path)
        print("[MLModelManager] Model loaded.")

    @classmethod
    async def cleanup(cls):
        """Async cleanup resources if needed"""
        await asyncio.sleep(0.05)
        cls._model = None
        print("[MLModelManager] Model cleaned up.")

    @classmethod
    def predict(cls, input_data):
        if cls._model is None:
            raise ValueError("Model not loaded. Call load_model() first.")
        return cls._model.predict([input_data]).tolist()
