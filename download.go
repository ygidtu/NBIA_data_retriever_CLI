package main

import (
	"archive/tar"

	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"


	"github.com/rs/zerolog/log"

	"fmt"

	"io/ioutil"

	"net/url"

	"strings"

	"net/http"
)

type FileInfo struct {
	Url         string
	Collection  string
	PatientId   string
	StudyUID    string
	SeriesUID   string
	Size        int64
	NumOfImages int
	Date        string
	Total 		[]string
}

func (info *FileInfo) Get() {
	log.Info().Msgf("Getting %s", info.Url)
	resp, err := http.Post(info.Url, "application/x-www-form-urlencoded; charset=ISO-8859-1", bytes.NewReader([]byte("")))

	if err != nil {
		log.Error().Msgf("%v", err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Msgf("%v", err)
	}

	data := strings.Split(string(content), "|")
	info.Total = data

	if len(data) < 11 {
		log.Error().Msgf("%v less than 11 elements", data)
	}

	info.Collection = data[0]
	info.PatientId = data[1]
	info.StudyUID = data[2]
	info.SeriesUID = data[3]

	if size, err := strconv.ParseInt(data[6], 10, 64); err != nil {
		log.Error().Msgf("%v", err)
	} else {
		info.Size = int64(size)
	}

	info.Date = data[11]

	log.Printf("%v", info)
}

func (info *FileInfo) GetOutput(output string) string {
	outputDir := filepath.Join(output, info.Collection, info.PatientId, info.StudyUID, info.Date, info.SeriesUID)

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err = os.MkdirAll(outputDir, 0755); err != nil {
			log.Error().Msgf("%v", err)
		}
	}

	return outputDir
}


// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func Untar(dst string, r io.Reader) error {

	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}

func (info *FileInfo) Download(output string) {

	log.Debug().Msgf("%v", info)

	outputFile := info.GetOutput(output)
	
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true},
		ResponseHeaderTimeout: timeout,
	}

	if proxy != "" {
		log.Info().Msgf(proxy)
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			log.Error().Msgf("%v", err)
		}

		tr = &http.Transport{
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			Proxy:                 http.ProxyURL(proxyURL),
			ResponseHeaderTimeout: timeout,
		}
	}

	// custom the request form
	form := url.Values{}
	form.Add("Range", fmt.Sprintf("bytes=0-"))
	form.Add("hasAnnotation", "false")
	form.Add("includeAnnotation", "true")
	form.Add("seriesUid", info.SeriesUID)
	form.Add("sopUids", "")
	form.Add("userId", "")
	form.Add("password", "")

	req, err := http.NewRequest("POST", baseURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Error().Msgf("%v", err)
	}
	// custom the request header
	req.Header.Add("password", "")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=ISO-8859-1")
	req.Header.Add("Connection", "Keep-Alive")

	log.Info().Msgf("Download %s to %s", info.SeriesUID, outputFile)

	client := &http.Client{Transport: tr, Timeout: timeout}

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Msgf("%v", err)
		os.Exit(1)
	}

	err = Untar(outputFile, resp.Body)

	if err != nil {
		log.Error().Msgf("%v", err)
	}
}

func (info *FileInfo) ToJson(output string) {
	rankingsJson, _ := json.MarshalIndent(info, "", "    ")
	err := ioutil.WriteFile(fmt.Sprintf("%s.json", info.GetOutput(output)), rankingsJson, 0644)

	if err != nil {
		log.Error().Msgf("%v", err)
	}
}

func (info *FileInfo) ToString() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%d\t%d\t%s", info.Url, info.Collection, info.PatientId, info.StudyUID, info.SeriesUID, info.Size, info.NumOfImages, info.Date)
}
