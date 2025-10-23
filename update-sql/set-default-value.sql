CREATE TABLE notes_new(id integer PRIMARY KEY AUTOINCREMENT,title varchar(512) NOT NULL UNIQUE, datelog integer, content text DEFAULT "", url text DEFAULT "", flags text DEFAULT "", reminder_ticks integer DEFAULT 0,timestamp integer DEFAULT CURRENT_TIMESTAMP, readonly integer DEFAULT 0, format_tag blob, alert_count integer DEFAULT 0, pixbuf_dict blob, time_spent integer DEFAULT 0, last_text_mark blob, language text DEFAULT "", file_ext text DEFAULT "");

UPDATE notes SET content="" WHERE content IS NULL;
UPDATE notes SET url="" WHERE url IS NULL;
UPDATE notes SET flags="" WHERE flags IS NULL;
UPDATE notes SET language="" WHERE language IS NULL;
UPDATE notes SET file_ext="" WHERE file_ext IS NULL;

INSERT INTO notes_new SELECT * FROM notes;

DROP TABLE notes;
ALTER TABLE notes_new RENAME TO notes;

INSERT INTO notes(title, content) VALUES("update schema success", "See app log for more details")
-- you better restart the app after