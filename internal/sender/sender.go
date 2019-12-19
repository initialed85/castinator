package sender

import (
	"fmt"
	"github.com/initialed85/castinator/internal/interfaces"
	"net"
)

func GetSender(addr *net.UDPAddr, srcAddr *net.UDPAddr) (conn *net.UDPConn, err error) {
	network := interfaces.GetNetwork(addr.String())

	conn, err = net.DialUDP(network, srcAddr, addr)
	if err != nil {
		err = fmt.Errorf("failed to DialUDP because %v", err)
		return
	}

	return
}
