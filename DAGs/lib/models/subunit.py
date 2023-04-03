from sqlalchemy import PrimaryKeyConstraint, String, Column

from lib.models.subunit_base import SubunitBase

class Subunit(SubunitBase):
    __table_args__ = (
        PrimaryKeyConstraint("organization_id"),
        {"schema": "public"},
    )

    organization_id = Column("organization_id", String(25), nullable=False)

    def __init__(self, *initial_dictionary, **kwargs):
        super(Subunit, self).__init__(*initial_dictionary, **kwargs)

