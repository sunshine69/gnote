package forms

import (
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
	isBackward bool
	searchBox *gtk.SearchEntry
	replaceBox *gtk.Entry
	m1 *gtk.TextMark
	m2 *gtk.TextMark
	curIter *gtk.TextIter
}

func (ns *NoteSearch) NoteFindIcase(o *gtk.CheckButton) {	ns.isIcase = o.GetActive() }

func (ns *NoteSearch) CommandFilter(o *gtk.CheckButton) {
	ns.isCmdFilter = o.GetActive()
	lastCmd, _ := GetConfig("last_cmd_filter", "perl -pe 's///'")
	ns.searchBox.SetText(lastCmd)
}

func (ns *NoteSearch) NoteFindBackward(o *gtk.CheckButton) {	ns.isBackward = o.GetActive()}

func (ns *NoteSearch) FindText() bool {
	buf := ns.np.buff
	keyword, _ := ns.searchBox.GetText()
	searchFlag := gtk.TEXT_SEARCH_TEXT_ONLY
	var foundIter1, foundIter2 *gtk.TextIter
	var ok bool = true
	var output = false

	if ns.isCmdFilter {//run external command and replace the note/selection with output
		text, startI, endI := ns.np.GetSelection()
		fmt.Printf("SELECTION %v\n", text)
		_tmpF, _ := ioutil.TempFile("", "browser")
		_tmpF.Write([]byte(text))
		cmdText := fmt.Sprintf("%s %s", keyword ,_tmpF.Name())
		fmt.Printf("Command: %v\n", cmdText)
		cmd := exec.Command("sh", "-c", cmdText)
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("DEBUG E %v\n", err)
		} else{
			fmt.Printf("DEBUG 1 %s\n", stdoutStderr)
			buf := ns.np.buff
			buf.SelectRange(startI, endI)
			buf.DeleteSelection(true, true)
			buf.InsertAtCursor(string(stdoutStderr))
		}
		os.Remove(_tmpF.Name())
		//Not sure why the curIter is invalid after running. Need to get back otherwise crash
		ns.curIter =  buf.GetIterAtMark(buf.GetInsert())
		SetConfig("last_cmd_filter", cmdText)
		return false //stop other actions
	} else {
		if ns.isIcase {
			searchFlag = gtk.TEXT_SEARCH_CASE_INSENSITIVE
		}
		if ns.isBackward {
			if ns.m1 != nil {
				buf.PlaceCursor(buf.GetIterAtMark(ns.m1))
				ns.curIter = buf.GetIterAtMark(buf.GetInsert())
			}
			foundIter1, foundIter2, ok = ns.curIter.BackwardSearch(keyword, searchFlag, nil)
		} else {
			if ns.m2 != nil {
				buf.PlaceCursor(buf.GetIterAtMark(ns.m2))
				ns.curIter = buf.GetIterAtMark(buf.GetInsert())
			}
			foundIter1, foundIter2, ok = ns.curIter.ForwardSearch(keyword, searchFlag, nil)
		}
		if ok {
			ns.np.textView.ScrollToIter(foundIter1, 0, true, 0, 0)
			buf.SelectRange(foundIter1, foundIter2)
			ns.m1 , ns.m2 = buf.CreateMark("s1", foundIter1, false), buf.CreateMark("s2", foundIter2, false)
			output = true
		} else {
			if !ok {
				MessageBox("Search text not found. Will reset iter")
				if ns.isBackward {
					ns.curIter = buf.GetEndIter()
					ns.m1, ns.m2 = nil, nil
				} else {
					ns.curIter = buf.GetStartIter()
					ns.m1, ns.m2 = nil, nil
				}
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
		"NoteFindBackward": ns.NoteFindBackward,
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