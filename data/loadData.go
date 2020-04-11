package data

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	current_time        = time.Now().Local().AddDate(0, 0, -1)
	URLTOLOAD    string = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_daily_reports/"
	currentDay   string = current_time.Format("01-02-2006")
)

func LoadData() ([][]string, error) {
	resp, err := http.Get(URLTOLOAD + currentDay + ".csv")

	if err != nil {
		log.Println("Error to load data, please check source")
		return nil, err
	}

	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)

	dataReader, err := reader.ReadAll()
	if err != nil {
		log.Println("Error to convert data csv, please contact to administrator")
		return nil, err
	}

	var dataToJson DataJSON
	var allJson []DataJSON

	for id, row := range dataReader {
		if id == 0 {
			continue
		}
		dataToJson.Country = row[3]
		i64, _ := strconv.ParseInt(row[7], 10, 32)
		dataToJson.Confirmed = int32(i64)
		i64, _ = strconv.ParseInt(row[8], 10, 32)
		dataToJson.Deaths = int32(i64)
		i64, _ = strconv.ParseInt(row[9], 10, 32)
		dataToJson.Recovered = int32(i64)

		allJson = append(allJson, dataToJson)
	}

	jsonData, err := json.Marshal(allJson)
	if err != nil {
		log.Println("Error to convert data to json")
	}

	if err := Transform(&jsonData, currentDay); err != nil {
		log.Println("Error to create json")
	}

	return nil, nil
}
