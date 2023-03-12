package clients

import (
	"cryptogen-retrieve/gateways/repositories"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
)

const NUM_REQUESTS = 1
const CRYPTO_MAP_URL = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/map"

type ClientRequests struct {
}

func ClientsRequests() *ClientRequests {
	return &ClientRequests{}
}

var _ CryptoRequests = (*ClientRequests)(nil)

func (c ClientRequests) RequestCryptoTypes(nb chan NonBlocking) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", CRYPTO_MAP_URL, nil)

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "5000")
	//q.Add("convert", "USD")

	godotenv_err := godotenv.Load()

	if godotenv_err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	nb <- NonBlocking{
		Response: resp,
		Error:    err,
	}
}

func (c ClientRequests) RequestCryptoDetails(nb chan NonBlocking) []*interface{} {
	//TODO implement me
	panic("implement me")
}

func (c ClientRequests) SaveDataToRepository(nb chan NonBlocking, wg *sync.WaitGroup) {
	for get := range nb {
		if get.Error != nil {
			log.Println(get.Error)
		}

		log.Println(get.Response.Status)
		respBody, _ := io.ReadAll(get.Response.Body)

		m := CryptoRequest{}
		err := json.Unmarshal(respBody, &m)

		if err != nil {
			fmt.Errorf("Error during unmarshall %x", err)
		}

		redisRepository := repositories.NewRedisRepo()
		redisErr := redisRepository.SaveDataList(m.Data)

		if err != nil {
			fmt.Errorf("Error during saving data in redis %x", redisErr)
		}

		wg.Done()
	}
}

func (c ClientRequests) SaveDataToFile(nb chan NonBlocking, wg *sync.WaitGroup) {
	for get := range nb {
		if get.Error != nil {
			log.Println(get.Error)
		}

		log.Println(get.Response.Status)
		respBody, _ := io.ReadAll(get.Response.Body)

		dir := "data_to_file"

		os.Mkdir(dir, 0777)

		fileName := path.Join(dir, "cryptogen-data.json")

		os.WriteFile(fileName, respBody, 0666)

		m := CryptoRequest{}
		err := json.Unmarshal(respBody, &m)

		if err != nil {
			fmt.Errorf("Error during unmarshall %x", err)
		}

		wg.Done()
	}
}
