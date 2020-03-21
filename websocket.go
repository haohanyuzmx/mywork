package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan *Message
	Register   chan *Client
	Unregister chan *Client
}

// Client is a websocket client
type Client struct {
	Recipient       string
	ID     int
	Socket *websocket.Conn
	Send   chan []byte
	Roomid int
}

// Message is return msg
type Message struct {
	Roomid    string `json:"roomid,omitempty"`
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

// Manager define a ws server manager
var Manager = ClientManager{
	Broadcast:  make(chan *Message),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[*Client]bool),
}

// Start is  项目运行前, 协程开启start -> go Manager.Start()
func (manager *ClientManager) Start() {
	for {
		log.Println("<---管道通信--->")

		select {
		case conn := <-Manager.Register:
			log.Printf("新用户加入:%v", conn.ID)
			Manager.Clients[conn] = true
			jsonMessage, _ := json.Marshal(&Message{Roomid:"0",Content: "Successful connection to socket service"})
			Manager.Send(jsonMessage, conn)
		case conn := <-Manager.Unregister:
			log.Printf("用户离开:%v", conn.ID)
			if _, ok := Manager.Clients[conn]; ok {
				close(conn.Send)
				delete(Manager.Clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "A socket has disconnected"})
				Manager.Send(jsonMessage, conn)
			}

		case message := <-Manager.Broadcast:
			jsonMessage, _ := json.Marshal(&Message{Roomid:message.Roomid,Recipient:message.Recipient,Sender:message.Sender,Content: message.Content})
			for conn := range Manager.Clients {
				select {
				case conn.Send <- jsonMessage:
				default:
					close(conn.Send)
					delete(Manager.Clients, conn)
				}
			}
		}
	}
}

// Send is to send ws message to ws client
func (manager *ClientManager) Send(message []byte, ignore *Client) {

	for conn := range manager.Clients {
		if conn == ignore { //向除了自己的socket 用户发送
			conn.Send <- message
		}
	}
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()

	for {

		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		txt:=Message{
			Roomid:   strconv.Itoa(c.Roomid),
			Sender:    strconv.Itoa(c.ID),
			Recipient: c.Recipient,
			Content:   string(message),
		}
		Manager.Broadcast <- &txt
	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Printf("发送到到客户端的信息:%s", string(message))
			var mm Message
			json.Unmarshal(message,&mm)
			if mm.Roomid=="0" {
				if mm.Recipient=="-1" {
					c.Socket.WriteMessage(websocket.TextMessage, message)
				}else if mm.Recipient==strconv.Itoa(c.ID) {
					c.Socket.WriteMessage(websocket.TextMessage, message)
				}

			}else if strconv.Itoa(c.Roomid)==mm.Roomid{
				c.Socket.WriteMessage(websocket.TextMessage, message)
			}

		}
	}

}

//TestHandler socket 连接 中间件 作用:升级协议,用户验证,自定义信息等
func TestHandler(c *gin.Context) {
	conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	//可以添加用户信息验证
	userID := c.Query("nameid")
	roomid:=c.DefaultQuery("roomid","0")
	recip:=c.DefaultQuery("recipient","-1")
	intname,_:=strconv.Atoi(userID)
	introom,_:=strconv.Atoi(roomid)
	client := &Client{
		Recipient:recip,
		ID:     intname,
		Socket: conn,
		Send:   make(chan []byte),
		Roomid: introom,
	}
	Manager.Register <- client
	go client.Read()
	go client.Write()
}
func main()  {
	go Manager.Start()
	r:=gin.Default()
	r.GET("",TestHandler)
	r.Run()
}

