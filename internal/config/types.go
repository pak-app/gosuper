package config

type Config struct {
	Supervisor SupervisorConfig         `yaml:"supervisor"`
	Services   map[string]ServiceConfig `yaml:"services"`
}

type SupervisorConfig struct {
	Name         string            `yaml:"name"`
	LogDir       string            `yaml:"log_dir"`
	RestartDelay string            `yaml:"restart_delay"`
	StopTimeout  string            `yaml:"stop_timeout"`
	Env          map[string]string `yaml:"env"`
}

type ServiceConfig struct {
	Command       []string          `yaml:"command"`
	Dir           string            `yaml:"dir"`
	RestartWindow string            `yaml:"restart_window"`
	Stdout        string            `yaml:"stdout"`
	Stderr        string            `yaml:"stderr"`
	Env           map[string]string `yaml:"env"`

	RestartLimit int  `yaml:"restart_limit"`
	Autostart    bool `yaml:"autostart"`
	Autorestart  bool `yaml:"autorestart"`
}
