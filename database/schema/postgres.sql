CREATE TABLE IF NOT EXISTS starcitizen_pledgedata (
    id serial PRIMARY KEY,
    pledge_timestamp timestamp NOT NULL,
    funding bigint,
    citizens integer,
    fleet integer,
    CONSTRAINT unique_timestamp UNIQUE (pledge_timestamp)
);