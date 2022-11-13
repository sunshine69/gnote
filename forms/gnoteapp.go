package forms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	u "github.com/sunshine69/golang-tools/utils"
)

// GnoteApp - struct
type GnoteApp struct {
	Builder         *gtk.Builder
	MainWindow      *gtk.Window
	model           *gtk.ListStore
	treeView        *gtk.TreeView
	selectedID      *[]int
	curNoteWindowID map[int]*NotePad
	searchBox       *gtk.SearchEntry
}

// ShowMain - show main window to do something. Meant to be called from NotePad
func (app *GnoteApp) ShowMain() {
	app.MainWindow.Present()
}

// RowActivated - Process when a treeview list row activated. Pop up a note window with the id
func (app *GnoteApp) RowActivated(treeView *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
	// note_id = model.get_value(model.get_iter(path), 0)
	model, _ := treeView.GetModel()
	iter, _ := model.ToTreeModel().GetIter(path)
	_v, _ := model.ToTreeModel().GetValue(iter, 0)
	v, _ := _v.GoValue()
	nID := v.(int)
	var _np *NotePad
	var ok bool
	if _np, ok = app.curNoteWindowID[nID]; !ok {
		_np := NewNotePad(nID)
		_np.app = app
		app.curNoteWindowID[nID] = _np
	} else {
		_np.w.Present()
	}
}

// ResultListKeyPress - evt
func (app *GnoteApp) ResultListKeyPress(w *gtk.TreeView, ev *gdk.Event) {
	keyEvent := &gdk.EventKey{ev}
	// fmt.Printf("DEBUG KEY %v\n", keyEvent.KeyVal() )
	if keyEvent.KeyVal() == 65535 { //Delete key
		for _, id := range *app.selectedID {
			fmt.Printf("ID %v\n", id)
			sql := fmt.Sprintf("DELETE FROM notes WHERE ID = '%d';", id)
			if e := DbConn.Unscoped().Exec(sql).Error; e != nil {
				fmt.Printf("ERROR %v\n", e)
			}
		}
		app.doSearch()
	}
}

// TreeSelectionChanged - evt
func (app *GnoteApp) TreeSelectionChanged(s *gtk.TreeSelection) {
	// Returns glib.List of gtk.TreePath pointers
	ListStore := app.model

	rows := s.GetSelectedRows(ListStore)
	items := make([]int, 0, rows.Length())

	for l := rows; l != nil; l = l.Next() {
		path := l.Data().(*gtk.TreePath)
		iter, _ := ListStore.GetIter(path)
		value, _ := ListStore.GetValue(iter, 0)
		str, _ := value.GoValue()
		items = append(items, str.(int))
	}
	app.selectedID = &items
	// fmt.Println(app.selectedID)
}

// Export a selections of notes into json file
func (app *GnoteApp) DoExportNotes() {
	if len(*app.selectedID) == 0 {
		return
	}
	outNoteList := []Note{}
	for _, _id := range *app.selectedID {
		_n := Note{}
		DbConn.First(&_n, _id)
		if _n.ID != -1 {
			outNoteList = append(outNoteList, _n)
		}
	}
	outputJson := u.JsonDump(outNoteList, "  ")
	dlg, _ := gtk.FileChooserDialogNewWith2Buttons(
		"choose file", nil, gtk.FILE_CHOOSER_ACTION_SAVE,
		"Open", gtk.RESPONSE_OK, "Cancel", gtk.RESPONSE_CANCEL,
	)
	defer dlg.Destroy()
	dlg.SetDefaultResponse(gtk.RESPONSE_OK)
	filter, _ := gtk.FileFilterNew()
	filter.SetName("my-export.json")
	filter.AddPattern("*.json")
	dlg.SetFilter(filter)
	response := dlg.Run()
	if response == gtk.RESPONSE_OK {
		filename := dlg.GetFilename()
		os.WriteFile(filename, []byte(outputJson), 0644)
	}
}

// Import notes from previously exported json file
func (app *GnoteApp) DoImportNotes() {
	var inputByte []byte
	var err error
	dlg, _ := gtk.FileChooserDialogNewWith2Buttons(
		"choose file", nil, gtk.FILE_CHOOSER_ACTION_OPEN,
		"Open", gtk.RESPONSE_OK, "Cancel", gtk.RESPONSE_CANCEL,
	)
	defer dlg.Destroy()
	dlg.SetDefaultResponse(gtk.RESPONSE_OK)
	filter, _ := gtk.FileFilterNew()
	filter.SetName("my-export.json")
	filter.AddPattern("*.json")
	dlg.SetFilter(filter)
	response := dlg.Run()
	if response == gtk.RESPONSE_OK {
		filename := dlg.GetFilename()
		inputByte, err = os.ReadFile(filename)
		if u.CheckErrNonFatal(err, "ReadFile DoImportNotes") != nil {
			MessageBox("ERROR ReadFile DoImportNotes")
			return
		}
	}
	if response == gtk.RESPONSE_CANCEL {
		return
	}
	inputNotes := []Note{}
	err = json.Unmarshal(inputByte, &inputNotes)
	if u.CheckErrNonFatal(err, "Unmarshal DoImportNotes") != nil {
		MessageBox("ERROR Unmarshal DoImportNotes")
		return
	}
	for _, note := range inputNotes {
		_n := Note{}
		DbConn.First(&_n, Note{Title: note.Title})
		if _n.ID == 0 {
			fmt.Println("new note")
			note.ID = 0
			DbConn.Create(&note)
		} else {
			fmt.Printf("update note new: %s - %d\nold: %s - %d\n", note.Title, note.Timestamp, _n.Title, _n.Timestamp)
			if note.Timestamp >= _n.Timestamp { // Update as input is newer or equal (rare)
				_n = note
				DbConn.Save(&_n)
			} else {
				fmt.Println("Input is older than current. Do not update")
				fmt.Printf("Existing note ID %d, title: '%s', TS: %d\n", _n.ID, _n.Title, _n.Timestamp)
				fmt.Printf("Input note note title: '%s', TS: %d\n", note.Title, note.Timestamp)
			}
		}
	}
}

// Create a note from current clipboard
func (app *GnoteApp) DoCreateNoteFromClipboard() {
	clip, err := gtk.ClipboardGet(gdk.SELECTION_PRIMARY)
	if u.CheckErrNonFatal(err, "ClipboardGet SELECTION_CLIPBOARD") != nil {
		MessageBox("ERROR Get clipboard " + err.Error())
		return
	}
	content, err := clip.WaitForText()
	if u.CheckErrNonFatal(err, "ClipboardGet WaitForText") != nil {
		MessageBox("ERROR Get clipboard WaitForText" + err.Error())
		return
	}
	note := Note{}
	note.NewNote(map[string]interface{}{
		"content": content,
	})
	DbConn.Save(&note)
}

// Change Passphrase
func (app *GnoteApp) DoChangePassphrase() {
	currentPass := InputDialog("title", "Enter Current Passphrase", "label", "Enter CURRENT passphrase to decode key.", "password-mask", '*')
	newPass := InputDialog("title", "Enter New Passphrase", "label", "Enter NEW passphrase to decode key.", "password-mask", '*')
	fullDbPath := os.Getenv("DBPATH")
	dbPath := strings.Split(fullDbPath, "?")[0]
	keyFile := dbPath + ".key"
	u.CheckErrNonFatal(ChangePassphrase(currentPass, newPass, keyFile), "Call ChangePassphrase")
}

// Set config Key -> Value. We no longer be able to open the db directly using sqlite3 command if teh file is encrypted thus we need this
func (app *GnoteApp) DoSetConfig() {
	key := InputDialog("title", "Set config KEY", "label", "Enter the config Key\nKey could be: select_limit, list_flags, webnote_url... See the config.go file for complete list")
	if key == "" {
		return
	}
	val := InputDialog("title", "Set config VALUE", "label", "Enter the config Value\nSee the config.go file for details")
	if val == "" {
		return
	}
	if err := SetConfig(key, val); err == nil {
		MessageBox("Config key: " + key + " Value: " + val + " have been set")
	} else {
		MessageBox("Some errors. " + err.Error())
	}
}

// NewNoteFromFile -
func (app *GnoteApp) NewNoteFromFile(o *gtk.FileChooserButton) {
	np := NewNoteFromFile(o.GetFilename())
	app.curNoteWindowID[np.ID] = np
	np.app = app
}

func (app *GnoteApp) DoVacuum() {
	if e := DbConn.Exec("VACUUM").Error; e != nil {
		MessageBox(fmt.Sprintf("ERROR VACUUM %v", e))
	}
}

func (app *GnoteApp) DoResetDB() {
	dbPath := os.Getenv("DBPATH")
	dbNewPath := fmt.Sprintf("%v.backup", dbPath)
	msg := `WARNING
This will rename your current database to the file %s.
It will exit the application to allow you to start it again to initialise the DB.
Are you sure to do that? Type 'yes'. otherwise type 'no' or hit enter.
	`
	confirm := InputDialog("title", "Confirmation required", "label", fmt.Sprintf(msg, dbNewPath))
	if confirm == "yes" {
		e := os.Rename(dbPath, dbNewPath)
		if e != nil {
			MessageBox(fmt.Sprintf("Error renaming db file. You may need to do it manualy. The file path is '%s'", dbPath))
		} else {
			MessageBox("Completed. You can click OK to shuttdown the app")
			os.Exit(0)
		}
	}
}

// DoUpdateResource -
func (app *GnoteApp) DoUpdateResource() {
	RestoreAssetsAll("./")
	MessageBox("Resource is updated. You need to restart the program to take effect")
}

// InitApp -
func (app *GnoteApp) InitApp() {
	Builder := app.Builder

	signals := map[string]interface{}{
		"ShowAbout": app.showAbout,
		// "OpenPref":             app.openPref,
		"NewNote": app.newNote,
		// "OpenDbfile":           app.openDBFile,
		"DoExit":               app.doExit,
		"DoClearSearchbox":     app.doClearSearchbox,
		"DoSearch":             app.doFullTextSearch,
		"RowActivated":         app.RowActivated,
		"ResultListKeyPress":   app.ResultListKeyPress,
		"TreeSelectionChanged": app.TreeSelectionChanged,
		"NewNoteFromFile":      app.NewNoteFromFile,
		"DoResetDB":            app.DoResetDB,
		"DoVacuum":             app.DoVacuum,
		"DoUpdateResource":     app.DoUpdateResource,
		"ChangePassphrase":     app.DoChangePassphrase,
		"SetConfig":            app.DoSetConfig,
		"DoExportNotes":        app.DoExportNotes,
		"DoImportNotes":        app.DoImportNotes,
		"NewNoteFromClipboard": app.DoCreateNoteFromClipboard,
		"DoSyncNotes":          app.DoSyncNotesFromWebnote,
	}

	Builder.ConnectSignals(signals)

	window := GetWindow(Builder, "window")
	window.Connect("delete-event", app.doExit)

	window.SetTitle("gnote")
	window.Connect("destroy", app.doExit)

	statusBar := GetStatusBar(Builder, "status_bar")
	statusBar.Push(1, "Welcome to gnote")

	app.MainWindow = window

	wT := GetTreeView(Builder, "treeview")
	app.treeView = wT
	app.model, _ = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
	wT.SetModel(app.model)
	wT.SetHeadersVisible(true)
	renderer, _ := gtk.CellRendererTextNew()
	col1, _ := gtk.TreeViewColumnNewWithAttribute("Title", renderer, "text", 1)
	col2, _ := gtk.TreeViewColumnNewWithAttribute("Date Log", renderer, "text", 2)
	col3, _ := gtk.TreeViewColumnNewWithAttribute("Last update", renderer, "text", 3)
	col1.SetResizable(true)
	col2.SetResizable(true)
	col3.SetResizable(true)
	col1.SetMaxWidth(200)
	col1.SetMinWidth(20)
	col2.SetMinWidth(20)
	wT.AppendColumn(col1)
	wT.AppendColumn(col2)
	wT.AppendColumn(col3)
	selection, _ := wT.GetSelection()
	selection.SetMode(gtk.SELECTION_MULTIPLE)
	// wT.SetSearchColumn(0)
	// window.SetPosition(gtk.WIN_POS_CENTER)
	// window.SetGravity(gdk.GDK_GRAVITY_NORTH_EAST)

	app.curNoteWindowID = make(map[int]*NotePad)
	app.searchBox = GetSearchEntry(Builder, "searchBox")

	wSize, _ := GetConfig("main_window_size", "300x291")
	_size := strings.Split(wSize, "x")
	w, _ := strconv.Atoi(_size[0])
	h, _ := strconv.Atoi(_size[1])
	window.SetDefaultSize(w, h)

	// window.Move(3000, 0)
	window.ShowAll()
}

// looks like handlers can literally be any function or method
// func (app *GnoteApp) openPref() {
// 	fmt.Println("Open Pref")
// 	Builder, err := gtk.BuilderNewFromFile("glade/gnote-editpref.glade")
// 	if err != nil {
// 		panic(err)
// 	}
// 	GetWindow(Builder, "edit_pref").Show()
// }

// func (app *GnoteApp) openDBFile() {
// 	fmt.Println("Open DB File")
// }

func (app *GnoteApp) DoSyncNotesFromWebnote() {
	duration := InputDialog("title", "Duration", "label", "Enter the time duration you want to sync: eg. 48h will sync all notes created in the last 48 hours", "default", "48h")
	client, csrfToken, webnoteURL := LoginToWebnote()
	if client == nil {
		return
	}
	resp, err := client.Get(fmt.Sprintf("%s/get_notes_titles?duration=%s&gorilla.csrf.Token=%s", webnoteURL, duration, csrfToken))
	if u.CheckErrNonFatal(err, "DoSyncNotesFromWebnote get_notes_titles") != nil {
		return
	}
	respText, _ := ioutil.ReadAll(resp.Body)
	// Even struct Note and webnote Note has small diff but we only care about common fields thus using Note here does work
	webnotes := []Note{}
	err = json.Unmarshal(respText, &webnotes)
	if u.CheckErrNonFatal(err, "Unmarshal") != nil {
		return
	}
	needToGetContentIds := []string{}
	for _, webnote := range webnotes {
		mynote := Note{}
		DbConnNew := DbConn.Where(&Note{Title: webnote.Title}).First(&mynote)
		if DbConnNew.RecordNotFound() {
			fmt.Printf("rec not found for title '%s', going to get it\n", webnote.Title)
			needToGetContentIds = append(needToGetContentIds, fmt.Sprintf("%d", webnote.ID))
		} else {
			fmt.Printf("TS mynote %d - webnote %d -  %s - %s\n", mynote.Timestamp, webnote.Timestamp, mynote.Title, webnote.Title)
			if mynote.Timestamp < webnote.Timestamp {
				fmt.Printf("we are older, so get webnote id %d\n", webnote.ID)
				needToGetContentIds = append(needToGetContentIds, fmt.Sprintf("%d", webnote.ID))
			}
		}
	}
	if len(needToGetContentIds) == 0 {
		fmt.Println("Nothing new to sync")
		return
	}
	data := fmt.Sprintf("(%s)", strings.Join(needToGetContentIds, ","))
	fmt.Printf("Will ask to get these IDs %s\n", data)
	getURL := fmt.Sprintf("%s/get_notes_by_id?ids=%s&gorilla.csrf.Token=%s", webnoteURL, url.QueryEscape(data), csrfToken)
	resp, err = client.Get(getURL)
	if u.CheckErrNonFatal(err, "DoSyncNotesFromWebnote get_notes_by_id") != nil {
		return
	}
	respText, _ = ioutil.ReadAll(resp.Body)

	webnotes1 := []Note{}
	err = json.Unmarshal(respText, &webnotes1)
	if u.CheckErrNonFatal(err, "Unmarshal") != nil {
		return
	}
	for _, note := range webnotes1 {
		newnote := Note{}
		noteExp := Note{ // Not have ID or title (unique field). We use this as expression. GOORM sucks!
			Datelog:   note.Datelog,
			Flags:     note.Flags,
			Content:   note.Content,
			URL:       note.URL,
			Timestamp: note.Timestamp,
			TimeSpent: note.TimeSpent,
			ReminderTicks: note.ReminderTicks,
		}
		DbconNew := DbConn.Where(&Note{Title: note.Title}).First(&newnote)
		if DbconNew.RecordNotFound() {
			newnote = noteExp
		} else {
			myID := newnote.ID
			newnote = noteExp
			newnote.ID = myID
		}
		newnote.Title = note.Title
		// Why DbconNew here does not do the Save? FFS goORM
		DbConn.Save(&newnote)
	}
	fmt.Printf("Have synced %d notes\n", len(webnotes1))
}

func (app *GnoteApp) newNote() *NotePad {
	np := NewNotePad(-1)
	np.app = app
	return np
}

func (app *GnoteApp) doExit() {
	w, h := app.MainWindow.GetSize()
	windowSize := fmt.Sprintf("%dx%d", w, h)
	fmt.Printf("save side - %dx%d\n", w, h)
	if e := SetConfig("main_window_size", windowSize); e != nil {
		fmt.Printf("ERROR save side - %v\n", e)
	}
	for _, np := range app.curNoteWindowID {
		np.saveBtnClick()
	}
	gtk.MainQuit()
}

func (app *GnoteApp) showAbout() {
	MessageBox("Gnote - A note and text processing system by stevek")
}

func (app *GnoteApp) doClearSearchbox() {
	fmt.Println("doClearSearchbox")
}

func (app *GnoteApp) doFullTextSearch() {
	fmt.Println("doFullTextSearch")
	b := app.Builder
	w := GetSearchEntry(b, "searchBox")
	keyword, _ := w.GetText()
	var sql string
	fmt.Printf("keyword: '%s'\n", keyword)
	if keyword == "" || strings.HasPrefix(keyword, "f:") || strings.HasPrefix(keyword, "flags:") || strings.HasPrefix(keyword, "F:") || strings.HasPrefix(keyword, "FLAGS:") || strings.HasPrefix(keyword, " ") {
		app.doSearch()
		return
	} else {
		sql = fmt.Sprintf("SELECT rowid FROM note_fts WHERE note_fts MATCH '%s' ORDER BY datelog DESC LIMIT 200;", keyword)
	}

	rows, e := DbConn.Raw(sql).Rows()
	if e != nil {
		fmt.Printf("ERROR - exec sql\n")
	}
	defer rows.Close()
	app.model.Clear()

	rowid, id, count := 0, 0, 0
	var title string
	var datelog, lastUpdate int64
	for rows.Next() {
		rows.Scan(&rowid)
		_note := Note{}
		if e = DbConn.First(&_note, rowid).Error; e != nil {
			fmt.Printf("Failt to get note id %d\n", rowid)
			break
		}
		id, title, datelog, lastUpdate = _note.ID, _note.Title, _note.Datelog, _note.Timestamp
		// fmt.Printf("row: %v - %v %v\n", id, title, datelog)
		_dateLogStr := u.NsToTime(datelog).Format(DateLayout)
		_lastUpdateStr := u.NsToTime(lastUpdate).Format(DateLayout)
		iter := app.model.Append()
		if e := app.model.Set(iter,
			[]int{0, 1, 2, 3},
			[]interface{}{id, title, _dateLogStr, _lastUpdateStr}); e != nil {
			fmt.Printf("ERROR appending data to model %v\n", e)
		}
		count = count + 1
	}
	s := GetStatusBar(app.Builder, "status_bar")
	s.Pop(1)
	s.Push(1, fmt.Sprintf("Found %d notes", count))
}

func (app *GnoteApp) doSearch() {
	fmt.Println("doSearch")
	b := app.Builder
	w := GetSearchEntry(b, "searchBox")
	searchFlags := false
	keyword, _ := w.GetText()
	q := ""
	keyword = strings.TrimSpace(keyword)
	tokens := []string{}
	if strings.HasPrefix(keyword, "F:") || strings.HasPrefix(keyword, "f:") {
		tokens = strings.Split(keyword[2:], ":")
		searchFlags = true
	} else if strings.HasPrefix(keyword, "FLAGS:") || strings.HasPrefix(keyword, "flags:") {
		tokens = strings.Split(keyword[6:], ":")
		searchFlags = true
	}
	if searchFlags {
		_l := len(tokens)
		for i, t := range tokens {
			if i == _l-1 {
				q = fmt.Sprintf("%v (flags LIKE '%%%v%%') ORDER BY datelog DESC LIMIT 200;", q, t)
			} else {
				q = fmt.Sprintf("%v (flags LIKE '%%%v%%') AND ", q, t)
			}
		}
		q = fmt.Sprintf("SELECT id, title, datelog, timestamp from notes WHERE %v", q)
	} else {
		tokens := strings.Split(keyword, " + ")
		_l := len(tokens)
		if _l == 0 {
			tokens = strings.Split(keyword, " & ")
			_l = len(tokens)
		}
		for i, t := range tokens {
			t = strings.TrimSpace(t)
			if i == _l-1 {
				q = fmt.Sprintf("%v (title LIKE '%%%v%%' OR content LIKE '%%%v%%') ORDER BY datelog DESC LIMIT 200;", q, t, t)
			} else {
				q = fmt.Sprintf("%v (title LIKE '%%%v%%' OR content LIKE '%%%v%%') AND ", q, t, t)
			}
		}
		q = fmt.Sprintf("SELECT id, title, datelog, timestamp from notes WHERE %v", q)
	}
	fmt.Println(q)
	rows, e := DbConn.Raw(q).Rows()
	if e != nil {
		fmt.Printf("ERROR - exec sql\n")
	}
	defer rows.Close()

	app.model.Clear()

	id, count := 0, 0
	var title string
	var datelog, lastUpdate int64
	for rows.Next() {
		rows.Scan(&id, &title, &datelog, &lastUpdate)
		// fmt.Printf("row: %v - %v %v\n", id, title, datelog)
		_dateLogStr := u.NsToTime(datelog).Format(DateLayout)
		_lastUpdateStr := u.NsToTime(lastUpdate).Format(DateLayout)
		iter := app.model.Append()
		if e := app.model.Set(iter,
			[]int{0, 1, 2, 3},
			[]interface{}{id, title, _dateLogStr, _lastUpdateStr}); e != nil {
			fmt.Printf("ERROR appending data to model %v\n", e)
		}
		count = count + 1
	}
	s := GetStatusBar(app.Builder, "status_bar")
	s.Pop(1)
	s.Push(1, fmt.Sprintf("Found %d notes", count))
}
