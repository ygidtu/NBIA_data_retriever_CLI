package main

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)


func ToJson(files []*FileInfo, output string) {
	rankingsJson, _ := json.MarshalIndent(files, "", "    ")
	err := ioutil.WriteFile(output, rankingsJson, 0644)

	if err != nil {
		log.Error().Msgf("%v", err)
	}
}
