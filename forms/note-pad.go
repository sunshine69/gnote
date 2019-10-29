package forms

import (
	"strconv"
	"strings"
	"github.com/gotk3/gotk3/gdk"
	"time"
	"github.com/araddon/dateparse"
	"fmt"
	"github.com/gotk3/gotk3/gtk"
)

//NotePad - GUI related
type NotePad struct {
	w *gtk.Window
	builder *gtk.Builder
	textView *gtk.TextView
	buff *gtk.TextBuffer
	Note
}

//Load - Load note data and set the widget with data
func (np *NotePad) Load(id int) {
	if id < 0 { return }
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
		w.SetText(fmt.Sprintf("%v", np.Datelog))

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

		_w, e = b.GetObject("content")
		if e != nil {
			fmt.Printf("ERROR Can not load widget\n")
			return
		}
		np.textView = _w.(*gtk.TextView)
		np.buff, e = np.textView.GetBuffer()
		if e != nil {
			fmt.Printf("ERROR Can not load widget\n")
			return
		}
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
		np.textView.GrabFocus()
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
	}
	builder.ConnectSignals(signals)
	_widget, e := builder.GetObject("content")
	if e != nil {
		fmt.Printf("ERROR get content\n")
	}
	vWidget := _widget.(*gtk.TextView)
	vWidget.SetWrapMode(gtk.WRAP_WORD)

	np.Load(id)
	_o, _ := np.builder.GetObject("bt_close")
	b := _o.(*gtk.Button)
	b.SetLabel("Close")

	wSize, _ := GetConfig("window_size")
	_size := strings.Split(wSize, "x")
	w, _ := strconv.Atoi(_size[0])
	h, _ := strconv.Atoi(_size[1])
	np.w.SetDefaultSize(w, h)

	np.w.ShowAll()
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
	NewNoteSearch(np)
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

//SaveNote - save current note
func (np *NotePad) SaveNote() {
	b := np.builder
	_widget, e := b.GetObject("title")
	if e != nil {
		fmt.Printf("ERROR get title entry\n")
	}
	widget := _widget.(*gtk.Entry)
	np.Title, e = widget.GetText()
	if e != nil {fmt.Printf("ERROR get title entry text\n")}

	_widget, e = b.GetObject("datelog")
	if e != nil {
		fmt.Printf("ERROR get datelog\n")
	}
	widget = _widget.(*gtk.Entry)
	_datelogStr, e := widget.GetText()
	if e != nil {
		fmt.Printf("ERROR get entry datelog\n")
	} else {
		np.Datelog, e = dateparse.ParseLocal(_datelogStr)
		if e != nil {
			fmt.Printf("ERROR can not parse datelog. Use Now\n")
			np.Datelog = time.Now()
		}
	}
	_widget, e = b.GetObject("flags")
	if e != nil {
		fmt.Printf("ERROR get flags\n")
	}
	widget = _widget.(*gtk.Entry)
	np.Flags, e = widget.GetText()
	if e != nil {
		fmt.Printf("ERROR get entry flags\n")
	}
	_widget, e = b.GetObject("url")
	if e != nil {
		fmt.Printf("ERROR get url\n")
	}
	widget = _widget.(*gtk.Entry)
	np.URL, e = widget.GetText()
	if e != nil {
		fmt.Printf("ERROR get entry url\n")
	}

	_widget, e = b.GetObject("content")
	if e != nil {
		fmt.Printf("ERROR get content\n")
	}
	vWidget := _widget.(*gtk.TextView)
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
		_l := len(np.Content)
		if _l > 64 {_l = 64}
		np.Title = np.Content[0:_l]
	}
	if e = DbConn.Save(&np.Note).Error; e != nil {
		fmt.Printf("ERROR can not save note - %v\n", e)
	} else {
		fmt.Printf("INFO Note saved\n")
		_o, _ := np.builder.GetObject("bt_close")
		b := _o.(*gtk.Button)
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
	_w, _ := np.builder.GetObject("content")
	w := _w.(*gtk.TextView)
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