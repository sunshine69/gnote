package forms

import (
	"io/ioutil"
	"strconv"
	"strings"
	"github.com/gotk3/gotk3/gdk"
	"time"
	"fmt"
	"github.com/gotk3/gotk3/gtk"
)

//NotePad - GUI related
type NotePad struct {
	app *GnoteApp
	w *gtk.Window
	builder *gtk.Builder
	textView *gtk.TextView
	buff *gtk.TextBuffer
	wTitle *gtk.Entry
	wFlags *gtk.Entry
	wDateLog *gtk.Entry
	wURL *gtk.Entry
	Note
}

//ShowMainWindowBtnClick -
func (np *NotePad) ShowMainWindowBtnClick(o *gtk.Button) {
	np.app.ShowMain()
}

//Load - Load note data and set the widget with data
func (np *NotePad) Load(id int) {
	if id < 0 {
		np.Datelog = time.Now()
		np.wDateLog.SetText(np.Datelog.Format("02-01-2006 15:04:05 MST"))
		return
	}

	if e := DbConn.FirstOrInit(&np.Note, Note{ID: id}).Error; e != nil {
		fmt.Printf("INFO Can not find that note ID %d\n", id)
		return
	} else {
		b := np.builder
		_w, e := b.GetObject("title")
		if e != nil {
			fmt.Printf("ERROR Can not load widget\n")
			return
		}
		w := _w.(*gtk.Entry)
		w.SetText(np.Title)

		_w, e = b.GetObject("datelog")
		if e != nil {
			fmt.Printf("ERROR Can not load widget\n")
			return
		}
		w = _w.(*gtk.Entry)
		w.SetText(np.Datelog.Format("02-01-2006 15:04:05 MST"))

		_w, e = b.GetObject("flags")
		if e != nil {
			fmt.Printf("ERROR Can not load widget\n")
			return
		}
		w = _w.(*gtk.Entry)
		w.SetText(np.Flags)

		_w, e = b.GetObject("url")
		if e != nil {
			fmt.Printf("ERROR Can not load widget\n")
			return
		}
		w = _w.(*gtk.Entry)
		w.SetText(np.URL)

		np.buff.SetText(np.Content)
		np.textView.SetEditable(!(np.Readonly == 1))
		np.buff.Connect("changed", np.TextChanged)

		_w, e = b.GetObject("bt_toggle_rw")
		if e != nil {
			fmt.Printf("ERROR Can not load widget\n")
			return
		}
		wTB := _w.(*gtk.ToggleButton)
		wTB.SetActive((np.Readonly == 1))
	}

}

//NewNotePad - Create new  NotePad
func NewNotePad(id int) *NotePad {
	np := &NotePad{}
	builder, err := gtk.BuilderNewFromFile("glade/note.glade")
	np.builder = builder
	if err != nil {
		panic(err)
	}
	_obj, err := builder.GetObject("notepad")
	if err != nil {
		panic(err)
	}
	np.w = _obj.(*gtk.Window)
	np.NewNote(map[string]interface{}{})
	fmt.Printf("Empty note created %v\n", np.Title)

	signals := map[string]interface{} {
		"SaveBtnClick": np.saveBtnClick,
		"CloseBtnClick": np.closeBtnClick,
		"ToggleReadOnly": np.ToggleReadOnly,
		"TextChanged": np.TextChanged,
		"KeyPressed": np.KeyPressed,
		"ShowMainWindowBtnClick": np.ShowMainWindowBtnClick,
	}
	builder.ConnectSignals(signals)
	_widget, e := builder.GetObject("content")
	if e != nil {
		fmt.Printf("ERROR get content\n")
	}
	vWidget := _widget.(*gtk.TextView)
	vWidget.SetWrapMode(gtk.WRAP_WORD)
	np.textView = vWidget
	np.buff, _ = vWidget.GetBuffer()

	_w, e := builder.GetObject("title")
	if e != nil {
		fmt.Printf("ERROR get title\n")
	}
	np.wTitle = _w.(*gtk.Entry)
	_w, e = builder.GetObject("flags")
	if e != nil {
		fmt.Printf("ERROR get flags\n")
	}
	np.wFlags = _w.(*gtk.Entry)
	_w, e = builder.GetObject("url")
	if e != nil {
		fmt.Printf("ERROR get url\n")
	}
	np.wURL = _w.(*gtk.Entry)

	_w, e = builder.GetObject("datelog")
	if e != nil {
		fmt.Printf("ERROR get datelog\n")
	}
	np.wDateLog = _w.(*gtk.Entry)

	np.Load(id)
	_o, _ := np.builder.GetObject("bt_close")
	b := _o.(*gtk.Button)
	b.SetLabel("Close")

	wSize, _ := GetConfig("window_size")
	_size := strings.Split(wSize, "x")
	w, _ := strconv.Atoi(_size[0])
	h, _ := strconv.Atoi(_size[1])
	np.w.SetDefaultSize(w, h)

	if ! np.textView.HasGrab() { np.textView.GrabFocus() }
	np.w.ShowAll()
	return np
}

//NewNoteFromFile -
func NewNoteFromFile(filename string) *NotePad {
	ct, e := ioutil.ReadFile(filename)
	if e !=nil {
		MessageBox("Can not open file for reading")
		return nil
	}
	np := NewNotePad(-1)
	np.buff.SetText(string(ct))
	np.wTitle.SetText(filename)
	return np
}

//SaveWindowSize -
func (np *NotePad) SaveWindowSize() {
	w,h := np.w.GetSize()
	windowSize := fmt.Sprintf("%dx%d", w, h)
	fmt.Printf("save side - %dx%d\n", w, h)
	if e := SetConfig("window_size", windowSize); e != nil {
		fmt.Printf("ERROR save side - %v\n", e)
	}
}

//NoteSearch - Search text in the note
func (np *NotePad) NoteSearch() {
	ns := NewNoteSearch(np)
	ns.w.Show()
}

//KeyPressed - handle key board
func (np *NotePad) KeyPressed(o interface{}, ev *gdk.Event) {
	keyEvent := &gdk.EventKey{ev}
	// if keyEvent.KeyVal() == 65535 {//Delete key	}

	if keyEvent.State() & gdk.GDK_CONTROL_MASK > 0 { //Control key pressed
		switch keyEvent.KeyVal() {
		case gdk.KeyvalFromName("s"):
			np.SaveNote()
		case gdk.KeyvalFromName("f"):
			np.NoteSearch()
		}

	}
}

//TextChanged - Marked as changed
func (np *NotePad) TextChanged() {
	_o, _ := np.builder.GetObject("bt_close")
	b := _o.(*gtk.Button)
	b.SetLabel("Cancel")
}

//FetchDataFromGUI - populate the Note data from GUI widget. Prepare to save to db or anything else
func (np *NotePad) FetchDataFromGUI() {
	b := np.builder
	var e error
	widget := GetEntry(b, "title")
	np.Title, e = widget.GetText()
	if e != nil {fmt.Printf("ERROR get title entry text\n")}

	widget = GetEntry(b, "datelog")
	_datelogStr, e := widget.GetText()
	if e != nil {
		fmt.Printf("ERROR get entry datelog\n")
	} else {
		np.Datelog, e = time.Parse("02-01-2006 15:04:05 MST",_datelogStr)
		if e != nil {
			fmt.Printf("ERROR can not parse datelog. Use Now\n")
			np.Datelog = time.Now()
		}
	}

	widget = GetEntry(b, "flags")
	np.Flags, e = widget.GetText()
	if e != nil {
		fmt.Printf("ERROR get entry flags\n")
	}

	widget = GetEntry(b, "url")
	np.URL, e = widget.GetText()
	if e != nil {
		fmt.Printf("ERROR get entry url\n")
	}

	vWidget := GetTextView(b, "content")
	textBuffer, e := vWidget.GetBuffer()
	if e != nil {
		fmt.Printf("ERROR get text buffer content\n")
	} else {
		startIter := textBuffer.GetStartIter()
        endIter := textBuffer.GetEndIter()
		np.Content, e = textBuffer.GetText(startIter, endIter, true)
		if e != nil {
			fmt.Printf("ERROR can get content\n")
		}
	}

	np.Timestamp = time.Now()
	if np.Title == "" {
		np.Title = ChunkString(np.Content, 64)[0]
	}
}

//SaveToWebnote - save to webnote store
func (np *NotePad) SaveToWebnote() {
	np.FetchDataFromGUI()

}

//SaveNote - save current note
func (np *NotePad) SaveNote() {
	np.FetchDataFromGUI()
	if e := DbConn.Save(&np.Note).Error; e != nil {
		fmt.Printf("ERROR can not save note - %v\n", e)
	} else {
		fmt.Printf("INFO Note saved\n")
		b := GetButton(np.builder, "bt_close")
		b.SetLabel("Close")
	}
}

func (np *NotePad) saveBtnClick() {
	np.SaveNote()
	np.SaveWindowSize()
	np.w.Destroy()
}

func (np *NotePad) closeBtnClick() {
	np.w.Destroy()
}

//ToggleReadOnly - set content readonly mode
func (np *NotePad) ToggleReadOnly(bt *gtk.ToggleButton) {
	state := bt.GetActive()
	if state {
		np.Readonly = 1
	} else {
		np.Readonly = 0
	}
	w := GetTextView(np.builder, "content")
	w.SetEditable(! (np.Readonly == 1))
}

//GetSelection - Get the current selection and return start_iter, end_iter, text
//To be used in various places
func (np *NotePad) GetSelection() (string, *gtk.TextIter, *gtk.TextIter) {
	buff, _ := np.textView.GetBuffer()
	if st, en, ok := buff.GetSelectionBounds(); ok {
		if selectedText, e := buff.GetText(st, en, true); e == nil {
			return selectedText, st, en
		} else {
			fmt.Printf("ERROR gettext %v\n", e)
			return "", st, en
		}
	}
	return "", nil, nil
}