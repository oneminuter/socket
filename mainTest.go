package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"socket/util"
	"strconv"
	"strings"
	"time"
)

func main() {
	//httpServer()
	socketServer()
}

func httpServer() {
	http.HandleFunc("/html/", httpHandler)
	http.ListenAndServe(":8080", nil)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("http server start...")
	bytes, err := ioutil.ReadFile("html/index.html")
	util.CheckErr(err)
	w.Write(bytes)
}

func socketServer() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8888")
	if util.CheckErr(err) {
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if util.CheckErr(err) {
		return
	}
	log.Println("socket server starting ...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go socketHandler(conn)
	}
}

func socketHandler(conn net.Conn) {
	fmt.Println(conn.RemoteAddr(), "connect...")
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(2 * time.Minute))
	request := make([]byte, 128)
	for {
		len, _ := conn.Read(request)
		if len == 0 {
			log.Println("len == 0")
			return
		} else if strings.TrimSpace(string(request[:len])) == "timestamp" {
			daytime := strconv.FormatInt(time.Now().Unix(), 10)
			conn.Write([]byte(daytime))
		} else {
			daytime := time.Now().String()
			conn.Write([]byte(daytime))
		}
		request = make([]byte, 128)
	}
}
