DROP TRIGGER podcasts_tsvectorupdate;
DROP FUNCTION podcasts_search_trigger;
DROP INDEX podcasts_tsv_idx;
ALTER TABLE podcasts DROP COLUMN tsv;
