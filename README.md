```sh
time CGOWRAP_CC=clang CC=cgowrap $(go env GOROOT)/pkg/tool/linux_amd64/cgo -debug-gcc -- $(pkg-config --cflags glib-2.0 gio-2.0) ./glib.go &> /tmp/cgo.out
```
