package data

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	current_time        = time.Now().Local().AddDate(0, 0, -1)
	URLTOLOAD    string = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_daily_reports/"
	currentDay   string = current_time.Format("01-02-2006")
	allJson      []DataJSON
)

func LoadData() ([]byte, error) {

	// Check if file exist
	_, errInfo := os.Stat(currentDay + ".json")

	if os.IsNotExist(errInfo) {
		//Execute function for create a file
		if errorRequest := requestToCSV(); errorRequest != nil {
			log.Println("Error to execute function for create a json file from csv")
			return nil, errorRequest
		}

		if err := loadJSON(); err != nil {
			log.Println("Error into load json file function")
		}
	} else {
		//Load local file
		if err := loadJSON(); err != nil {
			log.Println("Error into load json file function")
		}
	}

	encode, errEncode := json.Marshal(&allJson)
	if errEncode != nil {
		log.Println("Error to encode data json to bytes")
		return nil, errEncode
	}

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

	var dataToJson DataJSON

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

	return nil
}

func loadJSON() error {
	jsonFile, err := os.Open(currentDay + ".json")

	if err != nil {
		log.Println("Error to load json file")
		return err
	}

	defer jsonFile.Close()

	bytesValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(bytesValue, &allJson)

	return nil
	//return allJson, nil
}
