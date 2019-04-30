package main

import (
	"net"

	"github.com/premeidoworks/kanata/service/kanata_discovery"
)

func StartServer(discovery kanata_discovery.Discovery, listen string) (closeFunc func(), err error) {
	addr, err := net.ResolveTCPAddr("tcp", listen)
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				logErr(err)
				break
			}
			logInfo("connection income: ", conn)
			go clientIncoming(conn, discovery)
		}
	}()

	return func() {
		err := listener.Close()
		if err != nil {
			logErr(err)
		}
		//FIXME wait for all connection closed
	}, nil
}

func clientIncoming(conn *net.TCPConn, discovery kanata_discovery.Discovery) {
	//TODO
}
