package chat

import "log"

// Message : 채팅 메시지 구조체
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

type Hub struct {
	// 등록된 클라이언트들 (채팅방 접속자)
	clients map[*Client]bool

	// 메시지 전송 통로 (누가 말을 하면 이리로 들어옴)
	broadcast chan Message

	// 클라이언트 등록 요청 채널
	register chan *Client

	// 클라이언트 해제 요청 채널
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run : 채팅방을 계속 돌리면서 이벤트를 감지함
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client connected. Total: %d", len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client disconnected. Total: %d", len(h.clients))
			}

		case message := <-h.broadcast:
			// 메시지가 오면 모든 사람에게 뿌림
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
