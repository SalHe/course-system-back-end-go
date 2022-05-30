package flags

import "flag"

var (
	ConfigPath = flag.String("config", "./config.yml", "config file path")
)

func Parse() {
	flag.Parse()
}
