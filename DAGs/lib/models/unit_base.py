from sqlalchemy import (
    Boolean,
    DateTime,
    Integer,
    PrimaryKeyConstraint,
    String,
    Column,
)
from pandas import DataFrame

from lib.models.brreg_base import BrregBase


class UnitBase(BrregBase):
    __abstract__ = True
    __tablename__ = "units"
    __table_args__ = (
        PrimaryKeyConstraint("organization_id"),
    )

    organization_id = Column("organization_id", String(25), nullable=False)
    name = Column("name", String(255), nullable=False)
    registered_date = Column("registered_date", DateTime)
    tax_registered = Column("tax_registered", Boolean)
    founded_date = Column("founded_date", DateTime)
    foretaksregisteret = Column("foretaksregisteret", Boolean)
    stiftelsesregisteret = Column("stiftelsesregisteret", Boolean)
    frivillighetsregisteret = Column("frivillighetsregisteret", Boolean)
    last_accounting_year = Column("last_accounting_year", Integer)
    bankruptcy = Column("bankruptcy", Boolean)
    liquidation = Column("liquidation", Boolean)
    forced_liquidation = Column("forced_liquidation", Boolean)
    organization_type_code = Column(
        "organization_type_code", String(25)
    )
    postal_address_country = Column(
        "postal_address_country", String(100),
    )
    postal_address_country_code = Column(
        "postal_address_country_code", String(8),
    )
    postal_address_code = Column("postal_address_code", String(25))
    postal_address_place = Column("postal_address_place", String(100))
    postal_address = Column("postal_address", String(255))
    postal_address_county = Column("postal_address_county", String(100))
    postal_address_county_number = Column(
        "postal_address_county_number", String(25)
    )
    business_address_country = Column(
        "business_address_country", String(100)
    )
    business_address_country_code = Column(
        "business_address_country_code", String(8)
    )
    business_address_code = Column("business_address_code", String(25))
    business_address_place = Column(
        "business_address_place", String(100)
    )
    business_address = Column("business_address", String(255))
    business_address_county = Column(
        "business_address_county", String(100)
    )
    business_address_county_number = Column(
        "business_address_county_number", String(25)
    )

    @staticmethod
    def get_dropped_columns() -> list:
        """
        Get a list of columns dropped from the original dataframe.
        """
        return [
            "antallAnsatte",
            "maalform",
            "links",
            "hjemmeside",
            "frivilligMvaRegistrertBeskrivelser",
            "naeringskode2.hjelpeenhetskode",
            "naeringskode3.hjelpeenhetskode",
            "naeringskode1.beskrivelse",
            "naeringskode1.kode",
            "naeringskode2.beskrivelse",
            "naeringskode2.kode",
            "naeringskode3.beskrivelse",
            "naeringskode3.kode",
            "organisasjonsform.beskrivelse",
            "organisasjonsform.links",
            "institusjonellSektorkode.kode",
            "institusjonellSektorkode.beskrivelse",
            "institusjonellSektorkode.kode",
            "institusjonellSektorkode.beskrivelse",
            "overordnetEnhet",
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
            "stiftelsesdato": "founded_date",
            "registrertIForetaksregisteret": "foretaksregisteret",
            "registrertIStiftelsesregisteret": "stiftelsesregisteret",
            "registrertIFrivillighetsregisteret": "frivillighetsregisteret",
            "sisteInnsendteAarsregnskap": "last_accounting_year",
            "konkurs": "bankruptcy",
            "underAvvikling": "liquidation",
            "underTvangsavviklingEllerTvangsopplosning": "forced_liquidation",
            "organisasjonsform.kode": "organization_type_code",
            "postadresse.land": "postal_address_country",
            "postadresse.landkode": "postal_address_country_code",
            "postadresse.postnummer": "postal_address_code",
            "postadresse.poststed": "postal_address_place",
            "postadresse.adresse": "postal_address",
            "postadresse.kommune": "postal_address_county",
            "postadresse.kommunenummer": "postal_address_county_number",
            "forretningsadresse.land": "business_address_country",
            "forretningsadresse.landkode": "business_address_country_code",
            "forretningsadresse.postnummer": "business_address_code",
            "forretningsadresse.poststed": "business_address_place",
            "forretningsadresse.adresse": "business_address",
            "forretningsadresse.kommune": "business_address_county",
            "forretningsadresse.kommunenummer": "business_address_county_number",
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
        new_df["business_address"] = new_df["business_address"].astype("string")
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
        new_df["business_address"] = new_df["business_address"].str.replace(
            "['", "", regex=False
        )
        new_df["business_address"] = new_df["business_address"].str.replace(
            "']", "", regex=False
        )
        new_df["business_address"] = new_df["business_address"].str.replace(
            "'", "", regex=False
        )
        new_df["business_address"] = new_df["business_address"].str.replace(
            "[]", "", regex=False
        )

        return new_df
