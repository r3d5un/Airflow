from datetime import timedelta
from pathlib import Path

from pendulum import datetime

from airflow import DAG
from airflow.models import Variable
from airflow.operators.python import PythonOperator
from airflow.providers.postgres.operators.postgres import PostgresOperator

from brreg import (
    download_brreg_gzip,
    extract_brreg_gzip,
    create_brreg_parquet,
    load_brreg_parquet,
)
from brreg.brreg_config import BrregConfig
from lib.models.unit import Unit
from lib.models.unit_staging import UnitStaging
from utils import get_database_engine
from utils.database_credentials import DatabaseCredentials


brreg_config = BrregConfig(
    api_url="https://data.brreg.no/enhetsregisteret/api/enheter/lastned",
    header={
        "Accept": "application/vnd.brreg.enhetsregisteret.enhet.v1+gzip;charset=UTF-8",
    },
    storage_directory=Path("/storage"),
    gz_file="alle_enheter.json.gz",
    json_file="alle_enheter.json",
    parquet_file="units.parquet",
)

database_credentials = DatabaseCredentials(
    database_type=Variable.get("warehouse_database_type"),
    host=Variable.get("warehouse_host"),
    port=Variable.get("warehouse_port"),
    user=Variable.get("warehouse_user"),
    password=Variable.get("warehouse_password"),
    database=Variable.get("warehouse_database"),
)

with DAG(
    dag_id="brreg_units",
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
    template_searchpath=["/opt/airflow/dags/brreg/templates"],
):
    extract = PythonOperator(
        task_id="download_units",
        python_callable=download_brreg_gzip,
        op_kwargs={"config": brreg_config},
    )

    unzip = PythonOperator(
        task_id="extract_brreg_gzip",
        python_callable=extract_brreg_gzip,
        op_kwargs={"config": brreg_config},
    )

    transform = PythonOperator(
        task_id="create_parquet",
        python_callable=create_brreg_parquet,
        op_kwargs={"config": brreg_config, "db_model": UnitStaging},
    )

    create_staging = PythonOperator(
        task_id="create_staging_table",
        python_callable=UnitStaging.__table__.create,
        op_kwargs={
            "bind": get_database_engine(database_credentials),
            "checkfirst": True,
        },
    )

    create_public = PythonOperator(
        task_id="create_public_table",
        python_callable=Unit.__table__.create,
        op_kwargs={
            "bind": get_database_engine(database_credentials),
            "checkfirst": True,
        },
    )

    load = PythonOperator(
        task_id="load",
        python_callable=load_brreg_parquet,
        op_kwargs={
            "config": brreg_config,
            "db_model": UnitStaging,
            "database_credentials": database_credentials,
        },
    )

    switcharoo = PostgresOperator(
        task_id="switcharoo",
        sql="switcharoo.sql",
        postgres_conn_id="warehouse",
        params={
            "tablename": Unit.__table__.name,
            "target_schema": Unit.__table__.schema,
            "staging_schema": UnitStaging.__table__.schema,
            "index_names": [index.name for index in Unit.__table__.indexes],
        },
    )

    cleanup = PythonOperator(
        task_id="drop_staging",
        python_callable=UnitStaging.__table__.drop,
        op_kwargs={
            "bind": get_database_engine(database_credentials),
            "checkfirst": True,
        },
    )

    analyze = PostgresOperator(
        task_id="analyze",
        sql="analyze.sql",
        postgres_conn_id="warehouse",
        params={
            "tablename": Unit.__table__.name,
            "schema": Unit.__table__.schema,
        },
    )

    extract >> unzip >> transform >> load >> switcharoo >> [cleanup, analyze]  # type: ignore

    [create_staging, create_public] >> load  # type: ignore
