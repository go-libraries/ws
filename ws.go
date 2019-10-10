package ws

import (
    "net/http"
    "sync"
    "strings"
    "github.com/gorilla/websocket"
)


type Protocol struct {
    // Register requests from the clients.
    register chan *Connection

    // UnRegister requests from clients.
    unRegister chan *Connection

    //all connections, It's mapping O(1)
    ConnectionsMap  map[string]*sync.Map

    //use rw mutex
    rwm *sync.RWMutex

    //send and received message
    Msg chan IMessage
    num uint64 //count

    //log
    PLog ILog

    //upgrade check logic
    CheckOrigin func(r *http.Request) bool
    //upgrade error handler
    UpErrorHandler func(res http.ResponseWriter, r *http.Request, status int, reason error)

    //config
    Config IConfig
}

//router upgrade /ws/{room}
//router http   /room/xxx/cancel
//router http   /room/xxx/register
func (w *Protocol) ServeHTTP(rw http.ResponseWriter, r *http.Request)  {

    res := strings.Split(r.URL.Path, "/")
    l := len(res)
    if l < 2 {
        Response(rw, DefaultResponse{
            Code:"500",
            Msg:"protocol not support",
        })
        return
    }

    room := ""
    if res[0] == "ws" {
        if l == 2 {
            room = res[1]
        }
        if !w.roomExists(room) {
            Response(rw, DefaultResponse{
                Code:"500",
                Msg:"room id not exists",
            })
            return
        }
        w.registerWs(rw, r, room)
        Response(rw, SuccessResponse)
        return
    }

    if res[0] == "room" {
        room = res[1]
        if room == "" {
            Response(rw, DefaultResponse{
                Code:"500",
                Msg:"room id not exists",
            })
            return
        }
        isRegister := true
        if l > 2 {
            v := res[2]
            if v == "cancel" {
                isRegister = false
            }
        }
        if isRegister {
            w.registerRoom(room)
        }else{
            w.unRegisterRoom(room)
        }
        Response(rw, SuccessResponse)
        return
    }

    Response(rw, DefaultResponse{
        Code:"500",
        Msg:"not support this router",
    })
    return
}

func (w *Protocol) unRegisterRoom(room string) {
    w.rwm.Lock()
    defer w.rwm.Unlock()

    if ok:=w.roomExists(room); ok {
        m := w.ConnectionsMap[room]
        m.Range(func(key, value interface{}) bool {
            v1 := value.(*Connection)
            v1.Close()
            m.Delete(v1.GetRoom())
            return true
        })
    }
}
//lock room
func (w *Protocol)  registerRoom(room string) {
    w.rwm.Lock()
    defer w.rwm.Unlock()

    if ok:=w.roomExists(room); !ok {
        w.ConnectionsMap[room] = new(sync.Map)
    }
}

func (w *Protocol) roomExists(room string) bool{
    _,ok := w.ConnectionsMap[room]

    return ok
}

func (w *Protocol) registerWs(rw http.ResponseWriter, r *http.Request, room string)  {

    uniqueKey := r.Header.Get("Sec-WebSocket-Key")
    if uniqueKey == "" {
        //todo:
    }

    up := &websocket.Upgrader{
        ReadBufferSize:w.Config.GetReadBufferSize(),
        WriteBufferSize:w.Config.GetWriteBufferSize(),
    }

    if w.UpErrorHandler != nil{
        up.Error = w.UpErrorHandler
    }else{
        up.Error = w.upErrorHandler
    }

    if w.CheckOrigin != nil {
        up.CheckOrigin = w.CheckOrigin
    }else{
        up.CheckOrigin = w.checkAllowOrigin
    }

    con, err := up.Upgrade(rw, r, rw.Header())
    if err != nil {
        w.PLog.Println("handler err with message" + err.Error())
        rw.Write([]byte("fail to upGrader"))
        rw.WriteHeader(500)
        return
        //panic("handler err with message" + err.Error())
    }

    var wsConn  = &Connection {
        UniqueKey:uniqueKey,
        Conn:con,
        send: make(chan []byte, 256),
        room:room,
        CLog:w.PLog,

        maxMessageSize:w.Config.GetMaxMessageSize(),
        pongWait:w.Config.GetPongWaitTime(),
        writeWait:w.Config.GetWriteWaitTime(),
    }

    Wsp.Online(wsConn)

    go wsConn.read()
    go wsConn.write()
}

func (w *Protocol) Send(msg IMessage) {
    w.Msg <- msg
}
//send message with msg
func (w *Protocol) send(msg IMessage) {
    all := w.All(msg.GetRoom())
    bs,length,err := msg.GetMessage()

    if err != nil {
        w.PLog.Println(err)
        return
    }

    if length < 1 || len(bs) == 0 {
        w.PLog.Println("message is nil")
        return
    }

    for _,v := range all{
        v.Send(bs)
    }
}

//conn connection,write lock
func (w *Protocol) Online(conn *Connection) {
    w.register <- conn
}

//read lock
func (w *Protocol) All(room string) []*Connection {
    seen := make([]*Connection, 0)

    if w.roomExists(room) {
        m := w.ConnectionsMap[room]
        m.Range(func(ki, vi interface{}) bool {
            v := vi.(*Connection)
            seen = append(seen, v)
            return true
        })
    }

    return seen
}

//write lock
func (w *Protocol) OffLine(conn *Connection) {
    w.unRegister <- conn
}

//run
// 1 catch client in/out
// 2 catch message
func (w *Protocol) Run()  {
    for {
        select {
        case client := <-w.register:
            room := client.GetRoom()
            w.num = w.num + 1
            w.ConnectionsMap[room].Store(client.GetUniqueKey(), client)
            //w.Connections[client.GetUniqueKey()] = client
        case client := <-w.unRegister:
            w.num = w.num - 1
            room := client.GetRoom()
            w.ConnectionsMap[room].Delete(client.GetUniqueKey())
            //delete(w.Connections, client.GetUniqueKey())
        case msg := <-w.Msg:
            w.send(msg)
        }
    }
}

func (w *Protocol) GetNumber() uint64 {
    return w.num
}

func (w *Protocol) upErrorHandler(res http.ResponseWriter, req *http.Request, status int, reason error) {
    w.PLog.Println("handler err with message" + reason.Error())
    res.Write([]byte("fail to upGrader"))
    res.WriteHeader(status)
}


func (w *Protocol) checkAllowOrigin(r *http.Request) bool {
        return true
}