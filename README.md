## gnote

It is a golang port of my pretty old pnote using python gtk2.

As gtk2 and python2 is fading I got the option to port to pygtk3 or something else. Recently I just learn golang and like it so this is the result.

Many features is removed as I rarely use them however some new features is added.

I use it to quickly take notes, draft some code snippet and idea, search for error text in a log file (big that load into less or vim is very slow)

And a generic encryption/decyption text as well.

This uses gotk3 GUI and stole many code snippets from the examples there.

## Linux binary

- Requires gtk3 which is probably already install on standard ubuntu system.

[linux-amd64](https://xvt-public-repo.s3-ap-southeast-2.amazonaws.com/pub/devops/gnote-linux-amd64.tar.xz)

Extract and run the `gnote-linux-amd64` executable.

## Windows binary build

I did test build on windows and here is the bundle. Download
[this](https://xvt-public-repo.s3-ap-southeast-2.amazonaws.com/pub/devops/gnote-bundle-windows-amd64.7z)
and extract them. Run the `gnote-windows-amd64.exe` inside the bin directory. You can create a
shortcut in the desktop manually to point to the binary path.

For later update you just need to download the [gnote-windows-amd64.exe alone](https://xvt-public-repo.s3-ap-southeast-2.amazonaws.com/pub/devops/gnote-windows-amd64.7z)
and extract, save to the same location.

## Build

Checkout git repo and execute:

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


