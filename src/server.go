package main

import (
	"fmt"
	"log"
	"net"
)

type Session struct {
	id string
	conn net.Conn
}

func (s *Session) SessionId() string {
	return s.id
}

func (s *Session) WriteLine(str string) error {
	_, err := s.conn.Write([]byte(str + "\r\n"))
	return err
}

var nxtsid = 1
func generateSessionId() string {
	var sid = nxtsid
	nxtsid++
	return fmt.Sprintf("%d", sid)
}

func handleConn(conn net.Conn, inputChannel chan SessionEvent) error{
	buf := make([]byte, 4096)
	session := &Session{generateSessionId(), conn}
//	user := &User{name:  generateName(), session: session}

	inputChannel <- SessionEvent{session, &SessionCreatedEvent{},}

	for {
		n, err := conn.Read(buf)
		if err != nil {
			inputChannel <- SessionEvent{session, &SessionDisconnectEvent{}}
			return err
		}
		if n == 0 {
			log.Println("Terminating connection")
			inputChannel <- SessionEvent{session, &SessionDisconnectEvent{}}
			break
		}
		msg := string(buf[0 : n-2])
		log.Println("ACK :", msg)

		inputChannel <- SessionEvent{session, &SessionInputEvent{msg}}
	}

	return nil
}

func startServer(eventChannel chan SessionEvent) error {
	println("Initialising GRID")
	ln, err := net.Listen("tcp", ":77")

	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}
		go func() {
			if err := handleConn(conn, eventChannel); err != nil {
				log.Println("Error handling connection", err)
				return
			}
		}()
	}
}

