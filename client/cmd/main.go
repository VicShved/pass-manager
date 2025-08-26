package main

import "github.com/VicShved/pass-manager/client/internal/tui"

//go:generate env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=v0.1 -X 'main.BuildTime=$(date )'" -o  client_linux main.go
// go:generate env GOOS=darwin GOARCH=amd64 go build -ldflags "-X config.Version=v0.1 -X 'main.BuildTime=$(date )'" -o client_mac main.go
// go:generate env GOOS=windows GOARCH=amd64 go build -ldflags "-X config.Version=v0.1 -X 'main.BuildTime=$(date )'" -o client_win.exe main.go

var (
	Version   string
	BuildTime string
)

func main() {
	tui.StartTui(Version, BuildTime)
}
