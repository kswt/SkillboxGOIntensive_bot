package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"bytes"

	"math/rand"
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

type ActiveUserT struct{
		chat_id int //ID чата
		buddy_id int // ID собеседника. -1 - собеседника нет
		buddy_chat_id int /* ID чата с ним*/
	}

const debug  = true
//const debug = false

const baseTelegramUrl = "https://api.telegram.org"
const getUpdatesUri = "getUpdates"
const telegramToken = "1010543526:AAG8UPCvGaHqxF4KLjZPuMpabOkR2k_3xCs"
const sendMessageUrl = "sendMessage"

const keywordStart = "/start"


func main() {
	delay := 5
if debug{delay=1}
	var offset int = 0
	activeUsers := map[int]*ActiveUserT{}	 
	for {
		// установить счетчик времени, каждые 5 секунд получать обнолвения из telegram и отправлять ответы при совпажении по ключевым словам
		// сделать от 10 до 20 ключевых слов с разными реакциями и проверять совпадения в цикле (strings.Contains)
		// сделать боту возможность переписки пар людей (пересылка сообщений между парами, анонимный чат)
		time.Sleep(time.Duration(delay) * time.Second)

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
					text = "Бот приветствует тебя, " + item.Message.From.FirstName + " " + item.Message.From.LastName
				case strings.Contains(item.Message.Text, "Здравствуй"):
					text = "Здравствуй, " + item.Message.From.FirstName + " " + item.Message.From.LastName
				case strings.Contains(item.Message.Text, "надоел"):
					text = "Подсказка: чтобы закончить общение, воспользуйся командой /end"
				case strings.Contains(item.Message.Text, "скучно"):
					text = "Подсказка: чтобы сменить собеседника, воспользуйся последовательностью команд: /end /begin"
				case item.Message.Text == "пока" || item.Message.Text == "прощай":
					text = "Подсказка: не забудь воспользоваться командой /end чтобы выйти из чата"
				case strings.Contains(item.Message.Text, "дурак"):
					text = "Просим воздержаться от ругани"
				}

				//from_id = item.Message.From.Id
				if activeUsers[item.Message.From.Id] != nil {	// Если пользователь активен
					if activeUsers[item.Message.From.Id].buddy_id!=-1{ // И если у него есть собеседник
						sendMessage(activeUsers[item.Message.From.Id].buddy_chat_id, item.Message.Text) // отправим ему текст сообщения
					}
				}

				switch{ // Служебные команды бота, сообщения с которыми мы не будем пересылать собеседнику
				case item.Message.Text == "/start" :
					text = "Этот бот - отличная возможность найти себе анонимного собеседника.\nНапиши /begin чтобы подобрать себе собеседника и /end для завершения общения\n/count - количество активных пользователей в данный момент"
				case item.Message.Text == "/begin" :
					activeUsers[item.Message.From.Id] = &ActiveUserT{buddy_id:-1,buddy_chat_id:0, chat_id:item.Message.Chat.Id}

					text = "Отлично. Мы добавили тебя в список активных пользователей"

					activeUsers_arr := []int{}
					for key, _ :=range activeUsers{
						activeUsers_arr = append(activeUsers_arr, key)
					}

					rand.Seed(time.Now().UnixNano())
					rand.Shuffle(len(activeUsers_arr),func(i,j int){
						activeUsers_arr[i], activeUsers_arr[j] = activeUsers_arr[j], activeUsers_arr[i]})

					breakflag := false
					for _, v := range(activeUsers_arr){ // ищем собеседника
						if v != item.Message.From.Id {// главное - не выйти на самого себя
							if activeUsers[v].buddy_id == -1 { // если потенциальный собеседник без пары
								activeUsers[v].buddy_id = item.Message.From.Id // прописываем себя его парой
								activeUsers[v].buddy_chat_id = item.Message.Chat.Id // и добавляем ID своего чата

								activeUsers[item.Message.From.Id].buddy_id = v // прописывам собеседника парой себе
								activeUsers[item.Message.From.Id].buddy_chat_id = activeUsers[v].chat_id // и ID его чата

								sendMessage(activeUsers[item.Message.From.Id].buddy_chat_id, "Собеседник найден!") // уведомляем собеседника

								text = text + "\n Теперь у тебя есть собеседник"
								breakflag = true
								break
							}
						}
					}
					if breakflag == false {text = text + ", но пока не нашли тебе собеседника. Мы уведомим тебя сразу, как он появится"}




					/*
					делаем не более, чем len(activeUsers) попыток рандомно выбрать собеседника. Останавливаемся когда:
					-id не равен id самого человека
					-buddy_id != -1

					Как только находим собеседника, прописываемся и у него в структуре. а также пишем ему, что он теперь не один. И пишем нашему изначальному собеседнику. 
					в противном случае пишем, что собеседника пока нет.
					*/
				case item.Message.Text == "/end" :
					if activeUsers[item.Message.From.Id] != nil{
						buddy_id := activeUsers[item.Message.From.Id].buddy_id
						if buddy_id != -1 { //Если у пользователя был собеседник
							sendMessage(activeUsers[item.Message.From.Id].buddy_chat_id, "Твой собеседник покинул чат") // уведомим его о факте отключения
							activeUsers[buddy_id].buddy_id=-1 // удалим сведения о пользователе в структуре собеседника
						}
					}
					delete (activeUsers, item.Message.From.Id) // и удалим структуру самого пользователя 
					text = "Мы удалили тебя из списка активных собеседников"
				case item.Message.Text == "/count" :
					text = strconv.Itoa(len(activeUsers))
				}

				if text != "" {
					sendMessage(item.Message.Chat.Id, text)
				}

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

	sendMessage := SendMessageResponseT{}

	tojson := map[string]string{"chat_id": strconv.Itoa(chatId), "text": text}
	bytesRepresentation, err := json.Marshal(tojson)
	if err != nil {
		return sendMessage, err
	}
if debug {fmt.Println(string(bytesRepresentation))}

	resp, err :=  http.Post(url, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return sendMessage, err
	}
	response := make([]byte, 0)
	for true {
		bs := make([]byte, 1024)
		n, err := resp.Body.Read(bs)
		response = append(response, bs[:n]...)

		if n == 0 || err != nil{
			break
		}
	}

	err = json.Unmarshal(response, &sendMessage)
	if err != nil {
		return sendMessage, err
	}
if debug {fmt.Println(string(response))}
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
