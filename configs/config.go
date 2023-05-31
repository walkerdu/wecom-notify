package configs

import (
	"github.com/walkerdu/wecom-backend/pkg/wecom"
)

// 企业微信配置
type WeComConfig struct {
	AgentConfig wecom.AgentConfig `json:"agent_config"`
	Addr        string            `json:"addr"`
}

type Config struct {
	WeCom WeComConfig `json:"we_com"`
}
