#!/bin/sh

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
    ProductName=$( sw_vers | grep ProductName | sed 's/ //g' | cut -f2 -d: )
    ProductVersion=$( sw_vers | grep ProductVersion | sed 's/ //g' | cut -f2 -d: )
    TARBALL_NAME="gnote-${ProductName}-${ProductVersion}-${ARCH}.tgz"
fi

go build --tags "${GO_TAG}" -ldflags='-s -w' -o gnote

tar czf $TARBALL_NAME gnote 

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
