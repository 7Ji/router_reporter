package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

func listen_process(conn net.Conn, iface_addr *string) {
	defer conn.Close()
	for {
		msg := *iface_addr + "\n"
		_, err := conn.Write([]byte(msg))
		if err != nil {
			log.Println("Failed to continue writing", err)
			return
		}
		time.Sleep(time.Second)
	}
}

func listen_worker(listener net.Listener, iface_addr *string) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept TCP incoming stream,", err)
		} else {
			defer connection.Close()
			log.Println("Accepted TCP incoming stream,", connection.RemoteAddr())
			go listen_process(connection, iface_addr)
		}
	}
}

func ip_updater(iface_name string, iface_addr *string) {
	for {
		iface, err := net.InterfaceByName(iface_name)
		if err != nil {
			log.Println("Failed to get interface with name", iface_name, err)
		} else {
			addrs, err := iface.Addrs()
			if err != nil {
				log.Println("Failed to get addresses on interface", iface_name, err)
			}
			for _, addr := range addrs {
				addr_string := addr.String()
				if !strings.ContainsAny(addr_string, ":abcdef") {
					*iface_addr = addr_string
					break
				}
			}
		}
		time.Sleep(time.Second)
	}
}

func main() {
	flag_iface := flag.String("iface", "pppoe-wan", "Interface to get ip from")
	flag_listen := flag.String("listen", ":7777", "[host]:port to listen on")
	flag.Parse()
	listener, err := net.Listen("tcp4", *flag_listen)
	if err != nil {
		log.Println("Failed to listen on", *flag_listen, err)
	}
	defer listener.Close()
	var iface_addr string
	go ip_updater(*flag_iface, &iface_addr)
	go listen_worker(listener, &iface_addr)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	select {
	case <-interrupt:
		log.Println("Interrupt received, exiting")
		return
	}
}
