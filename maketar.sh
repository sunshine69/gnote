#!/bin/bash

# Build and create a tarball with name reflecting which current system to build

ARCH=$(uname -m)
OS=$(uname -s)
GO_TAG="icu json1 fts5 secure_delete"

if [ "$OS" = "Linux" ]; then
    DISTRO_NAME=$(grep '^NAME=' /etc/os-release | sed 's/ //g;s/"//g' | cut -f2 -d=)
    DISTRO_VER=$(grep '^VERSION_ID=' /etc/os-release | sed 's/ //g;s/"//g' | cut -f2 -d=)
    TARBALL_NAME="gnote-${DISTRO_NAME}-${DISTRO_VER}-${ARCH}.tgz"
    REDHAT_SUPPORT_PRODUCT_VERSION=$(grep '^REDHAT_SUPPORT_PRODUCT_VERSION=' /etc/os-release | sed 's/ //g;s/"//g' | cut -f2 -d=)
    if [ "$REDHAT_SUPPORT_PRODUCT_VERSION" = "8" ]; then
        GO_TAG="${GO_TAG} pango_1_42 gtk_3_22"
    fi
elif [ "$OS" = "Darwin" ]; then
    ProductName=$( sw_vers | grep ProductName | sed 's/ //g; s/\t//g' | cut -f2 -d: )
    ProductVersion=$( sw_vers | grep ProductVersion | sed 's/ //g; s/\t//g' | cut -f2 -d: )
    TARBALL_NAME="gnote-${ProductName}-${ProductVersion}-${ARCH}.tgz"
elif [[ "$OS" =~ MINGW64 ]]; then
    go build -ldflags="-s -w -H=windowsgui" --tags "json1 fts5 secure_delete"  -o gnote-windows-amd64.exe gnote.go
    if [ "$1" == "" ]; then
        echo "Enter your mingw64 root dir: "
        read MINGW64_ROOT_DIR
    else
        MINGW64_ROOT_DIR=$1
    fi
    if [ "$MINGW64_ROOT_DIR" != "" ]; then
        MINGW64_ROOT_OPT="-mingw64-root ${MINGW64_ROOT_DIR}"
    fi
    ./gnote-windows-amd64.exe -create-win-bundle $MINGW64_ROOT_OPT
    pushd .
    cd ..

    zip -r gnote-windows-bundle.zip gnote-windows-bundle
    echo "Output bundle file: $(pwd)/gnote-windows-bundle.zip"
    rm -rf gnote-windows-bundle
    popd
    exit 0
fi

go build --tags "${GO_TAG}" -ldflags='-s -w' -o gnote

mkdir gnote.app
cp -a gnote gnote.app/
tar czf $TARBALL_NAME gnote.app

echo Tar ball pkg is $TARBALL_NAME

#docker run --rm -v $(pwd):/work --entrypoint /usr/local/go/bin/go --workdir /work golang-ubuntu1804-build:latest build --tags "json1 fts5 secure_delete" -ldflags='-s -w' -o gnote-ubuntu1804-${ARCH}

#mv gnote-${ARCH} ~/Public/

#go build --tags "json1 fts5 secure_delete" -ldflags='-s -w' -o ~/Public/gnote-linux-amd64

#cd ~/Public
#echo "Quit current gnote to allow update. Build windows version and manually sync to Public. Then hit enter"
#read _junk

#cp -a gnote-${ARCH} ~/gnote/gnote
#for f in gnote-bundle-windows-amd64.7z gnote-ubuntu1804-${ARCH}; do
#	aws s3 cp $f s3://xvt-public-repo/pub/devops/ --profile xvt_aws
#done
