# NBIA data retriever CLI

> A simple replacement of NBIA data retriever.

---

https://www.cancerimagingarchive.net/wp-content/uploads/CMB-AML_October-2023-manifest.tcia

## Command line usage

```bash
Usage: nbia_cli_darwin [global options]                                       
                                                                              
Global options:                                                               
        -i, --input   Path to tcia file
        -o, --output  Output directory, or output file when --meta enabled (default: downloads)            
        -x, --proxy   Proxy
        -t, --timeout Due to limitation of target server, please set this time out value as big as possible (default: 1200000)
        -p, --process Start how many download at same time (default: 1)
        -m, --meta    Get Meta info of all files
        -u, --username Username for control data
        -w, --passwd   Password for control data
        -v, --version Show version
            --debug   Show debug info
            --help    Show this help
```

---

### [Update 2020.12.24]

Add `--username` and `--passwd`, maybe this is usedful retrive the restriced data.

>> I do not have an account for  NBIA, therefore, this is no tested yet.

---

### [Update 2019.09.17]

Just noticed original NBIA add tar wrapper of real dcm files

Now I add a tar wrapper to decompress the dcm files.
At the same time, I cannot check the download progress of single file anymore.
Therefore, I use a json file to record information of single seriesUID, and mark the relevant file of the seriesUID has been downloaded.

---

Issues with NBIA data retriever:

- Cannot resume download, if there is any error occurs, have to download all files from the beginning
- Swing is kind of heavy, and cannot run it in server

---
Advantages:

- Proxy like `socks5://127.0.0.1:1080` or `http://127.0.0.1:1080`
- Resume download
- Command line

---

Known issues:

- The `public.cancerimagingarchive.net/nbia-download/servlet` use `POST` to transfer data from server to local
, the connection may be terminated even before the download is complete. Therefore, **PLEASE** set timeout as huge as possible
- progress bar is a mess when using multiple process
- I do not have a account of NBIA, therefore this program could not handle the restricted data for now.
