package main

import (
	"github.com/parnurzeal/gorequest"
	"net/http"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"strings"
	"errors"
)

type RegionData []struct {
	Buildconfig   string `json:"buildconfig"`
	Buildid       string `json:"buildid"`
	Cdnconfig     string `json:"cdnconfig"`
	Keyring       string `json:"keyring,omitempty"`
	Region        string `json:"region"`
	Versionsname  string `json:"versionsname"`
	Productconfig string `json:"productconfig,omitempty"`
}


func makeRequest(url string) (*http.Response, string, []error) {
	return gorequest.New().Get(url).End()
}

func getRemoteVersionFile(code string, file string) (string, RegionData, error) {
	url := fmt.Sprintf("http://us.patch.battle.net:1119/%s/%s", code, file)

	resp, body, errs := makeRequest(url)

	if errs != nil {
		for e := 0; e < len(errs); e++ {
			fmt.Errorf("error: %s", errs[e])
		}

		errorMap := make([]map[string]string, 1)
		errorMap[0] = make(map[string]string)
		errorMap[0]["error"] = "Something happened internally"
		return errorMap[0]["error"], RegionData{}, errors.New("Something internally happened")
	}

	if resp.StatusCode == 404 || resp.StatusCode == 500 {
		errorMap := make([]map[string]string, 1)
		errorMap[0] = make(map[string]string)
		errorMap[0]["error"] = "Something happened or game wasn't found"
		return errorMap[0]["error"], RegionData{}, errors.New("Something happened or game wasn't found")
	}

	body = strings.TrimSpace(body)

	var result RegionData
	err := mapstructure.Decode(parseVersionFile(body), &result)
	if err != nil {
		panic(err)
	}

	return string(body), result, nil
}

func parseVersionFile(file string) []map[string]string {
	lines := strings.Split(file,"\n")
	keys := strings.Split(lines[0], `|`)
	keysList := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		keyList := strings.Split(keys[i], `!`)

		keysList[i] = strings.ToLower(keyList[0])
	}

	data := make([]map[string]string, len(lines) - 1)
	for i := 1; i < len(lines); i++ {
		local := make(map[string]string)

		lineData := strings.Split(lines[i], `|`)

		for x := 0; x < len(keysList); x++ {
			if len(lineData[x]) > 0 {
				local[keysList[x]] = lineData[x]
			}
		}

		data[i - 1] = local
	}

	return data
}

func Version(code string, file string) (string, RegionData, error) {
	return getRemoteVersionFile(code, file)
}