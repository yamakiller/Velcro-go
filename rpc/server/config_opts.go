package server

// ConnConfigOption 是一个配置rpc connector 的函数
type ConnConfigOption func(config *ConnConfig)

func Configure(options ...ConnConfigOption) *ConnConfig {
	config := defaultConfig()
	for _, option := range options {
		option(config)
	}

	return config
}

// WithKleepalive 设置心跳时间
func WithKleepalive(kleepalive int32) ConnConfigOption {

	return func(config *ConnConfig) {
		config.Kleepalive = kleepalive
	}
}

func WithPool(f RpcPool) ConnConfigOption {
	return func(config *ConnConfig) {
		config.Pool = f
	}
}
