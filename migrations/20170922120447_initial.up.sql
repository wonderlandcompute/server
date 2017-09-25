CREATE TABLE points (
    id              SERIAL NOT NULL,
    project         VARCHAR(40),
    status          SMALLINT,

    coordinate      TEXT NOT NULL DEFAULT '',
    metric_value    TEXT NOT NULL DEFAULT '',
    metadata        TEXT NOT NULL DEFAULT '',

    created         timestamp without time zone default (now() at time zone 'utc'),
    creator         VARCHAR(40),

    PRIMARY KEY (project, id)
);
