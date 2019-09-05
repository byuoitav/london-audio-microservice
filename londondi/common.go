package londondi

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/byuoitav/common/pooled"
	"github.com/fatih/color"
)

//TIMEOUT .
const TIMEOUT = 5

var pool = pooled.NewMap(45*time.Second, 400*time.Millisecond, getConnection)

//ReadWrite .
type ReadWrite int

//ReadWrite .
const (
	Read ReadWrite = 1 + iota
	Write
)

//GetConnection .
func getConnection(key interface{}) (pooled.Conn, error) {
	address, ok := key.(string)
	if !ok {
		return nil, fmt.Errorf("key must be a string")
	}
	log.Printf("%s", color.HiCyanString("[connection] getting connection to address on device %s", address))

	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return nil, err
	}

	// read the NOKEY line
	pconn := pooled.Wrap(conn)

	return pconn, nil
}
