package companynames

import (
	"bufio"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type Company struct {
	Name  string
	Email string
}

var ErrCannotReadFile = errors.New("failed to read file")

func FromTextFile(fileSystem fs.FS, fileName string) ([]Company, error) {
	file, err := fileSystem.Open(fileName)
	if err != nil {
		return nil, ErrCannotReadFile
	}
	defer file.Close()

	GetCompanies := func(file io.Reader) []Company {
		var companies []Company

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			companies = append(companies, Company{
				Name: scanner.Text(),
			})
		}

		return companies
	}
	companies := GetCompanies(file)

	return companies, nil
}

func ConstructURLWithParams(baseURL string, params map[string]string) (string, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	query := parsedURL.Query()
	
	for key, value := range params {
		query.Set(key, value)
	}
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

func GetWebsiteContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func ExtractURLs(data string) []string {
	const URLRegex = `(?i)\b((?:https?://|www\d{0,3}[.]|[a-z0-9.\-]+[.][a-z]
		{2,4}/)(?:[^\s()<>]+|\(([^\s()<>]+|(\([^\s()<>]+\)))*\))+(?:\(([^\s
		()<>]+|(\([^\s()<>]+\)))*\)|[^\s` + "`" + `!()\[\]{};:'".,<>?«»“”‘’]))`

	URLs := ExtractDataFromText(data, URLRegex)

	return URLs
}

const facebookURL = `https://www.facebook.com/`

func GetAllFacebookLink(urls []string) string {
	for _, URL := range urls {
		if strings.HasPrefix(URL, facebookURL) {
			trimmedURL := strings.TrimPrefix(URL, facebookURL)

			splitURL := strings.Split(trimmedURL, "/")
			if len(splitURL) > 1 {
				return facebookURL + splitURL[0]
			}

			return ""
		}
	}

	return ""
}

func AddAboutLinkSuffix(url string) string {
	return url + "/about"
}

func ExtractEmailsFromText(text string) []string {
	const EmailRegex = `[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}`
	emails := ExtractDataFromText(text, EmailRegex)

	return emails
}

func ExtractDataFromText(text string, regex string) []string {
	return regexp.MustCompile(regex).FindAllString(text, -1)
}

func GetFirstEntryInList(data []string) string {
	if len(data) > 0 {
		return data[0]
	}

	return ""
}
