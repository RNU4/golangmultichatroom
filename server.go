package main

import (
	static "github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"net/http"
	"strings"
	"fmt"
)

func main() {
	chatmap := make(map[string][]string)

	r := gin.Default()
	m := melody.New()
	r.Use(static.Serve("/", static.LocalFile("./public", true)))

	r.GET("/channel/:name", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./public/channel/channel.html")
	})

	//websocket setup
	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	r.GET("/channel/:name/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	// On new connection
	m.HandleConnect(func(s *melody.Session) {
		if (len(chatmap[s.Request.URL.Path])>0){
			fmt.Printf("%v\n", chatmap[s.Request.URL.Path])
			s.Write([]byte("$load "+strings.Join(chatmap[s.Request.URL.Path],"(%nl)")))
		}
	})

	//Handle message
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		str := strings.ReplaceAll(string(msg), "(%nl)", "(%del)")
		m.BroadcastFilter([]byte(str), func(q *melody.Session) bool {
			return q.Request.URL.Path == s.Request.URL.Path
		})
		chatmap[s.Request.URL.Path] = append(chatmap[s.Request.URL.Path],str)
	})

	r.Run(":5000")
}