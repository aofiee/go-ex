package telegram

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aofiee/go-ex/app/helper"
	model "github.com/aofiee/go-ex/app/model"
)

//telegramDatabase type struct
type telegramDatabase struct {
	TID          int    `json:"tid" sql:"AUTO_INCREMENT"`
	MemberID     int    `json:"member_id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	Created      int    `json:"created"`
	Text         string `json:"text"`
}

//ChatRoomDatabase struct
type ChatRoomDatabase struct {
	ID        int    `json:"cid" sql:"AUTO_INCREMENT"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

//GetUpdatesResponse struct
type GetUpdatesResponse struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
			} `json:"from"`
			Chat struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date     int    `json:"date"`
			Text     string `json:"text"`
			Entities []struct {
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				Type   string `json:"type"`
			} `json:"entities"`
		} `json:"message"`
	} `json:"result"`
}

//RecieverTelegramResponse struct
type RecieverTelegramResponse struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

type telegrameURLEndPoint string

var configuration = model.Configuration{}

//EndPoint telegrameURLEndPoint
var EndPoint telegrameURLEndPoint = "https://api.telegram.org/bot"

var getUpdatesResponse = GetUpdatesResponse{}

func init() {
	log.Println("init telegram")
}

func (c telegrameURLEndPoint) String() string {
	return string(c)
}

//GetUpdateURL func
func (c *GetUpdatesResponse) GetUpdateURL() string {
	url := EndPoint.String()
	var config = configuration.GetConfiguration()
	url += config.TelegramToken + "/getUpdates"
	return url
}

//RecieverTelegramProcess func
func RecieverTelegramProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	var result = []byte(body)
	str := fmt.Sprintf("%s", result)
	log.Println(str)
	var chatResponse = RecieverTelegramResponse{}
	err := json.Unmarshal(body, &chatResponse)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(chatResponse.Message.Text)
	chatResponse.echoMessage()
}
func (r *RecieverTelegramResponse) messageFromToString() string {
	return string(r.Message.From.ID)
}
func (r *RecieverTelegramResponse) echoMessage() {
	url := EndPoint.String()
	var config = configuration.GetConfiguration()
	var fromID string
	fromID = strconv.Itoa(r.Message.From.ID)

	url += config.TelegramToken + "/sendMessage?chat_id=" + fromID + "&text=" + r.Message.Text

	log.Println(url)
	spaceClient := http.Client{
		Timeout: time.Second * 200,
	}
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, errPost := spaceClient.Do(req)
	if errPost != nil {
		log.Fatal(errPost)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	var f interface{}
	_ = json.Unmarshal([]byte(body), &f)
	m := f.(map[string]interface{})
	log.Println("------------------------")
	log.Println(m)
	log.Println("------------------------")

}

//GetUpdateTelegramProcess func
func GetUpdateTelegramProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var getUpdatesResponse = GetUpdatesResponse{}
	var url = getUpdatesResponse.GetUpdateURL()
	log.Println(url)
	spaceClient := http.Client{
		Timeout: time.Second * 200,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	// var f interface{}
	// _ = json.Unmarshal([]byte(body), &f)
	// m := f.(map[string]interface{})
	// log.Println("------------------------")
	// log.Println(m)
	// log.Println("------------------------")

	//json output
	_ = json.Unmarshal([]byte(body), &getUpdatesResponse)

	if getUpdatesResponse.Ok == true {
		log.Println("i am ok")
		if len(getUpdatesResponse.Result) != 0 {
			log.Println(len(getUpdatesResponse.Result))
			getUpdatesResponse.insertNewMember(getUpdatesResponse)
		}

	}
	b, _ := json.Marshal(getUpdatesResponse)
	helper.SetHeader(w, b)
	//json output

}

func (c *GetUpdatesResponse) insertNewMember(member GetUpdatesResponse) bool {
	log.Println("------------------------")
	log.Println(member)
	log.Println("------------------------")
	var dbconnection = configuration.BuildDBConnectionString()
	log.Println("open connection....")
	db, err := sql.Open("mysql", dbconnection)
	if err != nil {
		log.Fatal(err)
		return false
	}
	for _, m := range member.Result {
		var telegramDB = telegramDatabase{
			MemberID:     m.Message.From.ID,
			IsBot:        m.Message.From.IsBot,
			FirstName:    m.Message.From.FirstName,
			LastName:     m.Message.From.LastName,
			Username:     m.Message.From.Username,
			LanguageCode: m.Message.From.LanguageCode,
			Created:      m.Message.Date,
			Text:         m.Message.Text,
		}
		resCountMember, errCheckIsHaveMember := db.Prepare("SELECT `member_id` FROM `telegram_members` WHERE `member_id` = ?")
		if errCheckIsHaveMember != nil {
			log.Fatal(errCheckIsHaveMember)
		}
		var memberID string
		errQueryRow := resCountMember.QueryRow(telegramDB.MemberID).Scan(&memberID)
		if errQueryRow != nil {
			res, errIn := db.Exec(`INSERT INTO telegram_members 
			( member_id, is_bot, first_name, last_name, username, language_code, created, text) 
			VALUES 
			(?,?,?,?,?,?,?,?)`,
				telegramDB.MemberID,
				telegramDB.IsBot,
				telegramDB.FirstName,
				telegramDB.LastName,
				telegramDB.Username,
				telegramDB.LanguageCode,
				telegramDB.Created,
				telegramDB.Text)
			if errIn != nil {
				log.Fatal(errIn)
			}
			log.Println(res)
			var chatID = ChatRoomDatabase{
				ID:        m.Message.Chat.ID,
				FirstName: m.Message.Chat.FirstName,
				LastName:  m.Message.Chat.LastName,
				Username:  m.Message.Chat.Username,
				Type:      m.Message.Chat.Type,
			}
			resCountChatID, errCheckIsHaveChatID := db.Prepare(`SELECT chat_id FROM chat_rooms WHERE chat_id = ?`)
			if errCheckIsHaveChatID != nil {
				log.Fatal(errCheckIsHaveChatID)
			}
			var resChatID string
			errQueryRow = resCountChatID.QueryRow(chatID.ID).Scan(&resChatID)
			if errQueryRow != nil {
				res, err := db.Exec(`INSERT INTO chat_rooms (chat_id, first_name, last_name, username, type) VALUES (?,?,?,?,?)`,
					chatID.ID,
					chatID.FirstName,
					chatID.LastName,
					chatID.Username,
					chatID.Type)
				if err != nil {
					log.Fatal(err)
				}
				log.Println(res)
			}

		}
	}

	defer db.Close()
	return true
}
