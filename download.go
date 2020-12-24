package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"

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

// FileInfo is struct handle all the information of one tcia file
type FileInfo struct {
	URL         string
	Collection  string
	PatientID   string
	StudyUID    string
	SeriesUID   string
	Size        int64
	NumOfImages int
	Date        string
	Total       []string
}

// Get main API to decode the file information
func (info *FileInfo) Get() {
	log.Info().Msgf("Getting %s", info.URL)
	resp, err := http.Post(info.URL, "application/x-www-form-urlencoded; charset=ISO-8859-1", bytes.NewReader([]byte("")))

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
	info.PatientID = data[1]
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

// GetOutput construct the output directory
func (info *FileInfo) GetOutput(output string) string {
	outputDir := filepath.Join(output, info.Collection, info.PatientID, info.StudyUID, info.Date, info.SeriesUID)

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err = os.MkdirAll(outputDir, 0755); err != nil {
			log.Error().Msgf("%v", err)
		}
	}

	return outputDir
}

// Download is real function to downlaod file
func (info *FileInfo) Download(output, username, password string) error {

	log.Debug().Msgf("%v", info)

	outputFile := info.GetOutput(output)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true},
		ResponseHeaderTimeout: timeout,
	}

	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return err
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
	form.Add("userId", username)
	form.Add("password", password)

	req, err := http.NewRequest("POST", baseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	// custom the request header
	req.Header.Add("password", "")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=ISO-8859-1")
	req.Header.Add("Connection", "Keep-Alive")

	log.Info().Msgf("Download %s to %s", info.SeriesUID, outputFile)

	client := &http.Client{Transport: tr, Timeout: timeout}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	return UnTar(outputFile, resp.Body)
}

// ToJSON is used to save downlaod file information and log the download progress
func (info *FileInfo) ToJSON(output string) {
	rankingsJSON, _ := json.MarshalIndent(info, "", "    ")
	err := ioutil.WriteFile(fmt.Sprintf("%s.json", info.GetOutput(output)), rankingsJSON, 0644)

	if err != nil {
		log.Error().Msgf("%v", err)
	}
}

// ToString is used to convert file info as string for human
func (info *FileInfo) ToString() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%d\t%d\t%s", info.URL, info.Collection, info.PatientID, info.StudyUID, info.SeriesUID, info.Size, info.NumOfImages, info.Date)
}
