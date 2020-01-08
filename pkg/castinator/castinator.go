package castinator

import (
	"fmt"
	"github.com/initialed85/castinator/internal/handler"
	"github.com/initialed85/castinator/pkg/interfaces"
	"github.com/initialed85/castinator/pkg/listener"
	"github.com/initialed85/castinator/pkg/sender"
	"log"
	"net"
	"os"
)

type Castinator struct {
	leftListener, leftSender, rightListener, rightSender         *net.UDPConn
	leftListenAddr, leftSendAddr, rightListenAddr, rightSendAddr *net.UDPAddr
	leftIntfc, rightIntfc                                        *net.Interface
	leftSrcAddr, rightSrcAddr                                    *net.UDPAddr
}

func New(leftIntfcName, leftUDPListenAddr, leftUDPSendAddr, rightIntfcName, rightUDPListenAddr, rightUDPSendAddr string) (c Castinator, err error) {
	if leftIntfcName == "" {
		return c, fmt.Errorf("leftIntfcName empty")
	}

	if leftUDPListenAddr == "" {
		return c, fmt.Errorf("leftUDPListenAddr empty")
	}

	if leftUDPSendAddr == "" {
		return c, fmt.Errorf("leftUDPSendAddr empty")
	}

	if rightIntfcName == "" {
		return c, fmt.Errorf("rightIntfcName empty")
	}

	if rightUDPListenAddr == "" {
		return c, fmt.Errorf("rightUDPListenAddr empty")
	}

	if rightUDPSendAddr == "" {
		return c, fmt.Errorf("rightUDPSendAddr empty")
	}

	//
	// left listener
	//

	c.leftListenAddr, c.leftIntfc, c.leftSrcAddr, err = interfaces.GetAddressesAndInterfaces(leftIntfcName, leftUDPListenAddr)
	if err != nil {
		return c, err
	}

	c.leftListener, err = listener.GetListener(c.leftListenAddr, c.leftIntfc)
	if err != nil {
		return c, err
	}

	//
	// left sender
	//

	c.leftSendAddr, err = interfaces.GetAddress(leftUDPSendAddr)
	if err != nil {
		return c, err
	}

	c.leftSender, err = sender.GetSender(c.leftSendAddr, c.leftSrcAddr)
	if err != nil {
		return c, err
	}

	//
	// right listener
	//

	c.rightListenAddr, c.rightIntfc, c.rightSrcAddr, err = interfaces.GetAddressesAndInterfaces(rightIntfcName, rightUDPListenAddr)
	if err != nil {
		return c, err
	}

	c.rightListener, err = listener.GetListener(c.rightListenAddr, c.rightIntfc)
	if err != nil {
		return c, err
	}

	//
	// right sender
	//

	c.rightSendAddr, err = interfaces.GetAddress(rightUDPSendAddr)
	if err != nil {
		return c, err
	}

	c.rightSender, err = sender.GetSender(c.rightSendAddr, c.rightSrcAddr)
	if err != nil {
		return c, err
	}

	log.Printf("leftIntfc = %+v", c.leftIntfc)
	log.Printf("leftListenAddr = %+v", c.leftListenAddr)
	log.Printf("leftListener = %+v", c.leftListener)
	log.Printf("leftSrcAddr = %+v", c.leftSrcAddr)
	log.Printf("leftSendAddr = %+v", c.leftSendAddr)
	log.Printf("leftSender = %+v", c.leftSender)

	log.Printf("rightIntfc = %+v", c.rightIntfc)
	log.Printf("rightListenAddr = %+v", c.rightListenAddr)
	log.Printf("rightListener = %+v", c.rightListener)
	log.Printf("rightSrcAddr = %+v", c.rightSrcAddr)
	log.Printf("rightSendAddr = %+v", c.rightSendAddr)
	log.Printf("rightSender = %+v", c.rightSender)

	return c, nil
}

func (c *Castinator) Start() {
	go handler.Handle(c.leftListener, c.rightSender, c.leftSender)
	go handler.Handle(c.rightListener, c.leftSender, c.rightSender)
}

func (c *Castinator) Stop() {
	os.Exit(0)
}
