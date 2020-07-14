package main

import "fmt" // Можно ссылку на гитхаб?

/*
const FirstConst = "My first constant"

var FirstStringVar string
var FirstIntVar int
var FirstInt8Var int8
var FirstInt64Var int64
var FirstFloatVar float64
var FirstBoolVar bool
*/

func main() {
	subscribers := 0
	fmt.Println("Введите количество подписчиков")
	fmt.Scan(&subscribers)

	groups := 3
	fmt.Println("Введите количество групп")
	fmt.Scan(&groups)

	fmt.Println("Подписчиков: ", subscribers)
	fmt.Println("Групп", groups)

	IsGroupChatAllowed := (subscribers % groups) == 0

	if IsGroupChatAllowed {
		fmt.Println("Alowed")
	} else {
		fmt.Println("NotAllowed")
	}
}
