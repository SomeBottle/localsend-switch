package constants

// 网络处理相关常量
const (
	// LocalSend 默认的 IPv4 组播地址
	LocalSendDefaultMulticastIPv4 = "224.0.0.167"
	// LocalSend 默认的组播端口
	LocalSendDefaultMulticastPort = "53317"
	// 组播数据读取时字节缓冲区大小
	MulticastReadBufferSize = 65536 // 64 KiB
	// 组播数据读取超时时间
	MulticastReadTimeout = 2 // 秒
	// 重试监听组播的间隔时间
	MulticastListenRetryInterval = 3 // 秒
	// TCP 最大连接数
	MaxTCPConnections = 255 * 255
	// TCP 心跳时间
	TCPHeartbeatInterval = 30 // 秒
	// TCP 读取超时时间，防止连接过长时间阻塞
	TCPReadTimeout = 2 // 秒
	// TCP 服务重启间隔时间
	TCPServerRestartInterval = 3 // 秒
	// 读取 TCP 数据时字节缓冲区大小
	TCPSocketReadBufferSize = 1024 * 1024 // 1 MiB
)
