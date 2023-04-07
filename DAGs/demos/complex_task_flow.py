from datetime import datetime, timedelta

from airflow import DAG
from airflow.operators.empty import EmptyOperator

with DAG(
    dag_id="complex_task_flow",
    schedule_interval=None,
    catchup=False,
    tags=["demo", "complex_task_flow", "task flow", "tutorial"],
    default_args={
        "owner": "oyvind.kristiansen@dfo.no",
        "retries": 5,
        "retries_delay": timedelta(seconds=60),
        "start_date": datetime(2022, 9, 22),
        "depends_on_past": False,
    },
) as dag:
    task_list_numbers = ["5", "6", "7", "8", "9", "10"]

    task_1 = EmptyOperator(task_id="task_1")
    task_2 = EmptyOperator(task_id="task_2")
    task_3 = EmptyOperator(task_id="task_3")
    task_4 = EmptyOperator(task_id="task_4")

    (
        task_1
        >> [task_2, task_3]
        >> task_4
        >> [
            EmptyOperator(task_id=f"task_{task_list_number}")
            for task_list_number in task_list_numbers
        ]
    )
