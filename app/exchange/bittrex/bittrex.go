package bittrex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aofiee/go-ex/app/helper"
	"github.com/aofiee/go-ex/app/model"
)

type bittrexURLEndpoint string

//Endpoint const
const Endpoint bittrexURLEndpoint = "https://bittrex.com/api/v1.1"

var configuration = model.Configuration{}

//BittrexSymbol string
type BittrexSymbol string

func init() {
	log.Println("init bittrex")
}
func (c BittrexSymbol) String() string {
	return string(c)
}

func (c bittrexURLEndpoint) String() string {
	return string(c)
}

//GetTicker BittrexSymbol  reciever
func (c *BittrexSymbol) GetTicker(APIKey string, APISecret string) (url string, hmacEncode string) {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	url = Endpoint.String()
	url += "/public/getticker?apikey=" + APIKey + "&nonce=" + timestamp
	enCryptsha512 := sha512.New()
	enCryptsha512.Write([]byte(url))
	hmac512 := hmac.New(sha512.New, []byte(APISecret))
	hmacEncode = base64.StdEncoding.EncodeToString(hmac512.Sum(nil))
	// log.Println(hmacEncode)
	return url, hmacEncode
}

func (c *BittrexSymbol) getTickerEndpointSlug(w http.ResponseWriter, r *http.Request) string {
	var aPath = "/v1/bittrex/get/ticker/"
	var slug string
	if strings.HasPrefix(r.URL.Path, aPath) {
		slug = r.URL.Path[len(aPath):]
	}
	return slug
}

//GetTickerEnpoint func
func GetTickerEnpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var config = configuration.GetConfiguration()
	var sym BittrexSymbol
	var url, hmacEndCode = sym.GetTicker(config.BittrexAPIKey, config.BittrexSecretKey)
	url += "&market=" + sym.getTickerEndpointSlug(w, r)
	// log.Println("xxxxxxxx" + url)
	// log.Println(hmacEndCode)

	spaceClient := http.Client{
		Timeout: time.Second * 200,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apisign", hmacEndCode)
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var result = []byte(body)
	str := fmt.Sprintf("%s", result)
	log.Println(str)

	var f interface{}
	_ = json.Unmarshal([]byte(body), &f)
	log.Println(f)

	m := f.(map[string]interface{})

	b, _ := json.Marshal(m)
	helper.SetHeader(w, b)
}
