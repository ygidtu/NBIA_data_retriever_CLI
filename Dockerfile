FROM golang:latest
COPY . /src
WORKDIR /src
RUN env GOOS=linux GOARCH=amd64 go build -o nbia_cli_linux_amd64 .
ENTRYPOINT ["/src/nbia_cli_linux_amd64"]
