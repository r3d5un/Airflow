from datetime import datetime, timedelta

from airflow import DAG
from airflow.operators.python import PythonOperator


def python_operator_print_function(param1, param2):
    print(f"This is the first parameter: {param1}")
    print(f"This is the second parameter: {param2}")
    print(f"{param1} {param2}!")


with DAG(
        dag_id="python_operator_demo",
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
        python_callable=python_operator_print_function,
        op_kwargs={"param1": "Hello", "param2": "World"},
    )
