CREATE TYPE point_status AS ENUM ('pending', 'pulled', 'running', 'failed', 'completed');


CREATE TABLE points (
    id          SERIAL NOT NULL,
    project     VARCHAR(40),
    status point_status,

    coordinate      TEXT,
    metric_value    DOUBLE PRECISION,
    metadata        TEXT,

    PRIMARY KEY (project, id)
);
