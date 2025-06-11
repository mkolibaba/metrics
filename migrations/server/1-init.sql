CREATE TABLE IF NOT EXISTS gauge
(
    id    varchar primary key not null,
    value double precision    not null
);

CREATE TABLE IF NOT EXISTS counter
(
    id    varchar primary key not null,
    delta bigint              not null
);
