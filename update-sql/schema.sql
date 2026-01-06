CREATE TABLE app_configs(id integer PRIMARY KEY AUTOINCREMENT,key varchar(128) UNIQUE, val text);

CREATE TABLE notes(id integer PRIMARY KEY AUTOINCREMENT,title varchar(512) NOT NULL UNIQUE, datelog integer, content text DEFAULT "", url text DEFAULT "", flags text DEFAULT "", reminder_ticks integer DEFAULT 0,timestamp integer DEFAULT CURRENT_TIMESTAMP, readonly integer DEFAULT 0, format_tag blob, alert_count integer DEFAULT 0, pixbuf_dict blob, time_spent integer DEFAULT 0, last_text_mark blob, language text DEFAULT "", file_ext text DEFAULT "");

CREATE VIRTUAL TABLE note_fts USING fts5(title, datelog, content, content='notes', content_rowid='id');

CREATE TRIGGER notes_ai AFTER INSERT ON notes BEGIN
  INSERT INTO note_fts(rowid, title, datelog, content) VALUES (new.id, new.title, new.datelog, new.content);
END;
CREATE TRIGGER notes_ad AFTER DELETE ON notes BEGIN
  INSERT INTO note_fts(note_fts, rowid, title, datelog, content) VALUES('delete', old.id, old.title, old.datelog, old.content);
END;
CREATE TRIGGER notes_au AFTER UPDATE ON notes BEGIN
  INSERT INTO note_fts(note_fts, rowid, title, datelog, content) VALUES('delete', old.id, old.title, old.datelog, old.content);
  INSERT INTO note_fts(rowid, title, datelog, content) VALUES (new.id, new.title, new.datelog, new.content);
END;