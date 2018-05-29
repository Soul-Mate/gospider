package cmd

import "flag"

type Cmd struct {
	Path string
}

func (cmd *Cmd) config() {
	useg := "Assign your config path. Otherwise, The default value will be used."
	flag.StringVar(&cmd.Path, "config", "", useg)
}

func InitCmd() *Cmd {
	cmd := new(Cmd)
	cmd.config()
	flag.Parse()
	return cmd
}
