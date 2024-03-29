import logging

from sqlalchemy import create_engine

from utils.database_credentials import DatabaseCredentials

logger = logging.getLogger(__name__)


def get_database_engine(database_credentials: DatabaseCredentials):
    """
    Creates a database engine for the given database credentials.
    """
    logger.info(f"Creating database engine for {database_credentials}")
    connection_string = (
        f"{database_credentials.database_type}://{database_credentials.user}:"
        f"{database_credentials.password}@{database_credentials.host}:"
        f"{database_credentials.port}/{database_credentials.database}"
        f"?sslmode=disable"
    )
    return create_engine(connection_string)
