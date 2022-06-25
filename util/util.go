package util

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
)

type LocalAddressInfo struct {
	LocalMacAddress []byte
	LocalIpAddress  []byte
	Index           int
}

// convert ip of string typr to byte type one.
func IpToByte(ip string) []byte {
	var convertedIp []byte
	for _, v := range strings.Split(ip, ".") {
		i, _ := strconv.ParseUint(v, 10, 8)
		convertedIp = append(convertedIp, byte(i))
	}
	return convertedIp
}

func GetLocalUIpAddress(ifname string) (LocalAddressInfo, error) {
	var localif LocalAddressInfo
	nif, err := net.InterfaceByName(ifname)
	if err != nil {
		return localif, err
	}
	localif.LocalMacAddress = nif.HardwareAddr
	localif.Index = nif.Index
	addrs, err := nif.Addrs()
	if err != nil {
		return localif, err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				localif.LocalIpAddress = ipnet.IP.To4()
			}
		}
	}
	return localif, nil
}

func ToByteArr(value any) []byte {
	rv := reflect.ValueOf(value)
	var arr []byte
	for i := 0; i < rv.NumField(); i++ {
		b := rv.Field(i).Interface().([]byte)
		arr = append(arr, b...)
	}
	return arr
}

func PrintByteArr(arr []byte) string {
	var str string
	for _, v := range arr {
		str += fmt.Sprintf("%x ", v)
	}
	return str
}
