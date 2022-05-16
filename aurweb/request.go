package aurweb

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

const AUR = "https://aur.archlinux.org/rpc?v=5&"

func Search(field, query string) ([]Result, error) {
	v := url.Values{}
	v.Set("type", "search")
	v.Set("by", field)
	v.Set("arg", query)

	response, err := http.Get(AUR + v.Encode())
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(response.Body)
	search := SearchResponse{}
	search.Results = make([]Result, 0, 0)

	err = decoder.Decode(&search)
	if err != nil {
		return nil, err
	}

	if search.Type == "error" {
		return nil, errors.New(search.Error)
	}
	return search.Results, nil
}

func Info(arg []string) ([]Result, error) {
	v := url.Values{}
	v.Set("type", "info")
	for _, a := range arg {
		v.Add("arg[]", a)
	}

	response, err := http.Get(AUR + v.Encode())
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(response.Body)
	search := SearchResponse{}
	search.Results = make([]Result, 0, 0)

	err = decoder.Decode(&search)
	if err != nil {
		return nil, err
	}

	if search.Type == "error" {
		return nil, errors.New(search.Error)
	}
	return search.Results, nil
}
