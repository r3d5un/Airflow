from datetime import datetime, timedelta

from airflow import DAG
from airflow.operators.bash import BashOperator

with DAG(
    dag_id="catchup_demo",
    schedule_interval="*/15 * * * *",
    catchup=True,
    tags=["demo", "backfill_demo", "tutorial"],
    default_args={
        "owner": "oyvind.kristiansen@dfo.no",
        "retries": 5,
        "retries_delay": timedelta(seconds=60),
        "start_date": datetime(2022, 10, 1),
        "depends_on_past": False,
    },
) as dag:
    execution_date = BashOperator(
        task_id="execution", bash_command="echo 'Execution date: {{ execution_date }}'"
    )

    interval_task = BashOperator(
        task_id="data_interval",
        bash_command=(
            "echo ' Interval: {{ data_interval_start}} - {{ data_interval_end }}'"
        ),
    )

    logical_date = BashOperator(
        task_id="logical_date", bash_command="echo 'Logical date: {{ ds }}'"
    )

    execution_date >> interval_task >> logical_date
