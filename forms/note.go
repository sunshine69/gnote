package forms

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Note - data structure
type Note struct {
	ID            int    `db:"id"`
	Title         string `db:"title"`
	Datelog       int64  `db:"datelog"`
	Content       string `db:"content"`
	URL           string `db:"url"`
	Flags         string `db:"flags"`
	ReminderTicks int64  `db:"reminder_ticks"`
	Timestamp     int64  `db:"timestamp"`
	Readonly      int8   `db:"readonly"`
	FormatTag     []byte `db:"format_tag"`
	AlertCount    int8   `db:"alert_count"`
	PixbufDict    []byte `db:"pixbuf_dict"`
	TimeSpent     int    `db:"time_spent"`
	LastTextMark  []byte `db:"last_text_mark"`
	Language      string `db:"language"`
	FileExt       string `db:"file_ext"`
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

	if lang, ok := in["language"]; ok {
		n.Language = lang.(string)
	} else {
		n.Language = "markdown"
	}

	if ext, ok := in["file-ext"]; ok {
		n.FileExt = ext.(string)
	} else {
		n.FileExt = ".md"
	}

	if _, e := DbConn.NamedExec(`INSERT INTO notes (title, datelog, content,url, flags , reminder_ticks, timestamp, readonly, format_tag , alert_count , pixbuf_dict , time_spent , last_text_mark , language , file_ext ) VALUES(:title, :datelog, :content,:url, :flags , :reminder_ticks,:timestamp, :readonly, :format_tag , :alert_count , :pixbuf_dict , :time_spent , :last_text_mark , :language , :file_ext) ON CONFLICT(title) DO UPDATE SET title=excluded.title, datelog=excluded.datelog, content=excluded.content, readonly=excluded.readonly, timestamp=excluded.timestamp, url=excluded.url, flags=excluded.flags, reminder_ticks=excluded.reminder_ticks, format_tag=excluded.format_tag, alert_count=excluded.alert_count, pixbuf_dict=excluded.pixbuf_dict,language=excluded.language, file_ext=excluded.file_ext`, n); e != nil {
		fmt.Printf("ERROR saving note - %v\n", e)
	}
}

// Update - Update existing note. Currently not need as the above already populate most data
func (n *Note) Update(in map[string]interface{}) {
	if e := DbConn.Get(n, "SELECT * FROM notes WHERE id=$1", n.ID); e != nil && errors.Is(e, sql.ErrNoRows) {
		fmt.Printf("INFO Can not find the note to update - %v\n", e)
		return
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
	if _, e := DbConn.NamedExec(`UPDATE note SET content=:content, url=:url, flags=:flags, readonly=:readonly, alert_count=:alert_count, time_spent=:time_spent WHERE id=:id`, n); e != nil {
		fmt.Printf("[ERROR] saving note - %s\n", e.Error())
	}
	return
}

func (n *Note) String() string { return n.Title }

// Delete - Delete note
func (n *Note) Delete() {
	if _, e := DbConn.Exec(`DELETE FROM notes WHERE id=:id`, n); e != nil {
		fmt.Printf("[ERROR] delete note - %s\n", e.Error())
	} else {
		*n = Note{}
	}
}
