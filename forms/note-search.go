package forms

import (
	"fmt"
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

func (ns *NoteSearch) NoteFindText(o *gtk.Button) {
	buf := ns.np.buff
	keyword, _, _ := ns.np.GetSelection()
	if keyword != "" {
		ns.searchBox.SetText(keyword)
	}
	keyword, _ = ns.searchBox.GetText()
	searchFlag := gtk.TEXT_SEARCH_TEXT_ONLY
	var foundIter1, foundIter2 *gtk.TextIter
	var ok bool = true
	// if ns.isIcase {
	// 	searchFlag = gtk.TEXT_SEARCH_CASE_INSENSITIVE
	// }
	if ns.isBackward {
		if ns.m1 != nil {
			fmt.Printf("keyword '%s'\n", keyword)
			buf.PlaceCursor(buf.GetIterAtMark(ns.m1))
			ns.curIter = buf.GetIterAtMark(buf.GetInsert())
		}
		foundIter1, foundIter2, ok = ns.curIter.BackwardSearch(keyword, searchFlag, buf.GetStartIter())
	} else {
		if ns.m2 != nil {
			fmt.Printf("keyword '%s'\n", keyword)
	  		buf.PlaceCursor(buf.GetIterAtMark(ns.m2))
			  ns.curIter = buf.GetIterAtMark(buf.GetInsert())
		}
		foundIter1, foundIter2, ok = ns.curIter.ForwardSearch(keyword, searchFlag, buf.GetEndIter())
	}
	if ok {
		ns.np.textView.ScrollToIter(foundIter2, 0, true, 0, 0)
		buf.SelectRange(foundIter1, foundIter2)
		ns.m1 , ns.m2 = buf.CreateMark ("", foundIter1, true), buf.CreateMark("", foundIter2, true)
	} else {
		if !ok {
			d := gtk.MessageDialogNew(ns.w, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE, "Text not found")
			d.Run()
			d.Destroy()
		}
	}
}

func (ns *NoteSearch) NoteReplaceText(o *gtk.Button) {
	fmt.Println("Do replace")
}

func (ns *NoteSearch) NoteReplaceAll(o *gtk.Button) {
	fmt.Println("Do replace all")
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

	_obj, err := builder.GetObject("note_search")
	if err != nil {
		panic(err)
	}
	ns.w = _obj.(*gtk.Window)

	_obj, err = builder.GetObject("text_ptn")
	if err != nil {
		panic(err)
	}
	ns.searchBox = _obj.(*gtk.SearchEntry)

	_obj, err = builder.GetObject("replace_text")
	if err != nil {
		panic(err)
	}
	ns.replaceBox = _obj.(*gtk.Entry)
	text, _, _ := ns.np.GetSelection()
	if text != "" {
		ns.searchBox.SetText(text)
	}

	buf := ns.np.buff
	ns.curIter =  np.buff.GetIterAtMark(np.buff.GetInsert())
	ns.m1 , ns.m2 = buf.CreateMark("start", ns.curIter, true), buf.CreateMark("end", buf.GetEndIter(), true)

	ns.w.ShowAll()
	return ns
}