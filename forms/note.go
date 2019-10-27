package forms

import (
	"fmt"
	// "log"
	"time"
	// "github.com/gotk3/gotk3/gtk"
	// "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/araddon/dateparse"
)

//Note - data structure
type Note struct {
	ID int `gorm:"primary_key"`
	Title string `gorm:"type:varchar(512);not null;unique_index"`
	Datelog time.Time
	Content string `gorm:"type:text"`
	URL string
	Flags string
	ReminderTicks int64
	Timestamp time.Time
	Readonly int8 `gorm:"default 0"`
	FormatTag []byte
	Econtent []byte
	AlertCount int8 `gorm:"type:int;default 0"`
	PixbufDict []byte
	TimeSpent int `gorm:"type:int;default 0"`
}

//NewNote - Create a new note object
func (n *Note) NewNote(in map[string]interface{}) {
	ct, ok := in["content"].(string)
	if !ok {
		fmt.Printf("INFO. content is empty\n")
		ct = ""
	}
	titleText, ok := in["title"].(string)
	if !ok {
		fmt.Printf("INFO No title provided, parse from content\n")
		if ct != ""{
			_l := len(ct)
			if _l >= 64 {_l = 64}
			titleText = ct[0:_l]
			n.Content = ct
		} else {
			fmt.Printf("INFO No content and title provided. Not creating note\n")
			return
		}
	}
	n.Content = ct
	n.Title = titleText

	if dateString, ok := in["datelog"]; ok {
		d, e := dateparse.ParseLocal(dateString.(string))
		if e != nil {
			fmt.Printf("ERROR Parse date string. Set to now - %v\n", e)
			d = time.Now()
		}
		n.Datelog = d
	} else {
		n.Datelog = time.Now()
	}

	n.Timestamp = time.Now()
	if e := DbConn.Save(n).Error; e != nil {
		fmt.Printf("ERROR saving note - %v\n", e)
	} else {
		n.Update(in)
	}
}

//Update - Update existing note
func (n *Note) Update(in map[string]interface{}) {
	if e := DbConn.Find(n, Note{ID: n.ID}).Error; e != nil {
		fmt.Printf("INFO Can not find the note to update - %v\n", e)
	}
	titleText, ok := in["title"].(string)
	if ok {
		n.Title = titleText
	}
	for k, v := range(in) {
		switch k {
		case "content":
			n.Content = v.(string)
		case "url":
			n.URL = v.(string)
		case "flags":
			n.Flags = v.(string)
		case "readonly":
			n.Readonly = v.(int8)
		case "alert_count":
			n.AlertCount = v.(int8)
		case "time_spent":
			n.TimeSpent = v.(int)
		}
	}
	n.Timestamp = time.Now()
	if e := DbConn.Save(n).Error; e != nil {
		fmt.Printf("ERROR saving note - %v\n", e)
	}
}

func (n *Note) String() string {return n.Title}

//Delete - Delete note
func (n *Note) Delete() {
	DbConn.Unscoped().Delete(&n)
	*n = Note{}
}