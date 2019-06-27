package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

type (
	// Config information.
	Config struct {
		url  string
		mode string
	}
)

var config Config

func main() {
	app := cli.NewApp()
	app.Name = "youtub-dl mp3"
	app.Usage = "TODO: Example Text"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url,u",
			Usage: "url: vidio link ",
		},
		cli.StringFlag{
			Name:  "mode,m",
			Usage: "mode: mp3, others",
		},
	}

	app.Run(os.Args)
}

func run(c *cli.Context) error {
	config = Config{
		url:  c.String("url"),
		mode: c.String("mode"),
	}

	switch config.mode {
	case "mp3":
		return mp3()
	default:
		return mp3()
	}
}

func mp3() error {
	fmt.Println("link:", config.url)
	fmt.Println("mode:", config.mode)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command("youtube-dl", "-x", "--audio-format", "mp3", "--audio-quality", "0", "-o", "~/Downloads/mp3/%(playlist)s/%(title)s.%(ext)s", config.url)

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()

	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}
	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()
	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

	return nil
}
