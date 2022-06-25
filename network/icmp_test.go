package network

import "testing"

func TestSendArpICMP(t *testing.T) {
	SendArpICMP("172.22.79.120", "eth0")
}
