package forms

import (
	"regexp"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
	"github.com/alecthomas/chroma/lexers"
	"github.com/gomarkdown/markdown"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/browser"
)

//NotePad - GUI related
type NotePad struct {
	app             *GnoteApp
	w               *gtk.Window
	builder         *gtk.Builder
	textView        *gtk.TextView
	buff            *gtk.TextBuffer
	wTitle          *gtk.Entry
	wFlags          *gtk.Entry
	wDateLog        *gtk.Entry
	wURL            *gtk.Entry
	tabCount        int
	StartUpdateTime time.Time
	Note
}

//ShowMainWindowBtnClick -
func (np *NotePad) ShowMainWindowBtnClick(o *gtk.Button) {
	np.app.ShowMain()
}

//Load - Load note data and set the widget with data
func (np *NotePad) Load(id int) {
	if id < 0 {//Datelog only constructed in here and never be updated for teh life of the note.
		np.Datelog = time.Now().UnixNano()
		np.wDateLog.SetText(nsToTime(np.Datelog).Format(DateLayout))
		np.StartUpdateTime = time.Now()
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
		w.SetText(nsToTime(np.Datelog).Format(DateLayout))

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
		np.StartUpdateTime = time.Now()
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
	// fmt.Printf("Empty note created %v\n", np.Title)

	signals := map[string]interface{}{
		"SaveBtnClick":             np.saveBtnClick,
		"CloseBtnClick":            np.closeBtnClick,
		"ToggleReadOnly":           np.ToggleReadOnly,
		"TextChanged":              np.TextChanged,
		"KeyPressed":               np.KeyPressed,
		"ShowMainWindowBtnClick":   np.ShowMainWindowBtnClick,
		"SendBtnClick":             np.SaveToWebnote,
		"HighlightBtnClick":        np.HighlightBtnClick,
		"AppendUpdateMarkBtnClick": np.AppendUpdateMarkBtnClick,
		"SearchNoteFromPad":        np.SearchNoteFromPad,
		"EndUpdateMarkBtnClick":    np.EndUpdateMarkBtnClick,
		"InsertFileToNote":         np.InsertFileToNote,
		"EncryptContent":           np.EncryptContent,
		"DecryptContent":           np.DecryptContent,
		"NoteSearchText": 			np.NoteSearchText,
		"ClearURL":					np.ClearURL,
		"ClearFlagsBtnClick":		np.ClearFlagsBtnClick,
		"NewLinkNote":				np.NewLinkNote,
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

	wSize, _ := GetConfig("window_size", "429x503")
	_size := strings.Split(wSize, "x")
	w, _ := strconv.Atoi(_size[0])
	h, _ := strconv.Atoi(_size[1])
	np.w.SetDefaultSize(w, h)

	if !np.textView.HasGrab() {
		np.textView.GrabFocus()
	}
	np.w.Connect("delete-event", func() bool {
		delete(np.app.curNoteWindowID, np.ID)
		return false
	})
	np.w.ShowAll()
	return np
}

func (np *NotePad) NewLinkNote() {
	newNp := np.app.newNote()
	if np.buff.GetHasSelection(){
		text, _, _ := np.GetSelection()
		newNp.buff.SetText(text)
	}
}

func (np *NotePad) ClearFlagsBtnClick() {
	wFlag := GetEntry(np.builder, "flags")
	wFlag.SetText("")
	wFlag.GrabFocus()
 }

func (np *NotePad) NoteSearchText() { np.NoteSearch() }

func (np *NotePad) ClearURL() {
	wURL := GetEntry(np.builder, "url")
	wURL.SetText("")
	wURL.GrabFocus()
 }

func (np *NotePad) DecryptContent() {
	key := InputDialog("title", "Password required", "label", "Enter passphrase to encrypt: ", "password-mask", '*')
	eCt, startI, endI := np.GetSelection()
	eCt = strings.TrimPrefix(eCt, "ENC:")
	eCt = strings.TrimSuffix(eCt, ":ENC")
	ct, e := Decrypt(eCt, key)
	if e != nil {
		MessageBox("Decrypt error. Please check password")
	} else{
		np.buff.SelectRange(startI, endI)
		np.buff.DeleteSelection(true, true)
		np.buff.InsertAtCursor(ct)
	}
}

func (np *NotePad) EncryptContent() {
	key := InputDialog("title", "Password required", "label", "Enter passphrase to encrypt: ", "password-mask", '*')
	ct, startI, endI := np.GetSelection()
	eCt := Encrypt(ct, key)
	eCt = fmt.Sprintf("ENC:%s:ENC", eCt)
	np.buff.SelectRange(startI, endI)
	np.buff.DeleteSelection(true, true)
	np.buff.InsertAtCursor(eCt)
}

func (np *NotePad) EndUpdateMarkBtnClick() {
	np.SaveNote()
	durationInsec := time.Now().Unix() - np.StartUpdateTime.Unix()
	np.TimeSpent = np.TimeSpent + int(durationInsec)
	dur, _ := time.ParseDuration(fmt.Sprintf("%ds", durationInsec))
	total, _ := time.ParseDuration(fmt.Sprintf("%ds", np.TimeSpent))
	text := fmt.Sprintf("\n---\nEnd Update %s. Time Spent: %s\nTotal time spent: %s\n", time.Now().Format(DateLayout), dur.String(), total.String())

	endI := np.buff.GetEndIter()
	np.buff.PlaceCursor(endI)
	np.buff.InsertAtCursor(text)
	np.textView.GrabFocus()
}

func (np *NotePad) InsertFileToNote(o *gtk.FileChooserButton) {
	ct, _ := ioutil.ReadFile(o.GetFilename())
	buf := np.buff
	buf.InsertAtCursor(string(ct))
}

func (np *NotePad) SearchNoteFromPad() {
	buf := np.buff
	if buf.GetHasSelection() {
		text, _, _ := np.GetSelection()
		if len(text) < 64 {
			np.app.searchBox.SetText(text)
			np.app.doFullTextSearch()
			np.app.MainWindow.Present()
		}
	}
}

func (np *NotePad) AppendUpdateMarkBtnClick() {
	text := fmt.Sprintf("---\nUpdate %s\n", time.Now().Format(DateLayout))
	endI := np.buff.GetEndIter()
	np.buff.PlaceCursor(endI)
	np.buff.InsertAtCursor(text)
	np.StartUpdateTime = time.Now()
	np.textView.GrabFocus()
}

//NewNoteFromFile -
func NewNoteFromFile(filename string) *NotePad {
	ct, e := ioutil.ReadFile(filename)
	if e != nil {
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
	w, h := np.w.GetSize()
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
func (np *NotePad) KeyPressed(o interface{}, ev *gdk.Event) bool {
	keyEvent := &gdk.EventKey{ev}
	// fmt.Printf("Key val: %v\n", keyEvent.KeyVal())
	if keyEvent.State()&gdk.GDK_CONTROL_MASK > 0 { //Control key pressed
		switch keyEvent.KeyVal() {
		case gdk.KeyvalFromName("T"): //All tab clear
			np.tabCount = 0
		case gdk.KeyvalFromName("t"): //reduce one tab level
			if np.tabCount > 0 {
				np.tabCount = np.tabCount - 1
			}
		case gdk.KeyvalFromName("s"):
			np.SaveNote()
		case gdk.KeyvalFromName("f"): //Find & Replace
			np.NoteSearch()
		case gdk.KeyvalFromName("b"): //Open in browser
			_t, _ := np.buff.GetText(np.buff.GetStartIter(), np.buff.GetEndIter(), true)
			md := []byte(_t)
			output := markdown.ToHTML(md, nil, nil)
			browser.OpenReader(strings.NewReader(string(output)))
		case gdk.KeyvalFromName("q"):
			np.w.Close()
		case gdk.KeyvalFromName("h"):
			helpTxt := `Keyboard shortcut of the notepad
Ctrl + s - Save note (not closing after save)
Ctrl + T - Clear all tabs count. When you press tab key it wil auto indent the level. Press this key to clear it
Ctrl + t - Reduce one tab level.
Ctrl + f - Show search and replace text. Finding text pattern and many useful features.
Ctrl + b - Show the content in a web browser. This will convert the markdown text into html if your note content is a markdown format text.
Ctrl + q - Close this note window.
			`
			MessageBox(helpTxt)
		}
	}
	switch keyEvent.KeyVal() {
	case 65293: // Enter key not sure what name is
		if np.tabCount > 0 {
			_str := ""
			for i := 1; i <= np.tabCount; i++ {
				_str = fmt.Sprintf("%s\t", _str)
			}
			_str = fmt.Sprintf("\n%s", _str)
			np.buff.InsertAtCursor(_str)
		} else {
			np.buff.InsertAtCursor("\n")
		}
		return true
	case gdk.KEY_Tab:
		np.tabCount = np.tabCount + 1
	}
	return false
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
	if e != nil {
		fmt.Printf("ERROR get title entry text\n")
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

	np.Timestamp = time.Now().UnixNano()
	if np.Title == "" {
		np.Title = ChunkString(np.Content, 64)[0]
	}
}

//SaveToWebnote - save to webnote store
func (np *NotePad) SaveToWebnote() {
	// SetupDefaultConfig()
	np.SaveNote()
	if WebNoteUser == "" {
		msg := `
		This feature allow user to save the note into a webnote.
		You need to have a webnote user account.
		Contact the author if you are interested.`
		MessageBox(msg)
		WebNoteUser = InputDialog("title", "Input required", "label", "Enter webnote username: ")
	}
	if WebNotePassword == "" {
		WebNotePassword = InputDialog(
			"title", "Password requried", "password-mask", '*', "label", "Enter webnote password. If you need OTP token, enter it at the end of the password separated with ':'")
	}
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	otpCode := ""
	otpPtn, _ := regexp.Compile(`\:([\d]+)$`)
	_otpCode := otpPtn.FindStringSubmatch(WebNotePassword)
	if len(_otpCode) > 0 {
		otpCode = _otpCode[1]
	}

	if WebNoteUser == "" || WebNotePassword == "" {
		MessageBox("No username or password. Aborting ...")
		return
	}
	data := url.Values{
		"action":      {"do_login"},
		"path":        {"action=show_frontpage"},
		"username":    {WebNoteUser},
		"password":    {WebNotePassword},
		"totp_number": {otpCode},
		"login":       {"Login"},
	}
	resp, err := client.PostForm("https://note.inxuanthuy.com", data)
	if err != nil {
		fmt.Printf("ERROR - CRITICAL login to webnote %v", err)
		WebNotePassword = ""
		WebNoteUser = ""
	}
	defer resp.Body.Close()
	respText, _ := ioutil.ReadAll(resp.Body)
	if strings.HasPrefix(string(respText), "ERROR - Response is None") {
		MessageBox("Error login to webnote. Check your password")
		WebNotePassword = ""
		WebNoteUser = ""
		return
	}

	data = url.Values{
		"action":     {"save_newnote"},
		"title":      {np.Title},
		"id":         {"0"},
		"datelog":    {nsToTime(np.Datelog).Format(DateLayout)} ,
		"timestamp":  {fmt.Sprintf("%d", np.Timestamp)},
		"flags":      {np.Flags},
		"content":    {np.Content},
		"url":        {np.URL},
		"ngroup":     {"default"},
		"permission": {"0"},
		"is_ajax":    {"1"},
		"raw_editor": {"1"},
		"savenote":   {"Save"},
	}
	resp, err = client.PostForm("https://note.inxuanthuy.com", data)
	if err != nil {
		fmt.Printf("ERROR - CRITICAL save to webnote %v", err)
		panic(2)
	}
	defer resp.Body.Close()
	respText, _ = ioutil.ReadAll(resp.Body)
	if string(respText) != "OK note saved" {
		SetConfig("webnote_user", WebNoteUser)
		browser.OpenReader(strings.NewReader(string(respText)))
	}
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
		np.app.curNoteWindowID[np.ID] = 1
	}
}

func (np *NotePad) saveBtnClick() {
	np.SaveNote()
	np.SaveWindowSize()
	np.w.Close()
}

func (np *NotePad) closeBtnClick() {
	np.w.Close()
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
	w.SetEditable(!(np.Readonly == 1))
}

//GetSelection - Get the current selection and return start_iter, end_iter, text
//To be used in various places
func (np *NotePad) GetSelection() (string, *gtk.TextIter, *gtk.TextIter) {
	buff := np.buff
	if buff.GetHasSelection() {
		// fmt.Printf("GetSelection\n")
		if st, en, ok := buff.GetSelectionBounds(); ok {
			if selectedText, e := buff.GetText(st, en, true); e == nil {
				return selectedText, st, en
			} else {
				fmt.Printf("ERROR gettext %v\n", e)
				return "", st, en
			}
		}
	} else {
		// fmt.Printf("Get whole note content\n")
		startI := buff.GetStartIter()
		endI := buff.GetEndIter()
		o, _ := buff.GetText(startI, endI, true)
		return o, startI, endI
	}
	return "", nil, nil
}

//HighlightBtnClick -
func (np *NotePad) HighlightBtnClick() {
	fmt.Printf("Start Highlight\n")
	buf := np.buff
	var someSourceCode string
	var startI, endI *gtk.TextIter
	if buf.GetHasSelection() {
		someSourceCode, startI, endI = np.GetSelection()
	} else {
		startI = buf.GetStartIter()
		endI = buf.GetEndIter()
		someSourceCode, _ = buf.GetText(startI, endI, true)
	}
	lexer := lexers.Analyse(someSourceCode)
	lexerStr := ""
	if lexer != nil {
		c := lexer.Config()
		fmt.Printf("Lexer detected type: %v\n", c.Name)
		lexerStr = c.Name
	} else {
		lexerStr = InputDialog("title", "Input required", "label", "Enter the language string for highlighter:", "default", "python")
	}
	formattedSource, err := ChromaHighlight(someSourceCode, lexerStr)
	if err == nil {
		buf.Delete(startI, endI)
		buf.InsertMarkup(startI, formattedSource)
	} else {
		fmt.Printf("%v\n", err)
	}
}
