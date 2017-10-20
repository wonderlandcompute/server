CREATE TABLE jobs (
  id           SERIAL NOT NULL,
  project      VARCHAR(40),
  status       SMALLINT,

  metadata     TEXT   NOT NULL             DEFAULT '',

  created      TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() AT TIME ZONE 'utc'),
  creator      VARCHAR(40),

  input        TEXT   NOT NULL             DEFAULT '',
  output       TEXT   NOT NULL             DEFAULT '',
  kind         TEXT   NOT NULL             DEFAULT '',

  PRIMARY KEY (id)
);

CREATE INDEX status_idx
  ON jobs (status);
CREATE INDEX access_idx
  ON jobs (project,kind);
