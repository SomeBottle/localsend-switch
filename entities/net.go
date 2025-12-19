package entities

// 网络处理相关实体

import (
	"net"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// PacketConn 封装了 IPv4 和 IPv6 的数据包连接，包括有 ReadFrom 和 Close 方法
type PacketConn struct{
	IPv4Conn *ipv4.PacketConn
	IPv6Conn *ipv6.PacketConn
}

func (pc *PacketConn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	if pc.IPv4Conn != nil {
		n, _, addr, err :=pc.IPv4Conn.ReadFrom(b)
		return n, addr, err
	}
	if pc.IPv6Conn != nil {
		n, _, addr, err := pc.IPv6Conn.ReadFrom(b)
		return n, addr, err
	}
	return 0, nil, nil
}

func (pc *PacketConn) Close() error {
	if pc.IPv4Conn != nil {
		if err := pc.IPv4Conn.Close(); err != nil {
			return err
		}
	}
	if pc.IPv6Conn != nil {
		if err := pc.IPv6Conn.Close(); err != nil {
			return err
		}
	}
	return nil
}