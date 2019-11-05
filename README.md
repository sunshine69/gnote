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
[this](https://xvt-public-repo.s3-ap-southeast-2.amazonaws.com/pub/devops/gnote-bundle-window-amd64.7z)
and extract them. Run the `gnote.exe` inside the bin directory. You can create a
shortcut in the desktop manually to point to the binary path.

For later update you just need to download the [gnote.exe
alone](https://xvt-public-repo.s3-ap-southeast-2.amazonaws.com/pub/devops/gnote-windows-amd64.exe)
and save to the same location.
