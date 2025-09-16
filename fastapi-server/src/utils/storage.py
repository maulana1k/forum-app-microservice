# app/utils/storage.py
import os
import json
import shutil

def save_model_version(local_path: str, version: str):
    model_dir = os.path.join(local_path, version)
    os.makedirs(model_dir, exist_ok=True)
    # Assume model already saved to model_dir
    meta = {"version": version}
    with open(os.path.join(local_path, "latest.json"), "w") as f:
        json.dump(meta, f)
