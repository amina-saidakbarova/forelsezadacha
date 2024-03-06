package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type SignUpStruct struct {
	Name          string
	TelegramLogin string
	Password      string
}

var SignUpSlice = []SignUpStruct{} // ? empty
func main() {
	r := gin.Default()

	r.Use(Cors)
	r.POST("/signup", SignUp)
	go Recovery()

	r.Run(":3434")

}
func Recovery() {
	ReadUser()
	botresult, err := tgbotapi.NewBotAPI("7009989148:AAHcaQ-gKvvFQoFOg8dlFyCG4GHomtOKedY")

	if err != nil {
		fmt.Printf("err: %v\n", err)

	}
	updates := tgbotapi.NewUpdate(0)
	allupdates, _ := botresult.GetUpdatesChan(updates)

	for update := range allupdates {
		if update.Message.IsCommand() {
			if update.Message.Command() == "reset" {
				for _, Item := range SignUpSlice {
					if Item.TelegramLogin == update.Message.Chat.UserName {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Enter new password")
						botresult.Send(msg)

					}

				}

			}
		} else {
			yast := false
			for i, item := range SignUpSlice {
					if item.TelegramLogin == update.Message.Chat.UserName{
						yast = true
						SignUpSlice[i].Password = update.Message.Text
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Enter new password")
						botresult.Send(msg)
					}
				}
				if !yast{
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Enter new password")
					botresult.Send(msg)

				}

			}
			WriteUser()

		}

	}


func SignUp(c *gin.Context) {
	var SignUpTemp SignUpStruct
	c.ShouldBindJSON(&SignUpTemp)

	if SignUpTemp.Name == "" || SignUpTemp.Password == "" || SignUpTemp.TelegramLogin == "" {
		c.JSON(404, "Empty field")
	} else {
		ReadUser()
		SignUpSlice = append(SignUpSlice, SignUpTemp)
		WriteUser()
	}
}
func WriteUser() {
	marsheledData, _ := json.Marshal(SignUpSlice)
	ioutil.WriteFile("app.json", marsheledData, 0644)
}
func ReadUser() {
	readedByte, _ := ioutil.ReadFile("app.json")
	json.Unmarshal(readedByte, &SignUpSlice)

}

func Cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://192.168.43.246:5500")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	}

	c.Next()
}
