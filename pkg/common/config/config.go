package config

var Config configStruct

type configStruct struct {
	Mysql struct {
		Address       string `yaml:"address"`
		Username      string `yaml:"username"`
		Password      string `yaml:"password"`
		Database      string `yaml:"database"`
		MaxOpenConn   int    `yaml:"maxOpenConn"`
		MaxIdleConn   int    `yaml:"maxIdleConn"`
		MaxLifeTime   int    `yaml:"maxLifeTime"`
		LogLevel      int    `yaml:"logLevel"`
		SlowThreshold int    `yaml:"slowThreshold"`
		Charset       string `yaml:"charset"`
		//字符集的排序规则或排序规则名称
		Collate string `yaml:"collate"`
	} `yaml:"mysql"`

	Redis struct {
		Address  []string `yaml:"address"`
		Username string   `yaml:"username"`
		Password string   `yaml:"password"`
	} `yaml:"redis"`

	RabbitMq struct {
		Address  string `yaml:"address"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"rabbitmq"`
}
