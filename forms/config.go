package forms

import (
	"fmt"
	"net/http/cookiejar"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	gorm.Model
	// Section string `gorm:"type:varchar(128);unique_index:section_key"`
	Key string `gorm:"type:varchar(128);unique_index:section_key"`
	Val string `gorm:"type:text"`
}

// DbConn - Global DB connection
var DbConn *gorm.DB

// SetupConfigDB - SetupDB. This is the initial point of config setup. Note init() does not work if it relies
// on DbConn as at the time the DBPATH is not yet available
func SetupConfigDB() {
	var err error
	dbPath := os.Getenv("DBPATH")
	// fmt.Printf("Use dbpath %v\n", dbPath)
	DbConn, err = gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("failed to connect database")
	}
	DbConn.AutoMigrate(&AppConfig{})
	DbConn.AutoMigrate(&Note{})
	setupSQL := `
CREATE VIRTUAL TABLE IF NOT EXISTS note_fts USING fts5(title, datelog, content, content='notes', content_rowid='id');

CREATE TRIGGER IF NOT EXISTS notes_ai AFTER INSERT ON notes BEGIN
  INSERT INTO note_fts(rowid, title, datelog, content) VALUES (new.id, new.title, new.datelog, new.content);
END;

CREATE TRIGGER IF NOT EXISTS notes_ad AFTER DELETE ON notes BEGIN
  INSERT INTO note_fts(note_fts, rowid, title, datelog, content) VALUES('delete', old.id, old.title, old.datelog, old.content);
END;

CREATE TRIGGER IF NOT EXISTS notes_au AFTER UPDATE ON notes BEGIN
  INSERT INTO note_fts(note_fts, rowid, title, datelog, content) VALUES('delete', old.id, old.title, old.datelog, old.content);
  INSERT INTO note_fts(rowid, title, datelog, content) VALUES (new.id, new.title, new.datelog, new.content);
END;
`
	DbConn.Exec(setupSQL)
	// DbConn.Exec("CREATE INDEX IF NOT EXISTS iTextContent ON notes(content COLLATE NOCASE);")

	// Example of loading a key dbpath
	// if err = DbConn.Find(&Config, AppConfig{Key: "dbpath"}).Error; err != nil {
	// 	log.Printf("Error can not load config table %v",err)
	// }
	// value := Config.Val
	DateLayout, _ = GetConfig("date_layout")
	WebNoteUser, _ = GetConfig("webnote_user")
	CreateDataNoteLangFileExt()
	CreateDataNoteListOfLanguageSupport()
}

// SetupDefaultConfig - Setup/reset default configuration set
func SetupDefaultConfig() {
	DbConn.Unscoped().Exec("DELETE FROM app_configs;")

	configSet := map[string]string{
		"config_created":    "",
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
		if e := DbConn.Create(&AppConfig{Key: key, Val: val}).Error; e != nil {
			fmt.Printf("ERROR %v\n", e)
		}
	}
}

// GetConfig - by key and return value. Give second arg as default value.
func GetConfig(key ...string) (string, error) {
	var cfg = AppConfig{}
	err := DbConn.Find(&cfg, AppConfig{Key: key[0]}).Error
	if err != nil {
		if len(key) == 2 {
			return key[1], nil
		} else {
			return "", err
		}
	} else {
		return cfg.Val, err
	}
}

// SetConfig - Set a config key with value
func SetConfig(key, val string) error {
	var cfg = AppConfig{}
	if e := DbConn.FirstOrInit(&cfg, AppConfig{Key: key}).Error; e != nil {
		return e
	}
	cfg.Val = val
	if e := DbConn.Save(&cfg).Error; e != nil {
		return e
	}
	return nil
}

// DeleteConfig - delete the config key
func DeleteConfig(key string) error {
	var cfg = AppConfig{}
	if e := DbConn.Find(&cfg, AppConfig{Key: key}).Error; e != nil {
		return e
	}
	return DbConn.Unscoped().Delete(&cfg).Error
}

// Populate some notes needed for data lookup - used by other part of the app
// Currently we store the language / file extention data but in the future we might store more
// This note is used to lookup Language => File Extention so we can save the note to file with correct extension in note-pad.go and note-search.go
func CreateDataNoteLangFileExt() {
	note := Note{}
	if e := DbConn.FirstOrInit(&note, Note{Title: "CreateDataNoteLangFileExt"}).Error; e != nil {
		fmt.Printf("INFO Can not create data note  CreateDataNoteLangFileExt %s\n", e.Error())
		return
	}
	defer DbConn.Save(&note)
	if note.Content == "" {
		// Fetch it so we do not waste memory by adding this resource to go-bindata
		jsonText, err := u.Curl("GET", "https://raw.githubusercontent.com/sunshine69/gnote/gtksourceview/CreateDataNoteLangFileExt.json", "", "", []string{})
		if u.CheckErrNonFatal(err, "CreateDataNoteLangFileExt GET") != nil {
			fmt.Println("Error fetching CreateDataNoteLangFileExt. You can manually search the note with title CreateDataNoteLangFileExt and insert the content yourself. The content is from the this repo project github")
			return
		}
		note.Content = jsonText
		note.Readonly = 1
	}
}

// Parse the language string that current gtksourceview support and save it to a note so we can check against it
// The list is created from this command on linux ls /usr/share/gtksourceview-3.0/language-specs/ | sed 's/.lang//'
func CreateDataNoteListOfLanguageSupport() {
	note := Note{}
	if e := DbConn.FirstOrInit(&note, Note{Title: "CreateDataNoteListOfLanguageSupport"}).Error; e != nil {
		fmt.Printf("INFO Can not create data note  CreateDataNoteListOfLanguageSupport %s\n", e.Error())
		return
	}
	defer DbConn.Save(&note)
	if note.Content == "" {
		// Fetch it so we do not waste memory by adding this resource to go-bindata
		jsonText, err := u.Curl("GET", "https://raw.githubusercontent.com/sunshine69/gnote/gtksourceview/CreateDataNoteListOfLanguageSupport.json", "", "", []string{})
		if u.CheckErrNonFatal(err, "CreateDataNoteListOfLanguageSupport GET") != nil {
			fmt.Println("Error fetching CreateDataNoteListOfLanguageSupport. You can manually search the note with title CreateDataNoteListOfLanguageSupport and insert the content yourself. The content is from the this repo project github")
			return
		}
		note.Content = jsonText
		note.Readonly = 1
	}
}
