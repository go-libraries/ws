package ws

import (
    "sync"
)
var Wsp *Protocol

func init()  {
    Wsp = &Protocol{

    }
    Wsp.ConnectionsMap = make(map[string]*sync.Map)

    Wsp.Msg = make(chan IMessage)
    Wsp.register = make(chan *Connection)
    Wsp.unRegister = make(chan *Connection)
    go Wsp.Run()
}