# Airflow

This is a reference configuration of the [Airflow Helm Chart](https://github.com/airflow-helm/charts) and a few [DAGs](https://airflow.apache.org/docs/apache-airflow/stable/core-concepts/dags.html) serving as demos.

## Prerequisites

The installation requires that the user has setup:

- A working Kubernetes install (this configuration assumes a single-node cluster)

- The [Helm Package Manager](https://helm.sh/)

- Python 3.8+

## Installing

- Add the [Airflow Helm Repository](https://github.com/airflow-helm/charts) to Helm.

```bash
helm repo add airflow-stable https://airflow-helm.github.io/charts
```

- Update repositories.

```bash
helm repo update
```

- Apply manifest files.

> It is highly recommended to generate new secrets/passwords and keys for the `manifest/airflow-secrets.yaml` file. Instructions can be found in the comments of the file.

```bash
kubectl apply -f manifest/
```

- Install Airflow.

```bash
helm upgrade --cleanup-on-fail --namespace airflow --install airflow airflow-stable/airflow -f airflow.yaml
```
