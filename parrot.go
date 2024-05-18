/* A small Go program to republish data from MQTT to MQTT.
 *
 * This program will exit on any error, so be sure to run it in an init system
 * or other process manager.
 *
 * This program can also be ran through the use of containers, use either
 * `docker` or `podman`:
 *
 * `podman run -e MQTT_HOST="tcp://127.0.0.1:1883" github.com/petspalace/parrot:latest`
 *
 * This program was made by:
 * - Simon de Vlieger <cmdr@supakeen.com>
 *
 * This program is licensed under the MIT license:
 *
 * Copyright 2022 Simon de Vlieger
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to
 * deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var logger = log.New(os.Stderr, "", log.LstdFlags)

/* Message passed along between *Loop and MessageLoop through a channel,
 * *Loop determines the data and where it goes. */
type MQTTMessage struct {
	Topic   string
	Payload string
	Retain  bool
}

/* Callback for an MQTT client which puts received messages into a channel. */
func MessageReadLoop(c MQTT.Client, rx, tx chan MQTTMessage) {
	for m := range rx {
		logger.Printf("Received message on topic %s: '%s'\n", m.Topic, m.Payload)
	}
}

/* Listens on a channel and submits messages through an MQTT client. */
func MessageSendLoop(c MQTT.Client, rx, tx chan MQTTMessage) {
	for m := range tx {
		topic := fmt.Sprintf("%s", m.Topic)

		if token := c.Publish(topic, 0, m.Retain, m.Payload); token.Wait() && token.Error() != nil {
			logger.Fatalln("MessageLoop could not publish message.")
		}

		logger.Printf("MessageLoop published topic='%s',payload='%s'\n", topic, m.Payload)
	}
}

func main() {
	hostFromEnv, hostExists := os.LookupEnv("MQTT_HOST")

	if !hostExists {
		logger.Fatalln("parrot needs `MQTT_HOST` set in the environment to a value such as `tcp://127.0.0.1:1883`.")
	}

	opts := MQTT.NewClientOptions().AddBroker(hostFromEnv).SetClientID("parrot")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := MQTT.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		logger.Panic(token.Error())
	}

	tx := make(chan MQTTMessage)
	rx := make(chan MQTTMessage)

	// Helper function to throw MQTT.Message into the received queue
	f := func(_ MQTT.Client, m MQTT.Message) {
		rx <- MQTTMessage{
			Topic:   m.Topic(),
			Payload: string(m.Payload()),
		}
	}

	if token := c.Subscribe("#", byte(0), f); token.Wait() && token.Error() != nil {
		logger.Panic(token.Error())
	}

	go MessageReadLoop(c, rx, tx)
	MessageSendLoop(c, rx, tx)

	c.Disconnect(250)

	time.Sleep(1 * time.Second)
}

// SPDX-License-Identifier: MIT
// vim: ts=4 sw=4 noet
