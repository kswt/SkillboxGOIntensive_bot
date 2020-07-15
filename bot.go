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
	}

//const debug  = true
const debug = false

const baseTelegramUrl = "https://api.telegram.org"
const getUpdatesUri = "getUpdates"
const telegramToken = "1010543526:AAG8UPCvGaHqxF4KLjZPuMpabOkR2k_3xCs"
const sendMessageUrl = "sendMessage"

const keywordStart = "/start"


func main() {
	delay := 5
if debug{delay=1}
	var offset int = 0
	active_chats := 0
	activeUsers := map[int]*ActiveUserT{} 
	for {
		// установить счетчик времени, каждые 5 секунд получать обнолвения из telegram и отправлять ответы при совпажении по ключевым словам
		// сделать от 10 до 20 ключевых слов с разными реакциями и проверять совпадения в цикле (strings.Contains)
		// сделать боту возможность переписки пар людей (пересылка сообщений между парами, анонимный чат)
		time.Sleep(time.Duration(delay) * time.Second)

if debug {fmt.Println(offset)}
		update, err := getUpdates(offset)

		if err != nil {
			fmt.Println(err.Error())

			return
		}


		for _, item := range(update.Result) {
			offset = item.UpdateId + 1
			from_id := item.Message.From.Id
			var text string
// Чат для живых людей. Сообщения от ботов игнорируются
			if item.Message.From.IsBot == false {
				send_to_buddy := true
				switch{ // Служебные команды бота, сообщения с которыми мы не будем пересылать собеседнику
				case item.Message.Text == "/start" :
					text = "Этот бот - отличная возможность найти себе анонимного собеседника.\nНапиши /begin чтобы подобрать себе собеседника и /end для завершения общения\n/users - количество активных пользователей в данный момент\n/chats - количество активных чатов"
					send_to_buddy = false
				case item.Message.Text == "/begin" :
					activeUsers[from_id] = &ActiveUserT{buddy_id:-1, chat_id:item.Message.Chat.Id}

					text = "Отлично. Мы добавили тебя в список активных пользователей"


					activeUsers_arr := []int{}
					for key, _ :=range activeUsers{
						activeUsers_arr = append(activeUsers_arr, key)
					}
if debug {fmt.Println(activeUsers_arr)}
					rand.Seed(time.Now().UnixNano())
					rand.Shuffle(len(activeUsers_arr),func(i,j int){
						activeUsers_arr[i], activeUsers_arr[j] = activeUsers_arr[j], activeUsers_arr[i]})
// При большом количестве пользователей блок кода выше станет узким местом. В целях оптимизации времени можно разбить пользователей на 2 отдельных объекта map: пользователи с собеседником и без, но до дедлайна я не успею.


					breakflag := false
					for _, v := range(activeUsers_arr){ // ищем собеседника
						if v != from_id {// главное - не выйти на самого себя
							if activeUsers[v].buddy_id == -1 { // если потенциальный собеседник без пары
								activeUsers[v].buddy_id = from_id // прописываем у собеседника себя
								activeUsers[from_id].buddy_id = v // прописывам у себя собеседника

								sendMessage(activeUsers[activeUsers[from_id].buddy_id].chat_id, "Собеседник найден!") // уведомляем собеседника

								text = text + "\nТеперь у тебя есть собеседник"
								active_chats++
								breakflag = true
								break
							}
						}
					}
					if breakflag == false {text = text + ", но пока не нашли тебе собеседника. Мы уведомим тебя сразу, как он появится"}
					send_to_buddy = false
				case item.Message.Text == "/end" :
					if activeUsers[from_id] != nil{
						buddy_id := activeUsers[from_id].buddy_id
						if buddy_id != -1 { //Если у пользователя был собеседник
							sendMessage(activeUsers[activeUsers[from_id].buddy_id].chat_id, "Твой собеседник покинул чат. Воспользуйся командой /begin,	чтобы найти нового") // уведомим его о факте отключения
							activeUsers[buddy_id].buddy_id=-1 // удалим сведения о пользователе в структуре собеседника
							active_chats--
						}
					}
					delete (activeUsers, from_id) // и удалим структуру самого пользователя 
					text = "Мы удалили тебя из списка активных собеседников"

					send_to_buddy = false
				case item.Message.Text == "/users" :
					text = strconv.Itoa(len(activeUsers))
					send_to_buddy = false
				case item.Message.Text == "/chats" :
					text = strconv.Itoa(active_chats)
					send_to_buddy = false
				}

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

				if  activeUsers[from_id] != nil && send_to_buddy {	// Если пользователь активен и текущее сообщение подлежит пересылке
					if activeUsers[from_id].buddy_id!=-1{ // И если у пользователя есть собеседник
						sendMessage(activeUsers[activeUsers[from_id].buddy_id].chat_id, item.Message.Text) // отправим ему текст сообщения
					}
				}


				if text != "" {
					sendMessage(item.Message.Chat.Id, text)
				}

			}
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
