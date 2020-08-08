package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func printHelp() {
	fmt.Println(`daemon - a cli tool for run a process as daemon

USAGE:
    $ daemon [option] [command]

OPTIONS:
    --help         	show help (default: false)
    --version      	print the version (default: false)
    --stdout				Specify the output file of stdout (default: nil)
    --stderr				Specify the output file of stderr (default: nil)

EXAMPLE:
    $ daemon --help
    $ daemon --version
    $ daemon node server.js
	`)
}

func ensureFile(file string) (*os.File, error) {
	if _, err := os.Stat(file); err != nil {
		if os.ErrNotExist == err {
			return os.Create(file)
		} else {
			return nil, err
		}
	} else {
		return os.Open(file)
	}
}

const VERSION = "v1.0.0"

/*
Usage:
	$ deamon [option] [command]
*/
func main() {
	var (
		help       bool
		version    bool
		stdoutFile string
		stderrFile string
	)

	flag.BoolVar(&help, "help", false, "print help information")
	flag.BoolVar(&version, "version", false, "print version information")
	flag.StringVar(&stdoutFile, "stdout", "", "the output file of stdout")
	flag.StringVar(&stderrFile, "stderr", "", "the output file of stderr")

	flag.Parse()

	if help {
		printHelp()
		os.Exit(0)
	}

	if version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if len(flag.Args()) < 1 {
		printHelp()
		os.Exit(1)
	}

	// $ daemon [command]
	if len(os.Args) < 2 {
		log.Fatalln("require 2 argument")
	}

	if os.Getppid() != 1 {
		filePath, _ := filepath.Abs(os.Args[0])
		args := flag.Args()

		cmd := exec.Command(filePath, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			log.Fatalln(err)
		}
		return
	} else {
		args := flag.Args()

		cmd := exec.Command(args[0], args[1:]...)

		if stdoutFile != "" {
			if file, err := ensureFile(stderrFile); err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else {
				cmd.Stdout = file
			}
		}

		if stderrFile != "" {
			if file, err := ensureFile(stderrFile); err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else {
				cmd.Stderr = file
			}
		}

		if err := cmd.Start(); err != nil {
			log.Fatalln(err)
		} else {
			os.Stderr.Write([]byte(fmt.Sprintf("[daemon]: isolated process pid '%d'\n", cmd.Process.Pid)))
		}

	}
}
