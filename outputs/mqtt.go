package outputs

import (
	"bitbucket.org/r_rudi/gostat/record"
	"encoding/json"
	"fmt"
	. "github.com/huin/mqtt"
	"io"
	"net"
	"strings"
)

type MQTT struct {
	conn        net.Conn
	isConnected bool
}

func (m MQTT) Header(r record.Record) {}

func (m MQTT) connect(conn io.Writer) error {
	msg := &Connect{
		Header: Header{
			DupFlag:  false,
			QosLevel: QosAtMostOnce,
			Retain:   false,
		},
		ProtocolName:    "MQIsdp",
		ProtocolVersion: 3,
		WillRetain:      false,
		WillFlag:        false,
		CleanSession:    false,
		WillQos:         QosAtMostOnce,
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
	topic_prefix := conf["mqtt"]["topic_prefix"]
	hostname := conf["root"]["hostname"]
	server := conf["mqtt"]["server"]
	port := conf["mqtt"]["port"]

	fmt.Println("hoge")

	if m.isConnected == false {
		conn, err := net.Dial("tcp", server+":"+port)
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

		m.conn = conn
		m.isConnected = true
	}

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
		topic := []string{topic_prefix, hostname, r.Tag}
		msg := &Publish{
			Header: Header{
				DupFlag:  false,
				QosLevel: QosAtMostOnce,
				Retain:   false,
			},
			TopicName: strings.Join(topic, "/"),
			MessageId: uint16(id),
			Payload:   BytesPayload(data),
		}
		if err := msg.Encode(m.conn); err != nil {
			fmt.Println(err)
			fmt.Println("send Publish failed")
			continue
		}
		/*
			puback, err := DecodeOneMessage(m.conn, nil)
			if err != nil {
				fmt.Println("recv PUB ACK failed")
				return err
			}
			if puback != nil {
				fmt.Println(puback)
			} // FIXME
			//		fmt.Println(puback)
		*/

	}
	/*
		if err := m.disconnect(m.conn); err != nil {
			fmt.Println("send DISCONNECT failed")
			return err
		}
	*/

	return nil
}
