package chat

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket 설정 (버퍼 크기 등)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 모든 오리진 허용 (개발 편의상)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client : 접속한 유저 정보
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan Message
	username string
}

// ServeWs : HTTP 요청을 웹소켓으로 업그레이드하고 클라이언트를 생성
func ServeWs(hub *Hub, c *gin.Context, username string) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan Message, 256), username: username}
	client.hub.register <- client

	// 읽기/쓰기 펌프 실행 (별도 고루틴)
	go client.writePump()
	go client.readPump()
}

// readPump : 유저가 보낸 메시지를 읽어서 Hub로 전달
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		// 메시지 포맷 정의
		var msg struct {
			Content string `json:"content"`
		}

		// 소켓에서 JSON 읽기
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Hub에 메시지 전달
		message := Message{
			Username: c.username,
			Content:  msg.Content,
		}
		c.hub.broadcast <- message
	}
}

// writePump : Hub에서 온 메시지를 유저에게 전달
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteJSON(message)
		}
	}
}
