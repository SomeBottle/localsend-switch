package constants

// 交换机制相关常量

const (
	// 发现信息存活时间 (包括传递过来的连接信息以及本地监听到的 UDP 发现包），单位秒
	DISCOVERY_INFO_DWELL_TIME_SECONDS = 5
	// 等候区大小，即本地停留的发现信息最大条目数，多余的会被丢弃
	LOUNGE_SIZE = 255 * 255
)
