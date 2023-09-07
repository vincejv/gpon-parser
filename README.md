# GPON Parser

Supports the following ONT models

* FiberHome HG6245D (Globe Telecom Philippines branded)
* FiberHome AN5506_04F1A (Globe Telecom Philippines branded)

### Compiling
    go build -ldflags "-s -w"

### ARM Build on Windows
    $env:GOARCH='arm'
    $env:GOOS='linux'