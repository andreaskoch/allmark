// +build ignore

// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
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
	verboseFlagIsSet      = flag.Bool("v", false, "Verbose mode")
	fmtFlagIsSet          = flag.Bool("fmt", false, "Format the source files")
	testFlagIsSet         = flag.Bool("test", false, "Execute all tests (go test")
	installFlagIsSet      = flag.Bool("install", false, "Force rebuild of everything (go install -a)")
	dependenciesFlagIsSet = flag.Bool("dependencies", false, "List all third-party dependencies")
	versionFlagIsSet      = flag.Bool("version", false, "Get the current version number of the repository")

	// working directory
	root = getWorkingDirectory()

	// packages to build
	buildPackages = []string{
		fmt.Sprintf("%s/cmd/allmark", ProjectNamespace),
	}

	// a regular expression matching all non-standard go library packages (e.g. github.com/..., allmark.io/... )
	nonStandardPackagePattern = regexp.MustCompile(`^\w+[\.-].+/`)

	// The git version pattern.
	gitVersionPattern = regexp.MustCompile(`\b\d\d\d\d-\d\d-\d\d-[0-9a-f]{7,7}\b`)
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	if *fmtFlagIsSet {
		format()
		return
	}

	if *testFlagIsSet {
		executeTests()
		return
	}

	if *installFlagIsSet {
		install()
		return
	}

	if *dependenciesFlagIsSet {
		listDependencies()
		return
	}

	if *versionFlagIsSet {
		printProjectVersionNumber()
		return
	}

	flag.PrintDefaults()
}

// Install all parts of allmark using go install.
func install() {

	for _, buildPackagage := range buildPackages {
		runGoCommand(os.Stdout, os.Stderr, root, "install", buildPackagage)
	}

}

// Execute all test in the current project.
func executeTests() {
	packages := getPackagesWithTests()

	for index, packageName := range packages {

		fmt.Printf("Testing package %02d of %v: %s\n", index+1, len(packages), packageName)
		runGoCommand(os.Stdout, os.Stderr, root, "test", packageName)
	}
}

// Format all packages of allmark using go fmt.
func format() {
	packages := getInternalPackages()

	for index, packageName := range packages {

		fmt.Printf("Formatting package %02d of %v: %s\n", index+1, len(packages), packageName)
		runGoCommand(os.Stdout, os.Stderr, root, "fmt", packageName)

	}
}

// List all third-party packages that allmark depends on.
func listDependencies() {
	thirdPartyPackages := getThirdPartyPackages()

	for _, dependency := range thirdPartyPackages {
		fmt.Println(dependency)
	}
}

// Print the current version number of the project.
func printProjectVersionNumber() {
	fmt.Println(gitVersion())
}

// Get all packages which tests in them.
func getPackagesWithTests() []string {
	isInternalPackageWithTests := func(packageName string) bool {
		isInternalPackage := strings.HasPrefix(packageName, ProjectNamespace)
		if !isInternalPackage {
			return false
		}

		return packageHasTests(packageName)
	}

	internalPackagesWithTests := getAllNonStandardLibraryPackages(isInternalPackageWithTests)
	return internalPackagesWithTests
}

// Get all internal packages used in this project.
func getInternalPackages() []string {

	isInternalPackage := func(packageName string) bool {
		return strings.HasPrefix(packageName, ProjectNamespace)
	}

	internalPackages := getAllNonStandardLibraryPackages(isInternalPackage)
	return internalPackages
}

// Get all third party packages used in this project.
func getThirdPartyPackages() []string {

	isThirdPartyPackage := func(packageName string) bool {
		return !strings.HasPrefix(packageName, ProjectNamespace)
	}

	thirdPartyPackages := getAllNonStandardLibraryPackages(isThirdPartyPackage)
	return thirdPartyPackages
}

// Get a sorted and unique list of all non-standard library packages used in this project that meet the supplied expression.
func getAllNonStandardLibraryPackages(inclusionExpression func(packageName string) bool) []string {

	// get all dependent packages (will include duplicates and standard library packages)
	allDependentPackages := make([]string, 0)
	for _, buildPackage := range buildPackages {
		output := new(bytes.Buffer)
		errors := new(bytes.Buffer)
		runGoCommand(output, errors, root, "list", "-f", `'{{ join .Deps ","}}'`, buildPackage)

		allDependentPackages = append(allDependentPackages, strings.Split(output.String(), ",")...)
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

// Check whether the package with the supplied name has tests or not.
func packageHasTests(packageName string) bool {
	packagePath := getPackagePathByName(packageName)
	testFilePattern := filepath.Join(packagePath, "*_test.go")
	matches, err := filepath.Glob(testFilePattern)
	if err != nil {
		log.Fatalf("Unable to find test files in %q. Error: %s", testFilePattern, err.Error())
	}

	packageContainsTestFiles := len(matches) > 0
	return packageContainsTestFiles
}

// Get the path of package from its name.
func getPackagePathByName(packageName string) string {
	return filepath.Join(root, "src", packageName)
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
func runGoCommand(stdout, stderr io.Writer, goPath string, args ...string) {

	commandName := "go"
	cmdName := fmt.Sprintf("%s %s", commandName, strings.Join(args, " "))

	// set the go path
	cmd := exec.Command(commandName, args...)
	cmd.Dir = goPath

	cmd.Env = cleanGoEnv()
	cmd.Env = setEnv(cmd.Env, GOPATH, goPath)
	cmd.Env = setEnv(cmd.Env, GOBIN, filepath.Join(goPath, "bin"))

	// execute the command
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if *verboseFlagIsSet {
		log.Printf("Running %s", cmdName)
	}

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error running %s: %v", cmdName, err)
	}
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

// gitVersion returns the git version of the git repo at root as a
// string of the form "yyyy-mm-dd-xxxxxxx", with an optional trailing
// '+' if there are any local uncomitted modifications to the tree.
func gitVersion() string {
	cmd := exec.Command("git", "rev-list", "--max-count=1", "--pretty=format:'%ad-%h'", "--date=short", "HEAD")
	cmd.Dir = root

	commandOutput, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error running git rev-list in %s: %v", root, err)
	}

	versionNumber := strings.TrimSpace(string(commandOutput))
	if m := gitVersionPattern.FindStringSubmatch(versionNumber); m != nil {
		versionNumber = m[0]
	} else {
		log.Fatalf("Failed to find git version in string %q", versionNumber)
	}

	cmd = exec.Command("git", "diff", "--exit-code")
	cmd.Dir = root
	if err := cmd.Run(); err != nil {
		versionNumber += "+"
	}

	return versionNumber
}
