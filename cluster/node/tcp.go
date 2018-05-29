package node

//import (
//	"net"
//	"log"
//)
//
//func InitTCP(address string, accept int) (err error) {
//	var (
//		addr   *net.TCPAddr
//		listen *net.TCPListener
//	)
//	if addr, err = net.ResolveTCPAddr("tcp", address); err != nil {
//		log.Printf("net.ResolveTCPAddr(\"tcp\", \"%s\") error(%v)", address, err)
//		return
//	}
//	if listen, err = net.ListenTCP("tpc", addr); err != nil {
//		log.Printf("net.ListenTCP(\"tcp\", \"%v\") error(%v)", addr, err)
//		return
//	}
//	// N CPU
//	for i := 0; i < accept; i++ {
//		go acceptTCP(listen)
//	}
//	return nil
//}
//
//func acceptTCP(listen *net.TCPListener) {
//	var (
//		err error
//		conn *net.TCPConn
//	)
//	conn ,err = listen.AcceptTCP()
//}
