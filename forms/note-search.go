package forms

import (
	// "fmt"
	"github.com/gotk3/gotk3/gtk"
)

//NoteSearch - GUI related
type NoteSearch struct {
	w *gtk.Window
	builder *gtk.Builder
	np *NotePad
	isIcase bool
	isRegexp bool
	isBackward bool
	searchBox *gtk.SearchEntry
	replaceBox *gtk.Entry
	m1 *gtk.TextMark
	m2 *gtk.TextMark
	curIter *gtk.TextIter
}

func (ns *NoteSearch) NoteFindIcase(o *gtk.CheckButton) {	ns.isIcase = o.GetActive() }

func (ns *NoteSearch) NoteFindRegexp(o *gtk.CheckButton) {	ns.isRegexp = o.GetActive()}

func (ns *NoteSearch) NoteFindBackward(o *gtk.CheckButton) {	ns.isBackward = o.GetActive()}

func (ns *NoteSearch) FindText() bool {
	buf := ns.np.buff
	keyword, _ := ns.searchBox.GetText()
	searchFlag := gtk.TEXT_SEARCH_TEXT_ONLY
	var foundIter1, foundIter2 *gtk.TextIter
	var ok bool = true
	var output = false

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
		"NoteFindRegexp": ns.NoteFindRegexp,
		"NoteFindText": ns.NoteFindText,
		"NoteReplaceText": ns.NoteReplaceText,
		"NoteReplaceAll": ns.NoteReplaceAll,
		"NoteFindBackward": ns.NoteFindBackward,
	}
	builder.ConnectSignals(signals)

	ns.w = GetWindow(builder, "note_search")

	ns.searchBox = GetSearchEntry(builder, "text_ptn")

	ns.replaceBox = GetEntry(builder, "replace_text")
	if ! np.textView.HasGrab() { np.textView.GrabFocus() } //Crash the following code if textview does not have pointer
	if ns.np.buff.GetHasSelection(){
		text, _, _ := ns.np.GetSelection()
		if text != "" {
			ns.searchBox.SetText(text)
		}
	}

	buf := ns.np.buff
	ns.curIter =  buf.GetIterAtMark(buf.GetInsert())
	ns.m1 , ns.m2 = nil, nil

	return ns
}