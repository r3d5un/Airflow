CREATE TABLE IF NOT EXISTS public.brreg_units (
    id VARCHAR(25) PRIMARY KEY,
    last_updated TIMESTAMP DEFAULT NOW(),
    brreg_unit JSONB NOT NULL
);
