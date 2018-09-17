package model

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	//mysql import
	_ "github.com/go-sql-driver/mysql"
)

//Configuration struct
type Configuration struct {
	Port             string `json:"Port"`
	DatabaseAddress  string `json:"DATABASE_ADDRESS"`
	DatabaseName     string `json:"DATABASE_NAME"`
	DatabaseUser     string `json:"DATABASE_USER"`
	DatabasePassword string `json:"DATABASE_PASSWORD"`
	BinanceAPIKey    string `json:"BINANCE_APIKEY"`
	BinanceSecretKey string `json:"BINANCE_SECRETKEY"`
	BittrexAPIKey    string `json:"BITTREX_APIKEY"`
	BittrexSecretKey string `json:"BITTREX_SECRETKEY"`
	KrakenAPIKey     string `json:"KRAKEN_APIKEY"`
	KrakenPrivateKey string `json:"KRAKEN_PRIVATEKEY"`
	BlockcypherToken string `json:"BLOCKCYPHER_TOKEN"`
	TelegramToken    string `json:"TELEGRAM_TOKEN"`
}

//AofExConfiguration variable of struct Configuration
var AofExConfiguration = Configuration{}

func init() {
	log.Println("init model")
	file, err := os.Open("./config/config.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&AofExConfiguration)
	if err != nil {
		log.Fatal(err)
		return
	}
	AofExConfiguration.initDatabase()
}

//GetConfiguration reciever
func (c *Configuration) GetConfiguration() Configuration {
	return AofExConfiguration
}

//BuildDBConnectionString reciever
func (c *Configuration) BuildDBConnectionString() string {
	return AofExConfiguration.DatabaseUser + ":" + AofExConfiguration.DatabasePassword + "@tcp(" + AofExConfiguration.DatabaseAddress + ")/" + AofExConfiguration.DatabaseName
}

func (c *Configuration) initDatabase() {

	var dbconnection = c.BuildDBConnectionString()
	log.Println("open connection....")
	db, err := sql.Open("mysql", dbconnection)

	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = db.Exec("SELECT * FROM `exchange_list`")

	if err != nil {
		log.Println("Create...exchange_list")
		c.createExchangeListTable()
		c.createTelegramUserTable()
		c.createTelegramChatTable()
	}

}

func (c *Configuration) createTelegramChatTable() {
	var dbconnection = c.BuildDBConnectionString()
	db, err := sql.Open("mysql", dbconnection)
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = db.Exec(`CREATE TABLE chat_rooms (
		cid int(11) NOT NULL,
		chat_id int(9) NOT NULL,
		first_name varchar(255) NOT NULL,
		last_name varchar(255) NOT NULL,
		username varchar(255) NOT NULL,
		type varchar(50) NOT NULL
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;`)

	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = db.Exec(`ALTER TABLE chat_rooms
	ADD PRIMARY KEY (cid);`)
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = db.Exec(`ALTER TABLE chat_rooms
	MODIFY cid int(11) NOT NULL AUTO_INCREMENT;`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("COMMIT;")
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Configuration) createTelegramUserTable() {
	var dbconnection = c.BuildDBConnectionString()
	db, err := sql.Open("mysql", dbconnection)
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = db.Exec(`CREATE TABLE telegram_members (tid int(11) NOT NULL,member_id int(11) NOT NULL,is_bot tinyint(1) NOT NULL DEFAULT '0',
	first_name varchar(255) NOT NULL,last_name varchar(255) NOT NULL,username varchar(50) NOT NULL,
	language_code varchar(5) NOT NULL,created int(10) NOT NULL,
	text varchar(255) NOT NULL) ENGINE=InnoDB DEFAULT CHARSET=utf8;`)
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = db.Exec("ALTER TABLE `telegram_members` ADD PRIMARY KEY (`tid`);")
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = db.Exec("ALTER TABLE `telegram_members` MODIFY `tid` int(11) NOT NULL AUTO_INCREMENT;")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("COMMIT;")
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Configuration) createExchangeListTable() {
	var dbconnection = c.BuildDBConnectionString()
	db, err := sql.Open("mysql", dbconnection)

	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = db.Exec(`CREATE TABLE exchange_list (exid int(11) NOT NULL,ex_name varchar(255) NOT NULL,
	ex_endpoint varchar(255) NOT NULL,
	ex_api_key varchar(255) NOT NULL,
	ex_secret_key varchar(255) NOT NULL) 
	ENGINE=InnoDB DEFAULT CHARSET=utf8;`)
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = db.Exec("ALTER TABLE `exchange_list` ADD PRIMARY KEY (`exid`);")
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = db.Exec("ALTER TABLE `exchange_list` MODIFY `exid` int(11) NOT NULL AUTO_INCREMENT;")
	if err != nil {
		log.Fatal(err)
		return
	}
	insert, err := db.Query(`INSERT INTO exchange_list (
		exid, ex_name, ex_endpoint, ex_api_key, ex_secret_key
		) VALUES (
		NULL, "Binance", "https://api.binance.com", "` + AofExConfiguration.BinanceAPIKey + `", 
		"` + AofExConfiguration.BinanceSecretKey + `");`)

	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("COMMIT;")
	if err != nil {
		log.Fatal(err)
	}
	defer insert.Close()
	defer db.Close()
}

//SetHeader reciever
func (c *Configuration) SetHeader(w http.ResponseWriter, b []byte) {
	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(b)
	// a := string(b)
	// log.Println(a)
	// log.Print(b)
	log.Print("Handle is ok")
}
