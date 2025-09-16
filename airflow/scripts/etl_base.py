import pandas as pd
import logging

logger = logging.getLogger(__name__)


def fetch_data_from_db(query: str, connection_string: str) -> pd.DataFrame:
    """
    Fetch raw data from DB.
    """
    import psycopg2
    import sqlalchemy

    engine = sqlalchemy.create_engine(connection_string)
    try:
        df = pd.read_sql(query, engine)
    except Exception as e:
        logger.error(f"Error fetching data: {e}")
        raise
    return df


def clean_text(df: pd.DataFrame, text_column: str) -> pd.DataFrame:
    """
    Basic text cleaning: lowercasing, removing extra spaces, etc.
    """
    df[text_column] = df[text_column].str.lower().str.strip()
    return df


def save_preprocessed(df: pd.DataFrame, path: str):
    """
    Save preprocessed dataset as CSV or parquet.
    """
    df.to_parquet(path, index=False)
    logger.info(f"Preprocessed data saved to {path}")
