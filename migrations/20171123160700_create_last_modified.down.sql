DO $$
BEGIN
  IF EXISTS(SELECT column_name
            FROM information_schema.columns
            WHERE table_name = 'jobs' AND column_name = 'last_modified')
  THEN
    ALTER TABLE jobs
      DROP last_modified;
  END IF;
END
$$