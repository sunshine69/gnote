package forms

import (
	"github.com/gotk3/gotk3/gdk"
	"regexp"
	"os"
	"io/ioutil"
	"os/exec"
	"fmt"
	"github.com/gotk3/gotk3/gtk"
)

//NoteSearch - GUI related
type NoteSearch struct {
	w *gtk.Window
	builder *gtk.Builder
	np *NotePad
	isIcase bool
	isCmdFilter bool
	isGetTempFileContent bool
	searchBox *gtk.SearchEntry
	replaceBox *gtk.Entry
	m1 *gtk.TextMark
	m2 *gtk.TextMark
	curIter *gtk.TextIter
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
	} else {
		ns.searchBox.SetText("")
		btnFind := GetButton(ns.builder, "find_btn")
		btnFind.SetLabel("Find")
	}
	ns.searchBox.GrabFocus()
}

func (ns *NoteSearch) GetTempFileContent(o *gtk.CheckButton) {
	ns.isGetTempFileContent = o.GetActive()
	ns.searchBox.GrabFocus()
}

func (ns *NoteSearch) FindText() bool {
	buf := ns.np.buff
	keyword, _ := ns.searchBox.GetText()
	searchFlag := gtk.TEXT_SEARCH_TEXT_ONLY
	var foundIter1, foundIter2 *gtk.TextIter
	var ok bool = true
	var output = false

	if ns.isCmdFilter {//run external command and replace the note/selection with output
		text, startI, endI := ns.np.GetSelection()
		replaceWith, _ := ns.replaceBox.GetText()
		outStr := ""
		if replaceWith == "" {
			_tmpF, _ := ioutil.TempFile("", "browser")
			_tmpF.Write([]byte(text))
			cmdText := fmt.Sprintf("%s %s", keyword ,_tmpF.Name())
			cmd := exec.Command("sh", "-c", cmdText)
			stdoutStderr, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("DEBUG E %v\n", err)
			} else{
				fmt.Printf("DEBUG 1 %s\n", stdoutStderr)
			}
			if ns.isGetTempFileContent {
				_outStr, e := ioutil.ReadFile(_tmpF.Name())
				if e != nil {
					MessageBox(fmt.Sprintf("ERROR %v\n", e))
					return false
				}
				outStr = string(_outStr)
			} else{
				outStr = string(stdoutStderr)
			}
			os.Remove(_tmpF.Name())
			SetConfig("last_cmd_filter", cmdText)
		} else {
			if ptn, e := regexp.Compile(keyword); e == nil {
				outStr = ptn.ReplaceAllString(text, replaceWith)
			} else {
				fmt.Printf("ERROR %v\n", e)
			}
		}
		buf := ns.np.buff
		buf.SelectRange(startI, endI)
		buf.DeleteSelection(true, true)
		buf.InsertAtCursor(outStr)
		//Not sure why the curIter is invalid after running. Need to get back otherwise crash
		ns.curIter =  buf.GetIterAtMark(buf.GetInsert())
		return false //stop other actions
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
			ns.m1 , ns.m2 = buf.CreateMark("s1", foundIter1, false), buf.CreateMark("s2", foundIter2, false)
			output = true
		} else {
			if !ok {
				MessageBox("Search text not found. Will reset iter")
				ns.curIter = buf.GetStartIter()
				ns.m1, ns.m2 = nil, nil
			}
		}
	}
	return output
}
//NoteFindText -
func (ns *NoteSearch) NoteFindText() {
	ns.FindText()
}
//NoteReplaceText -
func (ns *NoteSearch) NoteReplaceText(o *gtk.Button) {
	buf := ns.np.buff

	if buf.GetHasSelection() || ns.FindText() {
		buf.DeleteSelection(true, true)
		_rp := GetEntry(ns.builder, "replace_text")
		replaceText, _ := _rp.GetText()
		buf.InsertAtCursor(replaceText)
	}

}
//NoteReplaceAll -
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
	if keyEvent.State()&gdk.GDK_CONTROL_MASK > 0 { //Control key pressed
		switch keyEvent.KeyVal() {
		case gdk.KeyvalFromName("q"):
			ns.w.Close()
		}
	}
}

//NewNoteSearch - Create new  NotePad
func NewNoteSearch(np *NotePad) *NoteSearch {
	ns := &NoteSearch{np: np}
	builder, err := gtk.BuilderNewFromFile("glade/note-search.glade")
	if err != nil {
		panic(err)
	}
	ns.builder = builder
	signals := map[string]interface{} {
		"NoteFindIcase": ns.NoteFindIcase,
		"CommandFilter": ns.CommandFilter,
		"NoteFindText": ns.NoteFindText,
		"NoteReplaceText": ns.NoteReplaceText,
		"NoteReplaceAll": ns.NoteReplaceAll,
		"GetTempFileContent": ns.GetTempFileContent,
		"KeyPressed":		ns.KeyPressed,
	}
	builder.ConnectSignals(signals)

	ns.w = GetWindow(builder, "note_search")

	ns.searchBox = GetSearchEntry(builder, "text_ptn")

	ns.replaceBox = GetEntry(builder, "replace_text")

	//Crash the following code if textview does not have pointer
	if ! np.textView.HasGrab() { np.textView.GrabFocus() }

	buf := np.buff
	fmt.Println("Init curIter")
	ns.curIter =  buf.GetIterAtMark(buf.GetInsert())
	ns.m1 , ns.m2 = nil, nil

	return ns
}