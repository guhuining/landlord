package service

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"landlord/hub"
	"landlord/log"
	"landlord/model"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// CreateRoom 创建房间
func CreateRoom(c *gin.Context) {
	session := sessions.Default(c)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.MyLog.Printf("create room: %v", err)
		return
	}
	if session.Get("user_id") == nil {
		conn.WriteJSON(hub.Response{hub.ERROR, "用户未登录", nil})
		conn.Close()
		return
	}
	if hub.EmptyRooms.Number == 0 {
		conn.WriteJSON(hub.Response{hub.ERROR, "已达最大房间数", nil})
		conn.Close()
		return
	}
	room, err := hub.EmptyRooms.PopFront()
	if err != nil { // 此时出错代表在此语句执行和上一次检查之间最后一间房间被分配给其他人了
		conn.WriteJSON(hub.Response{hub.ERROR, "已达最大房间数", nil})
		conn.Close()
		return
	}
	conn.WriteJSON(hub.Response{hub.OK, "成功", nil})
	// 获取房间密码信息
	var request model.CreateRoomRequest
	err = conn.ReadJSON(&request)
	if err != nil {
		log.MyLog.Printf("read password: %v", err)
	}
	// 根据是否有密码，将房间放入不同的房间列表
	if request.Password == "" {
		hub.RoomWithoutPassword.PushBack(room)
	} else {
		room.Password = request.Password
		hub.RoomWithPassword.PushBack(room)
	}
	// 将房间状态改为等待中
	room.State = hub.WAITING
	// 将此用户加入进房间
	player := &hub.Player{
		UserID: session.Get("user_id").(int),
		Nickname: session.Get("nickname").(string),
		Conn: conn,
	}
	room.NewPlayer <- player

	for {
		req := hub.Request{}
		err = conn.ReadJSON(&req)
		if err != nil {  // 玩家异常退出
			room.PlayerExit <- player
			break
		}
		switch req.Data {
		default:

		}
	}
}
