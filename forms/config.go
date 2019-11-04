package forms

import (
	"os"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//DateLayout - global
var DateLayout string
//WebNotePassword
var WebNotePassword string
//WebNoteUser
var WebNoteUser string


//AppConfig - Application config struct
type AppConfig struct {
	gorm.Model
	// Section string `gorm:"type:varchar(128);unique_index:section_key"`
	Key string `gorm:"type:varchar(128);unique_index:section_key"`
	Val string
}

//DbConn - Global DB connection
var DbConn *gorm.DB

//SetupConfigDB - SetupDB. This is the initial point of config setup. Note init() does not work if it relies
//on DbConn as at the time the DBPATH is not yet available
func SetupConfigDB() {
	var err error
	dbPath := os.Getenv("DBPATH")
	fmt.Printf("Use dbpath %v\n", dbPath)
	DbConn, err = gorm.Open("sqlite3", dbPath)
	if err != nil {
	  panic("failed to connect database")
	}
	DbConn.AutoMigrate(&AppConfig{})
	DbConn.AutoMigrate(&Note{})
	DbConn.Exec("CREATE INDEX IF NOT EXISTS iTextContent ON notes(content COLLATE NOCASE);")

	// Example of loading a key dbpath
	// if err = DbConn.Find(&Config, AppConfig{Key: "dbpath"}).Error; err != nil {
	// 	log.Printf("Error can not load config table %v",err)
	// }
	// value := Config.Val
	DateLayout, _ = GetConfig("date_layout")
	WebNoteUser, _ = GetConfig("webnote_user")
}

//SetupDefaultConfig - Setup/reset default configuration set
func SetupDefaultConfig() {
	DbConn.Unscoped().Exec("DELETE FROM app_configs;")

	configSet := map[string]string{
		"config_created": "",
		"list_imap_account": "",
		"reminder_timer_interval" : "60",
		"checkmail" : "no",
		"clipboard_history_size" : "15",
		"run_startup_cmds" : "no",
		"last_font_desc" : "",
		"last_font_color" : "",
		"last_bgcolor" : "",
		"keywords" : "",
		"maxkwcount" : "20",
		"pnmain_win_pos" : "2202:54",
		"select_limit" : "250",
		"list_flags" : "TODO<|>IMPORTANT<|>URGENT",
		"recent_filter_cmd" : "",
		"maxcount_recent_filter_cmd" : "20",
		"webnote_password" : "",
		"window_size" : "429x503",
		"main_window_size" : "300x291",
		"default_font" : "None",
		"webnote_user": "msh.computing@gmail.com",
		"date_layout": "02-01-2006 15:04:05 MST",
	}
	for key, val := range(configSet) {
		fmt.Printf("Inserting %s - %s\n", key, val)
		if e := DbConn.Create(&AppConfig{Key: key, Val: val}).Error; e != nil {
			fmt.Printf("ERROR %v\n", e)
		}
	}
}

//GetConfig - by key and return value. Give second arg as default value.
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

//SetConfig - Set a config key with value
func SetConfig(key, val string) error {
	var cfg = AppConfig{}
	if e := DbConn.FirstOrInit(&cfg, AppConfig{Key: key}).Error; e != nil{
		return e
	}
	cfg.Val = val
	if e := DbConn.Save(&cfg).Error; e != nil {
		return e
	}
	return nil
}

//DeleteConfig - delete the config key
func DeleteConfig(key string) error {
	var cfg = AppConfig{}
	if e := DbConn.Find(&cfg, AppConfig{Key: key}).Error; e != nil {
		return e
	}
	return DbConn.Unscoped().Delete(&cfg).Error
}