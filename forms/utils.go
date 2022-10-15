package forms

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/quick"
	sourceview "github.com/linuxerwang/sourceview3"

	"path/filepath"

	"github.com/gotk3/gotk3/gtk"
	cp "github.com/otiai10/copy"
	u "github.com/sunshine69/golang-tools/utils"
)

// MessageBox - display a message
func MessageBox(msg string) {
	d := gtk.MessageDialogNew(nil, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE, msg)
	d.Run()
	d.Destroy()
}

// func ticks2Time(ticks string) time.Time {
// 	i, err := strconv.ParseInt(ticks, 10, 64)
// 	if err != nil {
// 		panic(err)
// 	}
// 	tm := time.Unix(i, 0)
// 	return tm
// }

// RestoreAssetsAll -
func RestoreAssetsAll(extractDir string) {
	for _, as := range AssetNames() {
		fmt.Printf("Restore %s\n", as)
		RestoreAssets(extractDir, as)
	}
}

//GUI helpers

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

func GetSourceView(b *gtk.Builder, id string) *sourceview.SourceView {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Tree view error: %s", e)
		return nil
	}

	view, _ := obj.(*sourceview.SourceView)
	return view
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

func GetButton(b *gtk.Builder, id string) (btn *gtk.Button) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Button error: %s", e)
		return nil
	}

	btn, _ = obj.(*gtk.Button)
	return
}

// func GetFileChooserButton(b *gtk.Builder, id string) (btn *gtk.FileChooserButton) {
// 	obj, e := b.GetObject(id)
// 	if e != nil {
// 		log.Printf("Button error: %s", e)
// 		return nil
// 	}

// 	btn, _ = obj.(*gtk.FileChooserButton)
// 	return
// }

// func GetComboBox(b *gtk.Builder, id string) (combobox *gtk.ComboBox) {
// 	obj, e := b.GetObject(id)
// 	if e != nil {
// 		log.Printf("ComboBox error: %s", e)
// 		return nil
// 	}

// 	combobox, _ = obj.(*gtk.ComboBox)
// 	return
// }

// func GetCheckButton(b *gtk.Builder, id string) (el *gtk.CheckButton) {
// 	obj, e := b.GetObject(id)
// 	if e != nil {
// 		log.Printf("CheckButton error: %s", e)
// 		return nil
// 	}

// 	el, _ = obj.(*gtk.CheckButton)
// 	return
// }

// func GetImage(b *gtk.Builder, id string) (el *gtk.Image) {
// 	obj, e := b.GetObject(id)
// 	if e != nil {
// 		log.Printf("Image error: %s", e)
// 		return nil
// 	}

// 	el, _ = obj.(*gtk.Image)
// 	return
// }

// func GetToggleButton(b *gtk.Builder, id string) (btn *gtk.ToggleButton) {
// 	obj, e := b.GetObject(id)
// 	if e != nil {
// 		log.Printf("Toggle button error: %s", e)
// 		return nil
// 	}

// 	btn, _ = obj.(*gtk.ToggleButton)
// 	return
// }

// func GetScrolledWindow(b *gtk.Builder, id string) (treeView *gtk.ScrolledWindow) {
// 	obj, e := b.GetObject(id)
// 	if e != nil {
// 		log.Printf("Scrolled window error: %s", e)
// 		return nil
// 	}

// 	treeView, _ = obj.(*gtk.ScrolledWindow)
// 	return
// }

// func GetSpinner(b *gtk.Builder, id string) (treeView *gtk.Spinner) {
// 	obj, e := b.GetObject(id)
// 	if e != nil {
// 		log.Printf("Spinner error: %s", e)
// 		return nil
// 	}

// 	treeView, _ = obj.(*gtk.Spinner)
// 	return
// }

// ChromaHighlight - Syntax highlighter using Chroma syntax
// highlighter: "github.com/alecthomas/chroma"
// informations above
func ChromaHighlight(inputString, lexer string) (out string, err error) {
	var buff bytes.Buffer
	writer := bufio.NewWriter(&buff)

	// Registrering pango formatter
	formatters.Register("pango", chroma.FormatterFunc(pangoFormatter))

	// Doing the job
	if err = quick.Highlight(writer, pangoPrepare(inputString), lexer, "pango", "github"); err != nil {
		return
	}
	writer.Flush()
	return pangoFinalize(string(buff.Bytes())), err
}

// pangoFormatter: is a part of "ChromaHighlight" function
func pangoFormatter(w io.Writer, style *chroma.Style, it chroma.Iterator) error {

	// Clear the background colour.
	var clearBackground = func(style *chroma.Style) *chroma.Style {
		builder := style.Builder()
		bg := builder.Get(chroma.Background)
		bg.Background = 0
		bg.NoInherit = true
		builder.AddEntry(chroma.Background, bg)
		style, _ = builder.Build()
		return style
	}

	closer, out := "", ""
	style = clearBackground(style)
	for token := it(); token != chroma.EOF; token = it() {
		entry := style.Get(token.Type)
		if !entry.IsZero() {
			closer, out = "", ""
			if entry.Bold == chroma.Yes {
				out += "<b>"
				closer = "</b>" + closer
			}
			if entry.Underline == chroma.Yes {
				out += "<u>"
				closer = "</u>" + closer
			}
			if entry.Italic == chroma.Yes {
				out += "<i>"
				closer = "</i>" + closer
			}
			if entry.Colour.IsSet() {
				out += fmt.Sprintf("<span foreground=\"#%02X%02X%02X\">", entry.Colour.Red(), entry.Colour.Green(), entry.Colour.Blue())
				closer = "</span>" + closer
			}
			if entry.Background.IsSet() {
				out += fmt.Sprintf("<span background=\"#%02X%02X%02X\">", entry.Background.Red(), entry.Background.Green(), entry.Background.Blue())
				closer = "</span>" + closer
			}
			fmt.Fprint(w, out)
		}
		fmt.Fprint(w, token.Value)
		if !entry.IsZero() {
			fmt.Fprint(w, closer)
		}
	}
	return nil
}

var pangoEscapeChar = [][]string{{"<", "&lt;", "lOwErThAnTmPrEpLaCeMeNt"}, {"&", "&amp;", "aMpErSaNdTmPrEpLaCeMeNt"}}

// prepare: sanitize input string to safely use with pango
func pangoPrepare(inString string) string {
	inString = strings.Replace(inString, pangoEscapeChar[1][0], pangoEscapeChar[1][2], -1)
	return strings.Replace(inString, pangoEscapeChar[0][0], pangoEscapeChar[0][2], -1)
}

// finalize: restore originals characters using markup replacement
func pangoFinalize(inString string) string {
	inString = strings.Replace(inString, pangoEscapeChar[1][2], pangoEscapeChar[1][1], -1)
	return strings.Replace(inString, pangoEscapeChar[0][2], pangoEscapeChar[0][1], -1)
}

// Python and msys shell is like s***t. File not found while file exists and etc etc.. FFS lets write it in golang
func CreateWinBundle() {
	srcRootDir := "c:/ansible_install"
	srcDir := srcRootDir + "/gnote"
	targetDir := srcRootDir + "/gnote-windows-bundle"
	mingw64Prefix := "c:/tools/msys64/mingw64"

	os.RemoveAll(targetDir)
	for _, _f := range []string{"/bin", "/lib", "/share"} {
		os.MkdirAll(targetDir+_f, 0755)
	}

	err := cp.Copy(mingw64Prefix+"/lib/gdk-pixbuf-2.0", targetDir+"/lib/gdk-pixbuf-2.0")
	fmt.Println(err)
	err = cp.Copy(mingw64Prefix+"/share/glib-2.0", targetDir+"/share/glib-2.0")
	fmt.Println(err)
	err = cp.Copy(mingw64Prefix+"/share/gtksourceview-3.0", targetDir+"/share/gtksourceview-3.0")
	fmt.Println(err)
	err = cp.Copy(mingw64Prefix+"/share/icons", targetDir+"/share/icons")
	fmt.Println(err)

	exeFiles, err := filepath.Glob(srcDir + "/gnote*.exe")
	u.CheckErr(err, "Glob")
	for _, _f := range exeFiles {
		cp.Copy(_f, targetDir+"/bin/"+filepath.Base(_f))
	}

	dllFilesByte, err := os.ReadFile(srcDir + "/dll_files.lst")
	u.CheckErr(err, "dll_files")
	dllFilesStr := string(dllFilesByte)
	dllFilesStr = strings.ReplaceAll(dllFilesStr, "\r\n", "\n")
	lines := strings.Split(dllFilesStr, "\n")
	for _, _f := range lines {
		if _f != "" {
			fmt.Printf("Copy %s/%s => %s/%s\n", mingw64Prefix, _f, targetDir+"/bin", _f)
			err = cp.Copy(mingw64Prefix+"/bin/"+_f, targetDir+"/bin/"+_f)
			fmt.Println(err)
		}
	}
	fmt.Println("Output folder " + targetDir)
}

func ChangePassphrase(old, new, keyFile string) error {
	if old == "" || new == "" || keyFile == "" {
		return fmt.Errorf("[ERROR] oldpass, newpass or keyfile is empty string")
	}
	keyEncData, err := os.ReadFile(keyFile)
	if u.CheckErrNonFatal(err, "keyEncData") != nil {
		return nil
	}
	key, err := u.Decrypt(string(keyEncData), old)
	if u.CheckErrNonFatal(err, "Decrypt keyEncData") == nil {
		keyEnc := u.Encrypt(key, new)
		err = os.WriteFile(keyFile, []byte(keyEnc), 0600)
		return u.CheckErrNonFatal(err, "WriteFile")
	} else {
		fmt.Println(string(keyEncData))
	}
	return nil
}
