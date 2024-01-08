package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// decodeTCIA is used to decode the tcia file
func decodeTCIA(path string) []*FileInfo {
	res := make([]*FileInfo, 0)

	f, err := os.Open(path)

	if err != nil {
		logger.Fatal(err)
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.ContainsAny(line, "=") {

			url_, err := makeURL(MetaUrl, map[string]interface{}{"SeriesInstanceUID": line})
			req, err := http.NewRequest("GET", url_, nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

			resp, err := client.Do(req)
			if err != nil {
				logger.Errorf("failed to do request: %v", err)
				continue
			}

			content, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if err != nil {
				logger.Errorf("failed to read response data: %v", err)
				continue
			}

			files := make([]*FileInfo, 0)
			err = json.Unmarshal(content, &files)
			if err != nil {
				logger.Errorf("failed to parse response data: %v - %s", err, content)
			}

			res = append(res, files...)
		}
	}

	return res
}

type FileInfo struct {
	NumberOfImages     string `json:"Number of Images"`
	SOPClassUID        string `json:"SOP Class UID"`
	Manufacturer       string `json:"Manufacturer"`
	DataDescriptionURI string `json:"Data Description URI"`
	LicenseURL         string `json:"License URL"`
	AnnotationSize     string `json:"Annotation Size"`
	Collection         string `json:"Collection"`
	StudyDescription   string `json:"Study Description"`
	SeriesUID          string `json:"Series UID"`
	StudyUID           string `json:"Study UID"`
	LicenseName        string `json:"License Name"`
	StudyDate          string `json:"Study Date"`
	SeriesDescription  string `json:"Series Description"`
	Modality           string `json:"Modality"`
	RdPartyAnalysis    string `json:"3rd Party Analysis"`
	FileSize           string `json:"File Size"`
	SubjectID          string `json:"Subject ID"`
	SeriesNumber       string `json:"Series Number"`
}

// GetOutput construct the output directory
func (info *FileInfo) getOutput(output string) string {
	outputDir := filepath.Join(output, info.SubjectID, info.StudyDate)

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err = os.MkdirAll(outputDir, 0755); err != nil {
			logger.Fatal(err)
		}
	}

	return outputDir
}

func (info *FileInfo) MetaFile(output string) string {
	return filepath.Join(info.getOutput(output), fmt.Sprintf("%s.json", info.SeriesUID))
}

func (info *FileInfo) DcimFiles(output string) string {
	return filepath.Join(info.getOutput(output), fmt.Sprintf("%s.zip", info.SeriesUID))
}

func (info *FileInfo) GetMeta(output string) error {
	f, err := os.OpenFile(info.MetaFile(output), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open meta file %s: %v", info.MetaFile(output), err)
	}
	content, err := json.MarshalIndent(info, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshall meta: %v", err)
	}
	_, err = f.Write(content)
	if err != nil {
		return err
	}
	return f.Close()
}

// Download is real function to download file
func (info *FileInfo) Download(output string) error {

	url_, err := makeURL(ImageUrl, map[string]interface{}{"SeriesInstanceUID": info.SeriesUID})
	req, err := http.NewRequest("GET", url_, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %v", err)
	}
	defer resp.Body.Close()

	if err != nil {
		return fmt.Errorf("failed to read response data: %v", err)
	}
	f, err := os.OpenFile(info.DcimFiles(output), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	bar := bytesBar(resp.ContentLength, info.SeriesUID)
	if fSize, err := strconv.Atoi(info.FileSize); err == nil {
		bar = bytesBar(int64(fSize), info.SeriesUID)
	}
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write data: %v", err)
	}
	return f.Close()
}
