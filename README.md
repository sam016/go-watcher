Watcher [![GoDoc](https://godoc.org/github.com/canthefason/go-watcher?status.svg)](https://godoc.org/github.com/canthefason/go-watcher) [![Build Status](https://travis-ci.org/canthefason/go-watcher.svg?branch=master)](https://travis-ci.org/canthefason/go-watcher)
=======

Watcher is a command line tool inspired by [fresh](https://github.com/pilu/fresh) and used for watching .go file changes, and restarting the app in case of an update/delete/add operation.

Most of the existing file watchers have a configuration burden, and even though Go has a really short build time, this configuration burden makes your binaries really hard to run right away. With Watcher, we aimed simplicity in configuration, and tried to make it as simple as possible.

Now ships with the awesome [Delve](https://github.com/go-delve/delve), debugger for Go

## Installation

Get the package with:

`go get github.com/sam016/go-watcher`

Install the binary under go/bin folder:

`go install github.com/sam016/go-watcher/watcher/cmd/watcher`

After this step, please make sure that your go/bin folder is appended to PATH environment variable.

## Usage

### Step-1

```bash
cd /path/to/myapp
```

### Step-2

Create a config file `watcher.yaml`

```yaml
watcher:
  run: hello
  watch: hello
package:
  path: "/go/src/hello"
  args:
  - --argXYZ=some-value
  - --argABC=another-value
  - --argEMPTY
delve-args:
  - --headless
  - --continue
  - --accept-multiclient
  - --api-version=2
  - --listen=:2345
  - --output=/tmp/__debug_bin
  - --log
  - --log-dest=debugger.log"
```

### Step-3

Start watcher:

`watcher -f watcher.yaml`

Watcher works like your native package binary. You can pass all your existing package arguments to the Watcher, which really lowers the learning curve of the package, and makes it practical.

### Current app usage

`myapp --argXYZ=some-value --argABC=another-value --argEMPTY`

### With watcher

Move the arguments to config file

```yaml
package:
  args:
    - --argXYZ=some-value
    - --argABC=another-value
    - --argEMPTY
```

```bash
watcher -f watcher.yaml
```

As you can see nothing changed between these two calls. When you run the command, Watcher starts watching folders recursively, starting from the current working directory. It only watches .go and .tmpl files and ignores hidden folders and _test.go files.

## Watcher in Docker

If you want to run Watcher in a containerized local environment, you can achieve this by using [sam016/go-watcher](https://hub.docker.com/r/sam016/go-watcher/) image in Docker Hub. There is an example project under [/docker-example](https://github.com/sam016/go-watcher/tree/dockerfile-gvm/docker-examples) directoy. Let's try to dockerize this example code first.

In our example, we are creating a server that listens to port 7000 and responds to all clients with "watcher is running" string. The most essential thing to run your code in Docker is, mounting your project volume to a container. In the containerized Watcher, our GOPATH is set to /go directory by default, so you need to mount your project to this GOPATH.

```bash
docker run \
  sam016/go-watcher\
  -v /path/to/hello:/go/src/hello \
  -p 7000:7000 \
  watcher -f watcher.yaml
```

To provide a more structured repo, we also integrated a docker-compose manifest file. That file already handles volume mounting operation that and exposes the port to the localhost. With docker-compose the only thing that you need to do from the root, invoking `docker-compose up

## Known Issues
On Mac OS X, when you make a tls connection, you can get a message like: x509: `certificate signed by unknown authority`

You can resolve this problem by setting CGO_ENABLED=0
https://github.com/golang/go/issues/14514
https://codereview.appspot.com/22020045

## TODO

[] Fix tests

## Author

* [Can Yucel](http://canthefason.com) **(Original)**
* [sam016](http://github.com/sam016)

## License

The MIT License (MIT) - see LICENSE.md for more details
