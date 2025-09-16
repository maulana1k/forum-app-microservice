from airflow import DAG
from airflow.operators.python import PythonOperator
from datetime import datetime
import subprocess
from config.settings import settings
import logging

logger = logging.getLogger(__name__)


def run_etl(model_name: str):
    """
    Run ETL pipeline for a specific model.
    """
    script_path = f"/opt/airflow/scripts/etl_{model_name}.py"
    logger.info(f"Running ETL for {model_name} using {script_path}")
    subprocess.run(["python3", script_path], check=True)


def run_training(model_name: str):
    """
    Run training script for a specific model.
    """
    script_path = f"/opt/airflow/scripts/train_{model_name}.py"
    logger.info(f"Running training for {model_name} using {script_path}")
    subprocess.run(["python3", script_path], check=True)


# --------------------------
# DAG Definition
# --------------------------
with DAG(
    dag_id="batch_train_models",
    start_date=datetime(2025, 9, 1),
    schedule_interval="@daily",
    catchup=False,
    tags=["ml", "training"],
    max_active_runs=1,
) as dag:

    for model_name in settings.MODELS.keys():
        # ETL task
        etl_task = PythonOperator(
            task_id=f"etl_{model_name}",
            python_callable=lambda name=model_name: run_etl(name),
            retries=1,
        )

        # Training task
        train_task = PythonOperator(
            task_id=f"train_{model_name}",
            python_callable=lambda name=model_name: run_training(name),
            retries=1,
        )

        # ETL -> Training dependency
        etl_task >> train_task  # type: ignore
