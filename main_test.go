package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	hf := http.HandlerFunc(handler)
	hf.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `Hello World!`
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestRouter(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)
	res, err := http.Get(mockServer.URL + "/hello")
	if err != nil {
		t.Fatal(err)
	}
	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	respString := string(b)
	expected := `Hello World!`
	if respString != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", respString, expected)
	}
}

func TestRouterForNonExistentRoute(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)
	res, err := http.Get(mockServer.URL + "/non-existent")
	if err != nil {
		t.Fatal(err)
	}
	if status := res.StatusCode; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestRouterForUnsupportedRouteMethod(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)
	res, err := http.Post(mockServer.URL+"/hello", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if status := res.StatusCode; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	respString := string(b)
	expected := ""
	if respString != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", respString, expected)
	}
}

func TestStaticFileServer(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)
	res, err := http.Get(mockServer.URL + "/assets/")
	if err != nil {
		t.Fatal(err)
	}
	// We want our status to be 200 (ok)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status should be 200, got %d", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("Wrong content type, expected %s, got %s", expectedContentType, contentType)
	}
}

func TestGetBirdsRoute(t *testing.T) {
	birds = []Bird{
		Bird{"Accipitriformes", "Sharp beaked birds of prey eg Hawks, Eagles"},
		Bird{"Apodiformes", "Tiny and underdeveloped feet eg humming birds"},
	}
	r := newRouter()
	mockServer := httptest.NewServer(r)
	res, err := http.Get(mockServer.URL + "/birds")
	if err != nil {
		t.Fatal(err)
	}
	// We want our status to be 200 (ok)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status should be 200, got %d", res.StatusCode)
	}

	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	respString := string(b)
	expected, err := json.Marshal(birds)
	if respString != string(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", respString, expected)
	}
}

func TestCreateBirdsRoute(t *testing.T) {
	r := newRouter()
	// initialize birds
	birds = []Bird{}
	mockServer := httptest.NewServer(r)
	birdSpecies := "Test Species"
	birdDesc := "Test description"
	res, err := http.PostForm(mockServer.URL+"/birds", map[string][]string{
		"species":     []string{birdSpecies},
		"description": []string{birdDesc},
	})
	if err != nil {
		t.Fatal(err)
	}
	// We want our status to be 200 (ok)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status should be %d, got %d", http.StatusOK, res.StatusCode)
	}

	if len(birds) != 1 {
		t.Errorf("Expected number of birds to be %v got %v", 1, len(birds))
	} else {
		bird := birds[0]
		if bird.Species != birdSpecies || bird.Description != birdDesc {
			t.Errorf("Expected bird species to be %s and description to be %s, got %s and %s", birdSpecies, birdDesc, bird.Species, bird.Description)
		}
	}

}
