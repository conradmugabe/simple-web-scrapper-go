package main

import (
	"os"

	companynames "github.com/conradmugabe/simple-web-scrapper-go/src"
	log "github.com/sirupsen/logrus"
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
		log.Info("Constructing google search url for ", company.Name)
		url, _ := companynames.ConstructURLWithParams("https://www.google.com/search", map[string]string{"q": company.Name})
		log.Info("Making a google search for ", company.Name)
		body, _ := companynames.GetWebsiteContent(url)
		log.Info("Extracting facebook url for ", company.Name)
		urls := companynames.ExtractURLs(body)
		companyFacebookURL := companynames.GetAllFacebookLink(urls)
		companyFacebookAboutPageURL := companynames.AddAboutLinkSuffix(companyFacebookURL)
		log.Info("Extracting email for ", company.Name, " from ", companyFacebookAboutPageURL)
		aboutPageBody, _ := companynames.GetWebsiteContent(companyFacebookAboutPageURL)
		emails := companynames.ExtractEmailsFromText(aboutPageBody)
		company.Email = companynames.GetFirstEntryInList(emails)
		log.Info("Email for ", company.Name, " is ", company.Email)
	}
}
