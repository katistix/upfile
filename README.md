# upfile

a simple file server you can just point to a folder and it will launch a nice clean minimalist fileserver.

goal: make it as simple as possible to just do `upfile .` and have that folder be accessible via this machine's ip or hostname on local network.

## install

requires [go](https://go.dev/dl/).

```sh
go install github.com/katistix/upfile@latest
```

## usage

```sh
upfile .              # serve the current folder on port 8000
upfile . --port 9000  # use a different port
```

open one of the printed URLs in any browser to browse and download files.
