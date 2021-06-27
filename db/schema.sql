CREATE TABLE urls (
    id               BIGSERIAL PRIMARY KEY,
    code             TEXT NOT NULL,
    url              TEXT NOT NULL,
    creation_date    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    update_date      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX ON urls (code);
CREATE INDEX ON urls (url);