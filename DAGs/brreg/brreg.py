from datetime import timedelta

from pendulum import datetime

from airflow import DAG
from airflow.providers.postgres.operators.postgres import PostgresOperator
from airflow.providers.cncf.kubernetes.operators.kubernetes_pod import (
    KubernetesPodOperator,
)

with DAG(
    dag_id="brreg",
    schedule_interval="@daily",
    catchup=False,
    max_active_runs=1,
    tags=["brreg", "etl", "units"],
    default_args={
        "owner": "airflow",
        "retries": 3,
        "retry_delay": timedelta(minutes=5),
        "start_date": datetime(2023, 3, 1),
        "depends_on_past": False,
    },
    template_searchpath=["/opt/airflow/dags/repo/dags/brreg/templates"],
):

    create_units_table = PostgresOperator(
        task_id="create_table_units",
        sql="create_table_units.sql",
        postgres_conn_id="warehouse",
    )

    create_subunits_table = PostgresOperator(
        task_id="create_table_subunits.sql",
        sql="create_table_subunits.sql",
        postgres_conn_id="warehouse",
    )

    analyze_units_table = PostgresOperator(
        task_id="analyze_table",
        sql="ANALYZE public.brreg_units",
        postgres_conn_id="warehouse",
    )

    analyze_subunits_table = PostgresOperator(
        task_id="analyze_subunits_table",
        sql="ANALYZE public.brreg_subunits",
        postgres_conn_id="warehouse",
    )

    ingest_brreg_units = KubernetesPodOperator(
        namespace="airflow",
        image="ghcr.io/r3d5un/brreg:latest",
        labels={"dataset": "units"},
        name="ingest-brreg-units",
        task_id="ingest_brreg_units",
        in_cluster=True,
        config_file=None,
        is_delete_operator_pod=True,
        get_logs=True,
        arguments=[
            "-workers",
            "25",
            # Hardcoded for convencience, do not put this into production
            "-dsn",
            "postgres://postgres:postgres@10.0.1.1:5432/warehouse?sslmode=disable",
            "-dataset",
            "units",
        ],
    )

    ingest_brreg_subunits = KubernetesPodOperator(
        namespace="airflow",
        image="ghcr.io/r3d5un/brreg:latest",
        labels={"dataset": "subunits"},
        name="ingest-brreg-subunits",
        task_id="ingest_brreg_subunits",
        in_cluster=True,
        config_file=None,
        is_delete_operator_pod=True,
        get_logs=True,
        arguments=[
            "-workers",
            "25",
            # Hardcoded for convencience, do not put this into production
            "-dsn",
            "postgres://postgres:postgres@10.0.1.1:5432/warehouse?sslmode=disable",
            "-dataset",
            "subunits",
        ],
    )

    create_units_table >> ingest_brreg_units >> analyze_units_table  # type: ignore
    create_subunits_table >> ingest_brreg_subunits >> analyze_subunits_table  # type: ignore
