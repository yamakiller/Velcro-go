package configs

type Config struct {
	LAddr             string `yaml:"laddr"`              // 监听地址
	VAddr             string `yaml:"vaddr"`              // 本服务虚拟地址
	Router            string `yaml:"router-uri"`         // 路由文件地址
	EncryptionEnabled bool   `yaml:"encryption-enabled"` // 是否开启通信加密
	OnlineOfNumber    int    `yaml:"online-number"`      // 在线人数限制
	RequestTimeout    int32  `yaml:"request-timeout-limit"`
}
