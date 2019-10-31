package forms

import (
	"log"
	"strconv"
	"time"
	"github.com/gotk3/gotk3/gtk"
)

//MessageBox - display a message
func MessageBox(msg string) {
	d := gtk.MessageDialogNew(nil, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE, msg)
	d.Run()
	d.Destroy()
}

func ticks2Time(ticks string) time.Time {
	i, err := strconv.ParseInt(ticks, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	return tm
}

//ChunkString -
func ChunkString(s string, chunkSize int) []string {
	var chunks []string
	runes := []rune(s)

	if len(runes) == 0 {
		return []string{s}
	}

	for i := 0; i < len(runes); i += chunkSize {
		nn := i + chunkSize
		if nn > len(runes) {
			nn = len(runes)
		}
		chunks = append(chunks, string(runes[i:nn]))
	}
	return chunks
}

func GetWindow(b *gtk.Builder, id string) (Window *gtk.Window) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Get Window error: %s", e)
		return nil
	}

	Window, _ = obj.(*gtk.Window)
	return
}

func GetDialog(b *gtk.Builder, id string) (Window *gtk.Dialog) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Get Dialog error: %s", e)
		return nil
	}

	Window, _ = obj.(*gtk.Dialog)
	return
}

func GetListStore(b *gtk.Builder, id string) (listStore *gtk.ListStore) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("List store error: %s", e)
		return nil
	}

	listStore, _ = obj.(*gtk.ListStore)
	return
}

func GetButton(b *gtk.Builder, id string) (btn *gtk.Button) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Button error: %s", e)
		return nil
	}

	btn, _ = obj.(*gtk.Button)
	return
}

func GetToggleButton(b *gtk.Builder, id string) (btn *gtk.ToggleButton) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Toggle button error: %s", e)
		return nil
	}

	btn, _ = obj.(*gtk.ToggleButton)
	return
}

func GetTreeView(b *gtk.Builder, id string) (treeView *gtk.TreeView) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Tree view error: %s", e)
		return nil
	}

	treeView, _ = obj.(*gtk.TreeView)
	return
}

func GetTextView(b *gtk.Builder, id string) (treeView *gtk.TextView) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Tree view error: %s", e)
		return nil
	}

	treeView, _ = obj.(*gtk.TextView)
	return
}

func GetLabel(b *gtk.Builder, id string) (treeView *gtk.Label) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Label error: %s", e)
		return nil
	}

	treeView, _ = obj.(*gtk.Label)
	return
}

func GetScrolledWindow(b *gtk.Builder, id string) (treeView *gtk.ScrolledWindow) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Scrolled window error: %s", e)
		return nil
	}

	treeView, _ = obj.(*gtk.ScrolledWindow)
	return
}

func GetSpinner(b *gtk.Builder, id string) (treeView *gtk.Spinner) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Spinner error: %s", e)
		return nil
	}

	treeView, _ = obj.(*gtk.Spinner)
	return
}

func GetEntry(b *gtk.Builder, id string) (treeView *gtk.Entry) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Entry error: %s", e)
		return nil
	}

	treeView, _ = obj.(*gtk.Entry)
	return
}

func GetSearchEntry(b *gtk.Builder, id string) (treeView *gtk.SearchEntry) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("SearchEntry error: %s", e)
		return nil
	}

	treeView, _ = obj.(*gtk.SearchEntry)
	return
}

func GetStatusBar(b *gtk.Builder, id string) (treeView *gtk.Statusbar) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("SearchEntry error: %s", e)
		return nil
	}

	treeView, _ = obj.(*gtk.Statusbar)
	return
}

func GetComboBox(b *gtk.Builder, id string) (combobox *gtk.ComboBox) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("ComboBox error: %s", e)
		return nil
	}

	combobox, _ = obj.(*gtk.ComboBox)
	return
}

func GetCheckButton(b *gtk.Builder, id string) (el *gtk.CheckButton) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("CheckButton error: %s", e)
		return nil
	}

	el, _ = obj.(*gtk.CheckButton)
	return
}

func GetImage(b *gtk.Builder, id string) (el *gtk.Image) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Image error: %s", e)
		return nil
	}

	el, _ = obj.(*gtk.Image)
	return
}

func GetMenuItem(b *gtk.Builder, id string) (el *gtk.MenuItem) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("MenuItem error: %s", e)
		return nil
	}

	el, _ = obj.(*gtk.MenuItem)
	return
}

func GetCheckMenuItem(b *gtk.Builder, id string) (el *gtk.CheckMenuItem) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("CheckMenuItem error: %s", e)
		return nil
	}

	el, _ = obj.(*gtk.CheckMenuItem)
	return
}
