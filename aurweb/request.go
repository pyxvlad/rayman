package aurweb

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

const AUR = "https://aur.archlinux.org/rpc"
const AURv5 = AUR + "?v=5&"

func Search(field, query string) (search SearchResponse, err error) {
	v := url.Values{}
	v.Set("type", "search")
	v.Set("by", field)
	v.Set("arg", query)

	var response *http.Response
	response, err = http.Get(AURv5 + v.Encode())
	if err != nil {
		return
	}

	decoder := json.NewDecoder(response.Body)
	search.Results = make([]Result, 0, 0)

	err = decoder.Decode(&search)
	if err != nil {
		return
	}

	if search.Type == "error" {
		err = errors.New(search.Error)
		return
	}
	return
}

func Info(arg []string) (info InfoResponse,err error) {
	v := url.Values{}
	v.Set("type", "info")
	for _, a := range arg {
		v.Add("arg[]", a)
	}

	response, err := http.Get(AURv5 + v.Encode())
	if err != nil {
		return
	}

	decoder := json.NewDecoder(response.Body)
	info.Results = make([]InfoResult, 0, 0)

	err = decoder.Decode(&info)
	return
}
