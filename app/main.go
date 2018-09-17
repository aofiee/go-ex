package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/aofiee/go-ex/app/helper"

	binance "github.com/aofiee/go-ex/app/exchange/binance"
	bittrex "github.com/aofiee/go-ex/app/exchange/bittrex"
	kraken "github.com/aofiee/go-ex/app/exchange/kraken"
	model "github.com/aofiee/go-ex/app/model"
	telegram "github.com/aofiee/go-ex/app/telegram"
)

var addr = flag.String("addr", ":1234", "http service address")
var configuration = model.Configuration{}

func main() {
	var config = configuration.GetConfiguration()
	fmt.Print(`Server listening on 0.0.0.0 port ` + config.Port + `
------------------------------------------------------------------------------
GO-EXCHANGGE MODULE Ver 0.1
------------------------------------------------------------------------------
Press CTRL+C to exit
`)
	log.Println(telegram.EndPoint)
	flag.Parse()
	http.HandleFunc("/", addHandler)
	http.HandleFunc("/v1/binance/get/depth/", binance.GetDepthEnpoint)
	http.HandleFunc("/v1/bittrex/get/ticker/", bittrex.GetTickerEnpoint)
	http.HandleFunc("/v1/kraken/get/depth/", kraken.GetDepthEnpoint)
	http.HandleFunc("/v1/telegram/getUpdate/", telegram.GetUpdateTelegramProcess)
	http.HandleFunc("/v1/telegram/reciever/", telegram.RecieverTelegramProcess)
	//https://api.telegram.org/bot/setWebhook?url=https://e2a99a10.ngrok.io/v1/telegram/reciever/
	if err := http.ListenAndServe(config.Port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	output := "GO-EXCHANGE MODULE Ver 0.1"
	outputByteSlice := []byte(output)
	helper.SetHeader(w, outputByteSlice)
}
