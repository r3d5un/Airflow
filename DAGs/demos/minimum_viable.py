from datetime import datetime, timedelta

from airflow import DAG
from airflow.operators.empty import EmptyOperator

with DAG(
    dag_id="minimum_viable",
    schedule_interval=None,
    catchup=False,
    tags=["demo", "minimum_viable", "tutorial"],
    default_args={
        "owner": "oyvind.kristiansen@dfo.no",
        "retries": 5,
        "retries_delay": timedelta(seconds=60),
        "start_date": datetime(2022, 9, 22),
        "depends_on_past": False,
    },
) as dag:
    start = EmptyOperator(task_id="start")
    end = EmptyOperator(task_id="end")

    start >> end  # type: ignore
