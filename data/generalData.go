package data

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
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
	globalJson   []GeneralJSON
)

func GeneralData() ([]byte, error) {

	_, errInfo := os.Stat("general.json")

	if os.IsNotExist(errInfo) {
		//Execute function for create a file
		if errorRequest := generateJson(); errorRequest != nil {
			log.Println("Error to execute function for create a json file from csv")
			return nil, errorRequest
		}

		if err := loadGeneralJSON(); err != nil {
			log.Println("Error into load json file function")
		}
	} else {
		//Load local file
		if err := loadGeneralJSON(); err != nil {
			log.Println("Error into load json file function")
		}
	}

	encode, errEncode := json.Marshal(&globalJson)
	if errEncode != nil {
		log.Println("Error to encode data json to bytes")
		return nil, errEncode
	}

	return encode, nil
}

func generateJson() error {

	wg := &sync.WaitGroup{}
	ch1 := make(chan *[]GeneralJSON)
	ch2 := make(chan *[]GeneralJSON)
	ch3 := make(chan *[]GeneralJSON)
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

	if err := TransformGeneral(&jsonData, currentDay); err != nil {
		log.Println("Error to create json")
	}

	return nil
}

func requestConfirmed(typeRequest string, wg *sync.WaitGroup, ch chan *[]GeneralJSON) {
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
		//return &globalJson
	}

	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)

	dataReader, errRead := reader.ReadAll()
	if errRead != nil {
		log.Println("Error to convert data csv, please contact to administrator")
		//return &globalJson
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
	//return &globalJson
}

func confirmToStruct(country string, number int, allData *[]GeneralJSON) {
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

func deathToStruct(country string, number int, allData *[]GeneralJSON) {
	var localData GeneralJSON
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

func recoverToStruct(country string, number int, allData *[]GeneralJSON) {
	var localData GeneralJSON
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

/* Funciones para refactor */
func validateExist(country string, allData *[]GeneralJSON) int {
	for key, val := range *allData {
		if val.Country == country {
			return key
		}
	}
	return -1
}

func loadGeneralJSON() error {
	jsonFile, err := os.Open("general.json")

	if err != nil {
		log.Println("Error to load json file")
		return err
	}

	defer jsonFile.Close()

	bytesValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(bytesValue, &globalJson)

	return nil
}
