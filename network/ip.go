package network

/*

【IPv4 Datagram Format】

------------------------------------------------------------------------------------------------
| Version | Internet Header Lengh | DSCP | ECN |                Total Lengh                    |
------------------------------------------------------------------------------------------------
|              Identification                  |   Flags   |          Fragment Offset          |
------------------------------------------------------------------------------------------------
|       Time To Live        |     Protocol     |               Header Checksum                 |
------------------------------------------------------------------------------------------------
|                                       Source Address                                         |
------------------------------------------------------------------------------------------------
|                                     Destination Address                                      |
------------------------------------------------------------------------------------------------
|                           Options                          |              Padding            |
------------------------------------------------------------------------------------------------
|                                                                                              |
|                                            Data                                              |
|                                                                                              |
------------------------------------------------------------------------------------------------

*/

type ProtocolType []byte

var (
	IP  ProtocolType = []byte{0x01}
	UDP ProtocolType = []byte{0x11}
	TCP ProtocolType = []byte{0x06}
)

var ProtocolTypeMap map[string]ProtocolType = map[string]ProtocolType{
	"IP":  IP,
	"UDP": UDP,
	"TCP": TCP,
}

type IPHeader struct {
	VersionAndHeaderLength []byte
	ServiceType            []byte
	TotalPacketLength      []byte
	PacketIdentification   []byte
	Flags                  []byte
	TTL                    []byte
	Protocol               []byte
	HeaderChecksum         []byte
	SourceIPAddr           []byte
	DestinationIPAddr      []byte
}

func NewIPHeader(sourceIP, destinationIp []byte, protocol string) IPHeader {
	ip := IPHeader{
		VersionAndHeaderLength: []byte{0x45},
		ServiceType:            []byte{0x00},
		TotalPacketLength:      []byte{0x00, 0x00},
		PacketIdentification:   []byte{0x00, 0x00},
		Flags:                  []byte{0x40, 0x00},
		TTL:                    []byte{0x40},
		Protocol:               ProtocolTypeMap[protocol],
		HeaderChecksum:         []byte{0x00, 0x00},
		SourceIPAddr:           sourceIP,
		DestinationIPAddr:      destinationIp,
	}
	return ip
}
