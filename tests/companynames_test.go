package companynames_test

import (
	"net/http"
	"net/http/httptest"
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
	assert.Nil(t, companies)
	assert.Equal(t, err, companynames.ErrCannotReadFile)
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
	assert.Nil(t, err)

	var want []companynames.Company
	assert.Equal(t, companies, want, "got %q, wanted %q", companies, want)
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
	assert.Nil(t, err)
	assert.Equal(t, len(companies), 2)
	assert.Equal(t, companies[0], companynames.Company{Name: "hello"})
	assert.Equal(t, companies[1], companynames.Company{Name: "world"})
}

func TestConstructURLWithParams(t *testing.T) {
	t.Run("constructs URL with params", func(t *testing.T) {
		baseURL := "http://example.com"
		params := map[string]string{
			"param1": "value 1",
			"param2": "value 2",
		}

		want := "http://example.com?param1=value+1&param2=value+2"

		got, err := companynames.ConstructURLWithParams(baseURL, params)
		assert.Nil(t, err)
		assert.Equal(t, got, want, "got %q, wanted %q", got, want)
	})
}

func TestGetWebsiteContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))

	defer server.Close()

	want := "hello world"

	got, err := companynames.GetWebsiteContent(server.URL)
	assert.Nil(t, err)
	assert.Equal(t, got, want, "got %q, wanted %q", got, want)
}

func TestExtractURLs(t *testing.T) {
	cases := []struct {
		name string
		data string
		want []string
	}{{
		name: "extracts URLs from string",
		data: "Lorem http://example.com neque, www.openai.com sapien. http://example.org/",
		want: []string{"http://example.com", "www.openai.com", "http://example.org/"},
	},
		{

			name: "returns empty array if no URLs found",
			data: "Lorem neque, sapien. Interdum.",
			want: nil,
		}}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := companynames.ExtractURLs(tt.data)
			assert.Equal(t, got, tt.want, "got %q, wanted %q", got, tt.want)
		})
	}
}

func TestGetAllFacebookLinks(t *testing.T) {
	cases := []struct {
		name string
		want string
		urls []string
	}{{
		name: "returns the first facebook URL in the list",
		want: "https://www.facebook.com/Test",
		urls: []string{"https://www.facebook.com/Test/amp&1",
			"https://www.facebook.com/Test2/amp&2",
			"https://www.facebook.com/Test3/amp&3"},
	}, {
		name: "returns empty array if no URLs found",
		want: "",
		urls: nil,
	}}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := companynames.GetAllFacebookLink(tt.urls)
			assert.Equal(t, got, tt.want, "got %q, wanted %q", got, tt.want)
		})
	}
}

func TestAddAboutLinkSuffix(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want string
	}{{
		name: "adds about link suffix",
		url:  "https://www.facebook.com/Test",
		want: "https://www.facebook.com/Test/about",
	}, {
		name: "adds about link suffix",
		url:  "",
		want: "/about",
	}}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := companynames.AddAboutLinkSuffix(tt.url)
			assert.Equal(t, got, tt.want, "got %q, wanted %q", got, tt.want)
		})
	}
}

func TestGetEmailsFromText(t *testing.T) {
	cases := []struct {
		name string
		data string
		want []string
	}{{
		name: "extracts email from string",
		data: "Lorem user@example.com neque, test@test.org sapien. @example.org/",
		want: []string{"user@example.com", "test@test.org"},
	}, {
		name: "returns empty array if no email found",
		data: "Lorem neque, sapien. Interdum.",
		want: nil,
	}}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := companynames.ExtractEmailsFromText(tt.data)
			assert.Equal(t, got, tt.want, "got %q, wanted %q", got, tt.want)
		})
	}
}

func TestGetFirstEntryInList(t *testing.T) {
	cases := []struct {
		name string
		data []string
		want string
	}{
		{
			name: "returns first entry in list",
			data: []string{"Test", "Test2", "Test3"},
			want: "Test",
		},
		{
			name: "returns empty string if list is empty",
			data: []string{},
			want: "",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := companynames.GetFirstEntryInList(tt.data)
			assert.Equal(t, got, tt.want, "got %q, wanted %q", got, tt.want)
		})
	}
}
