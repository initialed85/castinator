package interfaces

import (
	"fmt"
	"net"
	"strings"
)

func GetAddressesAndInterfaces(rawIntfc, rawAddr string) (addr *net.UDPAddr, intfc *net.Interface, srcAddr *net.UDPAddr, err error) {
	intfc, err = net.InterfaceByName(rawIntfc)
	if err != nil {
		err = fmt.Errorf("failed to get interface because %v", err)
		return
	}

	intfcAddrs, err := intfc.Addrs()
	if err != nil {
		err = fmt.Errorf("failed to get intfcAddrs because %v", err)
		return
	}

	addrNetwork := "udp4"
	if strings.Count(rawAddr, ":") > 1 {
		addrNetwork = "udp6"
	}

	addr, err = net.ResolveUDPAddr(addrNetwork, rawAddr)
	if err != nil {
		err = fmt.Errorf("failed to get addr because %v", err)
		return
	}

	srcAddr = &net.UDPAddr{}
	for _, v := range intfcAddrs {
		ipNet := v.(*net.IPNet)

		if addrNetwork == "udp4" && strings.Count(ipNet.IP.String(), ":") > 0 {
			continue
		} else if addrNetwork == "udp6" && strings.Count(ipNet.IP.String(), ":") == 0 {
			continue
		} else if ipNet.IP.IsLinkLocalUnicast() || ipNet.IP.IsInterfaceLocalMulticast() {
			continue
		}

		srcAddr.IP = ipNet.IP
		srcAddr.Zone = intfc.Name

		break
	}

	return
}
