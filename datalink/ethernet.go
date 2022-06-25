package datalink

type EthType []byte

var (
	IPv4 EthType = []byte{0x08, 0x00}
	ARP  EthType = []byte{0x08, 0x06}
	IPv6 EthType = []byte{0x86, 0xdd}
)

var EthTypeMap map[string]EthType = map[string]EthType{
	"IPv4": IPv4,
	"ARP":  ARP,
	"IPv6": IPv6,
}

type EthernetFrame struct {
	DstMacAddr []byte
	SrcMacAddr []byte
	Type       []byte
}

// constractor for ethernetframe
func NewEthernet(dstMacaddr, srcMacAddr []byte, ethType string) EthernetFrame {
	ethernet := EthernetFrame{
		DstMacAddr: dstMacaddr,
		SrcMacAddr: srcMacAddr,
		Type:       EthTypeMap[ethType],
	}
	return ethernet
}
