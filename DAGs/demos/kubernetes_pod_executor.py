from pendulum import datetime, duration
from airflow import DAG
from airflow.providers.cncf.kubernetes.operators.kubernetes_pod import (
    KubernetesPodOperator,
)

with DAG(
    dag_id="example_kubernetes_pod",
    schedule="@once",
    default_args={
        "owner": "airflow",
        "depends_on_past": False,
        "start_date": datetime(2022, 1, 1),
        "email_on_failure": False,
        "email_on_retry": False,
        "retries": 1,
        "retry_delay": duration(minutes=5),
    },
) as dag:
    KubernetesPodOperator(
        namespace="airflow",
        image="hello-world",
        labels={"environment": "development"},
        name="airflow-test-pod",
        task_id="task-one",
        in_cluster=True,
        config_file=None,
        is_delete_operator_pod=True,
        get_logs=True,
    )
