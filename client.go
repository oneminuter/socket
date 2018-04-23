package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"socket/util"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8888")
	if (util.CheckErr(err)) {
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if (util.CheckErr(err)) {
		return
	}

	_, err = conn.Write([]byte("timestamp"))
	if (util.CheckErr(err)) {
		return
	}

	result, err := ioutil.ReadAll(conn)
	if (util.CheckErr(err)) {
		return
	}
	fmt.Println(string(result))
	os.Exit(0)
}
