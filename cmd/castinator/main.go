package main

import (
	"fmt"
	"github.com/initialed85/castinator/pkg/handler"
	"github.com/initialed85/castinator/pkg/interfaces"
	"github.com/initialed85/castinator/pkg/listener"
	"github.com/initialed85/castinator/pkg/sender"
	"log"
	"os"
	"sync"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	if len(os.Args) < 5 {
		fmt.Println("usage: castinator [left interface] [left UDPv4/v6 address] [right interface] [right UDPv4/v6 address]")

		os.Exit(1)
	}

	leftAddr, leftIntfc, leftSrcAddr, err := interfaces.GetAddressesAndInterfaces(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	leftListener, err := listener.GetListener(leftAddr, leftIntfc)
	if err != nil {
		log.Fatal(err)
	}

	leftSender, err := sender.GetSender(leftAddr, leftSrcAddr)
	if err != nil {
		log.Fatal(err)
	}

	rightAddr, rightIntfc, rightSrcAddr, err := interfaces.GetAddressesAndInterfaces(os.Args[3], os.Args[4])
	if err != nil {
		log.Fatal(err)
	}

	rightListener, err := listener.GetListener(rightAddr, rightIntfc)
	if err != nil {
		log.Fatal(err)
	}

	rightSender, err := sender.GetSender(rightAddr, rightSrcAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("leftAddr = %+v\n", leftAddr)
	log.Printf("leftIntfc = %+v\n", leftIntfc)
	log.Printf("leftSrcAddr = %+v\n", leftSrcAddr)
	log.Printf("leftListener = %+v\n", leftListener)
	log.Printf("leftSender = %+v\n", leftSender)

	fmt.Println("")

	log.Printf("rightAddr = %+v\n", rightAddr)
	log.Printf("rightIntfc = %+v\n", rightIntfc)
	log.Printf("rightSrcAddr = %+v\n", rightSrcAddr)
	log.Printf("rightListener = %+v\n", rightListener)
	log.Printf("rightSender = %+v\n", rightSender)

	wg := sync.WaitGroup{}

	wg.Add(2)

	go handler.Handle(leftListener, rightSender, leftSender)
	go handler.Handle(rightListener, leftSender, rightSender)

	wg.Wait()
}
