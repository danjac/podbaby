DROP TRIGGER channels_tsvectorupdate;
DROP FUNCTION channels_search_trigger;
DROP INDEX channels_tsv_idx;
ALTER TABLE channels DROP COLUMN tsv;
