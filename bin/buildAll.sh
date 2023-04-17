#!/bin/bash

docker kill poulpe-build ; docker rm poulpe-build
docker run --name poulpe-build -v /app golang:1.20 go version &> /dev/null
docker cp . poulpe-build:/app
docker run --rm --volumes-from poulpe-build -w /app/ golang:1.20 bash -c 'echo "Fetching dependencies" && GOPROXY=https://proxy.golang.org go get -v && \
mkdir -p release
echo "Building"
export GOOS=linux && export GOARCH=386                         && go get && go build -buildvcs=false -o release/poulpe-$GOOS-$GOARCH &&            echo "done : "$GOOS"-"$GOARCH;
export GOOS=linux && export GOARCH=amd64                       && go get && go build -buildvcs=false -o release/poulpe-$GOOS-$GOARCH &&            echo "done : "$GOOS"-"$GOARCH;

export GOOS=linux && export GOARCH=arm && export GOARM=6       && go get && go build -buildvcs=false -o release/poulpe-$GOOS-$GOARM-$GOARCH &&     echo "done : $GOOS-$GOARM-$GOARCH";
export GOOS=linux && export GOARCH=arm64 && export GOARM=6     && go get && go build -buildvcs=false -o release/poulpe-$GOOS-$GOARM-$GOARCH &&     echo "done : $GOOS-$GOARM-$GOARCH";

export GOOS=linux && export GOARCH=arm && export GOARM=7       && go get && go build -buildvcs=false -o release/poulpe-$GOOS-$GOARM-$GOARCH &&     echo "done : $GOOS-$GOARM-$GOARCH";
export GOOS=linux && export GOARCH=arm64 && export GOARM=7     && go get && go build -buildvcs=false -o release/poulpe-$GOOS-$GOARM-$GOARCH &&     echo "done : $GOOS-$GOARM-$GOARCH";

export GOOS=darwin && export GOARCH=amd64                      && go get && go build -buildvcs=false -o release/poulpe-$GOOS-$GOARCH &&            echo "done : "$GOOS"-"$GOARCH;

export GOOS=windows && export GOARCH=386                       && go get && go build -buildvcs=false -o release/poulpe-$GOOS-$GOARCH.exe &&        echo "done : "$GOOS"-"$GOARCH;
export GOOS=windows && export GOARCH=amd64                     && go get && go build -buildvcs=false -o release/poulpe-$GOOS-$GOARCH.exe &&        echo "done : "$GOOS"-"$GOARCH;
echo "done all"
'
docker cp poulpe-build:/app/release ./
