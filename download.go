package main

import (
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

	"github.com/cheggaaa/pb/v3"
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

func (info *FileInfo) getOutput(output string) string {
	outputDir := filepath.Join(output, info.Collection, info.PatientId, info.Date)

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err = os.MkdirAll(outputDir, 0755); err != nil {
			log.Error().Msgf("%v", err)
		}
	}

	outputFile := fmt.Sprintf("%s/%s.dcm", outputDir, strings.Replace(info.SeriesUID, ".", "-", -1))

	return outputFile
}

func (info *FileInfo) Download(output string) {

	outputFile := info.getOutput(output)
	info.ToJson(outputFile)

	log.Info().Msgf("Download %s to %s", info.Url, outputFile)

	var start int64
	if stat, err := os.Stat(outputFile); os.IsNotExist(err) {
		start = 0
		_, err = os.Create(outputFile)
	} else {
		log.Debug().Msgf("%s exists", outputFile)
		start = stat.Size()
	}

	if start >= info.Size {
		log.Info().Msgf("Skip")
		return
	} else {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true},
			ResponseHeaderTimeout: timeout,
		}

		if proxy != "" {
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

		req, err := http.NewRequest("POST", baseURL, nil)
		if err != nil {
			log.Error().Msgf("%v", err)
		}
		// custom the request header
		req.Header.Add("password", "")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=ISO-8859-1")
		req.Header.Add("Connection", "Keep-Alive")

		// custom the request form
		form, _ := url.ParseQuery(req.URL.RawQuery)
		form.Add("Range", fmt.Sprintf("bytes=%d-%d", start, info.Size))
		form.Add("hasAnnotation", "false")
		form.Add("includeAnnotation", "true")
		form.Add("seriesUid", info.SeriesUID)
		form.Add("sopUids", "")
		form.Add("userId", "")
		form.Add("password", "")
		req.URL.RawQuery = form.Encode()

		client := &http.Client{Transport: tr, Timeout: timeout}

		resp, err := client.Do(req)
		if err != nil {
			log.Error().Msgf("%v", err)
		}

		writer, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		defer writer.Close()

		reader := io.LimitReader(resp.Body, info.Size-start)

		// start new bar
		bar := pb.Full.Start64(info.Size - start)
		// create proxy reader
		barReader := bar.NewProxyReader(reader)
		// copy from proxy reader
		_, err = io.Copy(writer, barReader)
		// finish bar
		bar.Finish()

		if err != nil {
			log.Error().Msgf("%v", err)
		}
	}

}

func (info *FileInfo) ToJson(outputFile string) {
	rankingsJson, _ := json.MarshalIndent(info, "", "    ")
	err := ioutil.WriteFile(fmt.Sprintf("%s.json", outputFile), rankingsJson, 0644)

	if err != nil {
		log.Error().Msgf("%v", err)
	}
}


func (info *FileInfo) ToString() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%d\t%d\t%s", info.Url, info.Collection, info.PatientId, info.StudyUID, info.SeriesUID, info.Size, info.NumOfImages, info.Date)
}
