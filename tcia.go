package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

/*
Decode .tcia file


downloadServerUrl=https://public.cancerimagingarchive.net/nbia-download/servlet/DownloadServlet
includeAnnotation=true
noOfrRetry=4
databasketId=manifest-1561766528292.tcia
manifestVersion=3.0
ListOfSeriesToDownload=
1.3.6.1.4.1.14519.5.2.1.2932.1975.258485580179562855968279889580
1.3.6.1.4.1.14519.5.2.1.2932.1975.828151202600553380765273057255
1.3.6.1.4.1.14519.5.2.1.2932.1975.833599136481996339747054915378
1.3.6.1.4.1.14519.5.2.1.2932.1975.683517573824842427464991915817
1.3.6.1.4.1.14519.5.2.1.2932.1975.139344037244157968824624956143
1.3.6.1.4.1.14519.5.2.1.2932.1975.140019191302805149397246344266
1.3.6.1.4.1.14519.5.2.1.2932.1975.108765081996899446170362622894
1.3.6.1.4.1.14519.5.2.1.2932.1975.319738355605489346002421981300
1.3.6.1.4.1.14519.5.2.1.2932.1975.177636765165982037627828115568
1.3.6.1.4.1.14519.5.2.1.2932.1975.317522988137423573172585740106
1.3.6.1.4.1.14519.5.2.1.2932.1975.704503971523330107813722872940
1.3.6.1.4.1.14519.5.2.1.2932.1975.171054842668611060585132309012
1.3.6.1.4.1.14519.5.2.1.2932.1975.806694139290428616868542160946
*/

const downloadLinePrefix = "downloadServerUrl="
const versionPrefix = "manifestVersion="
const downloadUrlConnect = "numberOfSeries=1&series1="

func DecodeTCIA(path string) []*FileInfo {
	res := make([]*FileInfo, 0, 0)

	f, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
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
			res = append(res, &FileInfo{Url: fmt.Sprintf("%s?%s%s", baseURL, downloadUrlConnect, line)})
		}
	}

	return res
}
