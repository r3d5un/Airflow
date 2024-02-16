CREATE TABLE IF NOT EXISTS public.brreg_subunits (
    id VARCHAR(25) PRIMARY KEY,
    parent_id VARCHAR(25) NOT NULL,
    last_updated TIMESTAMP DEFAULT NOW(),
    brreg_subunit JSONB NOT NULL
);
