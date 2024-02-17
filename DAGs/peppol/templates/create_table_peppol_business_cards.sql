CREATE TABLE IF NOT EXISTS public.peppol_business_cards (
    id VARCHAR(25) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    countrycode VARCHAR(2) NOT NULL,
    last_updated TIMESTAMP DEFAULT NOW(),
    business_card JSONB NOT NULL
);
