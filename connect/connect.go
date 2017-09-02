package connect

import (
	"bufio"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/fatih/color"
)

const TIMEOUT = 15

var CONNS = sync.Map{}

func GetConnection(address string) (*net.TCPConn, error) {

	color.Set(color.FgHiCyan)
	log.Printf("[connection] getting connection to address on device %s", address)
	color.Unset()

	//first see if the entry is in the map
	conn, ok := CONNS.Load(address)
	if !ok {

		color.Set(color.FgHiCyan)
		log.Printf("[connection] connection not found, connecting to %s", address)
		color.Unset()

		addr, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return nil, err
		}

		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			return nil, err
		}

		conn.SetDeadline(time.Now().Add(TIMEOUT * time.Second))
		CONNS.LoadOrStore(address, conn)

		return conn, nil
	}

	output, _ := conn.(net.TCPConn)
	return &output, nil
}

//@param conn - the connection in question
//@param msg - the message to be read or written
//@param act - connect.Read or connect.Write
//refreshes the connection by extending the deadline, then re-writes msg to the
//connection or re-reads until the first byte of msg is found
//@pre conn has connected successfully prior to this function call -- otherwise it will trigger a panic!
func HandleTimeout(conn *net.TCPConn, msg []byte, act ReadWrite) ([]byte, error) {

	//these three happen in any case
	errMsg := color.HiRedString("[connection] connection timed out, retrying...")

	if len(msg) == 0 {
		return []byte{}, errors.New("cannot write empty message to TCP connection")
	}

	log.Printf("%s", errMsg)
	conn.SetDeadline(time.Now().Add(TIMEOUT * time.Second))

	if act == Write {

		_, err := conn.Write(msg)
		return msg, err
	}

	reader := bufio.NewReader(conn)
	return reader.ReadBytes(msg[0])

}

type ReadWrite int

const (
	Read ReadWrite = 1 + iota
	Write
)
