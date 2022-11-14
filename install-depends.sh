#!/bin/bash

# Build and create a tarball with name reflecting which current system to build

ARCH=$(uname -m)
OS=$(uname -s)
GO_TAG="icu json1 fts5 secure_delete"

if [ "$OS" = "Linux" ]; then
    sudo apt-get update
    sudo apt-get -y install libgtk-3-0  libgtk-3-dev ca-certificates libgtksourceview-3.0-dev
elif [ "$OS" = "Darwin" ]; then
    brew install pkg-config gtk+3 adwaita-icon-theme
    brew install gtksourceview3
elif [[ "$OS" =~ MINGW64 ]]; then
    echo "Not supported Mingw64 yet"
    exit 1
fi