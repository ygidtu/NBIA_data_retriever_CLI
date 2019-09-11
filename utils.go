package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
)

func ToFile(files []*FileInfo, output string) {
	w, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Error().Msgf("write error %s", err)
	}
	defer w.Close()

	writer := bufio.NewWriter(w)
	for _, f := range files {
		fmt.Fprintln(writer, f.ToString())
	}
	writer.Flush()
}

func ToJson(files []*FileInfo, output string) {
	rankingsJson, _ := json.MarshalIndent(files, "", "    ")
	err := ioutil.WriteFile(output, rankingsJson, 0644)

	if err != nil {
		log.Error().Msgf("%v", err)
	}
}
