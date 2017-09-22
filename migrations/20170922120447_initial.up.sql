CREATE TABLE points (
    id              SERIAL NOT NULL,
    project         VARCHAR(40),
    status          SMALLINT,

    coordinate      TEXT NOT NULL DEFAULT '',
    metric_value    TEXT NOT NULL DEFAULT '',
    metadata        TEXT NOT NULL DEFAULT '',

    PRIMARY KEY (project, id)
);
