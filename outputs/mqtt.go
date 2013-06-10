package outputs

import (
	"bitbucket.org/r_rudi/gostat/record"
	"encoding/json"
	"fmt"
	. "github.com/huin/mqtt"
	"io"
	"net"
)

type MQTT struct{}

func (m MQTT) Header(r record.Record) {}

func (m MQTT) connect(conn io.Writer) error {
	msg := &Connect{
		Header: Header{
			DupFlag:  false,
			QosLevel: QosAtLeastOnce,
			Retain:   false,
		},
		ProtocolName:    "MQIsdp",
		ProtocolVersion: 3,
		WillRetain:      false,
		WillFlag:        false,
		CleanSession:    false,
		WillQos:         QosAtLeastOnce,
		KeepAliveTimer:  10,
		ClientId:        "gostat",
		WillTopic:       "",
		WillMessage:     "",
		UsernameFlag:    false,
		PasswordFlag:    false,
		Username:        "",
		Password:        "",
	}

	return msg.Encode(conn)
}

func (m MQTT) disconnect(conn io.Writer) error {
	msg := &Disconnect{
		Header: Header{
			DupFlag:  false,
			QosLevel: QosAtLeastOnce,
			Retain:   false,
		},
	}

	return msg.Encode(conn)
}

func (m MQTT) Emit(rs []record.Record, conf map[string]map[string]string) error {
	topic := conf["mqtt"]["topic"]
	server := conf["mqtt"]["server"]

	conn, err := net.Dial("tcp", server)
	if err != nil {
		fmt.Println("Could not connect server")
		return err
	}

	if err := m.connect(conn); err != nil {
		fmt.Println("send CONNECT failed")
		return err
	}
	msg, err := DecodeOneMessage(conn, nil)
	if err != nil {
		fmt.Println("recv CONNECT ACK failed")
		return err
	}
	if msg != nil {
	} // FIXME
	//	fmt.Println(m.getRet(msg))

	for id, r := range rs {
		value := make(map[string]string, 0)
		value["time"] = r.Time.UTC().String()
		value["tag"] = r.Tag

		for k, v := range r.Value {
			value[k] = v
		}

		data, err := json.Marshal(value)
		if err != nil {
			continue
		}
		msg := &Publish{
			Header: Header{
				DupFlag:  false,
				QosLevel: QosAtLeastOnce,
				Retain:   false,
			},
			TopicName: topic + "/" + r.Tag,
			MessageId: uint16(id),
			Payload:   BytesPayload(data),
		}

		if err := msg.Encode(conn); err != nil {
			fmt.Println("send Publish failed")
			continue
		}
		puback, err := DecodeOneMessage(conn, nil)
		if err != nil {
			fmt.Println("recv PUB ACK failed")
			return err
		}
		if puback != nil {
		} // FIXME
		//		fmt.Println(puback)

	}

	if err := m.disconnect(conn); err != nil {
		fmt.Println("send DISCONNECT failed")
		return err
	}

	return nil
}
