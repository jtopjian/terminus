// Copyright (c) 2015 Joe Topjian. All rights reserved.
// Copyright (c) 2014 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jtopjian/terminus/config"
)

var (
	externalFactsDir string
	httpAddr         string
	printVersion     bool
	debug            bool
)

var defaultExternalFacts = "/etc/terminus/facts.d"

func init() {
	log.SetFlags(0)
	flag.BoolVar(&debug, "debug", false, "print errors to stderr instead of ignoring them")
	flag.StringVar(&externalFactsDir, "external-facts-dir", defaultExternalFacts, "Path to external facts directory.")
	flag.StringVar(&httpAddr, "http", "", "HTTP service address (e.g., ':6060')")
	flag.BoolVar(&printVersion, "version", false, "print version and exit")
}

func main() {
	flag.Parse()

	// If there are arguments left over, use the first argument as a fact query.
	var path string
	if len(flag.Args()) > 0 {
		path = flag.Args()[0]
	}

	// Build a configuration struct with all known information
	c := config.Config{
		ExternalFactsDir: externalFactsDir,
		Path:             path,
		Debug:            debug,
	}

	if printVersion {
		fmt.Printf("terminus %s\n", Version)
		os.Exit(0)
	}

	if httpAddr != "" {
		http.Handle("/facts", httpHandler(factsHandler))
		log.Fatal(http.ListenAndServe(httpAddr, nil))
	}

	f := getFacts(c)
	facts, err := parseFacts(f, c)
	if err != nil {
		errorAndExit(err)
	}
	output, err := formatFacts(facts)
	if err != nil {
		errorAndExit(err)
	}

	fmt.Println(output)
	os.Exit(0)
}
