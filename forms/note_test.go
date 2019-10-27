package forms

import (
	"fmt"
	"testing"
	"os"
)

var p = fmt.Printf

func TestNote (t *testing.T) {
	os.Setenv("DBPATH", "test.db")
	SetupConfigDB()
	DbConn.Unscoped().Exec("DELETE from notes;")
	var n = &Note{}
	n.NewNote( map[string]interface{} {"title": "New note1", "content": "Content note 2"})
	if n.Title != "New note1" {p("FAIL Title New note1\n")}

	n = &Note{}; n.NewNote(map[string]interface{} {"content": "Content note 3"})
	if n.Title != "Content note 3" {p("FAIL Content note 3\n")}

	n = &Note{}; DbConn.Find(n, Note{Title: "New note1"})
	if n.Title != "New note1" {p("FAIL Can not find note Title New note1\n")}

	n = &Note{}; DbConn.Find(n, "content like ? OR title like ?", "%note 2%", "not found")
	p("Found %s\n", n.String())

	n = &Note{};  DbConn.Find(n, "content like '%note 3%' OR content like 'note 1'")
	p("Found %s\n", n.String())

	ns := []Note{}
	tokens := []string {"%note 3%", "%note 2%"}
	_l := len(tokens)
	q := ""
	for i, t := range(tokens) {
		if i == _l - 1 {
			q = fmt.Sprintf("%v content like '%v' ORDER BY ID", q, t)
		} else {
			q = fmt.Sprintf("%v content like '%v' OR ", q, t)
		}
	}
	q = fmt.Sprintf("SELECT * from notes WHERE %v", q)
	// fmt.Println(q)
	if e := DbConn.Exec(q).Find(&ns).Error; e != nil {
		fmt.Printf("ERROR %v\n", e)
	} else {
		for r, n := range(ns){
			fmt.Printf("Found note %v - %v\n", r, n.Title)
		}
	}
	n = &Note{}; n.NewNote( map[string]interface{} {"title": "note4", "content": "Content note 4"})
	if n.Title != "note4" {p("FAIL note4\n")}
	n.Delete()
	if n.Title != "" {p("FAIL 1 delete note4 - %v\n", n.Title)}
	DbConn.Find(n, Note{Title: "note4"})
	if n.Title != "" {p("FAIL 2 delete note4\n")}
}