FROM golang:latest
COPY . /src
WORKDIR /src
RUN env GOOS=linux GOARCH=amd64 go build -o nbia_cli_linux_amd64 main.go tcia.go download.go utils.go
ENTRYPOINT ["/src/nbia_cli_linux_amd64"]
