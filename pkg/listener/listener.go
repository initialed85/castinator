package listener

import (
	"fmt"
	"github.com/initialed85/castinator/internal/common"
	"github.com/initialed85/castinator/pkg/interfaces"
	"net"
	"syscall"
)

func GetListener(addr *net.UDPAddr, intfc *net.Interface) (conn *net.UDPConn, err error) {
	network := interfaces.GetNetwork(addr.String())

	if addr.IP.IsMulticast() {
		conn, err = net.ListenMulticastUDP(network, intfc, addr)
		if err != nil {
			err = fmt.Errorf("failed to ListenMulticastUDP because %v", err)
			return
		}
	} else {
		conn, err = net.ListenUDP(network, addr)
		if err != nil {
			err = fmt.Errorf("failed to ListenUDP because %v", err)
			return
		}
	}

	err = conn.SetReadBuffer(common.MaxDatagramSize)
	if err != nil {
		err = fmt.Errorf("failed to SetReadBuffer to %v because %v", common.MaxDatagramSize, err)
		return
	}

	file, err := conn.File()
	if err != nil {
		err = fmt.Errorf("failed to File for the socket because %v", err)
		return
	}

	fd := int(file.Fd())

	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		err = fmt.Errorf("failed to set SO_REUSADDR on the socket's file descriptor because %v", err)
		return
	}

	return
}
