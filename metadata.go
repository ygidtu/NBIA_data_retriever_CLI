package main

import (
	"fmt"
	"net/http"
	"os"
)

func saveResponseToFile(u, token, series, proxy, output string) error {
	url_, err := makeURL(u, map[string]interface{}{"SeriesUID": series})
	req, err := http.NewRequest("GET", url_, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	content, err := RespContent(req)
	if err != nil {
		return fmt.Errorf("failed to read response data: %v", err)
	}
	f, err := os.OpenFile(fmt.Sprintf("%s.json", output), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	_, err = f.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write data: %v", err)
	}
	return f.Close()
}

func RetrieveMeta(token, series, proxy, output string) error {
	err := saveResponseToFile(TagUrl, token, series, proxy, fmt.Sprintf("%s.json", output))
	if err != nil {
		return err
	}
	return saveResponseToFile(MetaUrl, token, series, proxy, fmt.Sprintf("%s.json", output))
}
