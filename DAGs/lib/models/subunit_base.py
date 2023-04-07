from sqlalchemy import (
    PrimaryKeyConstraint,
    Column,
    String,
    DateTime,
    Boolean,
    Integer,
)
from pandas import DataFrame

from lib.models.brreg_base import BrregBase


class SubunitBase(BrregBase):
    __abstract__ = True
    __tablename__ = "brreg_subunits"
    __table_args__ = PrimaryKeyConstraint("organization_id")

    organization_id = Column("organization_id", String(25), nullable=False)
    name = Column("name", String(255), nullable=False)
    registered_date = Column("registered_date", DateTime)
    tax_registered = Column("tax_registered", Boolean)
    number_of_employees = Column("number_of_employees", Integer)
    parent_organization_id = Column("parent_organization_id", String(25))
    start_date = Column("start_date", DateTime)
    organization_type_code = Column("organization_type_code", String(25))
    description = Column("description", String(255))
    industry_description = Column("industry_description", String(255))
    industry_code = Column("industry_code", String(25))
    address_country = Column("address_country", String(100))
    address_country_code = Column("address_country_code", String(25))
    address_code = Column("address_code", String(25))
    address_place = Column("address_place", String(100))
    address = Column("address", String(255))
    address_county = Column("address_county", String(100))
    address_county_code = Column("address_county_code", String(25))
    ownership_change_date = Column("ownership_change_date", DateTime)
    website = Column("website", String(255))
    postal_country = Column("postal_country", String(100))
    postal_country_code = Column("postal_country_code", String(25))
    postal_code = Column("postal_code", String(25))
    postal_place = Column("postal_place", String(100))
    postal_address = Column("postal_address", String(255))
    postal_county = Column("postal_county", String(100))
    postal_county_code = Column("postal_county_code", String(25))
    dissolution_date = Column("dissolution_date", DateTime)

    @staticmethod
    def get_dropped_columns() -> list:
        """
        Get a list of columns dropped from the original dataframe.
        """
        return [
            "links",
            "organisasjonsform.links",
            "naeringskode2.beskrivelse",
            "naeringskode2.kode",
            "naeringskode2.hjelpeenhetskode",
            "naeringskode3.beskrivelse",
            "naeringskode3.kode",
            "frivilligMvaRegistrertBeskrivelser",
            "naeringskode3.hjelpeenhetskode",
            "naeringskode1.hjelpeenhetskode",
        ]

    @staticmethod
    def get_renamed_columns() -> dict[str, str]:
        """
        Get a dict of renamed columns.
        """
        return {
            "organisasjonsnummer": "organization_id",
            "navn": "name",
            "registreringsdatoEnhetsregisteret": "registered_date",
            "registrertIMvaregisteret": "tax_registered",
            "antallAnsatte": "number_of_employees",
            "overordnetEnhet": "parent_organization_id",
            "oppstartsdato": "start_date",
            "organisasjonsform.kode": "organization_type_code",
            "organisasjonsform.beskrivelse": "description",
            "naeringskode1.kode": "industry_code",
            "naeringskode1.beskrivelse": "industry_description",
            "beliggenhetsadresse.land": "address_country",
            "beliggenhetsadresse.landkode": "address_country_code",
            "beliggenhetsadresse.postnummer": "address_code",
            "beliggenhetsadresse.poststed": "address_place",
            "beliggenhetsadresse.adresse": "address",
            "beliggenhetsadresse.kommune": "address_county",
            "beliggenhetsadresse.kommunenummer": "address_county_code",
            "datoEierskifte": "ownership_change_date",
            "hjemmeside": "website",
            "postadresse.land": "postal_country",
            "postadresse.landkode": "postal_country_code",
            "postadresse.postnummer": "postal_code",
            "postadresse.poststed": "postal_place",
            "postadresse.adresse": "postal_address",
            "postadresse.kommune": "postal_county",
            "postadresse.kommunenummer": "postal_county_code",
            "nedleggelsesdato": "dissolution_date",
        }

    @staticmethod
    def transform_df(df: DataFrame) -> DataFrame:
        """
        Performs basic transformations on a given Unit DataFrame. Returns a
        transformed DataFrame.
        """
        new_df = df.copy()
        # Original data consists of multiple address lines. These steps normalizes the
        # by casting the lists into strings, then removing the brackets and quotation
        # marks.
        new_df["postal_address"] = new_df["postal_address"].astype("string")
        new_df["address"] = new_df["address"].astype("string")
        new_df["postal_address"] = new_df["postal_address"].str.replace(
            "['", "", regex=False
        )
        new_df["postal_address"] = new_df["postal_address"].str.replace(
            "']", "", regex=False
        )
        new_df["postal_address"] = new_df["postal_address"].str.replace(
            "'", "", regex=False
        )
        new_df["postal_address"] = new_df["postal_address"].str.replace(
            "[]", "", regex=False
        )
        new_df["address"] = new_df["address"].str.replace("['", "", regex=False)
        new_df["address"] = new_df["address"].str.replace("']", "", regex=False)
        new_df["address"] = new_df["address"].str.replace("'", "", regex=False)
        new_df["address"] = new_df["address"].str.replace("[]", "", regex=False)

        return new_df
