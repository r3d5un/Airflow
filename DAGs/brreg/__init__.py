import json
import logging
import gzip
import shutil

import pandas as pd
import requests

from brreg.brreg_config import BrregConfig
from lib.models.brreg_base import BrregBase
from utils import get_database_engine
from utils.database_credentials import DatabaseCredentials


logger = logging.getLogger(__name__)


def download_brreg_gzip(config: BrregConfig) -> None:
    """
    Downloads a gzipped file from The Brunnøysund Register Centre.
    """
    logger.debug(f"Downloading file from The Brunnøysund Register Centre")
    response = requests.get(config.api_url, headers=config.header)
    if response.status_code != 200:
        raise Exception(
            f"Request to {config.api_url} failed with status code {response.status_code}"
        )

    logger.info(f"Writing response content to {config.gz_file}")
    with open(config.gz_file, "wb") as f:
        f.write(response.content)


def extract_brreg_gzip(config: BrregConfig) -> None:
    """
    Extracts a gzipped file from The Brunnøysundregisteret Register Centre.
    """
    logger.debug(f"Unzipping {config.gz_file} to {config.json_file}")
    with gzip.open(config.gz_file, "rb") as f_in, open(config.json_file, "wb") as f_out:
        shutil.copyfileobj(f_in, f_out)


def create_brreg_parquet(config: BrregConfig, db_model: BrregBase):
    """
    Reads the The Brunnøysund Registry Centre data, normalizes the JSON file into a
    DataFrame, performs some basic transformations and writes the data to a Parquet.
    """
    logger.debug(f"Opening {config.json_file}")
    with open(config.json_file, "r") as f:
        units = json.load(f)

    logger.debug(f"Normalize JSON to Pandas DataFrame")
    df = (
        pd.json_normalize(units)
        .drop(columns=db_model.get_dropped_columns())
        .rename(columns=db_model.get_renamed_columns())
    )
    logger.debug(f"Transforming and Cleaning DataFrame")
    df = db_model.transform_df(df)

    logger.info(f"Saving DataFrame to {config.parquet_file}")
    df.to_parquet(str(config.parquet_file), index=False)


def load_brreg_parquet(
    config: BrregConfig, db_model: BrregBase, database_credentials: DatabaseCredentials
):
    """
    Reads The Brunnøysund Registry Centre data from a Parquet file and loads it into
    the database.
    """
    logger.debug(f"Creating database engine")
    engine = get_database_engine(database_credentials)

    logger.info(
        f"Loading {config.parquet_file} to database {database_credentials.database}."
        f"{db_model.__table__.schema}.{db_model.__table__.name}"
    )
    pd.read_parquet(str(config.parquet_file)).to_sql(
        name=db_model.__table__.name,
        schema=db_model.__table__.schema,
        con=engine,
        index=False,
        if_exists="append",
    )
