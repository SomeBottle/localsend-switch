package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/somebottle/localsend-switchboard/constants"
	"github.com/somebottle/localsend-switchboard/entities"
	"github.com/somebottle/localsend-switchboard/services"
	"github.com/somebottle/localsend-switchboard/utils"
)

func main() {
	// 中断信号处理
	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// ------------ 先读取配置
	multicastAddr := os.Getenv("LOCALSEND_MULTICAST_ADDR")
	multicastPort := os.Getenv("LOCALSEND_MULTICAST_PORT")

	// 尝试从命令行读取配置
	flag.StringVar(&multicastAddr, "addr", multicastAddr, "Multicast address")
	flag.StringVar(&multicastPort, "port", multicastPort, "Multicast port")

	flag.Parse()

	// 没有配置就用默认值
	if multicastAddr == "" {
		multicastAddr = constants.LocalSendDefaultMulticastIPv4
		fmt.Println("Multicast address not provided, using default value: ", multicastAddr)
	}

	if multicastPort == "" {
		multicastPort = constants.LocalSendDefaultMulticastPort
		fmt.Println("Multicast port not provided, using default value: ", multicastPort)
	}

	// 检查是否为 IPv6 地址
	isIpv6, err := utils.IsIpv6(multicastAddr)
	if err != nil {
		fmt.Printf("Error parsing IP address: %v\n", err)
		return
	}
	fmt.Printf("Is IPv6: %v\n", isIpv6)
	// 获得首选出站 IP 地址
	selfIp, err := utils.GetOutboundIP()
	if err != nil {
		fmt.Printf("Error getting outbound IP address: %v\n", err)
		return
	}
	// 获得相应的网络接口
	outBoundInterface, err := utils.GetInterfaceByIP(selfIp)
	if err != nil {
		fmt.Printf("Error getting outbound network interface: %v\n", err)
		return
	}
	if outBoundInterface == nil {
		fmt.Printf("No network interface found for IP address: %s\n", selfIp.String())
		return
	}

	fmt.Printf("Outbound IP address: %s\n", selfIp.String())
	fmt.Printf("Using network interface: %s\n", outBoundInterface.Name)

	var network string
	if isIpv6 {
		network = "udp6"
	} else {
		network = "udp4"
	}

	// 数据传输通道
	udpPacketChan := make(chan entities.UDPPacketData)
	// ------------ 加入组播组，接收 LocalSend 的发现 UDP 包
	go services.ListenLocalSendMulticast(network, multicastAddr, multicastPort, outBoundInterface, sigCtx, udpPacketChan)

	// 测试接收数据
	for {
		select {
		case <-sigCtx.Done():
			fmt.Println("Shutting down gracefully...")
			return
		case udpPacket := <-udpPacketChan:
			fmt.Printf("Received UDP packet from %s:%d - Data: %s\n", udpPacket.SourceIP.String(), udpPacket.SourcePort, string(udpPacket.Data))
		}
	}

}
