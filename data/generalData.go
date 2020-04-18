package data

import (
	"encoding/csv"
	"encoding/json"
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

	for id, row := range dataReader {
		if id == 0 {
			continue
		}
		counter, _ := strconv.Atoi(row[len(row)-1])
		temporalToStruct(row[1], counter, &globalJson)
	}

	return nil
}

func temporalToStruct(country string, number int, allData *[]GeneralJSON) {
	var localData GeneralJSON
	if len(*allData) > 0 {
		if ok := validateExist(country, allData); ok != -1 {
			(*allData)[ok].Confirmed += int32(number)
		} else {
			localData.Country = country
			localData.Confirmed = int32(number)
			*allData = append(*allData, localData)
		}
	} else {
		localData.Country = country
		localData.Confirmed = int32(number)
		*allData = append(*allData, localData)
	}
}

func validateExist(country string, allData *[]GeneralJSON) int {
	for key, val := range *allData {
		if val.Country == country {
			return key
		}
	}
	return -1
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
