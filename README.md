# NBIA data retriever CLI

> A simple replacement of NBIA data retriever.

---

## Command line usage

```bash
SYNOPSIS:
    NBIA_data_retriever_CLI [--debug] [--help|-h] [--image-url <string>]
                            [--input|-i <string>] [--meta|-m]
                            [--meta-url <string>] [--output|-s <string>]
                            [--passwd|-w <string>] [--processes|-p <int>]
                            [--proxy|-x <string>] [--timeout|-t <int>]
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

    --passwd|-w <string>    password for control data (default: "")

    --processes|-p <int>    start how many download at same time (default: 1)

    --proxy|-x <string>     the proxy to use [http, socks5://user:passwd@host:port] (default: "")

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

