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
apt-get install libgtk-3-0  libgtk-3-dev ca-certificates libgtksourceview-3.0-dev
go build --tags "icu json1 fts5 secure_delete" -ldflags='-s -w'

# To avoid glibc incompatibility maybe try the below. Note it might cause link errors, if so use the previous build command.
go build --tags "icu json1 fts5 secure_delete osusergo netgo sqlite_stat4 sqlite_foreign_keys" -ldflags='-s -w'

# On RHEL8 is you get this error `could not determine kind of name for C.pango_attr_insert_hyphens_new` then add these tags into the build command `pango_1_42 gtk_3_22`

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

- MacOS

I used these commands to build on MacOS Monterey for amd64

```
brew install pkg-config gtk+3 adwaita-icon-theme
brew install gtksourceview3
go build --tags "json1 fts5 secure_delete osusergo netgo sqlite_stat4 sqlite_foreign_keys" -ldflags='-s -w'
```

To install it into the application folder so you can double click it in Finder and run it rather than run from terminal

```
# at terminal
mkdir ~/Applications/gnote.app
# Extract the binary downloaded from the release page
tar xf gnote-macos-amd64.tgz -C ~/Applications/gnote.app
# If you build it locally then just copy the binary to the fodler gnote.app above
```

Now you can use Finder and navigate to your Home/Applications you should see the app icon (actually it is the default icon). Double click the first time, it would extract these assets files like glade and icons. Then it will run the first setup.

I do not know how to make it appear in the normal Applications - I need to navigte to my home first and get inside the Applications folder. Not sure why.

ISSUES

The first time you run it it may give you the error like this

```
gnote[1532:17457] *** Terminating app due to uncaught exception 'NSInternalInconsistencyException', reason: 'NSWindow drag regions should only be invalidated on the Main Thread!'
```

I think as I display the dialog without opening a main window. On MAC they prohibit it. The setup process has been completed at that time thus the next time you run it will run fine.

I will have some time to think how to fix that problem.

## Text processing feature

gnote is a small source editor and runner. The note itself is using gtksourceview with all syntax highlighting
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

If you type `gopher-lua` then it will trigger to use the internal lua5.1 VM inpterpreter and will run the note (or selection content) as lua code. See [gopher-lua](https://github.com/yuin/gopher-lua) for more. I include some libraries as well, gluare (allow golang regexp syntax to use), gluahttp (http client), gopher-json, gluayaml (handle json and yaml), gluacrypto (some simple crypto func). Please refer to this site for examples of usage.

This allows us to use text processing features without any external command. For more about Lua programming language see [this](https://www.lua.org).

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

## New example - use python3 modules as jsonformatter or yaml formatter

Create a note with proper title and content below

```
from datetime import datetime, timedelta
import ipaddress
import time
import csv, re, sys, os, subprocess
from glob import glob
import yaml, json, sqlite3, uuid

def str_presenter(dumper, data):
    """configures yaml for dumping multiline strings
    Ref: https://stackoverflow.com/questions/8640959/how-can-i-control-what-scalar-form-pyyaml-uses-for-my-data
    usage:
        yaml.add_representer(str, str_presenter)
        yaml.representer.SafeRepresenter.add_representer(str, str_presenter)
    """
    if data.count('\n') > 0:  # check for multiline string
        return dumper.represent_scalar('tag:yaml.org,2002:str', data, style='|')
    return dumper.represent_scalar('tag:yaml.org,2002:str', data)

yaml.add_representer(str, str_presenter)
yaml.representer.SafeRepresenter.add_representer(str, str_presenter)

def json_dump(obj: dict, indent: int = 4):
    o = json.dumps(obj, indent=indent)
    print(o)
    return o

def json_dumps(jsontext: str, indent: int = 4):
    obj = json.loads(jsontext)
    return json_dump(obj, indent)

def yaml_load(ymltext: str) ->dict:
    return  yaml.load(ymltext, Loader=yaml.CLoader)

def yaml_dumps(ymltext: str) ->str:
    o = yaml_load(ymltext)
    os = yaml.dump(o)
    print(os)
    return os

def run_cmd(cmd,sendtxt=None, working_dir=".", args=[], shell=True, DEBUG=False, shlex=False):
    if DEBUG:
        cmd2 = re.sub('root:([^\s])', 'root:xxxxx', cmd) # suppress the root password printout
        print(cmd2)
    if sys.platform == "win32":
        args = cmd
    else:
        if shlex:
            import shlex
            args = shlex.split(cmd)
        else:
            args = cmd
    popen = subprocess.Popen(
            args,
            shell=shell,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            stdin=subprocess.PIPE,
            cwd=working_dir
        )
    if sendtxt: output, err = popen.communicate(bytearray(sendtxt, 'utf-8'))
    else: output, err = popen.communicate()
    code = popen.returncode
    if not code == 0 or DEBUG:
        output = "Command string: '%s'\n\n%s" % (cmd, output)
    return (output.strip().decode(encoding='utf-8'), code, err)

datatxt = '''
{"what": 233, "ever": 3222}
'''

json_dumps(datatxt)
```

Click the highlight button to detect src type, if it pops up select python.

You know, the `datatxt` is part where you paste your unformatted json or yaml. Then you can use the python func
`json_dumps()` or `yaml_dumps()` to print out the output. There are several functions we defined above.

After pasting text, and select which function to run (in here I use json_dumps); press Ctrl + f or click the
search button in the right note tool bar. Select cmd, new type the command `python3` and click `cmd` button. See
the output in the new note.

I import several modules just if we need it but not all of them is used in the sample. Get wild!

## New example - Use built in Lua VM

Create a note with content below

```
local re = require("re")
local http = require("http")

content = get_note("CreateDataNoteListOfLanguageSupport")
print(content)
```

First/Second line to allow you to use the module `re` using golang regexp syntax, but the command is the same as lua regex (string.find, gsub etc)

Next, you can get the data as string from a existing note. The func `get_note` will search a note with the title `CreateDataNoteListOfLanguageSupport` and get the content. Note that there are notes which is automatically created at startup if it does not exists to serve as a data point for some internal usage. At the moment there are 2 and you should not remove it as if you do, it will be re-created again at the next start, and in your session some feature might stop working.

Now you can use anything lua allows you with the content.

You can fetch url to get the content as well using http. See https://github.com/cjoudrey/gluahttp.

```
response, error_message = http.request("GET", "http://example.com", {
    query="page=1",
    timeout="30s",
    headers={
        Accept="*/*"
    }
})
print(response.body)
-- response has fields:  body (string), body_size (number), headers (table), cookies (table), status_code (number), url (string)
```

Function in lua to deal with notes have been implemented:
- get_note - give a string as note title - get one note
- search_notes - give an WHERE sql part, return all notes match that.
- update_notes - The full sql UPDATE statement

For these get functions, return a json string dump of the result. You can then use json.decode them into lua table and process.

More example ...

```
json = require("json")
-- List of usefull fields: title, datelog, content, url, flags, language, fileext, timestamp, readonly, reminderticks

o = search_notes("WHERE flags = ':TAO' ")

-- dump output to json file
file = io.open("notes-dump.json", "w")
file:write(o)
file:close()

-- print(o)

o1 = json.decode(o)

for _, v in ipairs(o1) do
  print(v.Title)
end
```

Lua split lines ..

```
--Returns a table splitting some string with a delimiter
--Changes to enhance the code from https://gist.github.com/jaredallard/ddb152179831dd23b230
function string:split(delimiter)
    local result = {}
    local from = 1
    local delim_from, delim_to = string.find(self, delimiter, from, true)
    while delim_from do
        if (delim_from ~= 1) then
            table.insert(result, string.sub(self, from, delim_from-1))
        end
        from = delim_to + 1
        delim_from, delim_to = string.find(self, delimiter, from, true)
    end
    if (from <= #self) then table.insert(result, string.sub(self, from)) end
    return result
end
-- data
mdata = [[qwe
asd


rzxc
]]
-- execution
lines = mdata:split("\n")
for _, s in ipairs(lines) do
  if not ((s == '') or (s == nil)) then
    print(s)
  end
end
```

Regex examples ...

```
local re = require("re")

-- quote
assert(re.quote("^$.?a") == [[\^\$\.\?a]])

local data = "Today is 21/10/2022"
local ptn = [[([\d]+/[\d]+/[\d]+)]]
i, j, s = re.find(data, ptn)
print(i, j, s)

a, b, s = re.find("abcd efgh i23", "i([0-9]+)")
print(a,b,s)

-- re wont work with lua regex ptn, only accept go regexp. Need to use raw string for the pattern
-- This will use normal lua string regex
print( string.find(data, "%d%d/%d%d/%d%d%d%d") )

-- find will return index start, end, capture (start from 1). sub extract index
-- capture return is optional
-- .sub wont use capture
s = "Deadline is 30/05/1999, firm"
date = "%d%d/%d%d/%d%d%d%d"
print(string.sub(s, string.find(s, date)))   --> 30/05/1999

date = [[([\d]+/[\d]+)/[\d]+]]
-- for this to work need to capture in the regex above. remember re.sub does not exists, u need to use string.sub
print(string.sub(s, re.find(s, date))) --> 30/05

```
