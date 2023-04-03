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
