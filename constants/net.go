package constants

// 网络处理相关常量
const (
	// LocalSend 默认的 IPv4 组播地址
	LocalSendDefaultMulticastIPv4 = "224.0.0.167"
	// LocalSend 默认的组播端口
	LocalSendDefaultMulticastPort = "53317"
	// UDP 连接读缓冲区大小
	ConnReadBufferSize = 1024 * 1024
	// 读取时字节缓冲区大小
	ReadBufferSize = 65536 // 64 KiB
	// 重试监听组播的间隔时间
	MulticastListenRetryInterval = 5 // 秒
)
