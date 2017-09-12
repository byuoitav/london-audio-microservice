package connect

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"syscall"
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
	}

	//cast to TCP connection and refresh
	output, _ := conn.(net.TCPConn)
	output.SetDeadline(time.Now().Add(TIMEOUT * time.Second))

	//check for broken pipe error (DSP reboot while microservice is still running)
	_, err := output.Write([]byte{0x00})
	if err != nil && err == syscall.EPIPE {
		output.Close()
		CONNS.Delete(address)
		err = addConnection(address)
		if err != nil {
			return nil, err
		}

		conn, _ = CONNS.Load(address)
		output, _ = conn.(net.TCPConn)
	}

	return &output, nil
}

//@param conn - the connection in question
//@param msg - the message to be read or written
//@param act - connect.Read or connect.Write
//refreshes the connection by extending the deadline, then re-writes msg to the
//connection or re-reads until the first byte of msg is found
//@pre conn has connected successfully prior to this function call -- otherwise it will trigger a panic!
func HandleTimeout(conn *net.TCPConn, msg []byte, method ReadWrite) ([]byte, error) {

	//these three happen in any case
	log.Printf("%s", color.HiRedString("[connection] connection timed out, retrying..."))

	if len(msg) == 0 {
		return []byte{}, errors.New("cannot write empty message to TCP connection")
	}

	conn.SetDeadline(time.Now().Add(TIMEOUT * time.Second))

	if method == Write {

		_, err := conn.Write(msg)
		return msg, err
	}

	reader := bufio.NewReader(conn)
	return reader.ReadBytes(msg[0])

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

func HandleBrokenPipe(address string) error {

	msg := fmt.Sprintf("[connection] handling broken pipe error with address %s...", address)
	log.Printf("%s", color.HiRedString("%s", msg))

	if conn, ok := CONNS.Load(address); !ok {

		msg := fmt.Sprintf("[connection] connection to %s not found. Adding to connection store...")
		log.Printf("%s", color.HiRedString("%s", msg))

		err := addConnection(address)
		if err != nil {
			msg = fmt.Sprintf("[connection] unable to add connection: %s", err.Error())
			log.Printf("%s", color.HiRedString("%s", msg))
			return errors.New(msg)
		}

	}

	conn.Close()

	return nil
}
