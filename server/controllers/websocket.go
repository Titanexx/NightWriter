package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"nightwriter/helpers"
	"nightwriter/models"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type PartHub struct {
	currIndex uint
	clients   map[uint]*websocket.Conn
	Updates   []string `json:"updates"`
	lock      sync.Mutex
}

func newPartHub() *PartHub {
	var g PartHub
	g.clients = make(map[uint]*websocket.Conn)
	return &g
}

func (partHub *PartHub) connectFirstConn(conn *websocket.Conn) bool {
	conn.WriteJSON(1)
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Error in waiting state", err)
		return false
	}
	partHub.Updates = append(partHub.Updates, string(msg))
	conn.WriteJSON(0)
	return true
}
func (partHub *PartHub) connectOtherConn(conn *websocket.Conn) bool {
	conn.WriteJSON(2)
	conn.WriteJSON(partHub.Updates)
	conn.WriteJSON(0)
	return true
}

func (partHub *PartHub) AddConn(conn *websocket.Conn) (idConn uint, err error) {
	partHub.lock.Lock()
	idConn = partHub.currIndex
	partHub.currIndex += 1
	if len(partHub.clients) == 0 {
		if partHub.connectFirstConn(conn) {
			partHub.clients[idConn] = conn
		} else {
			conn.Close()
			err = errors.New("can't add the 1st conn")
		}
		partHub.lock.Unlock()
	} else {
		partHub.clients[idConn] = conn
		partHub.lock.Unlock()
		partHub.connectOtherConn(conn)
	}
	return idConn, err
}

func (partHub *PartHub) BroadcastMsg(msg string, currentConn *websocket.Conn) {
	partHub.Updates = append(partHub.Updates, string(msg))
	for _, clientConn := range partHub.clients {
		if clientConn != currentConn {
			clientConn.WriteJSON(msg)
		}
	}
}

func (partHub *PartHub) RemoveConn(idConn uint) {
	partHub.lock.Lock()
	partHub.clients[idConn].Close()
	delete(partHub.clients, idConn)
	if len(partHub.clients) == 0 {
		partHub.Updates = nil
	}
	partHub.lock.Unlock()
}

var wsUpgrader = websocket.Upgrader{} // use default options
var hub = map[uint]map[uint]*PartHub{}

func WsHandlerPartByDocIDAndPartID(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	docId := helpers.GetUintIDInUrl(c, "docID")
	partId := helpers.GetUintIDInUrl(c, "partID")

	if models.IsDocAndPartExistByID(user, docId, partId) {
		if _, ok := hub[docId]; !ok {
			hub[docId] = map[uint]*PartHub{}
		}

		if _, ok := hub[docId][partId]; !ok {
			hub[docId][partId] = newPartHub()
		}
		w := c.Writer
		r := c.Request
		wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Failed to set websocket upgrade:", err)
			return
		}

		idConn, err := hub[docId][partId].AddConn(conn)
		if err != nil {
			return
		}
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			hub[docId][partId].BroadcastMsg(string(msg), conn)
		}
		hub[docId][partId].RemoveConn(idConn)
	} else {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
}
