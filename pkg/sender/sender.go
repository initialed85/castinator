package sender

import (
	"fmt"
	"net"
	"strings"
)

func GetSender(addr *net.UDPAddr, srcAddr *net.UDPAddr) (conn *net.UDPConn, err error) {
	network := "udp4"
	if strings.Count(addr.String(), ":") > 1 {
		network = "udp6"
	}

	conn, err = net.DialUDP(network, srcAddr, addr)
	if err != nil {
		err = fmt.Errorf("failed to DialUDP because %v", err)
		return
	}

	return
}
