package ws

import (
    "sync"
)
var Wsp *Protocol
var SuccessResponse  = DefaultResponse{
    Code:"200",
    Msg:"操作成功",
}
var ErrorResponse = DefaultResponse{
    Code:"200",
    Msg:"操作失败",
}
func init()  {
    Wsp = &Protocol{

    }
    Wsp.ConnectionsMap = make(map[string]*sync.Map)

    Wsp.Msg = make(chan IMessage)
    Wsp.register = make(chan *Connection)
    Wsp.unRegister = make(chan *Connection)
    go Wsp.Run()
}