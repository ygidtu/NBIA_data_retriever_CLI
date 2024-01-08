# NBIA data retriever CLI

> A simple replacement of NBIA data retriever.

---
## Installation

1. using the binary from this repository
2. compile from source

```bash
git clone https://github.com/ygidtu/NBIA_data_retriever_CLI.git
cd NBIA_data_retriever_CLI
go mod tidy   # prepare the dependencies

# check the build script arguments
python build.py --help

# build for current platform
python build.py

# build for most command platforms
python build.py --common

# build for specific platform and architecture
python build.py --platform linux --arch amd64

# build for all platforms
python build.py --all
```

3. using docker

   1. build docker image from source

    ```bash
    git clone https://github.com/ygidtu/NBIA_data_retriever_CLI.git
    cd NBIA_data_retriever_CLI
    docker build -t nbia .
    ```

    2. running docker

    ```bash
    docker run --rm -v $PWD:$PWD -w $PWD nbia --help
    ```


## Command line usage

```bash
SYNOPSIS:
    NBIA_data_retriever_CLI [--debug] [--help|-h] [--image-url <string>]
                            [--input|-i <string>] [--meta|-m]
                            [--meta-url <string>] [--output|-s <string>]
                            [--passwd <string>] [--processes|-p <int>]
                            [--prompt|-w] [--proxy|-x <string>] [--save-log]
                            [--token-url <string>] [--user|-u <string>]
                            [--version|-v] [<args>]

OPTIONS:
    --debug                 show more info (default: false)

    --help|-h               show help information (default: false)

    --image-url <string>    the api url to download image data (default: "https://services.cancerimagingarchive.net/nbia-api/services/v2/getImage")

    --input|-i <string>     path to input tcia file (default: "")

    --meta|-m               get Meta info of all files (default: false)

    --meta-url <string>     the api url get meta data (default: "https://services.cancerimagingarchive.net/nbia-api/services/v2/getSeriesMetaData")

    --output|-s <string>    Output directory, or output file when --meta enabled (default: "./")

    --passwd <string>       set password for control data in command line (default: "")

    --processes|-p <int>    start how many download at same time (default: 1)

    --prompt|-w             input password for control data (default: false)

    --proxy|-x <string>     the proxy to use [http, socks5://user:passwd@host:port] (default: "")

    --save-log              save debug log info to file (default: false)

    --token-url <string>    the api url of login token (default: "https://services.cancerimagingarchive.net/nbia-api/oauth/token")

    --user|-u <string>      username for control data (default: "nbia_guest")

    --version|-v            show version information (default: false)

```

This tool download files using the official API, the input file should be tcia file or just list of series instance ids.

 Basic usage
```bash
./NBIC_data_retriever_CLI -i path_to_tcia
```

---
`
Issues with NBIA data retriever:

- Cannot resume download, if there is any error occurs, have to download all files from the beginning
- Swing is kind of heavy, and cannot run it in server

---`
Advantages:

- Proxy like `socks5://127.0.0.1:1080` or `http://127.0.0.1:1080`
- Resume download
- Command line

