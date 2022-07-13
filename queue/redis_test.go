package queue

import (
	"encoding/json"
	"testing"
)

func TestNewRedis(t *testing.T) {
	var err error
	Q, err = NewRedis()
	if err != nil {
		t.Error(err)
	}
}

type Test_Struct struct {
	MsgId      string `json:"msg_id"`
	MsgContent string `json:"msg_content"`
}

func (i Test_Struct) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func TestRedis_Push(t *testing.T) {
	var test_msg Test_Struct
	err := Q.Push(test_msg)
	if err != nil {
		t.Error(err)
	}
}

func TestRedis_Size(t *testing.T) {
	t.Log(Q.Size())
}

func TestRedis_Pop(t *testing.T) {
	c, err := Q.Pop()
	if err != nil {
		t.Error(err)
	}
	t.Log(c)
}
