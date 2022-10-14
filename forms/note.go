package forms

import (
	"fmt"
	"time"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Note - data structure
type Note struct {
	// Do not embed gorm Model as we use our own ID as primary key
	// gorm.Model
	ID            int    `gorm:"primary_key"`
	Title         string `gorm:"type:varchar(512);not null;unique_index"`
	Datelog       int64  `gorm:"type:int"`
	Content       string `gorm:"type:text"`
	URL           string `gorm:"type:text"`
	Flags         string `gorm:"type:text"`
	ReminderTicks int64  `gorm:"type:int;default 0"`
	Timestamp     int64  `gorm:"type:int;default 0"`
	Readonly      int8   `gorm:"default 0"`
	FormatTag     []byte
	AlertCount    int8 `gorm:"type:int;default 0"`
	PixbufDict    []byte
	TimeSpent     int `gorm:"type:int;default 0"`
	LastTextMark  []byte
}

// NewNote - Create a new note object
func (n *Note) NewNote(in map[string]interface{}) {
	ct, ok := in["content"].(string)
	if !ok {
		// fmt.Printf("INFO. content is empty\n")
		ct = ""
	}
	titleText, ok := in["title"].(string)
	if !ok {
		// fmt.Printf("INFO No title provided, parse from content\n")
		if ct != "" {
			_l := len(ct)
			if _l >= 64 {
				_l = 64
			}
			titleText = ct[0:_l]
			n.Content = ct
		} else {
			// fmt.Printf("INFO No content and title provided. Not creating note\n")
			return
		}
	}
	n.Content = ct
	n.Title = titleText

	if dateData, ok := in["datelog"]; ok {
		switch v := dateData.(type) {
		case string:
			dateLog, e := time.Parse(DateLayout, v)
			if e != nil {
				fmt.Printf("ERROR can not parse date\n")
				n.Datelog = time.Now().UnixNano()
			} else {
				n.Datelog = dateLog.UnixNano()
			}
		case int64:
			n.Datelog = v
		}
	} else {
		n.Datelog = time.Now().UnixNano()
	}

	n.Timestamp = time.Now().UnixNano()

	if flags, ok := in["flags"]; ok {
		n.Flags = flags.(string)
	} else {
		n.Flags = ""
	}

	if url, ok := in["url"]; ok {
		n.URL = url.(string)
	} else {
		n.URL = ""
	}

	if readonly, ok := in["readonly"]; ok {
		n.Readonly = readonly.(int8)
	} else {
		n.Readonly = 0
	}

	if e := DbConn.Save(n).Error; e != nil {
		fmt.Printf("ERROR saving note - %v\n", e)
	}
}

// Update - Update existing note. Currently not need as the above already populate most data
func (n *Note) Update(in map[string]interface{}) {
	if e := DbConn.Find(n, Note{ID: n.ID}).Error; e != nil {
		fmt.Printf("INFO Can not find the note to update - %v\n", e)
	}
	titleText, ok := in["title"].(string)
	if ok {
		n.Title = titleText
	}
	for k, v := range in {
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
	if e := DbConn.Save(n).Error; e != nil {
		fmt.Printf("ERROR saving note - %v\n", e)
	}
}

func (n *Note) String() string { return n.Title }

// Delete - Delete note
func (n *Note) Delete() {
	DbConn.Unscoped().Delete(&n)
	*n = Note{}
}
