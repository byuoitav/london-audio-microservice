package connect

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/fatih/color"
	"golang.org/x/sync/syncmap"
)

const TIMEOUT = 5

var CONNS = syncmap.Map{}

type ReadWrite int

const (
	Read ReadWrite = 1 + iota
	Write
)

func GetConnection(address string) (*net.TCPConn, error) {

	log.Printf("%s", color.HiCyanString("[connection] getting connection to address on device %s", address))

	//first see if the entry is in the map
	conn, ok := CONNS.Load(address)
	if !ok {
		err := addConnection(address)
		if err != nil {
			return nil, err
		}

		conn, ok = CONNS.Load(address)
		if !ok {
			msg := "unable to add address to map"
			log.Printf("%s", color.HiRedString("[connection] %s", msg))
			return nil, errors.New(msg)
		}
	}

	//cast to TCP connection and refresh
	output, _ := conn.(net.TCPConn)
	output.SetDeadline(time.Now().Add(TIMEOUT * time.Second))

	//check for broken pipe error (DSP reboot while microservice is still running)
	//check in old event-router commit
	_, err := output.Write([]byte{0x00})
	if err != nil {

		log.Printf("%s", color.HiRedString("[connection] unable to write to connection: %s refreshing...", err.Error()))
		err = HandleStaleConnection(&output)
		if err != nil {
			msg := fmt.Sprintf("unable to refresh connection: %s", err.Error())
			log.Printf("%s", color.HiRedString("[connection] %s", msg))
			return nil, errors.New(msg)
		}

		err = addConnection(address)
		if err != nil {
			return nil, err
		}

		conn, _ = CONNS.Load(address)
		output, _ = conn.(net.TCPConn)
	}

	return &output, nil
}

func HandleStaleConnection(conn *net.TCPConn) error {

	if conn == nil {
		msg := "null connection"
		log.Printf("%s", color.HiRedString("[connection] %s", msg))
		return errors.New(msg)
	}

	if conn.RemoteAddr() == nil {
		msg := "no remote address"
		log.Printf("%s", color.HiRedString("[connection] %s", msg))
		return errors.New(msg)
	}

	address := conn.RemoteAddr().String()
	log.Printf("[connection] handling stale connection: %s", address)

	//close connection
	conn.Close()

	//remove connection from map
	//it will get added to the map on the next call to GetConnection
	CONNS.Delete(address)

	return nil
}

func addConnection(address string) error {

	log.Printf("[connection] adding connection to %s...", address)

	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	conn.SetDeadline(time.Now().Add(TIMEOUT * time.Second))

	CONNS.LoadOrStore(address, *conn)

	return nil
}
