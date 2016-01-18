ALTER TABLE bookmarks DROP CONSTRAINT bookmarks_user_id_fkey;

ALTER TABLE bookmarks

ADD CONSTRAINT bookmarks_user_id_fkey FOREIGN KEY (user_id)
      REFERENCES users (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE subscriptions DROP CONSTRAINT subscriptions_user_id_fkey;

ALTER TABLE subscriptions

ADD CONSTRAINT subscriptions_user_id_fkey FOREIGN KEY (user_id)
      REFERENCES users (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE;


