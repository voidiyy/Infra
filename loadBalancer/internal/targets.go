package internal

type TargetGroup struct {
	Name            string  `yaml:"name"`
	HandlePath      string  `yaml:"handlePath"`
	BalanceStrategy string  `yaml:"balanceStrategy"`
	Targets         *Target `yaml:"targets"`
}

type Target struct {
	URLs []string `yaml:"urls"`

	Balancer BalanceStrategy `yaml:"-"`
}
