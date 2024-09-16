-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TYPE IF EXISTS tender_type;
CREATE TYPE tender_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
);
DROP TYPE IF EXISTS tender_status;

CREATE TYPE tender_status AS ENUM (
    'Created',
    'Published',
    'Closed'
);

CREATE TABLE IF NOT EXISTS tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type tender_type,
    status tender_status NOT NULL DEFAULT 'Created',
    organization_id UUID,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
DROP TYPE IF EXISTS bid_status;

CREATE TYPE bid_status AS ENUM (
    'Created',
    'Published',
    'Closed'
);
DROP TYPE IF EXISTS author_type;

CREATE TYPE  author_type AS ENUM (
    'Organization',
    'User'
);

CREATE TABLE IF NOT EXISTS bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status bid_status NOT NULL DEFAULT 'Created',
    tender_id UUID REFERENCES tender(id) NOT NULL,
    author_type author_type NOT NULL,
    author_id UUID NOT NULL,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
