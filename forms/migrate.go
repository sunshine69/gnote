package forms

import (
	"strings"
	"log"
	"regexp"
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
	ptn := regexp.MustCompile(`^\d\d\-\d\d\-\d\d\d\d \d\d\:\d\d\:\d\d$`)
	ptn2 := regexp.MustCompile(`^\d\d\-\d\d\-\d\d\d\d \d\d\:\d\d$`)
	ptn3 := regexp.MustCompile(`(AEST|NZDT|NZST)$`)
	ptn4 := regexp.MustCompile(`(\d\d\/\d\d\/\d\d\d\d [\d]+\:\d\d:\d\d) ([pa]\.m\.)`)
	ptn5 := regexp.MustCompile(`(.*) \-\d\d\d\d`)

	for rows.Next() {
		var d time.Time
		rows.Scan(&note_id, &title, &datelog, &content, &flags, &url, &timestamp, &readonly)
		if ptn.FindStringSubmatch(datelog) != nil {
			datelog1 := fmt.Sprintf("%s AEST", datelog)
			d, e = time.Parse("02-01-2006 15:04:05 MST",datelog1)
			if e != nil {
				log.Fatalf("ERROR unexpected 1 %v\n", datelog1)
			}
		} else if ptn2.FindStringSubmatch(datelog) != nil {
			datelog1 := fmt.Sprintf("%s AEST", datelog)
			d, e = time.Parse("02-01-2006 15:04 MST",datelog1)
			if e != nil {
				log.Printf("ERROR unexpected 2 %v\n", datelog1)
				continue
			}
		} else if ptn3.FindStringSubmatch(datelog) != nil {
			d, e = time.Parse("Mon 2 Jan 2006 15:04:05 PM MST",datelog)
			if e != nil {
				log.Fatalf("ERROR unexpected 3 %v\n", datelog)
			}
		} else if matches := ptn4.FindStringSubmatch(datelog); matches != nil {
			_tmp := strings.ReplaceAll(matches[2], ".", "")
			_tmp = strings.ToUpper(_tmp)
			datelog1 := fmt.Sprintf("%s %s NZST", matches[1], _tmp)
			d, e = time.Parse("02/01/2006 15:04:05 PM MST",datelog1)
			if e != nil {
				log.Fatalf("ERROR unexpected 4 %v\n", datelog1)
			}
		} else if matches := ptn5.FindStringSubmatch(datelog); matches != nil {
			datelog1 := fmt.Sprintf("%s AEST", matches[1])
			d, e = time.Parse("Mon, 02 Jan 2006 15:04:05 MST",datelog1)
			if e != nil {
				log.Fatalf("ERROR unexpected 5 %v\n", datelog1)
			}
		} else {
			datelog1 := fmt.Sprintf("%s NZST", datelog)
			d, e = time.Parse("Mon Jan 2 15:04:05 2006 MST",datelog1)
			if e != nil {
				log.Fatalf("ERROR Should not reach here %v\n", datelog)
			}
		}

		note := Note{}
		note.NewNote(map[string]interface{} {
			"title": title,
			"datelog": d.UnixNano(),
			"content": content,
			"flags": flags,
			"URL": url,
			"readonly": readonly,
		} )
		count = count + 1
	}
	fmt.Printf("%v\n", errorList)
}