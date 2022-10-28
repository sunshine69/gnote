package forms

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/cjoudrey/gluahttp"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/kohkimakimoto/gluayaml"
	"github.com/sunshine69/gluare"
	u "github.com/sunshine69/golang-tools/utils"
	gopherjson "github.com/sunshine69/gopher-json"
	lua "github.com/yuin/gopher-lua"
)

// NoteSearch - GUI related
type NoteSearch struct {
	w                 *gtk.Window
	builder           *gtk.Builder
	np                *NotePad
	isIcase           bool
	isCmdFilter       bool
	isOutputToNewNote bool
	searchBox         *gtk.SearchEntry
	replaceBox        *gtk.Entry
	m1                *gtk.TextMark
	m2                *gtk.TextMark
	curIter           *gtk.TextIter
}

func (ns *NoteSearch) NoteFindIcase(o *gtk.CheckButton) {
	ns.isIcase = o.GetActive()
	ns.searchBox.GrabFocus()
}

func (ns *NoteSearch) CommandFilter(o *gtk.CheckButton) {
	ns.isCmdFilter = o.GetActive()
	if ns.isCmdFilter {
		lastCmd, _ := GetConfig("last_cmd_filter", "perl -pe 's///'")
		ns.searchBox.SetText(lastCmd)
		btnFind := GetButton(ns.builder, "find_btn")
		btnFind.SetLabel("Cmd")
		ns.replaceBox.SetText("<EXTERNAL_CMD_OUPUT>")

	} else {
		ns.searchBox.SetText("")
		btnFind := GetButton(ns.builder, "find_btn")
		btnFind.SetLabel("Find")
	}
	ns.searchBox.GrabFocus()
}

func (ns *NoteSearch) OutputToNewNote(o *gtk.CheckButton) {
	ns.isOutputToNewNote = o.GetActive()
	ns.searchBox.GrabFocus()
}

func GetNoteFromLua(L *lua.LState) int {
	title := L.ToString(1) /* get argument */
	note := Note{}
	DbConn.First(&note, Note{Title: title})
	L.Push(lua.LString(u.JsonDump(note, ""))) /* push result */
	return 1                                  /* number of results */
}

func SearchNotesFromLua(L *lua.LState) int {
	sqlWhereList := L.ToString(1) // lua supply arg like this [[ WHERE content LIKE '%text%' ]]
	if !strings.HasPrefix(sqlWhereList, "WHERE") {
		L.Push(lua.LString("ERROR - the arg is a string and should started with 'WHERE'"))
		return 1
	}
	sql := "SELECT id FROM notes " + sqlWhereList
	rows, err := DbConn.Raw(sql).Rows()
	if err != nil {
		fmt.Printf("ERROR - exec sql\n")
		L.Push(lua.LString(err.Error()))
		return 1
	}
	defer rows.Close()
	oNoteList := []Note{}
	for rows.Next() {
		_n, nid := Note{}, 0
		rows.Scan(&nid)
		DbConn.First(&_n, nid)
		oNoteList = append(oNoteList, _n)
	}
	L.Push(lua.LString(u.JsonDump(oNoteList, "")))
	return 1
}

func UpdateNotesFromLua(L *lua.LState) int {
	sql := L.ToString(1) // lua supply arg like this [[ UPDATE notes SET content = 'new content' WHERE xxx ]]
	if !strings.HasPrefix(sql, "UPDATE") {
		L.Push(lua.LString("ERROR - the arg is a string and should started with 'UPDATE'"))
		return 1
	}
	DbConn.Exec(sql)
	L.Push(lua.LString(fmt.Sprintf("OK %d rows affected", DbConn.RowsAffected)))
	return 1
}

func RunLuaFile(luaFileName string) string {
	old := os.Stdout // keep backup of the real stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("re", gluare.Loader)
	L.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)
	L.PreloadModule("yaml", gluayaml.Loader)
	L.PreloadModule("json", gopherjson.Loader)
	L.SetGlobal("get_note", L.NewFunction(GetNoteFromLua))
	L.SetGlobal("search_notes", L.NewFunction(SearchNotesFromLua))
	L.SetGlobal("update_notes", L.NewFunction(UpdateNotesFromLua))

	err := L.DoFile(luaFileName)
	if err := u.CheckErrNonFatal(err, "Lua DoFile"); err != nil {
		fmt.Print(err.Error())
	}

	w.Close()
	os.Stdout = old
	out := <-outC
	return out
}

func (ns *NoteSearch) FindText() bool {
	buf := ns.np.buff
	keyword, _ := ns.searchBox.GetText()
	searchFlag := gtk.TEXT_SEARCH_TEXT_ONLY
	var foundIter1, foundIter2 *gtk.TextIter
	var ok bool = true
	var output = false
	ns.curIter = buf.GetIterAtMark(buf.GetInsert())

	if ns.isCmdFilter { //run external command and replace the note/selection with output
		text, startI, endI := ns.np.GetSelection()
		replaceWith, e := ns.replaceBox.GetText()
		if e != nil {
			MessageBox(fmt.Sprintf("ERROR %v\n", e))
			return false
		}
		outStr := ""
		if replaceWith == "<EXTERNAL_CMD_OUPUT>" {
			_tmpF, _ := ioutil.TempFile("", fmt.Sprintf("gnote-*%s", ns.np.FileExt))
			_tmpF.Write([]byte(text))
			err := _tmpF.Close()
			u.CheckErrNonFatal(err, "run-command can not close tmp file")
			cmdText := fmt.Sprintf("%s %s", keyword, _tmpF.Name())

			commandList := strings.Fields(cmdText)

			if commandList[0] == "gopher-lua" {
				// Use internal lua VM to run the code
				outStr = RunLuaFile(_tmpF.Name())
			} else {
				cmd := exec.Command(commandList[0], commandList[1:]...)
				cmd.Env = append(os.Environ())
				stdoutStderr, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Printf("DEBUG E %v\n", err)
				}
				outStr = string(stdoutStderr)
			}
			os.Remove(_tmpF.Name())
			SetConfig("last_cmd_filter", keyword)

			if ns.isOutputToNewNote {
				_np := NewNotePad(-1)
				_np.app = ns.np.app
				_np.buff.SetText(outStr)
				_np.wTitle.SetText(fmt.Sprintf("note: %s result of cmd: %s ", ns.np.Title, keyword))
				return false //stop other actions
			}
		} else {
			if ptn, e := regexp.Compile(keyword); e == nil {
				var newTxt []string
				for _, lineTxt := range strings.Split(text, "\n") {
					newTxt = append(newTxt, ptn.ReplaceAllString(lineTxt, replaceWith))
				}
				outStr = strings.Join(newTxt, "\n")
			} else {
				MessageBox(fmt.Sprintf("ERROR %s\n", e.Error()))
			}

		}
		if !ns.isOutputToNewNote {
			buf := ns.np.buff
			buf.SelectRange(startI, endI)
			buf.DeleteSelection(true, true)
			buf.InsertAtCursor(outStr)
			//Not sure why the curIter is invalid after running. Need to get back otherwise crash

			return false //stop other actions
		}
	} else {
		if ns.isIcase {
			searchFlag = gtk.TEXT_SEARCH_CASE_INSENSITIVE
		}
		if ns.m2 != nil {
			buf.PlaceCursor(buf.GetIterAtMark(ns.m2))
			ns.curIter = buf.GetIterAtMark(buf.GetInsert())
		}
		foundIter1, foundIter2, ok = ns.curIter.ForwardSearch(keyword, searchFlag, nil)

		if ok {
			ns.np.textView.ScrollToIter(foundIter1, 0, true, 0, 0)
			buf.SelectRange(foundIter1, foundIter2)
			ns.m1, ns.m2 = buf.CreateMark("s1", foundIter1, false), buf.CreateMark("s2", foundIter2, false)
			output = true
		} else {
			if !ok {
				MessageBox("Search text not found. Will reset iter")
				buf.PlaceCursor(buf.GetStartIter())
				ns.curIter = buf.GetStartIter()
				ns.m1, ns.m2 = nil, nil
			}
		}
	}
	return output
}

// NoteFindText -
func (ns *NoteSearch) NoteFindText() {
	ns.FindText()
}

// NoteReplaceText -
func (ns *NoteSearch) NoteReplaceText(o *gtk.Button) {
	buf := ns.np.buff

	if buf.GetHasSelection() || ns.FindText() {
		buf.DeleteSelection(true, true)
		_rp := GetEntry(ns.builder, "replace_text")
		replaceText, _ := _rp.GetText()
		buf.InsertAtCursor(replaceText)
	}

}

// NoteReplaceAll -
func (ns *NoteSearch) NoteReplaceAll(o *gtk.Button) {
	buf := ns.np.buff

	for buf.GetHasSelection() || ns.FindText() {
		buf.DeleteSelection(true, true)
		_rp := GetEntry(ns.builder, "replace_text")
		replaceText, _ := _rp.GetText()
		buf.InsertAtCursor(replaceText)
	}
}

func (ns *NoteSearch) KeyPressed(o interface{}, ev *gdk.Event) {
	keyEvent := &gdk.EventKey{ev}
	if keyEvent.State()&gdk.CONTROL_MASK > 0 { //Control key pressed
		switch keyEvent.KeyVal() {
		case gdk.KeyvalFromName("q"):
			ns.w.Close()
		}
	}
}

func (ns *NoteSearch) ResetIter() {
	//Crash the following code if textview does not have pointer
	if !ns.np.textView.HasGrab() {
		ns.np.textView.GrabFocus()
	}
	buf := ns.np.buff
	// fmt.Println("Init curIter")
	ns.curIter = buf.GetIterAtMark(buf.GetInsert())
	ns.m1, ns.m2 = nil, nil
}

// NewNoteSearch - Create new  NotePad
func NewNoteSearch(np *NotePad) *NoteSearch {
	ns := &NoteSearch{np: np, isIcase: true}
	builder, err := gtk.BuilderNewFromFile("glade/note-search.glade")
	u.CheckErr(err, "BuilderNewFromFile note-search.glade")
	ns.builder = builder
	signals := map[string]interface{}{
		"NoteFindIcase":   ns.NoteFindIcase,
		"CommandFilter":   ns.CommandFilter,
		"NoteFindText":    ns.NoteFindText,
		"NoteReplaceText": ns.NoteReplaceText,
		"NoteReplaceAll":  ns.NoteReplaceAll,
		"OutputToNewNote": ns.OutputToNewNote,
		"KeyPressed":      ns.KeyPressed,
	}
	builder.ConnectSignals(signals)

	ns.w = GetWindow(builder, "note_search")

	ns.searchBox = GetSearchEntry(builder, "text_ptn")

	ns.replaceBox = GetEntry(builder, "replace_text")

	ns.ResetIter()

	ns.w.Connect("delete-event", func() bool {
		ns.np.noteSearch = nil
		return false
	})

	return ns
}
