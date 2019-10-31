package forms

import (
	"github.com/gotk3/gotk3/gtk"

)

//InputDialog - Args: title, label prompt, password-mask,
func InputDialog(opt ...interface{}) string {
    b, _ := gtk.BuilderNewFromFile("glade/input-dialog.glade")
    d := GetDialog(b, "input_dialog")
    entry := GetEntry(b, "input_entry")

    for i, v := range(opt) {
        if i % 2 == 0 {
            key := v.(string)
            switch key {
            case "title":
                d.SetTitle(opt[i+1].(string))
            case "label":
                l := GetLabel(b,"input_label")
                l.SetText(opt[i+1].(string))
            case "password-mask":
                entry.SetInvisibleChar(opt[i+1].(rune))
                entry.SetVisibility(false)
            }
        }
    }

    output := ""
    entry.Connect("activate", func (o *gtk.Entry) { d.Response(gtk.RESPONSE_OK) } )
    btok := GetButton(b, "bt_ok")
    btok.Connect("clicked", func (b *gtk.Button) { d.Response(gtk.RESPONSE_OK) } )

    btcancel := GetButton(b, "bt_cancel")
    btcancel.Connect("clicked", func (b *gtk.Button) { d.Response(gtk.RESPONSE_CANCEL) } )

    code := d.Run()
    if code == gtk.RESPONSE_OK {
        output, _ = entry.GetText()
    }

    d.Destroy()
    return output
}
