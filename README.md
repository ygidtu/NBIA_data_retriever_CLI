# NBIA data retriever CLI

> A simple replacement of NBIA data retriever.

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
