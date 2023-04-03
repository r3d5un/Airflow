from datetime import datetime, timedelta

from airflow import DAG
from airflow.operators.empty import EmptyOperator
from airflow.operators.bash import BashOperator

with DAG(
    dag_id="bash_template_demo",
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
    template_searchpath=["/opt/airflow/dags/repo/dags/demos"]
) as dag:
    start = EmptyOperator(task_id="start")
    echo_stdout = BashOperator(
        task_id="echo_stdout",
        bash_command="demo_bash_template.sh",
        params={"message": "Hello World!"},
    )
    end = EmptyOperator(task_id="end")

    start >> echo_stdout >> end
