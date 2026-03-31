CREATE TABLE IF NOT EXISTS urls (
    id          BIGSERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code   VARCHAR(10) NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT urls_original_url_unique UNIQUE (original_url),
    CONSTRAINT urls_short_code_unique   UNIQUE (short_code)
);