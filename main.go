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
    --help         show help (default: false)
    --version      print the version (default: false)

EXAMPLE:
    $ daemon --help
    $ daemon --version
    $ daemon node server.js
	`)
}

const VERSION = "v1.0.0"

/*
Usage:
	$ deamon [option] [command]
*/
func main() {
	var (
		help    bool
		version bool
	)

	flag.BoolVar(&help, "help", false, "print help information")
	flag.BoolVar(&version, "version", false, "print version information")

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
		if err := cmd.Start(); err != nil {
			log.Fatalln(err)
		} else {
			os.Stderr.Write([]byte(fmt.Sprintf("[daemon]: isolated process pid '%d'\n", cmd.Process.Pid)))
		}

	}
}
