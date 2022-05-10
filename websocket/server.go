package websocket

import (
	"fmt"
	"strconv"
	"time"
	"web/model"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

var wsupgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WSHandler(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %s", err.Error())
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("Message received: %s", msg)
		conn.WriteMessage(t, msg)
	}
}

func Logout(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var lastLogout time.Time
	model.DB.QueryRow(
		"SELECT last_logout FROM `govwa`.`activities` WHERE id = ? LIMIT 1",
		id).Scan(&lastLogout)
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %s", err.Error())
		return
	}
	msg := fmt.Sprintf("$('#%d_semaphore').attr('src','/img/0.png');", id)
	msg = msg + fmt.Sprintf("$('#%d_forceLogout').contents().unwrap();", id)
	msg = msg + fmt.Sprintf("$('.bi.bi-x-circle-fill.enabled.%d').attr('class','bi bi-x-circle-fill disabled %d');", id, id)
	msg = msg + fmt.Sprintf("date = new Date(\"%s\");", lastLogout)
	msg = msg + fmt.Sprintf("$('#%d_logout').text(date.toISOString());", id)
	fmt.Println(msg)
	conn.WriteMessage(1, []byte(msg))
}
