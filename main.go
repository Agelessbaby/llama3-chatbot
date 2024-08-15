package main

import (
	"github.com/Agelessbaby/llama3-chatbot/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	gin.SetMode(gin.DebugMode)
	engine.LoadHTMLFiles("view/chatview.html")
	engine.GET("/", func(context *gin.Context) {
		context.HTML(200, "chatview.html", nil)
	})
	engine.POST("/chat", handler.Chatwithllama)
	engine.Run("127.0.0.1:8080")
}
