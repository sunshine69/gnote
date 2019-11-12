#!/bin/sh

go build --tags "icu json1 fts5 secure_delete" -ldflags='-s -w' -o gnote-linux-amd64
rm -rf  /tmp/gnote-linux-amd64 >/dev/null 2>&1 ; mkdir -p /tmp/gnote-linux-amd64
cp -a glade icons gnote-linux-amd64 /tmp/gnote-linux-amd64/
CWD=$(pwd)
cd /tmp
tar Jcf ~/Public/gnote-linux-amd64.tar.xz gnote-linux-amd64
rm -rf gnote-linux-amd64
cd ~/Public
echo "Build windows version and manually sync to Public. Then hit enter"
read _junk
for f in gnote-windows-amd64.exe gnote-windows-amd64.7z gnote-linux-amd64.tar.xz; do
	aws s3 cp $f s3://xvt-public-repo/pub/devops/ --profile xvt_aws
done
