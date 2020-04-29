package war_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/websocket"
	"github.com/hiank/think/pb"
	war_pb "github.com/hiank/thinkend/war/proto"
)

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
