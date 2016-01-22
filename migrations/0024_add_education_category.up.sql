INSERT INTO categories (name) VALUES('Education');

SELECT add_subcategories('Education', ARRAY['Educational Technology', 'Higher Education', 'K-12', 'Language Courses', 'Training']);
