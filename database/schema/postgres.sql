CREATE TABLE IF NOT EXISTS PledgeData (
    ID serial PRIMARY KEY,
    TimeStamp timestamp NOT NULL,
    Funding bigint,
    Citizens integer,
    Fleet integer,
    CONSTRAINT unique_timestamp UNIQUE (TimeStamp)
);