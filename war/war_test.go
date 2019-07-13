package war_test

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/websocket"
	"github.com/hiank/think/pb"
	war_pb "github.com/hiank/thinkend/war/proto"
	"testing"
	"net/http"
	"net/url"
)

// func TestProtobufPerformance(t *testing.T) {


// 	msg := &war_pb.TestMessage{
// 	}
// 	p := &war_pb.TestParent{
// 		Msg : msg,
// 	}
// 	anyMsg, err := pb.Message2Any(p)
// 	if err != nil {
// 		t.Log(anyMsg, err)
// 	}
// 	anyBuf, err := pb.AnyEncode(anyMsg)
// 	if err != nil {
// 		t.Log(anyBuf, err)
// 	}


// 	curTime := time.Now().Nanosecond()

// }

func TestTeamEncode(t *testing.T) {

	// roles := make([]*proto1.Role, 2, 2)
	// for i:=0; i<2; i++ {

	// 	roles[i] = &proto1.Role {
	// 		Uid : uint64(1001),
	// 		ModelId : int32(1),
	// 		ModelLv: int32(1),
	// 		Cup : int32(1000),
	// 		Uname : *proto.String("hello"),
	// 	}
	// }

	// roles[1] = &proto1.Role{Uid : uint64(1002)}

	msg := &war_pb.TestMessage{
		Test : &war_pb.Test{Id: uint64(1002000000000000)},
	}
	p := &war_pb.TestParent{
		Msg : msg,
	}
	buf, err := proto.Marshal(p)
	t.Log(buf)
	t.Log(p)
	anyMsg, err := ptypes.MarshalAny(p)
	t.Log(anyMsg, err)
}

func dail(t *testing.T) (*websocket.Conn, *http.Response, error) {

	addr := "192.168.137.222:30250"
	t.Logf("address : %s\n", addr)

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	return websocket.DefaultDialer.Dial(u.String(), http.Header{"token": {"1022"}})
}


func TestConnect(t *testing.T) {

	conn, _, err := dail(t)
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	//NOTE: do war request
	anyMsg, err := ptypes.MarshalAny(&war_pb.S_War_Want{Type: war_pb.War_Type_TEAMWINNER})
	if err != nil {
		t.Errorf("marshal any error : %v\n", err)
		return
	}

	buf, err := pb.AnyEncode(anyMsg)
	if err != nil {
		t.Errorf("encode error : %v\n", err)
		return
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, buf); err != nil {

		t.Errorf("send want message error : %v\n", err)
		return
	}
	
	for {

		_, message, err := conn.ReadMessage()
		if err != nil {
			t.Logf("read error : %v\n", err)
			return
		}

		t.Logf("received: %s\n", message)
	}
}