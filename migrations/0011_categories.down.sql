DROP INDEX channel_categories_channel_id_idx;
DROP INDEX channel_categories_category_id_idx;
DROP INDEX categories_name_idx;
ALTER TABLE channel_categories DROP CONSTRAINT unique_channel_categories;
DROP TABLE channel_categories;
DROP TABLE categories;
