package data

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	current_time        = time.Now().Local().AddDate(0, 0, -1)
	URLTOLOAD    string = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_daily_reports/"
	currentDay   string = current_time.Format("01-02-2006")
	allJson      []DataJSON
)

func LoadData(country string) ([]byte, error) {

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

		if country != "" {
			filterByCountryLast(country, &allJson)
		}

	} else {
		//Load local file
		if err := loadJSON(); err != nil {
			log.Println("Error into load json file function")
		}

		if country != "" {
			filterByCountryLast(country, &allJson)
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

	//var dataToJson DataJSON

	for id, row := range dataReader {
		if id == 0 {
			continue
		}
		/*dataToJson.Country = row[3]
		i64, _ := strconv.ParseInt(row[7], 10, 32)
		dataToJson.Confirmed = int32(i64)
		i64, _ = strconv.ParseInt(row[8], 10, 32)
		dataToJson.Deaths = int32(i64)
		i64, _ = strconv.ParseInt(row[9], 10, 32)
		dataToJson.Recovered = int32(i64)

		allJson = append(allJson, dataToJson)*/
		confirmed, _ := strconv.Atoi(row[7])
		death, _ := strconv.Atoi(row[8])
		recovered, _ := strconv.Atoi(row[9])
		dataToStruct(confirmed, death, recovered, row[3], &allJson)
	}

	jsonData, err := json.Marshal(&allJson)
	if err != nil {
		log.Println("Error to convert data to json")
	}

	if err := Transform(&jsonData, currentDay); err != nil {
		log.Println("Error to create json")
	}

	return nil
}

func dataToStruct(confirmed int, deaths int, recovered int, country string, dataJson *[]DataJSON) {
	var dataDaily DataJSON
	if len(*dataJson) > 0 {
		if ok := validateExistDayly(country, dataJson); ok != -1 {
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

/*Refactor, se puede crear una funcion global que utilice generalData y locaData*/
func validateExistDayly(country string, allData *[]DataJSON) int {
	for key, val := range *allData {
		if val.Country == country {
			return key
		}
	}
	return -1
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
}

func filterByCountryLast(country string, jsonFilter *[]DataJSON) {
	var tmp []DataJSON
	for _, val := range *jsonFilter {
		if country == strings.ToLower(val.Country) {
			tmp := append(tmp, val)
			(*jsonFilter) = tmp
			break
		}
	}
}
