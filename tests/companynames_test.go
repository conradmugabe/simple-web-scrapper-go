package companynames_test

import (
	"errors"
	"reflect"
	"testing"
	"testing/fstest"

	companynames "github.com/conradmugabe/simple-web-scrapper-go/src"
)

func TestReadTextFileErrorsWhenFileNotFound(t *testing.T) {
	t.Parallel()
	fs := fstest.MapFS{
		"test.txt":  {Data: []byte("")},
		"test2.txt": {Data: []byte("hello world 2")},
		"test3.txt": {Data: []byte("hello world 3")},
	}
	fileName := "test4.txt"

	_, err := companynames.FromTextFile(fs, fileName)
	assertError(t, err, companynames.ErrCannotReadFile)
}

func TestReadTextFile(t *testing.T) {
	t.Parallel()
	fs := fstest.MapFS{
		"test.txt":  {Data: []byte("")},
		"test2.txt": {Data: []byte("hello world 2")},
		"test3.txt": {Data: []byte("hello world 3")},
	}
	fileName := "test.txt"

	companies, err := companynames.FromTextFile(fs, fileName)
	if err != nil {
		t.Fatal(err)
	}

	want := []companynames.Company{}
	if len(companies) != len(want) {
		t.Errorf("got %d, wanted %d", len(companies), len(want))
	}
}

func TestReadTestFileContent(t *testing.T) {
	t.Parallel()
	fs := fstest.MapFS{
		"test.txt":  {Data: []byte("hello\nworld")},
		"test2.txt": {Data: []byte("hello world 2")},
		"test3.txt": {Data: []byte("hello world 3")},
	}
	fileName := "test.txt"
	companies, _ := companynames.FromTextFile(fs, fileName)

	if len(companies) != 2 {
		t.Errorf("got %v, wanted %v", len(companies), 2)
	}

	assertCompany(t, companies[0], companynames.Company{Name: "hello"})
	assertCompany(t, companies[1], companynames.Company{Name: "world"})
}

func assertCompany(t *testing.T, got, want companynames.Company) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, wanted %+v", got, want)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if errors.Is(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
