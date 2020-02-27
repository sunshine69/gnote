#!/bin/sh

#go build --tags "icu json1 fts5 secure_delete" -ldflags='-s -w' -o ~/Public/gnote-linux-amd64

ARCH=$(uname -m)

docker run --rm -v $(pwd):/work --entrypoint /usr/local/go/bin/go --workdir /work golang-ubuntu1804-build:latest build --tags "json1 fts5 secure_delete" -ldflags='-s -w' -o gnote-ubuntu1804-${ARCH}

mv gnote-ubuntu1804-${ARCH} ~/Public/

#go build --tags "json1 fts5 secure_delete" -ldflags='-s -w' -o ~/Public/gnote-linux-amd64

cd ~/Public
echo "Quit current gnote to allow update. Build windows version and manually sync to Public. Then hit enter"
read _junk

cp -a gnote-ubuntu1804-${ARCH} ~/gnote/gnote
for f in gnote-bundle-windows-amd64.7z gnote-ubuntu1804-${ARCH}; do
	aws s3 cp $f s3://xvt-public-repo/pub/devops/ --profile xvt_aws
done
