package entities

import switchdata "github.com/somebottle/localsend-switch/generated/switchdata/v1"

// 发现相关实体

// SwitchMessage 包装了连接要递交的交换数据
type SwitchMessage struct {
	SourceAddr string
	Payload    *switchdata.DiscoveryMessage
}
