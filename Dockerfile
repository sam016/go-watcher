FROM golang:1.12
LABEL maintainer="devmaster@sam016.com"

RUN go get -v github.com/fatih/color
RUN go get -v gopkg.in/fsnotify.v1
RUN go get -v golang.org/x/sys/...
RUN go get -v gopkg.in/yaml.v2

COPY . /go/src/github.com/sam016/go-watcher

ENV VERSION 1.0.0

RUN go get -v github.com/sam016/go-watcher/watcher/...
RUN GOVERSION="$(go version)" go build \
    -o $GOPATH/bin/watcher \
    -ldflags "-X main.goversion='$GOVERSION' -X main.version=$VERSION -X main.commitID=000 -X main.buildTime=$(date +%s)" \
    github.com/sam016/go-watcher/watcher/cmd/watcher

RUN watcher --version

WORKDIR /go/src

ADD entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
CMD ["watcher"]
