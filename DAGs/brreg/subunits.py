from datetime import timedelta

from pendulum import datetime

from airflow import DAG
from airflow.operators.empty import EmptyOperator
from airflow.providers.postgres.operators.postgres import PostgresOperator

with DAG(
    dag_id="brreg_subunits",
    schedule_interval="@daily",
    catchup=False,
    max_active_runs=1,
    tags=["brreg", "etl", "subunits"],
    default_args={
        "owner": "airflow",
        "retries": 1,
        "retry_delay": timedelta(minutes=5),
        "start_date": datetime(2023, 1, 1),
        "depends_on_past": False,
    },
    template_searchpath=["/opt/airflow/dags/brreg/templates"],
):

    create_table = PostgresOperator(
        task_id="create_table",
        sql="create_subunits.sql",
        postgres_conn_id="warehouse",
    )

    analyze_table = PostgresOperator(
        task_id="analyze_table",
        sql="ANALYZE public.subunits",
        postgres_conn_id="warehouse",
    )

    etl = EmptyOperator(task_id="etl")

    create_table >> etl >> analyze_table  # type: ignore
