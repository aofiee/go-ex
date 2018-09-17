package kraken

/*
NOTE: All API URLs should use the domain api.kraken.com.

Public methods can use either GET or POST.

Private methods must use POST and be set up as follows:

HTTP header:

API-Key = API key
API-Sign = Message signature using HMAC-SHA512 of (URI path + SHA256(nonce + POST data)) and base64 decoded secret API key
POST data:

nonce = always increasing unsigned 64 bit integer
otp = two-factor password (if two-factor enabled, otherwise not required)
*/
import (
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

type krakenURLEndPoint string

const Endpoint krakenURLEndPoint = "https://api.kraken.com/0"

var configuration = model.Configuration{}

//KrakenSymbol string
type KrakenSymbol string

func init() {
	log.Println("init kraken")
	var sym KrakenSymbol = "XBTEUR"
	var url = sym.GetDepth(100)
	log.Println(url)

}
func (c KrakenSymbol) String() string {
	return string(c)
}

func (c krakenURLEndPoint) String() string {
	return string(c)
}

//GetDepth KrakenSymbol  reciever
func (c *KrakenSymbol) GetDepth(limit int) (url string) {
	url = Endpoint.String()
	l := strconv.Itoa(limit)
	url += "/public/Depth?pair=" + c.String() + "&count=" + l
	return url
}
func (c *KrakenSymbol) getDepthEndpointSlug(w http.ResponseWriter, r *http.Request) (slugSymbol string, limit int) {
	var aPath = "/v1/kraken/get/depth/"
	var slug string
	if strings.HasPrefix(r.URL.Path, aPath) {
		slug = r.URL.Path[len(aPath):]
		s := strings.Split(slug, "/")
		slug = s[0]
		limit = 10
		if len(s) != 0 {
			limit, _ = strconv.Atoi(s[1])
		}
	}
	return slug, limit
}

//GetDepthEnpoint func
func GetDepthEnpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var config = configuration.GetConfiguration()
	var sym KrakenSymbol
	var slug, limit = sym.getDepthEndpointSlug(w, r)
	log.Println(slug)
	log.Println(limit)
	sym = KrakenSymbol(slug)
	var url = sym.GetDepth(limit)
	spaceClient := http.Client{
		Timeout: time.Second * 200,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-Key", config.KrakenAPIKey)
	// req.Header.Set("API-Sign", config.KrakenAPIKey)
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
