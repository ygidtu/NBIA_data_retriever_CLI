env GOOS=darwin GOARCH=amd64 go build -o nbia_cli_darwin main.go tcia.go download.go

env GOOS=darwin GOARCH=amd64 go build -o nbia_cli_linux_amd64 main.go tcia.go download.go

env GOOS=darwin GOARCH=amd64 go build -o nbia_cli_win64.exe main.go tcia.go download.go

