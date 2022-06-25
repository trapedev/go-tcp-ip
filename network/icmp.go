package network

import (
	"fmt"
	"log"
	"syscall"

	"github.com/trapedev/go-tcp-ip/datalink"
	"github.com/trapedev/go-tcp-ip/util"
)

type ICMP struct {
	Type           []byte
	Code           []byte
	CheckSum       []byte
	Identification []byte
	SequenceNumber []byte
	Data           []byte
}

func NewICMP() ICMP {
	icmp := ICMP{
		Type:           []byte{0x08},
		Code:           []byte{0x00},
		CheckSum:       []byte{0x00, 0x00},
		Identification: []byte{0x00, 0x10},
		SequenceNumber: []byte{0x00, 0x01},
		Data:           []byte{0x01, 0x02},
	}
	icmpSum := util.SumByteArr(util.ToByteArr(icmp))
	icmp.CheckSum = util.Checksum(icmpSum)
	return icmp
}

func (*ICMP) Send(ifindex int, packet []byte) ICMP {
	address := syscall.SockaddrLinklayer{
		Protocol: syscall.ETH_P_IP,
		Ifindex:  ifindex,
	}

	sendfd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(convertShort(syscall.ETH_P_ALL)))
	if err != nil {
		log.Fatalf("create icmp sendfd err : %v\n", err)
	}
	err = syscall.Sendto(sendfd, packet, 0, &address)
	if err != nil {
		log.Fatalf("Send to err : %v\n", err)
	}
	fmt.Println("send icmp packet")

	for {
		recvBuf := make([]byte, 1500)
		_, _, err := syscall.Recvfrom(sendfd, recvBuf, 0)
		if err != nil {
			log.Fatalf("read err : %v", err)
		}
		if recvBuf[23] == 0x01 {
			return parseICMP(recvBuf[34:])
		}
	}
}

func parseICMP(packet []byte) ICMP {
	return ICMP{
		Type:           []byte{packet[0]},
		Code:           []byte{packet[1]},
		CheckSum:       []byte{packet[2], packet[3]},
		Identification: []byte{packet[4], packet[5]},
		SequenceNumber: []byte{packet[6], packet[7]},
		Data:           packet[8:],
	}
}

func SendArpICMP(destinationIp string, ifname string) {
	localif, err := util.GetLocalUIpAddress(ifname)
	if err != nil {
		log.Fatalf("get local ip address: err -> %v", err)
	}
	// create arp packet
	arp := NewArpRequest(localif, destinationIp)
	var sendArp []byte

	// ccreate ethernet packet
	sendArp = append(sendArp, util.ToByteArr(datalink.NewEthernet([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, localif.LocalMacAddress, "ARP"))...)
	sendArp = append(sendArp, util.ToByteArr(arp)...)
	// send arp
	arpreply := arp.Send(localif.Index, sendArp)
	fmt.Printf("ARP Reply : %s\n", util.PrintByteArr(arpreply.SenderMacAddress))

	var sendIcmp []byte
	// create icmp packet
	icmpPacket := NewICMP()
	// create ip header
	header := NewIPHeader(localif.LocalIpAddress, util.IpToByte(destinationIp), "IP")
	header.TotalPacketLength = util.UintTo2byte(util.ToByteLen(header) + util.ToByteLen(icmpPacket))
	// calc checksum
	ipsum := util.SumByteArr(util.ToByteArr(header))
	header.HeaderChecksum = util.Checksum(ipsum)

	sendIcmp = append(sendIcmp, util.ToByteArr(datalink.NewEthernet(arpreply.SenderMacAddress, localif.LocalMacAddress, "IPv4"))...)
	sendIcmp = append(sendIcmp, util.ToByteArr(header)...)
	sendIcmp = append(sendIcmp, util.ToByteArr(icmpPacket)...)
	// send icmp packet
	icmpreply := icmpPacket.Send(localif.Index, sendIcmp)
	if icmpreply.Type[0] == 0 {
		fmt.Printf("ICMP Reply is %d, OK!\n", icmpreply.Type[0])
	}
}
