from etl_base import fetch_data_from_db, clean_text, save_preprocessed
from config.settings import settings
import os

MODEL_NAME = "sentiment"
OUTPUT_PATH = f"/tmp/{MODEL_NAME}_dataset.parquet"


def run_sentiment_etl():
    """
    ETL pipeline for sentiment model data.
    """
    # --------------------------
    # 1. Fetch raw data
    # --------------------------
    query = "SELECT text, label FROM posts WHERE created_at > NOW() - INTERVAL '1 day';"
    df = fetch_data_from_db(query, connection_string=settings.DATABASE_URL)

    # --------------------------
    # 2. Clean / preprocess
    # --------------------------
    df = clean_text(df, "text")

    # --------------------------
    # 3. Save preprocessed batch
    # --------------------------
    os.makedirs(os.path.dirname(OUTPUT_PATH), exist_ok=True)
    save_preprocessed(df, OUTPUT_PATH)

    return OUTPUT_PATH
