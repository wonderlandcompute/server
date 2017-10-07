CREATE TABLE points (
  id           SERIAL NOT NULL,
  project      VARCHAR(40),
  status       SMALLINT,

  coordinate   TEXT   NOT NULL             DEFAULT '',
  metric_value TEXT   NOT NULL             DEFAULT '',
  metadata     TEXT   NOT NULL             DEFAULT '',

  created      TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() AT TIME ZONE 'utc'),
  creator      VARCHAR(40),

  input        TEXT   NOT NULL             DEFAULT '',
  output       TEXT   NOT NULL             DEFAULT '',
  kind         SMALLINT,

  PRIMARY KEY (project, id)
);

CREATE INDEX status_idx
  ON points (status);
