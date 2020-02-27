package ws

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type DefaultResponse struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (res DefaultResponse) Error() string {
	return res.Msg
}

func (res DefaultResponse) ToJsonBytes() []byte {
	b, err := json.Marshal(res)
	if err != nil {
		return []byte{}
	}

	return b
}

func (res DefaultResponse) GetData() interface{} {
	return res.Data
}

func (res DefaultResponse) ToJsonString() string {
	b := res.ToJsonBytes()
	return string(b[:])
}

func (res DefaultResponse) GetCode() int {
	code, err := strconv.Atoi(res.Code)

	if err != nil {
		code = 200
	}

	return code
}

func Response(res http.ResponseWriter, response IResponse) {
	code := response.GetCode()
	res.WriteHeader(code)
	res.Write(response.ToJsonBytes())
}
