package companynames

import (
	"bufio"
	"io"
	"io/fs"
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
