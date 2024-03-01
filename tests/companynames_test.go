package companynames_test

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

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

	companies, err := companynames.FromTextFile(fs, fileName)
	if assert.NotNil(t, err) {
		assert.Equal(t, err, companynames.ErrCannotReadFile)
	}

	assert.Nil(t, companies)
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

	companies, err := companynames.FromTextFile(fs, fileName)
	if err != nil {
		t.Fatal(err)
	}

	if len(companies) != 2 {
		t.Errorf("got %v, wanted %v", len(companies), 2)
	}

	assert.Equal(t, companies[0], companynames.Company{Name: "hello"})
	assert.Equal(t, companies[1], companynames.Company{Name: "world"})
}
