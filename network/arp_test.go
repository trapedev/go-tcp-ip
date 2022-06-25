package network

import "testing"

func TestSendArp(t *testing.T) {
	SendArp("eth0")
}
