package companynames_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"testing/fstest"

	companynames "github.com/conradmugabe/simple-web-scrapper-go/src"
)

func TestReadTextFileErrorsWhenFileNotFound(t *testing.T) {
	fs := fstest.MapFS{
		"test.txt":  {Data: []byte("")},
		"test2.txt": {Data: []byte("hello world 2")},
		"test3.txt": {Data: []byte("hello world 3")},
	}
	fileName := "test4.txt"
	_, err := companynames.CompanyNamesFromTextFile(fs, fileName)

	if err == nil {
		t.Errorf("got %v, wanted error", err)
	}
}

func TestReadTextFile(t *testing.T) {
	fs := fstest.MapFS{
		"test.txt":  {Data: []byte("")},
		"test2.txt": {Data: []byte("hello world 2")},
		"test3.txt": {Data: []byte("hello world 3")},
	}
	fileName := "test.txt"
	companies, err := companynames.CompanyNamesFromTextFile(fs, fileName)

	if err != nil {
		t.Fatal(err)
	}

	want := []companynames.Company{}
	if len(companies) != len(want) {
		t.Errorf("got %d, wanted %d", len(companies), len(want))
	}
}

func TestReadTestFileContent(t *testing.T) {
	fs := fstest.MapFS{
		"test.txt":  {Data: []byte("hello\nworld")},
		"test2.txt": {Data: []byte("hello world 2")},
		"test3.txt": {Data: []byte("hello world 3")},
	}
	fileName := "test.txt"
	companies, _ := companynames.CompanyNamesFromTextFile(fs, fileName)

	if len(companies) != 2 {
		t.Errorf("got %v, wanted %v", len(companies), 2)
	}

	assertPost(t, companies[0], companynames.Company{Name: "hello"})
	assertPost(t, companies[1], companynames.Company{Name: "world"})

}

func assertPost(t *testing.T, got companynames.Company, want companynames.Company) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, wanted %+v", got, want)
	}
}

func TestConstructURLWithParams(t *testing.T) {
	t.Run("constructs URL with params", func(t *testing.T) {
		baseURL := "http://example.com"
		params := map[string]string{
			"param1": "value 1",
			"param2": "value 2",
		}

		want := "http://example.com?param1=value+1&param2=value+2"

		got, _ := companynames.ConstructURLWithParams(baseURL, params)
		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	})
}

func TestGetWebsiteContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))

	defer server.Close()

	want := "hello world"

	got, _ := companynames.GetWebsiteContent(server.URL)
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
