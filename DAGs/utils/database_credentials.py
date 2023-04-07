from dataclasses import dataclass


@dataclass
class DatabaseCredentials:
    def __init__(
        self,
        database_type: str,
        host: str,
        port: int,
        user: str,
        password: str,
        database: str,
    ):
        self.database_type = database_type
        self.host = host
        self.port = port
        self.user = user
        self.password = password
        self.database = database

    def __repr__(self) -> str:
        return (
            f"<{self.__class__.__name__}(database_type:{self.database_type})"
            f"host:{self.host}, port:{self.port}, database:{self.database})>"
        )
