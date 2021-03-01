/*
 * Copyright (c) 2013 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package mqtt

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

func Test_Start(t *testing.T) {
	ops := NewClientOptions().SetClientID("Start").AddBroker(FVTTCP)
	c := NewClient(ops)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	c.Disconnect(250)
}

/* uncomment this if you have connection policy disallowing FailClientID
 func Test_InvalidConnRc(t *testing.T) {
	 ops := NewClientOptions().SetClientID("FailClientID").
		 AddBroker("tcp://" + FVT_IP + ":17003").
		 SetStore(NewFileStore("/tmp/fvt/InvalidConnRc"))

	 c := NewClient(ops)
	 _, err := c.Connect()
	 if err != ErrNotAuthorized {
		 t.Fatalf("Did not receive error as expected, got %v", err)
	 }
	 c.Disconnect(250)
 }
*/

// Helper function for Test_Start_Ssl
// func NewTLSConfig() *tls.Config {
// 	certpool := x509.NewCertPool()
// 	pemCerts, err := ioutil.ReadFile("samples/samplecerts/CAfile.pem")
// 	if err == nil {
// 		certpool.AppendCertsFromPEM(pemCerts)
// 	}

// 	cert, err := tls.LoadX509KeyPair("samples/samplecerts/client-crt.pem", "samples/samplecerts/client-key.pem")
// 	if err != nil {
// 		panic(err)
// 	}

// 	return &tls.Config{
// 		RootCAs:            certpool,
// 		ClientAuth:         tls.NoClientCert,
// 		ClientCAs:          nil,
// 		InsecureSkipVerify: true,
// 		Certificates:       []tls.Certificate{cert},
// 	}
// }

/* uncomment this if you have ssl setup
 func Test_Start_Ssl(t *testing.T) {
	 tlsconfig := NewTlsConfig()
	 ops := NewClientOptions().SetClientID("StartSsl").
		 AddBroker(FVT_SSL).
		 SetStore(NewFileStore("/tmp/fvt/Start_Ssl")).
		 SetTlsConfig(tlsconfig)

	 c := NewClient(ops)

	 _, err := c.Connect()
	 if err != nil {
		 t.Fatalf("Error on Client.Connect(): %v", err)
	 }

	 c.Disconnect(250)
 }
*/

func Test_Publish_1(t *testing.T) {
	ops := NewClientOptions()
	ops.AddBroker(FVTTCP)
	ops.SetClientID("Publish_1")

	c := NewClient(ops)
	token := c.Connect()
	if token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	c.Publish("test/Publish", 0, false, "Publish qo0")

	c.Disconnect(250)
}

func Test_Publish_2(t *testing.T) {
	ops := NewClientOptions()
	ops.AddBroker(FVTTCP)
	ops.SetClientID("Publish_2")

	c := NewClient(ops)
	token := c.Connect()
	if token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	c.Publish("/test/Publish", 0, false, "Publish1 qos0")
	c.Publish("/test/Publish", 0, false, "Publish2 qos0")

	c.Disconnect(250)
}

func Test_Publish_3(t *testing.T) {
	ops := NewClientOptions()
	ops.AddBroker(FVTTCP)
	ops.SetClientID("Publish_3")

	c := NewClient(ops)
	token := c.Connect()
	if token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	c.Publish("/test/Publish", 0, false, "Publish1 qos0")
	c.Publish("/test/Publish", 1, false, "Publish2 qos1")
	c.Publish("/test/Publish", 2, false, "Publish2 qos2")

	c.Disconnect(250)
}

func Test_Publish_BytesBuffer(t *testing.T) {
	ops := NewClientOptions()
	ops.AddBroker(FVTTCP)
	ops.SetClientID("Publish_BytesBuffer")

	c := NewClient(ops)
	token := c.Connect()
	if token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	payload := bytes.NewBufferString("Publish qos0")

	c.Publish("test/Publish", 0, false, payload)

	c.Disconnect(250)
}

func Test_Subscribe(t *testing.T) {
	pops := NewClientOptions()
	pops.AddBroker(FVTTCP)
	pops.SetClientID("Subscribe_tx")
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.AddBroker(FVTTCP)
	sops.SetClientID("Subscribe_rx")
	var f MessageHandler = func(client Client, msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	}
	sops.SetDefaultPublishHandler(f)
	s := NewClient(sops)

	sToken := s.Connect()
	if sToken.Wait() && sToken.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", sToken.Error())
	}

	s.Subscribe("/test/sub", 0, nil)

	pToken := p.Connect()
	if pToken.Wait() && pToken.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", pToken.Error())
	}

	p.Publish("/test/sub", 0, false, "Publish qos0")

	p.Disconnect(250)
	s.Disconnect(250)
}

func Test_Will(t *testing.T) {
	willmsgc := make(chan string, 1)

	sops := NewClientOptions().AddBroker(FVTTCP)
	sops.SetClientID("will-giver")
	sops.SetWill("/wills", "good-byte!", 0, false)
	sops.SetConnectionLostHandler(func(client Client, err error) {
		fmt.Println("OnConnectionLost!")
	})
	sops.SetAutoReconnect(false)
	c := NewClient(sops).(*client)

	wops := NewClientOptions()
	wops.AddBroker(FVTTCP)
	wops.SetClientID("will-subscriber")
	wops.SetDefaultPublishHandler(func(client Client, msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		willmsgc <- string(msg.Payload())
	})
	wops.SetAutoReconnect(false)
	wsub := NewClient(wops)

	if wToken := wsub.Connect(); wToken.Wait() && wToken.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", wToken.Error())
	}

	if wsubToken := wsub.Subscribe("/wills", 0, nil); wsubToken.Wait() && wsubToken.Error() != nil {
		t.Fatalf("Error on Client.Subscribe(): %v", wsubToken.Error())
	}

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	c.forceDisconnect()

	if <-willmsgc != "good-byte!" {
		t.Fatalf("will message did not have correct payload")
	}

	wsub.Disconnect(250)
}

func Test_CleanSession(t *testing.T) {
	clsnc := make(chan string, 1)

	sops := NewClientOptions().AddBroker(FVTTCP)
	sops.SetClientID("clsn-sender")
	sops.SetConnectionLostHandler(func(client Client, err error) {
		fmt.Println("OnConnectionLost!")
	})
	sops.SetAutoReconnect(false)
	c := NewClient(sops).(*client)

	wops := NewClientOptions()
	wops.AddBroker(FVTTCP)
	wops.SetClientID("clsn-tester")
	wops.SetCleanSession(false)
	wops.SetDefaultPublishHandler(func(client Client, msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		clsnc <- string(msg.Payload())
	})
	wops.SetAutoReconnect(false)
	wsub := NewClient(wops)

	if wToken := wsub.Connect(); wToken.Wait() && wToken.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", wToken.Error())
	}

	if wsubToken := wsub.Subscribe("clean", 1, nil); wsubToken.Wait() && wsubToken.Error() != nil {
		t.Fatalf("Error on Client.Subscribe(): %v", wsubToken.Error())
	}

	wsub.Disconnect(250)
	time.Sleep(2 * time.Second)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	if pToken := c.Publish("clean", 1, false, "clean!"); pToken.Wait() && pToken.Error() != nil {
		t.Fatalf("Error on Client.Publish(): %v", pToken.Error())
	}

	c.Disconnect(250)

	wsub = NewClient(wops)
	if wToken := wsub.Connect(); wToken.Wait() && wToken.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", wToken.Error())
	}

	select {
	case msg := <-clsnc:
		if msg != "clean!" {
			t.Fatalf("will message did not have correct payload")
		}
	case <-time.NewTicker(5 * time.Second).C:
		t.Fatalf("failed to receive publish")
	}

	wsub.Disconnect(250)

	wops.SetCleanSession(true)

	wsub = NewClient(wops)
	if wToken := wsub.Connect(); wToken.Wait() && wToken.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", wToken.Error())
	}

	wsub.Disconnect(250)
}

func Test_Binary_Will(t *testing.T) {
	willmsgc := make(chan []byte, 1)
	will := []byte{
		0xDE,
		0xAD,
		0xBE,
		0xEF,
	}

	sops := NewClientOptions().AddBroker(FVTTCP)
	sops.SetClientID("will-giver")
	sops.SetBinaryWill("/wills", will, 0, false)
	sops.SetConnectionLostHandler(func(client Client, err error) {
	})
	sops.SetAutoReconnect(false)
	c := NewClient(sops).(*client)

	wops := NewClientOptions().AddBroker(FVTTCP)
	wops.SetClientID("will-subscriber")
	wops.SetDefaultPublishHandler(func(client Client, msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %v\n", msg.Payload())
		willmsgc <- msg.Payload()
	})
	wops.SetAutoReconnect(false)
	wsub := NewClient(wops)

	if wToken := wsub.Connect(); wToken.Wait() && wToken.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", wToken.Error())
	}

	if wsubToken := wsub.Subscribe("/wills", 0, nil); wsubToken.Wait() && wsubToken.Error() != nil {
		t.Fatalf("Error on Client.Subscribe() %v", wsubToken.Error())
	}

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	c.forceDisconnect()

	if !bytes.Equal(<-willmsgc, will) {
		t.Fatalf("will message did not have correct payload")
	}

	wsub.Disconnect(250)
}

/**
"[...] a publisher is responsible for determining the maximum QoS a
message can be delivered at, but a subscriber is able to downgrade
the QoS to one more suitable for its usage.
The QoS of a message is never upgraded."
**/

/***********************************
 * Tests to cover the 9 QoS combos *
 ***********************************/

func wait(c chan bool) {
	fmt.Println("choke is waiting")
	<-c
}

// Pub 0, Sub 0

func Test_p0s0(t *testing.T) {
	topic := "/test/p0s0"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.AddBroker(FVTTCP)
	pops.SetClientID("p0s0-pub")
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.AddBroker(FVTTCP)
	sops.SetClientID("p0s0-sub")
	var f MessageHandler = func(client Client, msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)

	if token := s.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	if token := s.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Subscribe(): %v", token.Error())
	}

	if token := p.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}
	p.Publish(topic, 0, false, "p0s0 payload 1")
	p.Publish(topic, 0, false, "p0s0 payload 2")

	wait(choke)
	wait(choke)

	p.Publish(topic, 0, false, "p0s0 payload 3")
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)
}

// Pub 0, Sub 1

func Test_p0s1(t *testing.T) {
	topic := "/test/p0s1"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.AddBroker(FVTTCP)
	pops.SetClientID("p0s1-pub")
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.AddBroker(FVTTCP)
	sops.SetClientID("p0s1-sub")
	var f MessageHandler = func(client Client, msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	if token := s.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	if token := s.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Subscribe(): %v", token.Error())
	}

	if token := p.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}
	p.Publish(topic, 0, false, "p0s1 payload 1")
	p.Publish(topic, 0, false, "p0s1 payload 2")

	wait(choke)
	wait(choke)

	p.Publish(topic, 0, false, "p0s1 payload 3")
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)
}

// Pub 0, Sub 2

func Test_p0s2(t *testing.T) {
	topic := "/test/p0s2"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.AddBroker(FVTTCP)
	pops.SetClientID("p0s2-pub")
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.AddBroker(FVTTCP)
	sops.SetClientID("p0s2-sub")
	var f MessageHandler = func(client Client, msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	if token := s.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	if token := s.Subscribe(topic, 2, nil); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Subscribe(): %v", token.Error())
	}

	if token := p.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}
	p.Publish(topic, 0, false, "p0s2 payload 1")
	p.Publish(topic, 0, false, "p0s2 payload 2")

	wait(choke)
	wait(choke)

	p.Publish(topic, 0, false, "p0s2 payload 3")

	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)
}

// Pub 1, Sub 0

func Test_p1s0(t *testing.T) {
	topic := "/test/p1s0"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.AddBroker(FVTTCP)
	pops.SetClientID("p1s0-pub")
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.AddBroker(FVTTCP)
	sops.SetClientID("p1s0-sub")
	var f MessageHandler = func(client Client, msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	if token := s.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	if token := s.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Subscribe(): %v", token.Error())
	}

	if token := p.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}
	p.Publish(topic, 1, false, "p1s0 payload 1")
	p.Publish(topic, 1, false, "p1s0 payload 2")

	wait(choke)
	wait(choke)

	p.Publish(topic, 1, false, "p1s0 payload 3")

	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)
}

// Pub 1, Sub 1

func Test_p1s1(t *testing.T) {
	topic := "/test/p1s1"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.AddBroker(FVTTCP)
	pops.SetClientID("p1s1-pub")
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.AddBroker(FVTTCP)
	sops.SetClientID("p1s1-sub")
	var f MessageHandler = func(client Client, msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	if token := s.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	if token := s.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Subscribe(): %v", token.Error())
	}

	if token := p.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}
	p.Publish(topic, 1, false, "p1s1 payload 1")
	p.Publish(topic, 1, false, "p1s1 payload 2")

	wait(choke)
	wait(choke)

	p.Publish(topic, 1, false, "p1s1 payload 3")
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)
}

// Pub 1, Sub 2

func Test_p1s2(t *testing.T) {
	topic := "/test/p1s2"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.AddBroker(FVTTCP)
	pops.SetClientID("p1s2-pub")
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.AddBroker(FVTTCP)
	sops.SetClientID("p1s2-sub")
	var f MessageHandler = func(client Client, msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	if token := s.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	if token := s.Subscribe(topic, 2, nil); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Subscribe(): %v", token.Error())
	}

	if token := p.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}
	p.Publish(topic, 1, false, "p1s2 payload 1")
	p.Publish(topic, 1, false, "p1s2 payload 2")

	wait(choke)
	wait(choke)

	p.Publish(topic, 1, false, "p1s2 payload 3")

	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)
}

// Pub 2, Sub 0

func Test_p2s0(t *testing.T) {
	topic := "/test/p2s0"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.AddBroker(FVTTCP)
	pops.SetClientID("p2s0-pub")
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.AddBroker(FVTTCP)
	sops.SetClientID("p2s0-sub")
	var f MessageHandler = func(client Client, msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	if token := s.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	if token := s.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Subscribe(): %v", token.Error())
	}

	if token := p.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}
	p.Publish(topic, 2, false, "p2s0 payload 1")
	p.Publish(topic, 2, false, "p2s0 payload 2")
	wait(choke)
	wait(choke)

	p.Publish(topic, 2, false, "p2s0 payload 3")
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)
}

// Pub 2, Sub 1

func Test_p2s1(t *testing.T) {
	topic := "/test/p2s1"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.AddBroker(FVTTCP)
	pops.SetClientID("p2s1-pub")
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.AddBroker(FVTTCP)
	sops.SetClientID("p2s1-sub")
	var f MessageHandler = func(client Client, msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	if token := s.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	if token := s.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Subscribe(): %v", token.Error())
	}

	if token := p.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}
	p.Publish(topic, 2, false, "p2s1 payload 1")
	p.Publish(topic, 2, false, "p2s1 payload 2")

	wait(choke)
	wait(choke)

	p.Publish(topic, 2, false, "p2s1 payload 3")

	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)
}

// Pub 2, Sub 2

func Test_p2s2(t *testing.T) {
	topic := "/test/p2s2"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.AddBroker(FVTTCP)
	pops.SetClientID("p2s2-pub")
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.AddBroker(FVTTCP)
	sops.SetClientID("p2s2-sub")
	var f MessageHandler = func(client Client, msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	if token := s.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	if token := s.Subscribe(topic, 2, nil); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Subscribe(): %v", token.Error())
	}

	if token := p.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}
	p.Publish(topic, 2, false, "p2s2 payload 1")
	p.Publish(topic, 2, false, "p2s2 payload 2")

	wait(choke)
	wait(choke)

	p.Publish(topic, 2, false, "p2s2 payload 3")

	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)
}

func Test_PublishMessage(t *testing.T) {
	topic := "/test/pubmsg"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.AddBroker(FVTTCP)
	pops.SetClientID("pubmsg-pub")
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.AddBroker(FVTTCP)
	sops.SetClientID("pubmsg-sub")
	var f MessageHandler = func(client Client, msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		if string(msg.Payload()) != "pubmsg payload" {
			fmt.Println("Message payload incorrect", msg.Payload(), len("pubmsg payload"))
			t.Fatalf("Message payload incorrect")
		}
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	if token := s.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	if token := s.Subscribe(topic, 2, nil); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Subscribe(): %v", token.Error())
	}

	if token := p.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	text := "pubmsg payload"
	p.Publish(topic, 0, false, text)
	p.Publish(topic, 0, false, text)
	wait(choke)
	wait(choke)

	p.Publish(topic, 0, false, text)
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)
}

func Test_PublishEmptyMessage(t *testing.T) {
	topic := "/test/pubmsgempty"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.AddBroker(FVTTCP)
	pops.SetClientID("pubmsgempty-pub")
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.AddBroker(FVTTCP)
	sops.SetClientID("pubmsgempty-sub")
	var f MessageHandler = func(client Client, msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		if string(msg.Payload()) != "" {
			t.Fatalf("Message payload incorrect")
		}
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	if sToken := s.Connect(); sToken.Wait() && sToken.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", sToken.Error())
	}

	if sToken := s.Subscribe(topic, 2, nil); sToken.Wait() && sToken.Error() != nil {
		t.Fatalf("Error on Client.Subscribe(): %v", sToken.Error())
	}

	if pToken := p.Connect(); pToken.Wait() && pToken.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", pToken.Error())
	}

	p.Publish(topic, 0, false, "")
	p.Publish(topic, 0, false, "")
	wait(choke)
	wait(choke)

	p.Publish(topic, 0, false, "")
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)
}

// func Test_Cleanstore(t *testing.T) {
// 	store := "/tmp/fvt/cleanstore"
// 	topic := "/test/cleanstore"

// 	pops := NewClientOptions()
// 	pops.AddBroker(FVTTCP)
// 	pops.SetClientID("cleanstore-pub")
// 	pops.SetStore(NewFileStore(store + "/p"))
// 	p := NewClient(pops)

// 	var s *Client
// 	sops := NewClientOptions()
// 	sops.AddBroker(FVTTCP)
// 	sops.SetClientID("cleanstore-sub")
// 	sops.SetCleanSession(false)
// 	sops.SetStore(NewFileStore(store + "/s"))
// 	var f MessageHandler = func(client Client, msg Message) {
// 		fmt.Printf("TOPIC: %s\n", msg.Topic())
// 		fmt.Printf("MSG: %s\n", msg.Payload())
// 		// Close the connection after receiving
// 		// the first message so that hopefully
// 		// there is something in the store to be
// 		// cleaned.
// 		s.ForceDisconnect()
// 	}
// 	sops.SetDefaultPublishHandler(f)

// 	s = NewClient(sops)
// 	sToken := s.Connect()
// 	if sToken.Wait() && sToken.Error() != nil {
// 		t.Fatalf("Error on Client.Connect(): %v", sToken.Error())
// 	}

// 	sToken = s.Subscribe(topic, 2, nil)
// 	if sToken.Wait() && sToken.Error() != nil {
// 		t.Fatalf("Error on Client.Subscribe(): %v", sToken.Error())
// 	}

// 	pToken := p.Connect()
// 	if pToken.Wait() && pToken.Error() != nil {
// 		t.Fatalf("Error on Client.Connect(): %v", pToken.Error())
// 	}

// 	text := "test message"
// 	p.Publish(topic, 0, false, text)
// 	p.Publish(topic, 0, false, text)
// 	p.Publish(topic, 0, false, text)

// 	p.Disconnect(250)

// 	s2ops := NewClientOptions()
// 	s2ops.AddBroker(FVTTCP)
// 	s2ops.SetClientID("cleanstore-sub")
// 	s2ops.SetCleanSession(true)
// 	s2ops.SetStore(NewFileStore(store + "/s"))
// 	s2ops.SetDefaultPublishHandler(f)

// 	s2 := NewClient(s2ops)
// 	sToken = s2.Connect()
// 	if sToken.Wait() && sToken.Error() != nil {
// 		t.Fatalf("Error on Client.Connect(): %v", sToken.Error())
// 	}

// 	// at this point existing state should be cleared...
// 	// how to check?
// }

func Test_MultipleURLs(t *testing.T) {
	ops := NewClientOptions()
	ops.AddBroker("tcp://127.0.0.1:10000")
	ops.AddBroker(FVTTCP)
	ops.SetClientID("MultiURL")

	c := NewClient(ops)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	if pToken := c.Publish("/test/MultiURL", 0, false, "Publish qo0"); pToken.Wait() && pToken.Error() != nil {
		t.Fatalf("Error on Client.Publish(): %v", pToken.Error())
	}

	c.Disconnect(250)
}

// A test to make sure ping mechanism is working
func Test_ping1_idle5(t *testing.T) {
	ops := NewClientOptions()
	ops.AddBroker(FVTTCP)
	ops.SetClientID("p3i10")
	ops.SetConnectionLostHandler(func(c Client, err error) {
		t.Fatalf("Connection-lost handler was called: %s", err)
	})
	ops.SetKeepAlive(4 * time.Second)

	c := NewClient(ops)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}
	time.Sleep(8 * time.Second)
	c.Disconnect(250)
}

func Test_autoreconnect(t *testing.T) {
	ops := NewClientOptions()
	ops.AddBroker(FVTTCP)
	ops.SetClientID("auto_reconnect")
	ops.SetAutoReconnect(true)
	ops.SetOnConnectHandler(func(c Client) {
		t.Log("Connected")
	})
	ops.SetKeepAlive(2 * time.Second)

	c := NewClient(ops)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	time.Sleep(5 * time.Second)

	fmt.Println("Breaking connection")
	c.(*client).internalConnLost(fmt.Errorf("autoreconnect test"))

	time.Sleep(5 * time.Second)
	if !c.IsConnected() {
		t.Fail()
	}

	c.Disconnect(250)
}

func Test_cleanUpMids(t *testing.T) {
	ops := NewClientOptions()
	ops.AddBroker(FVTTCP)
	ops.SetClientID("auto_reconnect")
	ops.SetCleanSession(true)
	ops.SetAutoReconnect(true)
	ops.SetKeepAlive(10 * time.Second)

	c := NewClient(ops)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	token := c.Publish("/test/cleanUP", 2, false, "cleanup test")
	c.(*client).messageIds.Lock()
	fmt.Println("Breaking connection", len(c.(*client).messageIds.index))
	if len(c.(*client).messageIds.index) == 0 {
		t.Fatalf("Should be a token in the messageIDs, none found")
	}
	c.(*client).messageIds.Unlock()
	c.(*client).internalConnLost(fmt.Errorf("cleanup test"))

	time.Sleep(1 * time.Second)
	if !c.IsConnected() {
		t.Fail()
	}

	c.(*client).messageIds.Lock()
	if len(c.(*client).messageIds.index) > 0 {
		t.Fatalf("Should have cleaned up messageIDs, have %d left", len(c.(*client).messageIds.index))
	}
	c.(*client).messageIds.Unlock()

	// This test used to check that token.Error() was not nil. However this is not something that can
	// be done reliably - it is likely to work with a remote broker but less so with a local one.
	// This is because:
	// - If the publish fails in net.go while transmitting then an error will be generated
	// - If the transmit succeeds (regardless of whether the handshake completes then no error is generated)
	// If the intention is that an error should always be returned if the publish is incomplete upon disconnedt then
	// internalConnLost needs to be altered (if c.options.CleanSession && !c.options.AutoReconnect)
	//if token.Error() == nil {
	//t.Fatal("token should have received an error on connection loss")
	//}
	fmt.Println(token.Error())

	c.Disconnect(250)
}

// Test that cleanup happens properly on explicit Disconnect()
func Test_cleanUpMids_2(t *testing.T) {
	ops := NewClientOptions()
	ops.AddBroker(FVTTCP)
	ops.SetClientID("auto_reconnect")
	ops.SetCleanSession(true)
	ops.SetAutoReconnect(true)
	ops.SetKeepAlive(10 * time.Second)

	c := NewClient(ops)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	token := c.Publish("/test/cleanUP", 2, false, "cleanup test 2")
	if len(c.(*client).messageIds.index) == 0 {
		t.Fatalf("Should be a token in the messageIDs, none found")
	}
	fmt.Println("Disconnecting", len(c.(*client).messageIds.index))
	c.Disconnect(0)

	fmt.Println("Wait on Token")
	// We should be able to wait on this token without any issue
	token.Wait()

	if len(c.(*client).messageIds.index) > 0 {
		t.Fatalf("Should have cleaned up messageIDs, have %d left", len(c.(*client).messageIds.index))
	}
	if token.Error() == nil {
		t.Fatal("token should have received an error on connection loss")
	}
	fmt.Println(token.Error())
}

func Test_ConnectRetry(t *testing.T) {
	// Connect for publish - initially use invalid server
	cops := NewClientOptions().AddBroker("256.256.256.256").SetClientID("cr-pub").
		SetConnectRetry(true).SetConnectRetryInterval(time.Second / 2)
	c := NewClient(cops).(*client)
	connectToken := c.Connect()

	time.Sleep(time.Second) // Wait a second to ensure we are past SetConnectRetryInterval
	if connectToken.Error() != nil {
		t.Fatalf("Connect returned error (should be retrying) (%v)", connectToken.Error())
	}
	c.optionsMu.Lock() // Protect c.options.Servers so that servers can be added in test cases
	c.options.AddBroker(FVTTCP)
	c.optionsMu.Unlock()
	if connectToken.Wait() && connectToken.Error() != nil {
		t.Fatalf("Error connecting after valid broker added: %v", connectToken.Error())
	}
	c.Disconnect(250)
}

func Test_ConnectRetryPublish(t *testing.T) {
	topic := "/test/connectRetry"
	payload := "sample Payload"
	choke := make(chan bool)

	// subscribe to topic and wait for expected message (only received after connection successful)
	sops := NewClientOptions().AddBroker(FVTTCP).SetClientID("crp-sub")
	var f MessageHandler = func(client Client, msg Message) {
		if msg.Topic() != topic || string(msg.Payload()) != payload {
			t.Fatalf("Received unexpected message: %v, %v", msg.Topic(), msg.Payload())
		}
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	if token := s.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Connect(): %v", token.Error())
	}

	if token := s.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on Client.Subscribe(): %v", token.Error())
	}

	// Connect for publish - initially use invalid server
	memStore := NewMemoryStore()
	memStore.Open()
	pops := NewClientOptions().AddBroker("256.256.256.256").SetClientID("crp-pub").
		SetStore(memStore).SetConnectRetry(true).SetConnectRetryInterval(time.Second / 2)
	p := NewClient(pops).(*client)
	connectToken := p.Connect()
	p.Publish(topic, 1, false, payload)
	// Check publish packet in the memorystore
	ids := memStore.All()
	if len(ids) == 0 {
		t.Fatalf("Expected published message to be in store")
	} else if len(ids) != 1 {
		t.Fatalf("Expected 1 message to be in store")
	}
	packet := memStore.Get(ids[0])
	if packet == nil {
		t.Fatal("Failed to retrieve packet from store")
	}
	pp, ok := packet.(*packets.PublishPacket)
	if !ok {
		t.Fatalf("Message in store not of the expected type (%T)", packet)
	}
	if pp.TopicName != topic || string(pp.Payload) != payload {
		t.Fatalf("Stored message Packet contents not as expected (%v, %v)", pp.TopicName, pp.Payload)
	}
	time.Sleep(time.Second) // Wait a second to ensure we are past SetConnectRetryInterval
	if connectToken.Error() != nil {
		t.Fatalf("Connect returned error (should be retrying) (%v)", connectToken.Error())
	}

	// disconnecting closes the store (both in disconnect and in Connect which runs as a goRoutine).
	// As such we duplicate the store
	memStore2 := NewMemoryStore()
	memStore2.Open()
	memStore2.Put(ids[0], packet)

	// disconnect and then reconnect with correct server
	p.Disconnect(250)

	pops = NewClientOptions().AddBroker(FVTTCP).SetClientID("crp-pub").SetCleanSession(false).
		SetStore(memStore2).SetConnectRetry(true).SetConnectRetryInterval(time.Second / 2)
	p = NewClient(pops).(*client)
	if token := p.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error on valid Publish.Connect(): %v", token.Error())
	}

	if connectToken.Wait() && connectToken.Error() == nil {
		t.Fatalf("Expected connection error - got nil")
	}
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)
	memStore.Close()
}

func Test_ResumeSubs(t *testing.T) {
	topic := "/test/ResumeSubs"
	var qos byte = 1
	payload := "sample Payload"
	choke := make(chan bool)

	// subscribe to topic before establishing a connection, and publish a message after the publish client has connected successfully
	subMemStore := NewMemoryStore()
	subMemStore.Open()
	sops := NewClientOptions().AddBroker("256.256.256.256").SetClientID("resumesubs-sub").SetConnectRetry(true).
		SetConnectRetryInterval(time.Second / 2).SetResumeSubs(true).SetStore(subMemStore)

	s := NewClient(sops)
	sConnToken := s.Connect()

	subToken := s.Subscribe(topic, qos, nil)

	// Verify the subscribe packet exists in the memorystore
	ids := subMemStore.All()
	if len(ids) == 0 {
		t.Fatalf("Expected subscribe packet to be in store")
	} else if len(ids) != 1 {
		t.Fatalf("Expected 1 packet to be in store")
	}
	packet := subMemStore.Get(ids[0])
	if packet == nil {
		t.Fatal("Failed to retrieve packet from store")
	}
	sp, ok := packet.(*packets.SubscribePacket)
	if !ok {
		t.Fatalf("Packet in store not of the expected type (%T)", packet)
	}
	if len(sp.Topics) != 1 || sp.Topics[0] != topic || len(sp.Qoss) != 1 || sp.Qoss[0] != qos {
		t.Fatalf("Stored Subscribe Packet contents not as expected (%v, %v)", sp.Topics, sp.Qoss)
	}

	time.Sleep(time.Second) // Wait a second to ensure we are past SetConnectRetryInterval
	if sConnToken.Error() != nil {
		t.Fatalf("Connect returned error (should be retrying) (%v)", sConnToken.Error())
	}
	if subToken.Error() != nil {
		t.Fatalf("Subscribe returned error (should be persisted) (%v)", sConnToken.Error())
	}

	// test that the stored subscribe packet gets sent to the broker after connecting
	subMemStore2 := NewMemoryStore()
	subMemStore2.Open()
	subMemStore2.Put(ids[0], packet)

	s.Disconnect(250)

	// Connect to broker and test that subscription was resumed
	sops = NewClientOptions().AddBroker(FVTTCP).SetClientID("resumesubs-sub").
		SetStore(subMemStore2).SetResumeSubs(true).SetCleanSession(false).SetConnectRetry(true).
		SetConnectRetryInterval(time.Second / 2)

	var f MessageHandler = func(client Client, msg Message) {
		if msg.Topic() != topic || string(msg.Payload()) != payload {
			t.Fatalf("Received unexpected message: %v, %v", msg.Topic(), msg.Payload())
		}
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)
	s = NewClient(sops).(*client)
	if sConnToken = s.Connect(); sConnToken.Wait() && sConnToken.Error() != nil {
		t.Fatalf("Error on valid subscribe Connect(): %v", sConnToken.Error())
	}

	// publish message to subscribed topic to verify subscription
	pops := NewClientOptions().AddBroker(FVTTCP).SetClientID("resumesubs-pub").SetCleanSession(true).
		SetConnectRetry(true).SetConnectRetryInterval(time.Second / 2)
	p := NewClient(pops).(*client)
	if pConnToken := p.Connect(); pConnToken.Wait() && pConnToken.Error() != nil {
		t.Fatalf("Error on valid Publish.Connect(): %v", pConnToken.Error())
	}

	if pubToken := p.Publish(topic, 1, false, payload); pubToken.Wait() && pubToken.Error() != nil {
		t.Fatalf("Error on valid Client.Publish(): %v", pubToken.Error())
	}

	wait(choke)

	s.Disconnect(250)
	p.Disconnect(250)
}

func Test_ResumeSubsWithReconnect(t *testing.T) {
	topic := "/test/ResumeSubs"
	var qos byte = 1

	// subscribe to topic before establishing a connection, and publish a message after the publish client has connected successfully
	ops := NewClientOptions().SetClientID("Start").AddBroker(FVTTCP).SetConnectRetry(true).SetConnectRetryInterval(time.Second / 2).
		SetResumeSubs(true).SetCleanSession(false)
	c := NewClient(ops)
	sConnToken := c.Connect()
	sConnToken.Wait()
	if sConnToken.Error() != nil {
		t.Fatalf("Connect returned error (%v)", sConnToken.Error())
	}

	// Send subscription request and then immediately force disconnect (hope it will happen before sub sent)
	subToken := newToken(packets.Subscribe).(*SubscribeToken)
	sub := packets.NewControlPacket(packets.Subscribe).(*packets.SubscribePacket)
	sub.Topics = append(sub.Topics, topic)
	sub.Qoss = append(sub.Qoss, qos)
	subToken.subs = append(subToken.subs, topic)

	if sub.MessageID == 0 {
		sub.MessageID = c.(*client).getID(subToken)
		subToken.messageID = sub.MessageID
	}
	DEBUG.Println(CLI, sub.String())

	persistOutbound(c.(*client).persist, sub)
	//subToken := c.Subscribe(topic, qos, nil)
	c.(*client).internalConnLost(fmt.Errorf("reconnection subscription test"))

	// As reconnect is enabled the client should automatically reconnect
	subDone := make(chan bool)
	go func(t *testing.T) {
		subToken.Wait()
		if err := subToken.Error(); err != nil {
			t.Errorf("Connect returned error (should be retrying) (%v)", err)
		}
		close(subDone)
	}(t)
	// Wait for done or timeout
	select {
	case <-subDone:
	case <-time.After(4 * time.Second):
		t.Fatalf("Timed out waiting for subToken to complete")
	}

	c.Disconnect(250)
}
