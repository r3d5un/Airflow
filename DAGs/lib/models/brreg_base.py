from pandas import DataFrame

from lib.models import Base


class BrregBase(Base):
    __abstract__ = True
    __tablename__ = "brreg_units"

    @staticmethod
    def get_dropped_columns() -> list:
        """
        Get a list of columns dropped from the original dataframe.
        """
        return []

    @staticmethod
    def get_renamed_columns() -> dict:
        """
        Get a renamed columns from the original dataframe.
        """
        return {}

    @staticmethod
    def transform_df(df: DataFrame) -> DataFrame:
        """
        Perform any transformations then returns a new dataframe.
        """
        return df.copy()
