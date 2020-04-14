package data

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var (
	URLCONFIRMED string = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_global.csv"
	URLDEATHS    string = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_deaths_global.csv"
	URLRECOVERED string = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_recovered_global.csv"
	name         string = "time_series_covid19_global"
	globalJson   []GeneralJSON
)

func GeneralData() ([]byte, error) {
	err := requestConfirmed()

	if err != nil {
		fmt.Println("Error")
	}

	encode, errEncode := json.Marshal(&globalJson)
	if errEncode != nil {
		log.Println("Error to encode data json to bytes")
		return nil, errEncode
	}

	return encode, nil
}

func requestConfirmed() error {
	resp, err := http.Get(URLCONFIRMED)

	if err != nil {
		log.Println("Error to load data, please check source")
		return err
	}

	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)

	dataReader, err := reader.ReadAll()
	if err != nil {
		log.Println("Error to convert data csv, please contact to administrator")
		return err
	}

	temporal := make(map[string]int) //We can't access to memory address, because map is dinamic, so declare into function

	for id, row := range dataReader {
		if id == 0 {
			continue
		}
		counter, _ := strconv.Atoi(row[len(row)-1])
		if _, ok := temporal[row[1]]; ok {
			temporal[row[1]] += counter
		} else {
			temporal[row[1]] = counter
		}
	}

	if error := temporalToStruct(temporal); error != nil {
		log.Println(error)
		return error
	}

	return nil
}

func temporalToStruct(data map[string]int) error {
	if len(data) > 0 {
		for key, val := range data {
			var localData GeneralJSON
			localData.Country = key
			localData.Confirmed = int32(val)

			globalJson = append(globalJson, localData)
		}
	} else {
		return errors.New("Size data from global confirmed is empty")
	}

	return nil
}

/*func requestDeaths() {

}

func requestRecovered() {

}

func requestToGeneralCSV() error {

	jsonData, err := json.Marshal(allJson)
	if err != nil {
		log.Println("Error to convert data to json")
	}

	if err := Transform(&jsonData, currentDay); err != nil {
		log.Println("Error to create json")
	}

	return nil
}
*/
