from datetime import datetime, timedelta

from airflow import DAG
from airflow.models import Variable
from airflow.operators.python import PythonOperator


def print_airflow_variables():
    print(f"We are currently running in the {Variable.get('ENVIRONMENT')} environment")


with DAG(
    dag_id="airflow_variables_demo",
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
    print_task = PythonOperator(
        task_id="print_airflow_variables",
        python_callable=print_airflow_variables,
    )
