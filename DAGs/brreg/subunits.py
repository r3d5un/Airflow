from datetime import timedelta
from pathlib import Path
from airflow.operators.python import PythonOperator

from pendulum import datetime

from airflow import DAG
from airflow.models import Variable
from airflow.providers.postgres.operators.postgres import PostgresOperator

from brreg import (
    download_brreg_gzip,
    extract_brreg_gzip,
    create_brreg_parquet,
    load_brreg_parquet,
)
from brreg.brreg_config import BrregConfig
from lib.models.subunit import Subunit
from lib.models.subunit_staging import SubunitStaging
from utils import get_database_engine
from utils.database_credentials import DatabaseCredentials


brreg_config = BrregConfig(
    api_url="https://data.brreg.no/enhetsregisteret/api/underenheter/lastned",
    header={
        "Accept": "application/vnd.brreg.enhetsregisteret.underenhet.v1+gzip;charset=UTF-8",
    },
    storage_directory=Path("/storage"),
    gz_file="alle_underenheter.json.gz",
    json_file="alle_underenheter.json",
    parquet_file="subunits.parquet",
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
    extract = PythonOperator(
        task_id="download_subunits",
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
        op_kwargs={"config": brreg_config, "db_model": SubunitStaging},
    )

    create_staging = PythonOperator(
        task_id="create_staging",
        python_callable=SubunitStaging.__table__.create,
        op_kwargs={
            "bind": get_database_engine(database_credentials),
            "checkfirst": True,
        },
    )

    create_public = PythonOperator(
        task_id="create_public_table",
        python_callable=Subunit.__table__.create,
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
            "db_model": SubunitStaging,
            "database_credentials": database_credentials,
        },
    )

    switcharoo = PostgresOperator(
        task_id="switcharoo",
        sql="switcharoo.sql",
        postgres_conn_id="warehouse",
        params={
            "tablename": Subunit.__table__.name,
            "target_schema": Subunit.__table__.schema,
            "staging_schema": SubunitStaging.__table__.schema,
            "index_names": [index.name for index in Subunit.__table__.indexes],
        },
    )

    cleanup = PythonOperator(
        task_id="drop_staging",
        python_callable=SubunitStaging.__table__.drop,
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
            "tablename": Subunit.__table__.name,
            "schema": Subunit.__table__.schema,
        },
    )

    extract >> unzip >> transform >> load >> switcharoo >> [cleanup, analyze]  # type: ignore

    [create_staging, create_public] >> load  # type: ignore
