// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/kameshsampath/drone-gcloud-run/plugin"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// These are populated at build time
var (
	// Version is the version string at which the CLI is built.
	Version string
	// BuildDate is the date on which this CLI binary was built
	BuildDate string
	// Commit is the git commit from which this CLI binary was built.
	Commit string
	//BuiltBy is the release program that built this binary
	BuiltBy string
	//Os the Operating System for which this binary is built
	Os string
	//Arch the Architecture for which this binary is compatible
	Arch string
)

func main() {
	versionInfo := flag.Bool("version", false, "display version information")
	flag.Parse()

	if *versionInfo {
		fmt.Printf("Version:      %s\n", Version)
		fmt.Printf("Build Date:   %s\n", BuildDate)
		fmt.Printf("Git Revision: %s\n", Commit)
		fmt.Printf("Built-By: %s\n", BuiltBy)
		fmt.Printf("OS: %s\n", Os)
		fmt.Printf("Arch: %s\n", Arch)
		os.Exit(0)
	}

	logrus.SetFormatter(new(formatter))
	var args plugin.Args
	if err := envconfig.Process("", &args); err != nil {
		logrus.Fatalln(err)
	}

	switch args.Level {
	case "debug":
		logrus.SetFormatter(textFormatter)
		logrus.SetLevel(logrus.DebugLevel)
	case "trace":
		logrus.SetFormatter(textFormatter)
		logrus.SetLevel(logrus.TraceLevel)
	}

	if err := plugin.Exec(context.Background(), args); err != nil {
		logrus.Fatalln(err)
	}
}

// default formatter that writes logs without including timestamp
// or level information.
type formatter struct{}

func (*formatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message), nil
}

// text formatter that writes logs with level information
var textFormatter = &logrus.TextFormatter{
	DisableTimestamp: true,
}
