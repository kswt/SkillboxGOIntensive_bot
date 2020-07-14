package main


import (
	"fmt"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UpdateT struct {
	Ok bool `json:"ok"`
	Result []UpdateResultT `json:"result"`
}

type UpdateResultT struct {
	UpdateId int `json:"update_id"`
	Message UpdateResultMessageT `json:"message"`
}

type UpdateResultMessageT struct {
	MessageId int `json:"message_id"`
	From UpdateResultFromT `json:"from"`
	Chat UpdateResultChatT `json:"chat"`
	Date int `json:"date"`
	Text string `json:"text"`
}

type UpdateResultFromT struct {
	Id int `json:"id"`
	IsBot bool `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Username string `json:"username"`
	Language string `json:"language_code"`
}

type UpdateResultChatT struct {
	Id int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Username string `json:"username"`
	Type string `json:"type"`
}

type SendMessageResponseT struct {
	Ok bool `json:"ok"`
	Result UpdateResultMessageT `json:"result"`
}

const debug  = true

const baseTelegramUrl = "https://api.telegram.org"
const getUpdatesUri = "getUpdates"
const telegramToken = "1010543526:AAG8UPCvGaHqxF4KLjZPuMpabOkR2k_3xCs"
const sendMessageUrl = "sendMessage"

const keywordStart = "/start"

func main() {
	var offset int = 0
	for {
		// установить счетчик времени, каждые 5 секунд получать обнолвения из telegram и отправлять ответы при совпажении по ключевым словам
		// сделать от 10 до 20 ключевых слов с разными реакциями и проверять совпадения в цикле (strings.Contains)
		// сделать боту возможность переписки пар людей (пересылка сообщений между парами, анонимный чат)
		time.Sleep(5 * time.Second)

		//offset++
if debug {fmt.Println(offset)}
		update, err := getUpdates(offset)

		if err != nil {
			fmt.Println(err.Error())

			return
		}


		for _, item := range(update.Result) {
			var text string
			if item.Message.From.IsBot == false {
				switch{
				case item.Message.Text == "Привет" :
					text = "Привет, " + item.Message.From.FirstName + " " + item.Message.From.LastName
				case strings.Contains(item.Message.Text, "Здравствуй"):
					text = "Здравствуй, " + item.Message.From.FirstName + " " + item.Message.From.LastName
				case item.Message.Text == "/start" :
					text = "Этот бот - отличная возможность найти себе анонимного собеседника.\nНапиши /begin чтобы начать"
				}



				sendMessage(item.Message.Chat.Id, text)
			}
			offset = item.UpdateId + 1
		}
	}
}

func getUpdates(offset int) (UpdateT, error) {
	url := baseTelegramUrl + "/bot" + telegramToken + "/" + getUpdatesUri + "?offset=" + strconv.Itoa(offset)
	response := getResponse(url)

if debug{fmt.Println(string(response))}

	update := UpdateT{}
	err := json.Unmarshal(response, &update)
	if err != nil {
		return update, err
	}

	return update, nil
}

func sendMessage(chatId int, text string) (SendMessageResponseT, error) {
	url := baseTelegramUrl + "/bot" + telegramToken + "/" + sendMessageUrl
	url = url + "?chat_id=" + strconv.Itoa(chatId) + "&text=" + text
	response := getResponse(url)

	sendMessage := SendMessageResponseT{}
	err := json.Unmarshal(response, &sendMessage)
	if err != nil {
		return sendMessage, err
	}

	return sendMessage, nil
}

func getResponse(url string) []byte {
	response := make([]byte, 0)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)

		return response
	}

	defer resp.Body.Close()

	for true {
		bs := make([]byte, 1024)
		n, err := resp.Body.Read(bs)
		response = append(response, bs[:n]...)

		if n == 0 || err != nil{
			break
		}
	}

	return response
}
