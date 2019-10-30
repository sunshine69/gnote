package forms

import (
	"time"
	"fmt"
	"github.com/jinzhu/gorm"
)

//DoMigration - once off
func DoMigration() {
	oldDBCon, _ := gorm.Open("sqlite3", "/home/stevek/Documents/clt.db")
	DbConn, _ = gorm.Open("sqlite3", "/home/stevek/.gnote.db")
	rows, e := oldDBCon.Raw(`SELECT note_id, title, cast(datelog as text), content, flags, url, timestamp, readonly FROM lsnote;`).Rows()
	if e != nil {
		fmt.Printf("ERROR - exec sql\n")
	}
	defer rows.Close()
	var readonly int8
	var count, note_id int
	var title, content, flags, url, timestamp string
	var datelog string
	var errorList []int
	for rows.Next() {
		rows.Scan(&note_id, &title, &datelog, &content, &flags, &url, &timestamp, &readonly)
		d, e := time.Parse("02-01-2006 15:04:05",datelog)
		if e != nil {
			d, e = time.Parse("Mon 2 Jan 2006 15:04:05 PM MST",datelog)
			if e != nil {
				d, e = time.Parse("02-01-2006 15:04",datelog)
				if e != nil {
					errorList = append(errorList, note_id)
					fmt.Printf("Error parse time %v\n", d)
				}
			}
		}
		note := Note{}
		note.NewNote(map[string]interface{} {
			"title": title,
			"datelog": d,
			"content": content,
			"flags": flags,
			"URL": url,
			"readonly": readonly,
		} )
		count = count + 1
	}
	fmt.Printf("%v\n", errorList)
}