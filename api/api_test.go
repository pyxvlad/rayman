package api_test

import (
	"gitlab.com/rayone121/rayman/api"
	"testing"
)

func TestAPI(t *testing.T) {
	a, err := api.New()
	if err != nil {
		panic(err)
		return
	}
	err = a.Listen()
	if err != nil {
		panic(err)
		return
	}

}
