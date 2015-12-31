package api

import (
	"net/http"
	"testing"
)

type varHandler func(*http.Request) map[string]string

func mockGetVars(vars map[string]string) varHandler {

	return func(r *http.Request) map[string]string {
		return vars
	}
}

func TestGetInt64IfOk(t *testing.T) {

	vars := map[string]string{
		"id": "11",
	}
	getVars = mockGetVars(vars)

	url := "/channels/11/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	result, err := getInt64(req, "id")
	if err != nil {
		t.Fatal(err)
	}

	if result != 11 {
		t.Errorf("Invalid parameter:%d", result)
	}

}

func TestGetInt64IfEmpty(t *testing.T) {

	vars := make(map[string]string)
	getVars = mockGetVars(vars)

	url := "/channels/11/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = getInt64(req, "id")
	if err == nil {
		t.Fail()
	}

}

func TestGetInt64IfNotInt(t *testing.T) {
	vars := map[string]string{
		"id": "foo",
	}

	getVars = mockGetVars(vars)

	url := "/channels/foo/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = getInt64(req, "id")
	if err == nil {
		t.Fail()
	}

}
