-- remove categories we don't want

DROP FUNCTION add_categories(INTEGER, TEXT[]);

CREATE OR REPLACE FUNCTION add_categories(my_channel_id INTEGER, my_names TEXT[]) RETURNS VOID AS
$$
DECLARE
my_name TEXT;
my_category_id INTEGER;
BEGIN
    -- DELETE ALL CURRENT CATEGORIES
    DELETE FROM channels_categories WHERE channel_id=my_channel_id;
    FOREACH my_name IN ARRAY my_names
    LOOP
        SELECT id FROM categories WHERE name ILIKE TRIM(my_name) INTO my_category_id;
        IF found THEN
        BEGIN
            INSERT INTO channels_categories(channel_id, category_id) VALUES (my_channel_id, my_category_id);
            EXCEPTION WHEN unique_violation THEN
            -- do nothing
        END;
        END IF;
    END LOOP;
END;
$$
LANGUAGE plpgsql;

-- ADD SUBCATEGORIES
CREATE OR REPLACE FUNCTION add_subcategories(my_parent TEXT, my_names TEXT[]) RETURNS VOID AS
$$
DECLARE
my_name TEXT;
my_parent_id INTEGER;
BEGIN
    SELECT id FROM categories WHERE name=my_parent INTO my_parent_id;
    IF NOT found THEN
        RETURN;
    END IF;
    FOREACH my_name IN ARRAY my_names
    LOOP
        BEGIN
            INSERT INTO categories(name, parent_id) VALUES (my_name, my_parent_id);
            EXCEPTION WHEN unique_violation THEN
            -- do nothing
        END;
    END LOOP;
END;
$$
LANGUAGE plpgsql;
END;

-- create only specific categories
-- see https://gist.github.com/skattyadz/814315

DELETE FROM channels_categories;
DELETE FROM categories;

ALTER TABLE categories ADD COLUMN parent_id INTEGER NULL;
ALTER TABLE categories ADD CONSTRAINT categories_parent_id_fk FOREIGN KEY(parent_id) REFERENCES categories (id);


-- create categories

INSERT INTO categories (name) VALUES('Arts');
INSERT INTO categories (name) VALUES('Business');
INSERT INTO categories (name) VALUES('Comedy');
INSERT INTO categories (name) VALUES('Games & Hobbies');
INSERT INTO categories (name) VALUES('Government & Organizations');
INSERT INTO categories (name) VALUES('Health');
INSERT INTO categories (name) VALUES('Music');
INSERT INTO categories (name) VALUES('News & Politics');
INSERT INTO categories (name) VALUES('Religion & Spirituality');
INSERT INTO categories (name) VALUES('Science & Medicine');
INSERT INTO categories (name) VALUES('Society & Culture');
INSERT INTO categories (name) VALUES('Sports & Recreation');
INSERT INTO categories (name) VALUES('TV & Film');
INSERT INTO categories (name) VALUES('Technology');

SELECT add_subcategories('Arts', ARRAY['Design', 'Fashion & Beauty', 'Food', 'Literature', 'Performing Arts', 'Spoken Word', 'Visual Arts']);
SELECT add_subcategories('Business', ARRAY['Business News', 'Careers', 'Investing', 'Management & Marketing', 'Shopping']);
SELECT add_subcategories('Education', ARRAY['Educational Technology', 'Higher Education', 'K-12', 'Language Courses', 'Training']);
SELECT add_subcategories('Games & Hobbies', ARRAY['Automotive', 'Aviation', 'Hobbies', 'Other Games', 'Video Games']);
SELECT add_subcategories('Government & Organizations', ARRAY['Local', 'National', 'Non-Profit', 'Regional']);
SELECT add_subcategories('Health', ARRAY['Alternative Health', 'Fitness & Nutrition', 'Self-Help', 'Sexuality', 'Kids & Family']);
SELECT add_subcategories('Music', ARRAY['Alternative', 'Blues', 'Country', 'Easy Listening', 'Electronic', 'Folk', 'Freeform', 'Hip-Hop & Rap',
    'Inspirational', 'Jazz', 'Latin', 'Metal', 'New Age', 'Oldies', 'Pop', 'R&B & Urban', 'Reggae', 'Rock', 'Seasonal & Holiday',
    'Soundtracks', 'World']);
SELECT add_subcategories('Electronic',  ARRAY['Acid House', 'Ambient', 'Big Beat', 'Breakbeat', 'Disco', 'Downtempo', 'Drum ''n`'' Bass',
    'Garage', 'Hard House', 'House', 'IDM', 'Jungle', 'Progressive', 'Techno', 'Trance', 'Tribal', 'Trip Hop']);
SELECT add_subcategories('News & Politics', ARRAY['Conservative (Right)', 'Liberal (Left)']);
SELECT add_subcategories('Religion & Spirituality', ARRAY['Buddhism', 'Christianity', 'Hinduism', 'Islam', 'Judaism', 'Other', 'Spirituality']);
SELECT add_subcategories('Science & Medicine', ARRAY['Medicine', 'Natural Sciences', 'Social Sciences']);
SELECT add_subcategories('Society & Culture', ARRAY['Gay & Lesbian', 'History', 'Personal Journals', 'Philosophy', 'Places & Travel']);
SELECT add_subcategories('Sports & Recreation', ARRAY['Amateur', 'College & High School', 'Outdoor', 'Professional']);
SELECT add_subcategories('Technology', ARRAY['Gadgets', 'IT News', 'Podcasting', 'Software How-To']);

