package notice

import (
	"log"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	"github.com/gin-gonic/gin"
)

var HubIns = NewHub()

func init() {
	go HubIns.Run()
}

type NoticeRequest struct {
	Topic string `json:"topic" form:"topic" binding:"required"`
	Uid   string `json:"uid" form:"uid" binding:"required"`
}

type NoticeResponse struct {
	Rid     string `json:"rid"`
	Event   string `json:"event"`
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data"`
}

// serveWs handles websocket requests from the peer.
func ServeWs(ctx *gin.Context) {
	input := new(NoticeRequest)
	err := ctx.ShouldBind(input)
	if err != nil {
		core.WithErr(ctx, errors.ApiErrValidation)
		core.WriteResponse(ctx, err, nil)
		return
	}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:   HubIns,
		Topic: input.Topic,
		Uid:   input.Uid,
		conn:  conn,
		send:  make(chan []byte, 256),
	}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
