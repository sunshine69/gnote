package forms

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	u "github.com/sunshine69/golang-tools/utils"
)

func DoMigrationV1(oldDB, newDB string) {
	// oldDBCon, _ := gorm.Open("sqlite3", "/home/stevek/Documents/clt.db")
	oldDBCon, _ := gorm.Open("sqlite3", oldDB)
	newDbConn, _ := gorm.Open("sqlite3", newDB)
	// rows, e := oldDBCon.Raw(`SELECT id, title, datelog, content, flags, url, timestamp, readonly FROM notes;`).Rows()
	// if e != nil {
	// 	fmt.Printf("ERROR - exec sql\n")
	// }
	// defer rows.Close()
	oldDBCon.AutoMigrate(&Note{})
	oldDBCon.AutoMigrate(&AppConfig{})
	newDbConn.AutoMigrate(&Note{})
	newDbConn.AutoMigrate(&AppConfig{})
	oldRows, err := oldDBCon.Raw("select ID, datelog, title, flags, timestamp, readonly, content, url, reminder_ticks, timestamp, format_tag, alert_count, pixbuf_dict, time_spent from notes;").Rows()
	if err != nil {
		log.Fatalf("[ERROR] %v\n", err)
	}
	defer oldRows.Close()
	count := 0
	newDbConn.Begin().New()
	for oldRows.Next() {
		// if count > 500 {
		// 	break
		// }
		_newNote := Note{}
		oldRows.Scan(&_newNote.ID, &_newNote.Datelog, &_newNote.Title, &_newNote.Flags, &_newNote.Timestamp, &_newNote.Readonly, &_newNote.Content, &_newNote.URL, &_newNote.ReminderTicks, &_newNote.Timestamp, &_newNote.FormatTag, &_newNote.AlertCount, &_newNote.PixbufDict, &_newNote.TimeSpent)
		// fmt.Printf("note: %v\n", _newNote)
		newDbConn.Create(&_newNote)
		// newDbConn.Save(&_newNote)
		count++
	}
	newDbConn.Commit()

}

// DoMigration - once off - this is old
func DoMigration(oldDB, newDB string) {
	// oldDBCon, _ := gorm.Open("sqlite3", "/home/stevek/Documents/clt.db")
	oldDBCon, _ := gorm.Open("sqlite3", oldDB)
	DbConn, _ = gorm.Open("sqlite3", newDB)
	// rows, e := oldDBCon.Raw(`SELECT note_id, title, cast(datelog as text), content, flags, url, timestamp, readonly FROM lsnote;`).Rows()
	rows, e := oldDBCon.Raw(`SELECT id, title, cast(datelog as text), content, flags, url, timestamp, readonly FROM notes;`).Rows()
	if e != nil {
		fmt.Printf("ERROR - exec sql\n")
	}
	defer rows.Close()
	var readonly int8
	var count, note_id int
	var title, content, flags, url, timestamp string
	var datelog string
	var errorList []int
	ptn := regexp.MustCompile(`^\d\d\-\d\d\-\d\d\d\d \d\d\:\d\d\:\d\d$`)
	ptn2 := regexp.MustCompile(`^\d\d\-\d\d\-\d\d\d\d \d\d\:\d\d$`)
	ptn3 := regexp.MustCompile(`(AEST|NZDT|NZST)$`)
	ptn4 := regexp.MustCompile(`(\d\d\/\d\d\/\d\d\d\d [\d]+\:\d\d:\d\d) ([pa]\.m\.)`)
	ptn5 := regexp.MustCompile(`(.*) \-\d\d\d\d`)

	for rows.Next() {
		var d time.Time
		rows.Scan(&note_id, &title, &datelog, &content, &flags, &url, &timestamp, &readonly)
		_dData, e := strconv.ParseInt(datelog, 10, 64)
		if e == nil {
			d = u.NsToTime(_dData)
		} else {
			if ptn.FindStringSubmatch(datelog) != nil {
				datelog1 := fmt.Sprintf("%s AEST", datelog)
				d, e = time.Parse("02-01-2006 15:04:05 MST", datelog1)
				if e != nil {
					log.Fatalf("ERROR unexpected 1 %v\n", datelog1)
				}
			} else if ptn2.FindStringSubmatch(datelog) != nil {
				datelog1 := fmt.Sprintf("%s AEST", datelog)
				d, e = time.Parse("02-01-2006 15:04 MST", datelog1)
				if e != nil {
					log.Printf("ERROR unexpected 2 %v\n", datelog1)
					continue
				}
			} else if ptn3.FindStringSubmatch(datelog) != nil {
				d, e = time.Parse("Mon 2 Jan 2006 15:04:05 PM MST", datelog)
				if e != nil {
					log.Fatalf("ERROR unexpected 3 %v\n", datelog)
				}
			} else if matches := ptn4.FindStringSubmatch(datelog); matches != nil {
				_tmp := strings.ReplaceAll(matches[2], ".", "")
				_tmp = strings.ToUpper(_tmp)
				datelog1 := fmt.Sprintf("%s %s NZST", matches[1], _tmp)
				d, e = time.Parse("02/01/2006 15:04:05 PM MST", datelog1)
				if e != nil {
					log.Fatalf("ERROR unexpected 4 %v\n", datelog1)
				}
			} else if matches := ptn5.FindStringSubmatch(datelog); matches != nil {
				datelog1 := fmt.Sprintf("%s AEST", matches[1])
				d, e = time.Parse("Mon, 02 Jan 2006 15:04:05 MST", datelog1)
				if e != nil {
					log.Fatalf("ERROR unexpected 5 %v\n", datelog1)
				}
			} else {
				datelog1 := fmt.Sprintf("%s NZST", datelog)
				d, e = time.Parse("Mon Jan 2 15:04:05 2006 MST", datelog1)
				if e != nil {
					log.Fatalf("ERROR Should not reach here %v\n", datelog)
				}
			}
		}

		note := Note{}
		note.NewNote(map[string]interface{}{
			"title":    title,
			"datelog":  d.UnixNano(),
			"content":  content,
			"flags":    flags,
			"URL":      url,
			"readonly": readonly,
		})
		count = count + 1
		// if count == 100 { break }
	}
	fmt.Printf("%v\n", errorList)
}
