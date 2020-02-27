package ws

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

type Connection struct {
	*websocket.Conn

	UniqueKey string
	send      chan []byte
	room      string
	close     bool
	CLog      ILog
	Config    IConfig

	maxMessageSize int64
	pongWait       time.Duration
	writeWait      time.Duration
}

func (wsc *Connection) GetUniqueKey() string {
	return wsc.UniqueKey
}

func (wsc *Connection) GetRoom() string {
	return wsc.room
}

func (wsc *Connection) SetRoom(room string) {
	wsc.room = room
}

func (wsc *Connection) Send(b []byte) {
	wsc.send <- b
}

//read
func (wsc *Connection) read() {
	defer func() {
		wsc.Close()
	}()
	wsc.SetReadLimit(wsc.maxMessageSize)
	wsc.SetReadDeadline(time.Now().Add(wsc.pongWait))
	wsc.SetPongHandler(func(string) error { wsc.SetReadDeadline(time.Now().Add(wsc.pongWait)); return nil })
	for {
		msgType, message, err := wsc.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				wsc.CLog.Println(fmt.Sprintf("error: %v", err))
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		wsc.CLog.Println(fmt.Sprintf("receive message: %s", string(message[:])))

		switch msgType {
		case websocket.CloseMessage:
			wsc.CLog.Println("client send close:" + wsc.GetUniqueKey())
			return
		case websocket.TextMessage:
			//todo:预留
		case websocket.BinaryMessage:
			//todo:预留
			//default:
			//todo:忽略,select 移除default减少性能损失
		}
	}
}

//close and offline
func (wsc *Connection) Close() {
	//panic: close of closed channel
	if !wsc.close {
		wsc.close = true
		Wsp.OffLine(wsc)
		wsc.Conn.Close()
		close(wsc.send)
	}
}

func (wsc *Connection) write() {
	ticker := time.NewTicker(getPingPeriod(wsc.pongWait))
	defer func() {
		ticker.Stop()
		wsc.Close()
	}()
	for {
		select {
		case message, ok := <-wsc.send:
			wsc.SetWriteDeadline(time.Now().Add(wsc.writeWait))
			//fmt.Println(message)
			if !ok {
				// The hub closed the channel.
				wsc.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := wsc.WriteBytes(message)
			if err != nil {
				return
			}
			//心跳
		case <-ticker.C:
			wsc.SetWriteDeadline(time.Now().Add(wsc.writeWait))
			if err := wsc.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (wsc *Connection) WriteString(message string) error {
	return wsc.WriteBytes([]byte(message))
}

func (wsc *Connection) WriteBytes(message []byte) error {
	wsc.SetWriteDeadline(time.Now().Add(wsc.writeWait))

	w, err := wsc.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	_, err = w.Write(message)
	if err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}
