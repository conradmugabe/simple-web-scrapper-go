package companynames

import (
	"bufio"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"regexp"
)

type Company struct {
	Name string
}

func CompanyNamesFromTextFile(fileSystem fs.FS, fileName string) ([]Company, error) {
	file, err := fileSystem.Open(fileName)
	if err != nil {
		return nil, err
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
	const URLRegex = `(?i)\b((?:https?://|www\d{0,3}[.]|[a-z0-9.\-]+[.][a-z]{2,4}/)(?:[^\s()<>]+|\(([^\s()<>]+|(\([^\s()<>]+\)))*\))+(?:\(([^\s()<>]+|(\([^\s()<>]+\)))*\)|[^\s` + "`" + `!()\[\]{};:'".,<>?«»“”‘’]))`

	URLs := regexp.MustCompile(URLRegex).FindAllString(data, -1)
	return URLs
}
