package clients

import (
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

func Request(nb chan NonBlocking) {
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

func HandleResponse(nb chan NonBlocking, wg *sync.WaitGroup) {
	for get := range nb {
		if get.Error != nil {
			log.Println(get.Error)
		}

		log.Println("------DATA------")
		log.Println(get.Response.Status)
		log.Println("=======")
		respBody, _ := io.ReadAll(get.Response.Body)

		dir := "data_to_file"

		os.Mkdir(dir, 0777)

		fileName := path.Join(dir, "cryptogen-data.json")

		os.WriteFile(fileName, respBody, 0666)

		//fmt.Println(string(respBody))
		log.Println("------DATA------")

		wg.Done()
	}
}
