package user

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	fundation "gobang"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan *Message
	Register   chan *Client
	Unregister chan *Client
}

// Client is a websocket client
type Client struct {
	Recipient string
	ID        string
	Socket    *websocket.Conn
	Send      chan []byte
	Roomid    string
}
type Message struct {
	Roomid    string `json:"roomid"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
	Types     string `json:"types"` //types=1表示下棋，2表示说话
}

var Manager = ClientManager{
	Broadcast:  make(chan *Message),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[*Client]bool),
}

// Start 项目运行前, 协程开启start -> go Manager.Start()
func (manager *ClientManager) Start() {
	for {
		log.Println("<---管道通信--->")

		select {
		case conn := <-Manager.Register:
			log.Printf("新用户加入:%v", conn.ID)
			Manager.Clients[conn] = true
			jsonMessage, _ := json.Marshal(&Message{Roomid: "0", Content: "Successful connection to socket service"})
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
			jsonMessage, _ := json.Marshal(&Message{Roomid: message.Roomid, Recipient: message.Recipient, Sender: message.Sender, Content: message.Content})
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
//judge判读五子棋胜负
func judge(ma [16][16]int) int {
	var dir = [8][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}, {1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for i := 1; i <= 15; i++ {
		for j := 1; j <= 15; j++ {
			for t := 0; t <= 7; t++ {
				if (ma[i][j] != 1 && ma[i][j] != 2) || j+4*dir[t][1] > 15 || j+4*dir[t][1] < 1 || i+4*dir[t][0] > 15 || i+4*dir[t][0] < 1 {
					continue
				}
				x := dir[t][1]
				y := dir[t][0]
				if ma[i][j] == ma[i+y][j+x] && ma[i][j] == ma[i+2*y][j+2*x] && ma[i][j] == ma[i+3*y][j+3*x] && ma[i][j] == ma[i+4*y][j+4*x] {
					return ma[i][j]
				}
			}
		}
	}
	return 0
}

// Send is to send ws message to ws client
func (manager *ClientManager) Send(message []byte, ignore *Client) {
	for conn := range manager.Clients {
		if conn != ignore { //向除了自己的socket 用户发送
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
			log.Println(err, "121行")
			break
		}
		me := string(message)
		m := strings.Split(me, "+")
		if len(m)<2 {
			m=append(m,"2")
		}
		if m[0] == "2" {
			if m[1] == "认输" {
				m[1] = c.ID + "输了"
				m[0] = "2"
				conn := fundation.Pool.Get()
				conn.Do("ZRANGE", c.Roomid+"+chess", "0", "-1", "WITHSCORES")
				conn.Do("spop", c.Roomid)
				conn.Do("spop", c.Roomid)
				goto MESS
			}
		}
		if m[0] == "1" {
			step := strings.Split(m[1], ",")
			var a [2]int
			for i, i2 := range step {
				a[i], err = strconv.Atoi(i2)
			}
			if a[1] > 15 || a[1] < 1 || a[0] > 15 || a[0] < 1 {
				m[1] = "请不要乱下棋"
				m[0] = "2"
				goto MESS
			}
			conn := fundation.Pool.Get()
			if t, _ := conn.Do("zrank", c.Roomid+"+chess", m[1]); t != nil {
				m[1] = "此步已下"
				m[0] = "2"
				goto MESS
			}
			name, err := redis.Strings(conn.Do("smembers", c.Roomid))
			if err != nil {
				Manager.Unregister <- c
				c.Socket.Close()
				log.Println(err, "147行")
				break
			}
			isMem := 0
			for i, i2 := range name {
				if i2 == c.ID {
					conn.Do("zadd", c.Roomid+"+chess", strconv.Itoa(i+1), m[1])
					isMem = 1
				}
			}
			if isMem == 0 {
				m[1] = "请不要下棋,你是观战"
				m[0] = "2"
				goto MESS
			}
			allstep, err := redis.StringMap(conn.Do("ZRANGEBYSCORE", c.Roomid+"+chess", "-inf", "+inf", "withscores"))
			if err != nil {
				Manager.Unregister <- c
				c.Socket.Close()
				fmt.Println(err)
				break
			}
			var step1num, step2num = 0, 0
			var thechess [16][16]int
			for i, i2 := range allstep {
				xy := strings.Split(i, ",")
				x, _ := strconv.Atoi(xy[0])
				y, _ := strconv.Atoi(xy[1])
				thechess[x][y], _ = strconv.Atoi(i2)
				if i2 == "1" {
					step1num++
				} else {
					step2num++
				}
			}
			if step1num-step2num > 1 || step1num-step2num < 0 {

				conn.Do("zrem", c.Roomid+"+chess", m[1])
				m[1] = "不是你的回合"
				m[0] = "2"
				goto MESS
			}
			win := judge(thechess)
			if win != 0 {
				m[1] = name[win-1] + "获胜"
				m[0] = "2"
				conn.Do("ZRANGE", c.Roomid+"+chess", "0", "-1", "WITHSCORES")
				conn.Do("spop", c.Roomid)
				conn.Do("spop", c.Roomid)
				goto MESS
			}
			mjson, _ := json.Marshal(allstep)
			m[1] = string(mjson)
			fmt.Println(allstep, "\n", mjson, "\n", m[1])
			m[0] = "2"
		}
	MESS:
		txt := Message{
			Roomid:    c.Roomid,
			Sender:    c.ID,
			Recipient: c.Recipient,
			Content:   m[1],
			Types:     m[0],
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
			json.Unmarshal(message, &mm)
			if mm.Roomid == "0" {
				if mm.Recipient == "-1" {
					c.Socket.WriteMessage(websocket.TextMessage, message)
				} else if mm.Recipient == c.ID {
					c.Socket.WriteMessage(websocket.TextMessage, message)
				}

			} else if c.Roomid == mm.Roomid {
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
	//userID, err := c.Cookie("username")
	//var roomid string
	//if roomid, err = c.Cookie("roomid"); err != nil {
	//	roomid = c.DefaultQuery("roomid", "0")
	//}
	userID:=c.Query("username")
	roomid:=c.Query("roomid")
	recip := c.DefaultQuery("recipient", "-1")
	client := &Client{
		Recipient: recip,
		ID:        userID,
		Socket:    conn,
		Send:      make(chan []byte),
		Roomid:    roomid,
	}
	Manager.Register <- client
	go client.Read()
	go client.Write()
}
