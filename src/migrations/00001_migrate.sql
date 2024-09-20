-- +goose Up
--create types
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tender_status') THEN
        CREATE TYPE tender_status AS ENUM (
            'Created',
            'Published',
            'Closed'
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tender_type') THEN
        CREATE TYPE tender_type AS ENUM (
            'Construction',
            'Delivery',
            'Manufacture'
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bid_status') THEN
        CREATE TYPE bid_status AS ENUM (
            'Created',
            'Published',
            'Closed'
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'author_type') THEN
        CREATE TYPE author_type AS ENUM (
            'Organization',
            'User'
        );
    END IF;
END
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

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


CREATE TABLE IF NOT EXISTS tender_history (
    id SERIAL PRIMARY KEY,
    tender_id UUID REFERENCES tender(id) NOT NULL,
    data JSONB,
    version INT,
    updated_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementBegin

CREATE OR REPLACE FUNCTION update_tender_version() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = OLD.status THEN
        INSERT INTO tender_history (tender_id, data, version)
        VALUES (OLD.id, row_to_json(OLD), OLD.version);
        NEW.version := OLD.version + 1;
    ELSE
        NEW.version := OLD.version;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'tender_update_trigger') THEN
        CREATE TRIGGER tender_update_trigger
        BEFORE UPDATE ON tender
        FOR EACH ROW EXECUTE FUNCTION update_tender_version();
    END IF;
END
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION rollback_tender_to_version(p_tender_id UUID, p_version INT) RETURNS VOID AS $$
DECLARE
    v_old_data JSONB;
BEGIN
    SELECT data INTO v_old_data
    FROM tender_history
    WHERE tender_id = p_tender_id AND version = p_version
    LIMIT 1;
    UPDATE tender
    SET name = v_old_data->>'name',
        description = v_old_data->>'description',
        "type" = (v_old_data->>'type')::tender_type,
        version = p_version + 1
    WHERE id = p_tender_id;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

