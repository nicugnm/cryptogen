package clients

import (
	"cryptogen-retrieve/domain"
	"cryptogen-retrieve/gateways/repositories"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"
)

const CRYPTO_TYPE_MAP_URL = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/map"
const CRYPTO_DATA_MAP_URL = "https://sandbox-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest"

type ClientRequests struct {
}

func ClientsRequests() *ClientRequests {
	return &ClientRequests{}
}

var _ CryptoRequests = (*ClientRequests)(nil)

func (c ClientRequests) RequestCryptoTypes(nb chan CryptoTypeChannel, wg *sync.WaitGroup) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", CRYPTO_TYPE_MAP_URL, nil)

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

	nb <- CryptoTypeChannel{
		Response: resp,
		Error:    err,
	}

	close(nb)

	wg.Done()
}

func (c ClientRequests) RequestCryptoDetails(nbType chan CryptoTypeChannel, nbMetadata chan CryptoMetadataChannel, wg *sync.WaitGroup) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", CRYPTO_DATA_MAP_URL, nil)

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	listSymbols := ""

	for get := range nbType {
		if get.Error != nil {
			log.Println(get.Error)
		}

		log.Println(get.Response.Status)
		respBody, _ := io.ReadAll(get.Response.Body)

		m := CryptoTypeRequest{}
		err := json.Unmarshal(respBody, &m)

		if err != nil {
			fmt.Errorf("error during unmarshall %x", err)
		}

		for index, cryptoMetadata := range m.Data {
			// verify a string only contains letters, numbers, underscores and dashes
			if !regexp.MustCompile(`^[A-Za-z0-9_-]*$`).MatchString(cryptoMetadata.Symbol) {
				continue
			}

			// if is the last, we don't need to have "," on final
			// maximum 1000 crypto, request does not support more than 1000 symbols
			if index == 1000-1 {
				listSymbols += cryptoMetadata.Symbol
				break
			} else {
				listSymbols += cryptoMetadata.Symbol + ","
			}
		}

		/*dir := "data_to_file"

		os.Mkdir(dir, 0777)

		fileName := path.Join(dir, "list-symbols.json")

		os.WriteFile(fileName, []byte(listSymbols), 0666)

		fmt.Printf("List symbols: %s\n", listSymbols)*/
	}

	q := url.Values{}
	q.Add("symbol", listSymbols)

	godotenv_err := godotenv.Load()

	if godotenv_err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	nbMetadata <- CryptoMetadataChannel{
		Response: resp,
		Error:    err,
	}

	close(nbMetadata)

	wg.Done()
}

func (c ClientRequests) SaveDataToRepository(nb chan CryptoMetadataChannel, wg *sync.WaitGroup) {
	for get := range nb {
		if get.Error != nil {
			log.Println(get.Error)
		}

		log.Println(get.Response.Status)
		respBody, _ := io.ReadAll(get.Response.Body)

		m := CryptoMetadataRequest{}
		err := json.Unmarshal(respBody, &m)

		if err != nil {
			fmt.Errorf("Error during unmarshall %s", err)
		}

		/*fmt.Println()
		fmt.Printf("MaxSuply: %s\n", m.Data["BTC"].MaxSuply)
		fmt.Printf("TotalSuply: %s\n", m.Data["BTC"].TotalSuply)
		fmt.Printf("Quote: %s\n", m.Data["BTC"].Quote)

		fmt.Printf("Price: %f\n", m.Data["BTC"].Quote.USD.Price)
		fmt.Printf("LastUpdated: %s\n", m.Data["BTC"].Quote.USD.LastUpdated)
		fmt.Printf("MarketCap: %f\n", m.Data["BTC"].Quote.USD.MarketingCap)
		fmt.Printf("Volume: %f\n", m.Data["BTC"].Quote.USD.Volume)
		fmt.Printf("Percent 1h: %f\n", m.Data["BTC"].Quote.USD.PercentChange1h)*/

		var cryptoDataList []domain.CryptoDataMetadata

		for _, cryptoData := range m.Data {
			cryptoDataList = append(cryptoDataList, cryptoData)
		}

		redisRepository := repositories.NewRedisRepo()
		redisErr := redisRepository.SaveDataList(cryptoDataList)

		if redisErr != nil {
			fmt.Errorf("Error during saving data in redis %x", redisErr)
		}
	}

	wg.Done()
}
