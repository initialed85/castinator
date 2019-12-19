package handler

import (
	"fmt"
	"github.com/initialed85/castinator/internal/common"
	"log"
	"net"
	"time"
)

func Handle(listener, sender, otherSender *net.UDPConn) {
	for {
		buf := make([]byte, common.MaxDatagramSize)
		n, src, err := listener.ReadFromUDP(buf)
		if err != nil {
			log.Printf(fmt.Sprintf("failed to ReadFromUDP because %v", err))

			time.Sleep(time.Millisecond * 100)

			continue
		}

		log.Printf("%v -> %v - received %v bytes", src, listener.LocalAddr().String(), n)

		data := string(buf[:n])

		senderLocalAddr := sender.LocalAddr().(*net.UDPAddr)
		if senderLocalAddr.IP.Equal(src.IP) && senderLocalAddr.Port == src.Port {
			log.Printf("%v -> %v - skipped because src is sender", src, listener.LocalAddr().String())
			continue
		}

		otherSenderLocalAddr := otherSender.LocalAddr().(*net.UDPAddr)
		if otherSenderLocalAddr.IP.Equal(src.IP) && otherSenderLocalAddr.Port == src.Port {
			log.Printf("%v -> %v - skipped because src is otherSender", src, listener.LocalAddr().String())
			continue
		}

		err = sender.SetWriteDeadline(time.Now().Add(time.Second))
		if err != nil {
			log.Fatal(err)
		}

		_, err = sender.Write([]byte(data))
		if err != nil {
			log.Printf("failed to Write because %v", err)

			time.Sleep(time.Millisecond * 100)

			continue
		}

		log.Printf("%v -> %v - sent %v bytes to %v", src, listener.LocalAddr().String(), n, sender.RemoteAddr())
	}
}
