from dataclasses import dataclass
from pathlib import Path


@dataclass
class BrregConfig:
    def __init__(
        self,
        api_url: str,
        header: dict,
        storage_directory: Path,
        gz_file: str,
        json_file: str,
        parquet_file: str,
    ):
        self.api_url = api_url
        self.header = header
        self.storage_directory = Path(storage_directory)
        self.gz_file = self.storage_directory.joinpath(gz_file)
        self.json_file = self.storage_directory.joinpath(json_file)
        self.parquet_file = self.storage_directory.joinpath(parquet_file)
