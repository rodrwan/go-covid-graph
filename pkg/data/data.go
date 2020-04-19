package data

import (
	"encoding/csv"
	"io"
	"net/http"
	"strconv"
)

const (
	regionalConfirmedCasesURL = "https://raw.githubusercontent.com/MinCiencia/Datos-COVID19/master/output/producto4/2020-04-19-CasosConfirmados-totalRegional.csv"
)

type Case struct {
	Region string
	New    int
	Total  int
	Dead   int
}

func GetRegionalConfirmedCases() ([]*Case, error) {
	req, err := http.NewRequest("GET", regionalConfirmedCasesURL, nil)
	if err != nil {
		return nil, err
	}

	clt := http.Client{}
	resp, err := clt.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Parse the file
	r := csv.NewReader(resp.Body)
	cases := make([]*Case, 0)

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		n, _ := strconv.Atoi(record[1])
		t, _ := strconv.Atoi(record[2])
		d, _ := strconv.Atoi(record[4])

		cases = append(cases, &Case{
			Region: record[0],
			New:    n,
			Total:  t,
			Dead:   d,
		})
	}

	return cases, nil
}
