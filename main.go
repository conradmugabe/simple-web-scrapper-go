package main

import (
	"fmt"
	"os"

	companynames "github.com/conradmugabe/simple-web-scrapper-go/src"
)

const (
	fileName  = "companies.txt"
	directory = "."
)

func main() {
	fs := os.DirFS(directory)

	companies, err := companynames.FromTextFile(fs, fileName)
	if err != nil {
		panic(err)
	}

	for _, company := range companies {
		url, _ := companynames.ConstructURLWithParams("https://www.google.com/search", map[string]string{"q": company.Name})
		body, _ := companynames.GetWebsiteContent(url)
		urls := companynames.ExtractURLs(body)
		companyFacebookURL := companynames.GetAllFacebookLink(urls)
		companyFacebookAboutPageURL := companynames.AddAboutLinkSuffix(companyFacebookURL)
		aboutPageBody, _ := companynames.GetWebsiteContent(companyFacebookAboutPageURL)
		emails := companynames.ExtractEmailsFromText(aboutPageBody)
		company.Email = companynames.GetFirstEntryInList(emails)
		fmt.Println(company)
	}
}
