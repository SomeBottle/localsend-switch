package entities

import "net"

// UDPPacketData 表示一个 UDP 数据包的数据内容及其来源信息
type UDPPacketData struct {
	SourceIP   net.IP
	SourcePort int
	Data       []byte
}
