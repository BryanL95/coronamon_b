package data

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	URLTOLOAD string = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_daily_reports/"
	allJson   []DataJSON
)

func LoadData(country string) ([]byte, error) {

	// Check if file exist
	exist := fileExist(false)

	if !exist {
		//Execute function for create a file
		if errorRequest := requestToCSV(); errorRequest != nil {
			log.Println("Error to execute function for create a json file from csv")
			return nil, errorRequest
		}

		if err := loadJSON("", &allJson); err != nil {
			log.Println("Error into load json file function")
		}

		if country != "" {
			filterByCountry(country, &allJson)
		}

	} else {
		//Load local file
		if err := loadJSON("", &allJson); err != nil {
			log.Println("Error into load json file function")
		}

		if country != "" {
			filterByCountry(country, &allJson)
		}
	}

	encode, errEncode := json.Marshal(&allJson)
	if errEncode != nil {
		log.Println("Error to encode data json to bytes")
		return nil, errEncode
	}

	go func() {
		nameToDelete := fiveFormat + ".json"
		_, errInfo := os.Stat(nameToDelete)
		if !os.IsNotExist(errInfo) {
			delete(nameToDelete)
		}
	}()

	return encode, nil
}

func requestToCSV() error {
	resp, err := http.Get(URLTOLOAD + currentDay + ".csv")

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
		confirmed, _ := strconv.Atoi(row[7])
		death, _ := strconv.Atoi(row[8])
		recovered, _ := strconv.Atoi(row[9])
		dataToStruct(confirmed, death, recovered, row[3], &allJson)
	}

	jsonData, err := json.Marshal(&allJson)
	if err != nil {
		log.Println("Error to convert data to json")
	}

	if err := Transform(&jsonData, false); err != nil {
		log.Println("Error to create json")
	}

	return nil
}

func dataToStruct(confirmed int, deaths int, recovered int, country string, dataJson *[]DataJSON) {
	var dataDaily DataJSON
	if len(*dataJson) > 0 {
		if ok := validateExist(country, dataJson); ok != -1 {
			(*dataJson)[ok].Confirmed += int32(confirmed)
			(*dataJson)[ok].Deaths += int32(deaths)
			(*dataJson)[ok].Recovered += int32(recovered)
		} else {
			dataDaily.Country = country
			dataDaily.Confirmed = int32(confirmed)
			dataDaily.Deaths = int32(deaths)
			dataDaily.Recovered = int32(recovered)
			*dataJson = append(*dataJson, dataDaily)
		}
	} else {
		dataDaily.Country = country
		dataDaily.Confirmed = int32(confirmed)
		dataDaily.Deaths = int32(deaths)
		dataDaily.Recovered = int32(recovered)
		*dataJson = append(*dataJson, dataDaily)
	}
}
