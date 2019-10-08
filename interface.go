package ws

import "time"

type IMessage interface {
    GetMessage() (bs []byte,length int64, err error)
    GetRoom() string
}

type ILog interface {
    Println(v... interface{})
}

type IConfig interface {
    GetWriteWaitTime() (d time.Duration)
    GetReadWaitTime() (d time.Duration)
    GetConnectionWaitTime() (d time.Duration)
    GetMaxMessageSize() int64
    GetPongWaitTime() (d time.Duration)
    GetReadBufferSize() int
    GetWriteBufferSize() int
}