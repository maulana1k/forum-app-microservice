from airflow import DAG
from airflow.operators.python import PythonOperator
from datetime import datetime
import subprocess
import logging

logger = logging.getLogger(__name__)
MODEL_NAME = "sentiment"


def run_etl():
    subprocess.run(["python3", f"/opt/airflow/scripts/etl_{MODEL_NAME}.py"], check=True)


def run_training():
    subprocess.run(
        ["python3", f"/opt/airflow/scripts/train_{MODEL_NAME}.py"], check=True
    )


with DAG(
    dag_id=f"dag_{MODEL_NAME}",
    start_date=datetime(2025, 9, 1),
    schedule_interval="@daily",
    catchup=False,
    tags=["ml", "training"],
    max_active_runs=1,
) as dag:

    etl_task = PythonOperator(
        task_id=f"etl_{MODEL_NAME}", python_callable=run_etl, retries=1
    )

    train_task = PythonOperator(
        task_id=f"train_{MODEL_NAME}", python_callable=run_training, retries=1
    )

    etl_task >> train_task  # type: ignore
