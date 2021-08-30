# chanbus

![build.yaml](https://github.com/go-laeo/chanbus/actions/workflows/build.yaml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/go-laeo/chanbus.svg)](https://pkg.go.dev/github.com/go-laeo/chanbus) ![golangci.yaml](https://github.com/go-laeo/chanbus/actions/workflows/golangci-lint.yaml/badge.svg)

Package chanbus is a 1:n channel broadcasting for golang.

## Install

```shell
go get github.com/go-laeo/chanbus
```

## Example

```golang
// First 0 means no buffer of main channel, and second
// 0 means blocking on forwarding to derived channels.
ch := chanbus.New(1, 0)

// Send value to main channel.
ch.Send("say, hi")

// Derives new readonly channel
nch := ch.Derive(1)

for v := range nch {
    fmt.Println(v) // say, hi
}
```

## License

The MIT license.
