package ws

//default IMessage
type RoomMsg struct {
	msg  []byte
	room string
}

func (rm *RoomMsg) SetMsg(m []byte) {
	rm.msg = m
}

func (rm *RoomMsg) SetRoom(r string) {
	rm.room = r
}

func (rm *RoomMsg) GetMsg() []byte {
	return rm.msg
}

func (rm *RoomMsg) GetRoom() string {
	return rm.room
}

func (rm *RoomMsg) GetMessage() (bs []byte, length int64, err error) {
	return rm.msg, int64(len(rm.msg)), nil
}
