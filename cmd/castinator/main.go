package castinator

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"syscall"
)

const (
	maxDatagramSize = 65536
)

func getAddressesAndInterfaces(rawIntfc, rawAddr string) (addr *net.UDPAddr, intfc *net.Interface, srcAddr *net.UDPAddr, err error) {
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

func getListener(addr *net.UDPAddr, intfc *net.Interface) (conn *net.UDPConn, err error) {
	network := "udp4"
	if strings.Count(addr.String(), ":") > 1 {
		network = "udp6"
	}

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

	err = conn.SetReadBuffer(maxDatagramSize)
	if err != nil {
		err = fmt.Errorf("failed to SetReadBuffer to %v because %v", maxDatagramSize, err)
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

func getSender(addr *net.UDPAddr, srcAddr *net.UDPAddr) (conn *net.UDPConn, err error) {
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

func handle(listener, sender, otherSender *net.UDPConn) {
	for {
		buf := make([]byte, maxDatagramSize)
		n, src, err := listener.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(fmt.Sprintf("failed to ReadFromUDP because %v", err))

			return
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

		_, err = sender.Write([]byte(data))
		if err != nil {
			log.Fatalf("failed to Write because %v", err)

			return
		}

		log.Printf("%v -> %v - sent %v bytes to %v", src, listener.LocalAddr().String(), n, sender.RemoteAddr())
	}
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("usage: castinator [left interface] [left UDPv4/v6 address] [right interface] [right UDPv4/v6 address]")

		os.Exit(1)
	}

	leftAddr, leftIntfc, leftSrcAddr, err := getAddressesAndInterfaces(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	leftListener, err := getListener(leftAddr, leftIntfc)
	if err != nil {
		log.Fatal(err)
	}

	leftSender, err := getSender(leftAddr, leftSrcAddr)
	if err != nil {
		log.Fatal(err)
	}

	rightAddr, rightIntfc, rightSrcAddr, err := getAddressesAndInterfaces(os.Args[3], os.Args[4])
	if err != nil {
		log.Fatal(err)
	}

	rightListener, err := getListener(rightAddr, rightIntfc)
	if err != nil {
		log.Fatal(err)
	}

	rightSender, err := getSender(rightAddr, rightSrcAddr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("leftAddr = %+v\n", leftAddr)
	fmt.Printf("leftIntfc = %+v\n", leftIntfc)
	fmt.Printf("leftSrcAddr = %+v\n", leftSrcAddr)
	fmt.Printf("leftListener = %+v\n", leftListener)
	fmt.Printf("leftSender = %+v\n", leftSender)

	fmt.Println("")

	fmt.Printf("rightAddr = %+v\n", rightAddr)
	fmt.Printf("rightIntfc = %+v\n", rightIntfc)
	fmt.Printf("rightSrcAddr = %+v\n", rightSrcAddr)
	fmt.Printf("rightListener = %+v\n", rightListener)
	fmt.Printf("rightSender = %+v\n", rightSender)

	wg := sync.WaitGroup{}

	wg.Add(2)

	go handle(leftListener, rightSender, leftSender)
	go handle(rightListener, leftSender, rightSender)

	wg.Wait()
}
