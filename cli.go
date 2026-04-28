package main

import (
	"flag"
	"fmt"
	"io"
)

func RunCLI(args []string) error {
	fs := flag.NewFlagSet("goquarium", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	info := fs.Bool("info", false, "show app info and exit")
	classic := fs.Bool("classic", false, "use classic fish set")
	version := fs.Bool("version", false, "show version and exit")
	shortV := fs.Bool("v", false, "show version and exit")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *version || *shortV {
		fmt.Println(VersionString())
		return nil
	}
	if *info {
		fmt.Print(InfoText())
		return nil
	}
	anim := NewAnimation()
	return anim.Run(SetupAquarium, *classic)
}
