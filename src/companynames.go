package companynames

import (
	"bufio"
	"errors"
	"io"
	"io/fs"
)

type Company struct {
	Name string
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
