#!/bin/sh

go build --tags "icu json1 fts5 secure_delete" -ldflags='-s -w' -o ~/Public/gnote-linux-amd64
cd ~/Public
echo "Quit current gnote to allow update. Build windows version and manually sync to Public. Then hit enter"
read _junk
cp gnote-linux-amd64 ~/gnote/gnote
for f in gnote-windows-amd64.exe gnote-linux-amd64; do
	aws s3 cp $f s3://xvt-public-repo/pub/devops/ --profile xvt_aws
done
