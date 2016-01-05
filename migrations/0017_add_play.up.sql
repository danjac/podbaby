CREATE OR REPLACE FUNCTION add_play (my_podcast_id INTEGER, my_user_id INTEGER) returns INTEGER AS
$$
DECLARE
play_id INTEGER;
BEGIN

    BEGIN
        INSERT INTO plays (podcast_id, user_id) VALUES(my_podcast_id, my_user_id) RETURNING id INTO play_id;
    EXCEPTION WHEN unique_violation THEN
        SELECT id FROM plays WHERE podcast_id=my_podcast_id AND user_id=my_user_id INTO play_id;
        UPDATE plays SET created_at = NOW() WHERE id=play_id;
    END;

    RETURN play_id;

END;
$$
LANGUAGE plpgsql;
