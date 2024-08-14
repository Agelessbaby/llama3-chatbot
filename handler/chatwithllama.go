package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

var mylog *log.Logger

func init() {
	file, err := os.OpenFile("AIChatBot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		os.Exit(1)
	}
	mylog = log.New(file, "[AIChatBot]", log.LstdFlags|log.Lshortfile)
}

type Message struct {
	Message string `json:"message" binding:"gt=0"`
}

type Response struct {
	Reply string `json:"reply"`
}

func Chatwithllama(ctx *gin.Context) {
	var mes Message
	err := ctx.ShouldBindJSON(&mes)
	if err != nil {
		ctx.String(http.StatusBadRequest, "invalid message")
		return
	}
	mylog.Printf("ip:%s：send a request: %s\n", ctx.ClientIP(), mes.Message)

	data := map[string]interface{}{
		"model":  "llama3:8b",
		"prompt": mes.Message,
		"stream": false,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		mylog.Printf("Error marshalling JSON:", err)
		ctx.String(http.StatusInternalServerError, "internal server error")
		return
	}
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "internal server error")
		mylog.Printf("Error sending request: %s\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "internal server error")
		mylog.Println("Error reading llama response body: %s\n", err)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		ctx.String(http.StatusInternalServerError, "internal server error")
		mylog.Println("Error unmarshalling json")
		return
	}
	bs := result["response"]
	ans, _ := bs.(string)
	mylog.Printf("llama3 give back a response:%s", result["response"], ans)
	res := Response{Reply: ans}
	ctx.JSON(http.StatusOK, res)
}
