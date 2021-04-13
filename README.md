## gnote

It is a golang port of my pretty old pnote using python gtk2.

As gtk2 and python2 is fading I got the option to port to pygtk3 or something else. Recently I just learn golang and like it so this is the result.

Many features is removed as I rarely use them however some new features is added.

I use it to quickly take notes, draft some code snippet and idea, search for error text in a log file (big that load into less or vim is very slow)

And a generic encryption/decyption text as well.

This uses gotk3 GUI and stole many code snippets from the examples there.

## Linux binary

- Requires gtk3 which is probably already install on standard ubuntu system.

Download the tar ball from the github [release page](https://github.com/sunshine69/gnote/releases)

Extract and run the `gnote-linux-amd64` executable.

## Windows binary build

Download the bundle in the git hub [release page](https://github.com/sunshine69/gnote/releases)
and extract them. Run the `gnote-windows-amd64.exe` or gnote.exe if exists inside the bin directory. You can create a
shortcut in the desktop manually to point to the binary path.

For upgrading you can only download the gnote.exe in the realease and copy it to the bin dir of the bundle folder.

Double click the exe to run the program.

## Build

Checkout git repo and execute:

If you change glade files recently then update the bindata before build

```
go-bindata -pkg forms -o forms/bindata.go -nomemcopy glade icons
```

- Linux

```
apt-get install libgtk-3-0  libgtk-3-dev ca-certificates
go build --tags "icu json1 fts5 secure_delete" -ldflags='-s -w'
```

My build is using a docker image to build and save cached. Basically at build host I pull and run image ubuntu:18.04 and exec in to install the above command. Download go and extract it to /usr/local. Then commit it into the image `golang-ubuntu-build`.  Now let the jenkins run it will use this image, pull go pkgs and build and save it as cache for the next build.


- Windows

```
go build -ldflags="-s -w -H=windowsgui" --tags "json1 fts5 secure_delete"  -o gnote-windows-amd64.exe gnote.go
```

There is a simple ansible playbook to build it on a windows build agent. To setup the windows box see [https://github.com/gotk3/gotk3/wiki/Installing-on-Windows](https://github.com/gotk3/gotk3/wiki/Installing-on-Windows) basically:

```
PS C:\> choco install golang
PS C:\> choco install git
PS C:\> choco install msys2
PS C:\> mingw64
$ pacman -S mingw-w64-x86_64-gtk3 mingw-w64-x86_64-toolchain base-devel glib2-devel
$ echo 'export PATH=/c/Go/bin:$PATH' >> ~/.bashrc
$ echo 'export PATH=/c/Program\ Files/Git/bin:$PATH' >> ~/.bashrc
$ source ~/.bashrc
$ sed -i -e 's/-Wl,-luuid/-luuid/g' /mingw64/lib/pkgconfig/gdk-3.0.pc # This fixes a bug in pkgconfig
$ go get github.com/gotk3/gotk3/gtk
```

## Text processing feature

gnote is a smaller source editor and runner. The note itself is using gtksourceview with all syntax highlighting
as you type. Just create a note, type some code snipper and click the button `source code highlight`. If it can
not detect the code it will prompt you to type the language in.

In case it detects wrongly, you can force it to prompt by select some text (that is can not be detected) and
click the highlight button, it will prompt.

All external text processing is done using the `Search & Replace` feature. Within the note click the search button or `Ctrl + f`, the window will popup.

The normal operation is search a raw text and replace with a text.

If you type in the `search text` a regex then it will search and replace using the regex (golang regex which is
also PCRE compatible.

Click the `cmd` button to turn it into run a external command on the sections/or the whole content. In here you
can run `sed` or `perl` etc. You do not need to quote - example like type `sed s/^/#/g` to add `#` at the
begining of the line.

It works by taking the note content or selection, write it to a temporary file with the extention detected by
the highlight button and run the command you type using the file path as the last argument.

The output will be replace to the note content, or the selection.

If you click the `new` button, it will output to a new note instead.

So using this you can run a code snippet, such as python or go (anything really). An example note below is to process text using
go (assume you have go installed)

Create a note with the content below

```
package main

import (
	"fmt"
	"regexp"
)

func main() {
	content := data()
	// Regex pattern captures "key: value" pair from the content.
	pattern := regexp.MustCompile(`(?m)(?P<key>\w+):\s+(?P<value>\w+)`)
	// Template to convert "key: value" to "key=value" by
	// referencing the values captured by the regex pattern.
	template := "$key,$value\n"

	result := []byte{}

	// For each match of the regex in the content.
	for _, submatches := range pattern.FindAllStringSubmatchIndex(content, -1) {
		// Apply the captured submatches to the template and append the output
		// to the result.
		result = pattern.ExpandString(result, template, content, submatches)
	}
	fmt.Println(string(result))
}


func data() string {
return `
this is a pattern kety1: value1 then parse and make it as a map k2: thisIs value2

# comment line
	this option1: value1
	option2: value2
	# another comment line
	option3: value3

`
}
```

Save and click the hightlight button, it should detect the go language type.

Press `Ctrl + f` , in the window pop up click `cmd` and `new` radio button.

Type in the search text `go run`. It will run the note content as a go prog.

See the output in the new note.

Now you have the idea how to utilize this feature.

