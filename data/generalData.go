package data

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var (
	URLCONFIRMED string = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_global.csv"
	URLDEATHS    string = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_deaths_global.csv"
	URLRECOVERED string = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_recovered_global.csv"
	name         string = "time_series_covid19_global"
	globalJson   []DataJSON
)

func GeneralData(country string) ([]byte, error) {

	exist := fileExist(true)

	if !exist {
		//Execute function for create a file
		if errorRequest := generateJson(); errorRequest != nil {
			log.Println("Error to execute function for create a json file from csv")
			return nil, errorRequest
		}
		if err := loadJSON("general", &globalJson); err != nil {
			log.Println("Error into load json file function")
		}

		if country != "" {
			filterByCountry(country, &globalJson)
		}

	} else {
		//Load local file
		if err := loadJSON("general", &globalJson); err != nil {
			log.Println("Error into load json file function")
		}

		if country != "" {
			filterByCountry(country, &globalJson)
		}
	}

	encode, errEncode := json.Marshal(&globalJson)
	if errEncode != nil {
		log.Println("Error to encode data json to bytes")
		return nil, errEncode
	}

	go func() {
		nameToDelete := fiveFormat + "general.json"
		_, errInfo := os.Stat(nameToDelete)
		if !os.IsNotExist(errInfo) {
			delete(nameToDelete)
		}
	}()

	return encode, nil
}

func generateJson() error {

	wg := &sync.WaitGroup{}
	ch1 := make(chan *[]DataJSON)
	ch2 := make(chan *[]DataJSON)
	ch3 := make(chan *[]DataJSON)
	wg.Add(3)

	go requestConfirmed("c", wg, ch1)
	go requestConfirmed("d", wg, ch2)
	go requestConfirmed("r", wg, ch3)

	globalJson := <-ch1
	globalJson = <-ch2
	globalJson = <-ch3

	wg.Wait()

	jsonData, err := json.Marshal(globalJson)
	if err != nil {
		log.Println("Error to convert data to json")
	}

	if err := Transform(&jsonData, true); err != nil {
		log.Println("Error to create json")
	}

	return nil
}

func requestConfirmed(typeRequest string, wg *sync.WaitGroup, ch chan *[]DataJSON) {
	var resp *http.Response
	var err error

	switch typeRequest {
	case "c":
		resp, err = http.Get(URLCONFIRMED)
	case "d":
		resp, err = http.Get(URLDEATHS)
	case "r":
		resp, err = http.Get(URLRECOVERED)
	default:
		resp, err = http.Get(URLCONFIRMED)
	}

	if err != nil {
		log.Println("Error to load data, please check source")
	}

	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)

	dataReader, errRead := reader.ReadAll()
	if errRead != nil {
		log.Println("Error to convert data csv, please contact to administrator")
	}

	for id, row := range dataReader {
		if id == 0 {
			continue
		}
		counter, _ := strconv.Atoi(row[len(row)-1])
		if typeRequest == "c" {
			confirmToStruct(row[1], counter, &globalJson)
		} else if typeRequest == "d" {
			deathToStruct(row[1], counter, &globalJson)
		} else {
			recoverToStruct(row[1], counter, &globalJson)
		}

	}

	ch <- &globalJson
	wg.Done()
}

func confirmToStruct(country string, number int, allData *[]DataJSON) {
	var localData DataJSON
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

func deathToStruct(country string, number int, allData *[]DataJSON) {
	var localData DataJSON
	if len(*allData) > 0 {
		if ok := validateExist(country, allData); ok != -1 {
			(*allData)[ok].Deaths += int32(number)
		} else {
			localData.Country = country
			localData.Deaths = int32(number)
			*allData = append(*allData, localData)
		}
	} else {
		localData.Country = country
		localData.Deaths = int32(number)
		*allData = append(*allData, localData)
	}
}

func recoverToStruct(country string, number int, allData *[]DataJSON) {
	var localData DataJSON
	if len(*allData) > 0 {
		if ok := validateExist(country, allData); ok != -1 {
			(*allData)[ok].Recovered += int32(number)
		} else {
			localData.Country = country
			localData.Recovered = int32(number)
			*allData = append(*allData, localData)
		}
	} else {
		localData.Country = country
		localData.Recovered = int32(number)
		*allData = append(*allData, localData)
	}
}
