// +build linux

package main

import (
	"bufio"
	"os/exec"
	"strconv"
	"syscall"
)

type SSRClient struct {
	cmd *exec.Cmd
}

func (c *SSRClient) Start(cfg *ssrConfig, listenPort int) {
	cmd := exec.Command("python3", "shadowsocksr/shadowsocks/local.py", "-s", cfg.Host, "-p", strconv.Itoa(cfg.Port), "-k", cfg.Password, "-m", cfg.Method, "-O", cfg.Protocol, "-o", cfg.Obfs, "-G", cfg.ProtoParam, "-g", cfg.ObfsParam, "-l", strconv.Itoa(listenPort), "-v")
	cmd.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGTERM}
	c.cmd = cmd

	stdout, err := cmd.StderrPipe()
	check(err)

	err = cmd.Start()
	check(err)

	outputBuf := bufio.NewReader(stdout)
	for {
		output, _, err := outputBuf.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				check(err)
			}
		}

		logs(string(output))
	}
	logs("stop cmd")
}

func (c *SSRClient) Stop () {
	err := c.cmd.Process.Kill()
	check(err)
}

func runSSR(cfgChan chan *ssrConfig, listenPort int) {
	cfg := <-cfgChan
	for {
		client := new(SSRClient)
		go client.Start(cfg, listenPort)
		cfg = <-cfgChan
		logs("stop...")
		client.Stop()
	}
}
