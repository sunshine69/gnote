package forms

import (
	"errors"
	"fmt"
	"net/http/cookiejar"
	"os"
	"time"

	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/mutecomm/go-sqlcipher/v4"
	u "github.com/sunshine69/golang-tools/utils"
)

// DateLayout - global
var DateLayout string

// WebNotePassword
var WebNotePassword string

// WebNoteUser
var WebNoteUser string
var CookieJar *cookiejar.Jar

// AppConfig - Application config struct
type AppConfig struct {
	ID  int    `db:"id"`
	Key string `db:"key"`
	Val string `db:"val"`
}

// DbConn - Global DB connection
var DbConn *sqlx.DB

// SetupConfigDB - SetupDB. This is the initial point of config setup. Note init() does not work if it relies
// on DbConn as at the time the DBPATH is not yet available
func SetupConfigDB() {
	// dbPath := os.Getenv("DBPATH")
	// fmt.Printf("Use dbpath %v\n", dbPath)
	// DbConn = sqlx.MustConnect("sqlite3", dbPath)

	setupSQL := `
CREATE TABLE app_configs(id integer PRIMARY KEY AUTOINCREMENT,key varchar(128) UNIQUE, val text);

CREATE TABLE notes(id integer PRIMARY KEY AUTOINCREMENT,title varchar(512) NOT NULL UNIQUE, datelog integer, content text,url text, flags text, reminder_ticks integer,timestamp integer DEFAULT CURRENT_TIMESTAMP, readonly integer DEFAULT 0, format_tag blob, alert_count integer DEFAULT 0, pixbuf_dict blob, time_spent integer DEFAULT 0, last_text_mark blob, language text, file_ext text);

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
`
	if _, err := DbConn.Exec(setupSQL); err == nil {
		DateLayout, _ = GetConfig("date_layout")
		WebNoteUser, _ = GetConfig("webnote_user")
		CreateDataNoteLangFileExt()
		CreateDataNoteListOfLanguageSupport()
	}
}

// SetupDefaultConfig - Setup/reset default configuration set
func SetupDefaultConfig() {
	DbConn.Exec("DELETE FROM app_configs;")

	configSet := map[string]string{
		"config_created":    "1",
		"pnmain_win_pos":    "2202:54",
		"select_limit":      "250",
		"list_flags":        "TODO<|>IMPORTANT<|>URGENT",
		"recent_filter_cmd": "",
		"window_size":       "429x503",
		"main_window_size":  "300x291",
		"date_layout":       "02-01-2006 15:04:05 MST",
		"webnote_url":       "",
		"debug":             "false",
	}
	for key, val := range configSet {
		fmt.Printf("Inserting %s - %s\n", key, val)
		if _, e := DbConn.NamedExec("INSERT INTO app_configs(key, val) VALUES (:key, :val)", &AppConfig{Key: key, Val: val}); e != nil {
			fmt.Printf("ERROR %v\n", e)
		}
	}
}

// GetConfig - by key and return value. Give second arg as default value. This runs very first
func GetConfig(key ...string) (string, error) {
	DbConn = sqlx.MustConnect("sqlite3", os.Getenv("DBPATH"))
	var cfg = AppConfig{}
	// If use Get - sqlx requires strict col mapping to the struct u put into and the old db might have
	// different schema for app_configs. The QueryRow and scan wont have that problem
	row := DbConn.QueryRow("SELECT key, val FROM app_configs WHERE `key` = $1", key[0])
	if row.Err() != nil {
		// fmt.Printf("[DEBUG 0] %s\n", row.Err().Error())
		if len(key) == 2 {
			return key[1], nil
		} else {
			return "", row.Err()
		}
	} else {
		row.Scan(&cfg.Key, &cfg.Val)
		// fmt.Printf("[DEBUG 1] %v - Val %v\n", row.Err(), cfg)
		return cfg.Val, nil
	}
	
}

// SetConfig - Set a config key with value
func SetConfig(key, val string) error {
	_, err := DbConn.NamedExec("INSERT INTO app_configs(key, val) VALUES (:key, :val) ON CONFLICT(key) DO UPDATE SET val=excluded.val", AppConfig{Key: key, Val: val})
	return err
}

// DeleteConfig - delete the config key
func DeleteConfig(key string) error {
	if _, e := DbConn.NamedExec("DELETE FROM app_configs WHERE key = :key", AppConfig{Key: key}); e != nil {
		return e
	}
	return nil
}

// Populate some notes needed for data lookup - used by other part of the app
// Currently we store the language / file extention data but in the future we might store more
func CreateDataNote(title string, fetchDataUrl string) {
	prepareNote := func(note *Note) {
		// Fetch it so we do not waste memory by adding this resource to go-bindata
		jsonText, err := u.Curl("GET", fetchDataUrl, "", "", []string{})
		if u.CheckErrNonFatal(err, title+"CreateDataNote GET") != nil {
			fmt.Printf("Error fetching. You can manually search the note with title %s and insert the content yourself. The content is from the this repo project github", title)
			return
		}
		note.Content = jsonText
		datelog := time.Now().UnixNano()
		note.Readonly, note.Datelog, note.Timestamp = 1, datelog, datelog
	}

	note := Note{Title: title, Datelog: time.Now().Unix()}
	if e := DbConn.Get(&note, "SELECT * FROM notes WHERE title = :title", title); errors.Is(e, sql.ErrNoRows) {
		prepareNote(&note)
		_, err := DbConn.NamedExec("INSERT INTO notes(title, datelog, content, readonly, timestamp) VALUES(:title, :datelog, :content, :readonly, :timestamp) ON CONFLICT(title) DO UPDATE SET title=excluded.title, datelog=excluded.datelog, content=excluded.content, readonly=excluded.readonly, timestamp=excluded.timestamp", note)
		u.CheckErr(err, title+" INSERT")
		return
	}
}

// This note is used to lookup Language => File Extention so we can save the note to file with correct extension in note-pad.go and note-search.go
func CreateDataNoteLangFileExt() {
	CreateDataNote("CreateDataNoteLangFileExt", "https://raw.githubusercontent.com/sunshine69/gnote/gtksourceview/CreateDataNoteLangFileExt.json")
}

// Parse the language string that current gtksourceview support and save it to a note so we can check against it
// The list is created from this command on linux ls /usr/share/gtksourceview-3.0/language-specs/ | sed 's/.lang//'
func CreateDataNoteListOfLanguageSupport() {
	CreateDataNote("CreateDataNoteListOfLanguageSupport", "https://raw.githubusercontent.com/sunshine69/gnote/gtksourceview/CreateDataNoteListOfLanguageSupport.json")
}
