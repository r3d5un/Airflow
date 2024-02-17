from datetime import timedelta

from pendulum import datetime

from airflow import DAG
from airflow.providers.postgres.operators.postgres import PostgresOperator
from airflow.providers.cncf.kubernetes.operators.kubernetes_pod import (
    KubernetesPodOperator,
)

with DAG(
    dag_id="peppol",
    schedule_interval="@daily",
    catchup=False,
    max_active_runs=1,
    tags=["peppol", "etl"],
    default_args={
        "owner": "r3d5un",
        "retries": 3,
        "retry_delay": timedelta(minutes=5),
        "start_date": datetime(2023, 3, 1),
        "depends_on_past": False,
    },
    template_searchpath=["/opt/airflow/dags/repo/dags/peppol/templates"],
):

    create_peppol_table = PostgresOperator(
        task_id="create_table_peppol",
        sql="create_table_peppol_business_cards.sql",
        postgres_conn_id="warehouse",
    )

    analyze_peppol_table = PostgresOperator(
        task_id="analyze_table",
        sql="ANALYZE public.peppol_business_cards",
        postgres_conn_id="warehouse",
    )

    ingest_peppol = KubernetesPodOperator(
        namespace="airflow",
        image="ghcr.io/r3d5un/peppol:latest",
        labels={"dataset": "peppol"},
        name="ingest-peppol",
        task_id="ingest_peppol",
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
        ],
    )

    create_peppol_table >> ingest_peppol >> analyze_peppol_table  # type: ignore
