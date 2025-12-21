package constants

// 交换机制相关常量

const (
	// 发现信息存活时间 (包括传递过来的连接信息以及本地监听到的 UDP 发现包），单位秒
	DISCOVERY_INFO_DWELL_TIME_SECONDS = 5  
	// 分配给发现信息存储的最大内存空间，多余的会丢弃（按理说不太可能达到），单位字节
	LOUNGE_SIZE         = 50 * 1024 * 1024 
)