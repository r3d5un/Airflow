CREATE TABLE IF NOT EXISTS public.peppol_business_cards (
    id VARCHAR(128) PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    countrycode VARCHAR(2) NOT NULL,
    last_updated TIMESTAMP DEFAULT NOW(),
    business_card JSONB NOT NULL
);
