package main

import (
	"archive/tar"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	downloadLinePrefix = "downloadServerUrl="
	versionPrefix      = "manifestVersion="
	downloadURLConnect = "numberOfSeries=1&series1="
)

var baseURL string

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

		if strings.HasPrefix(line, downloadLinePrefix) {
			baseURL = strings.Replace(line, downloadLinePrefix, "", -1)
		} else if strings.HasPrefix(line, versionPrefix) {
			version := strings.Replace(line, versionPrefix, "", -1)

			if strings.HasSuffix(version, ".0") {
				version = strings.Replace(version, ".0", "", -1)
			}
			baseURL = fmt.Sprintf("%sV%s", baseURL, version)
		} else if !strings.ContainsAny(line, "=") {
			res = append(res, &FileInfo{URL: fmt.Sprintf("%s?%s%s", baseURL, downloadURLConnect, line)})
		}
	}

	return res
}

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
func (info *FileInfo) Get(token string) {
	logger.Infof("Getting %s", info.URL)

	req, err := http.NewRequest("POST", info.URL, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	content, err := RespContent(req)
	if err != nil {
		logger.Fatal(err)
	}

	data := strings.Split(string(content), "|")
	info.Total = data

	if len(data) < 11 {
		logger.Fatalf("%v less than 11 elements", data)
	}

	info.Collection = data[0]
	info.PatientID = data[1]
	info.StudyUID = data[2]
	info.SeriesUID = data[3]

	if size, err := strconv.ParseInt(data[6], 10, 64); err != nil {
		logger.Fatal(err)
	} else {
		info.Size = int64(size)
	}

	info.Date = data[11]

	logger.Debugf("%v", info)
}

// GetOutput construct the output directory
func (info *FileInfo) GetOutput(output string) string {
	outputDir := filepath.Join(output, info.Collection, info.PatientID, info.StudyUID, info.Date, info.SeriesUID)

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err = os.MkdirAll(outputDir, 0755); err != nil {
			logger.Fatal(err)
		}
	}

	return outputDir
}

// Download is real function to download file
func (info *FileInfo) Download(output, username, password string) error {
	outputFile := info.GetOutput(output)

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
		return fmt.Errorf("failed to create request: %v", err)
	}
	// custom the request header
	req.Header.Add("password", "")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=ISO-8859-1")
	req.Header.Add("Connection", "Keep-Alive")

	logger.Infof("Download %s to %s", info.SeriesUID, outputFile)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %v", err)
	}
	defer resp.Body.Close()

	return UnTar(outputFile, resp.Body)
}

// ToJSON is used to save download file information and log the download progress
func (info *FileInfo) ToJSON(output string) {
	rankingsJSON, _ := json.MarshalIndent(info, "", "    ")

	f, err := os.OpenFile(fmt.Sprintf("%s.json", info.GetOutput(output)), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logger.Fatalf("failed to open json file: %v", err)
	}
	_, err = f.Write(rankingsJSON)
	if err != nil {
		logger.Fatalf("failed to write json file: %v", err)
	}
	if err := f.Close(); err != nil {
		logger.Fatalf("failed to close json file: %v", err)
	}
}

// ToString is used to convert file info as string for human
func (info *FileInfo) ToString() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%d\t%d\t%s", info.URL, info.Collection, info.PatientID, info.StudyUID, info.SeriesUID, info.Size, info.NumOfImages, info.Date)
}

/*
UnTar takes a destination path and a reader; a tar reader loops over the tarfile
creating the file structure at 'dst' along the way, and writing any files
*/
func UnTar(dst string, r io.Reader) error {

	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return fmt.Errorf("unknown error while untar: %v", err)

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

		// if it's a dir, and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, os.ModePerm); err != nil {
					return fmt.Errorf("failed to create dir while untar: %v", err)
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file while untar: %v", err)
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return fmt.Errorf("failed to copy data while untar: %v", err)
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			return f.Close()
		}
	}
}

// ToJSON as name says
func ToJSON(files []*FileInfo, output string) {
	rankingsJSON, _ := json.MarshalIndent(files, "", "    ")
	f, err := os.OpenFile(output, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logger.Fatalf("failed to open json file: %v", err)
	}
	_, err = f.Write(rankingsJSON)
	if err != nil {
		logger.Fatalf("failed to write json file: %v", err)
	}
	if err := f.Close(); err != nil {
		logger.Fatalf("failed to close json file: %v", err)
	}

	if err != nil {
		logger.Fatal(err)
	}
}
