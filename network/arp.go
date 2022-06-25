package network

import (
	"fmt"
	"log"
	"syscall"

	"github.com/trapedev/go-tcp-ip/datalink"
	"github.com/trapedev/go-tcp-ip/util"
)

/*

【ARP Packet Format】

-------------------------------------------------------------------------
|          hardware type            |           protocol type           |
-------------------------------------------------------------------------
| hardware length | protocol length |             operation             |
-------------------------------------------------------------------------
|                          sender mac address                           |
-------------------------------------------------------------------------
|        sender mac address         |          sender ip address        |
-------------------------------------------------------------------------
|          sender ip address        |          target mac address       |
-------------------------------------------------------------------------
|                          target mac address                           |
-------------------------------------------------------------------------
|                           target ip address                           |
-------------------------------------------------------------------------

*/

type Arp struct {
	HardwareType     []byte
	ProtocolType     []byte
	HardwareLength   []byte
	ProtocolLength   []byte
	Operation        []byte
	SenderMacAddress []byte
	SenderIpAddress  []byte
	TargetMacAddress []byte
	TargetIpAddress  []byte
}

func NewArpRequest(localif util.LocalAddressInfo, targetIp string) Arp {
	return Arp{
		HardwareType:     []byte{0x00, 0x01}, // 0x0001 if ethernet
		ProtocolType:     []byte{0x08, 0x00}, // 0x0800 if IPv4
		HardwareLength:   []byte{0x06},
		ProtocolLength:   []byte{0x04},
		Operation:        []byte{0x00, 0x01},
		SenderMacAddress: localif.LocalMacAddress,
		SenderIpAddress:  localif.LocalIpAddress,
		TargetMacAddress: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		TargetIpAddress:  util.IpToByte(targetIp),
	}
}

func convertShort(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

func (*Arp) Send(ifindex int, packet []byte) Arp {
	address := syscall.SockaddrLinklayer{
		Protocol: syscall.ETH_P_ARP,
		Ifindex:  ifindex,
		Hatype:   syscall.ARPHRD_ETHER,
	}
	sendfd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(convertShort(syscall.ETH_P_ALL)))
	if err != nil {
		log.Fatalf("create sendfd err : %v\n", err)
	}
	defer syscall.Close(sendfd)

	err = syscall.Sendto(sendfd, packet, 0, &address)
	if err != nil {
		log.Fatalf("Send to err : %v", err)
	}

	for {
		recvBuf := make([]byte, 80)
		_, _, err := syscall.Recvfrom(sendfd, recvBuf, 0)
		if err != nil {
			log.Fatalf("read err : %v", err)
			if recvBuf[12] == 0x08 && recvBuf[13] == 0x06 {
				if recvBuf[20] == 0x00 && recvBuf[21] == 0x02 {
					return parseArpPacket(recvBuf[14:])
				}
			}
		}
	}
}

func parseArpPacket(packet []byte) Arp {
	return Arp{
		HardwareType:     []byte{packet[0], packet[1]},
		ProtocolType:     []byte{packet[2], packet[3]},
		HardwareLength:   []byte{packet[4]},
		ProtocolLength:   []byte{packet[5]},
		Operation:        []byte{packet[6], packet[7]},
		SenderMacAddress: []byte{packet[8], packet[9], packet[10], packet[11], packet[12], packet[13]},
		SenderIpAddress:  []byte{packet[14], packet[15], packet[16], packet[17]},
		TargetMacAddress: []byte{packet[18], packet[19], packet[20], packet[21], packet[22], packet[23]},
		TargetIpAddress:  []byte{packet[24], packet[25], packet[26], packet[27]},
	}
}

func SendArp(ifname string) {
	localif, err := util.GetLocalUIpAddress(ifname)
	if err != nil {
		log.Fatalf("could not get local ip address -> err:%v", err)
	}
	ethernet := datalink.NewEthernet([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, localif.LocalMacAddress, "ARP")
	arpReq := NewArpRequest(localif, "172.22.71.61")

	var sendArp []byte
	sendArp = append(sendArp, util.ToByteArr(ethernet)...)
	sendArp = append(sendArp, util.ToByteArr(arpReq)...)
	arpReply := arpReq.Send(localif.Index, sendArp)
	fmt.Printf("ARP Reply : %v\n", util.PrintByteArr(arpReply.SenderMacAddress))
}
