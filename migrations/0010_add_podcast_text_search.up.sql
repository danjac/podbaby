ALTER TABLE podcasts ADD COLUMN tsv tsvector;
CREATE INDEX podcasts_tsv_idx ON podcasts USING gin(tsv);

CREATE OR REPLACE FUNCTION podcasts_search_trigger() RETURNS trigger AS $$
BEGIN
  new.tsv :=
    SETWEIGHT(TO_TSVECTOR(COALESCE(new.title, '')), 'A') ||
    SETWEIGHT(TO_TSVECTOR(COALESCE(new.description, '')), 'D');
  RETURN new;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER podcasts_tsvectorupdate BEFORE INSERT OR UPDATE
  ON podcasts FOR EACH ROW EXECUTE PROCEDURE podcasts_search_trigger();
