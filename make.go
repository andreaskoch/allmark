// +build ignore

// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// The main project namespace
const ProjectNamespace = "allmark.io"

// GOPATH environment variable name
const GOPATH = "GOPATH"

// GOBIN environment variable name
const GOBIN = "GOBIN"

var (

	// command line flags
	verboseFlagIsSet = flag.Bool("v", false, "Verbose mode")
	installFlagIsSet = flag.Bool("install", false, "Force rebuild of everything (go install -a)")
	fmtFlagIsSet     = flag.Bool("fmt", false, "Format the source files")

	// working directory
	root = getWorkingDirectory()

	// packages to build
	buildPackages = []string{
		fmt.Sprintf("%s/cmd/allmark", ProjectNamespace),
	}

	// a regular expression matching all non-standard go library packages (e.g. github.com/..., allmark.io/... )
	nonStandardPackagePattern = regexp.MustCompile(`^\w+[\.-].+/`)
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	if *fmtFlagIsSet {
		format()
		return
	}

	if *installFlagIsSet {
		install()
		return
	}

	flag.PrintDefaults()
}

// Install all parts of allmark using go install.
func install() {

	for _, buildPackagage := range buildPackages {
		fmt.Println(runGoCommand(root, "install", buildPackagage))
	}

}

// Format all packages of allmark using go fmt.
func format() {
	packages := getInternalPackages()

	for index, packageName := range packages {

		fmt.Printf("Formatting package %02d of %v: %s\n", index+1, len(packages), packageName)
		output := runGoCommand(root, "fmt", packageName)
		if output != "" {
			fmt.Println(output)
		}

		index++
	}
}

// Get all internal packages used in this project.
func getInternalPackages() []string {

	isInternalPackage := func(packageName string) bool {
		return strings.HasPrefix(packageName, ProjectNamespace)
	}

	internalPackages := getAllNonStandardLibraryPackages(isInternalPackage)
	return internalPackages
}

// Get a sorted and unique list of all non-standard library packages used in this project that meet the supplied expression.
func getAllNonStandardLibraryPackages(inclusionExpression func(packageName string) bool) []string {

	// get all dependent packages (will include duplicates and standard library packages)
	allDependentPackages := make([]string, 0)
	for _, buildPackage := range buildPackages {
		output := runGoCommand(root, "list", "-f", `'{{ join .Deps ","}}'`, buildPackage)
		allDependentPackages = append(allDependentPackages, strings.Split(output, ",")...)
	}

	// sort the list
	sort.Strings(allDependentPackages)

	// unique
	packages := make([]string, 0)
	uniquePackages := make(map[string]int)

	for _, dep := range allDependentPackages {

		// skip packages we have already seen
		if _, exists := uniquePackages[dep]; exists {

			// increment
			uniquePackages[dep] = uniquePackages[dep] + 1
			continue
		}

		// skip standard library packages
		if isStandardLibraryPackage(dep) {
			continue
		}

		// skip all packages that don't meet the expression
		if !inclusionExpression(dep) {
			continue
		}

		packages = append(packages, dep)
		uniquePackages[dep] = uniquePackages[dep] + 1
	}

	return packages

}

// Check if the supplied package name is a standard library package.
func isStandardLibraryPackage(packageName string) bool {
	return nonStandardPackagePattern.MatchString(packageName) == false
}

// getWorkingDirectory returns the current working directory path or fails.
func getWorkingDirectory() string {
	root, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	return root
}

// Execute go in the specified go path with the supplied command arguments.
// Stdout will be returned as the output.
func runGoCommand(goPath string, args ...string) (output string) {

	commandName := "go"
	cmdName := fmt.Sprintf("%s %s", commandName, strings.Join(args, " "))

	// set the go path
	cmd := exec.Command(commandName, args...)
	cmd.Env = cleanGoEnv()
	cmd.Env = setEnv(cmd.Env, GOPATH, goPath)
	cmd.Env = setEnv(cmd.Env, GOBIN, filepath.Join(goPath, "bin"))

	// execute the command
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	if *verboseFlagIsSet {
		log.Printf("Running %s", cmdName)
	}

	outputBytes, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error running %s: %v", cmdName, err)
	}

	return string(outputBytes)
}

// cleanGoEnv returns a copy of the current environment with GOPATH and GOBIN removed.
func cleanGoEnv() (clean []string) {
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, GOPATH+"=") || strings.HasPrefix(env, GOBIN+"=") {
			continue
		}

		clean = append(clean, env)
	}

	return
}

// setEnv sets the given key & value in the provided environment.
// Each value in the env list should be of the form key=value.
func setEnv(env []string, key, value string) []string {
	for i, s := range env {
		if strings.HasPrefix(s, fmt.Sprintf("%s=", key)) {
			env[i] = envPair(key, value)
			return env
		}
	}
	env = append(env, envPair(key, value))
	return env
}

// Create an environment variable of the form key=value.
func envPair(key, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}
