package binance

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aofiee/go-ex/app/helper"
	"github.com/aofiee/go-ex/app/model"
)

type binanceURLEndpoint string

//Endpoint const
const Endpoint binanceURLEndpoint = "https://api.binance.com"

var configuration = model.Configuration{}

//BinaceSymbol string
type BinaceSymbol string

func init() {
	log.Println("init binance")
}
func (c BinaceSymbol) String() string {
	return string(c)
}

func (c binanceURLEndpoint) String() string {
	return string(c)
}

//GetDepth BinaceSymbol  reciever
func (c *BinaceSymbol) GetDepth(limit int, APIKey string) string {
	l := strconv.Itoa(limit)
	// log.Println("limit " + l)
	url := Endpoint.String()
	url += "/api/v1/depth?symbol="
	url += c.String()
	url += "&limit=" + l
	return url
}

//GetDepthEnpoint func
func GetDepthEnpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var config = configuration.GetConfiguration()

	var sym BinaceSymbol
	var limit int
	sym, limit = sym.getDepthEnpointSlug(w, r)

	// log.Println("the slug is: " + sym)
	var url = sym.GetDepth(limit, config.BinanceAPIKey)
	spaceClient := http.Client{
		Timeout: time.Second * 200,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-MBX-APIKEY", config.BinanceAPIKey)
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// var result = []byte(body)
	// str := fmt.Sprintf("%s", result)
	// log.Println(str)

	var f interface{}
	_ = json.Unmarshal([]byte(body), &f)
	// log.Println(f)

	m := f.(map[string]interface{})
	// log.Println(m)
	// log.Println(m["asks"])

	b, _ := json.Marshal(m)
	helper.SetHeader(w, b)
}

func (c *BinaceSymbol) getDepthEnpointSlug(w http.ResponseWriter, r *http.Request) (a BinaceSymbol, limit int) {
	const aPath = "/v1/binance/get/depth/"
	var slug string
	var symbol string
	var depthLegalRange = []int{5, 10, 20, 50, 100, 500, 1000}
	var checkRange = false
	if strings.HasPrefix(r.URL.Path, aPath) {
		slug = r.URL.Path[len(aPath):]
		s := strings.Split(slug, "/")
		l, err := strconv.Atoi(s[1])
		checkRange, _ = helper.InArray(l, depthLegalRange)
		limit = 5
		if checkRange == true {
			limit = l
		}

		if err != nil {
			log.Fatal(err)
		}
		symbol = s[0]
	}
	// fmt.Println("the slug is: ", symbol)
	// fmt.Println("the slug is: ", limit)
	a = BinaceSymbol(symbol)
	return a, limit
}
