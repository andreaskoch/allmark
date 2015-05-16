// +build ignore

// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program builds allmark.
//
// $ go run make.go -install
//
// View the README.md for further details.
//
// The output binaries go into the ./bin/ directory (under the goPath, where make.go is)
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

// GOPATH environment variable name.
const GOPATH_ENVIRONMENT_VARIABLE = "GOPATH"

// GOBIN environment variable name
const GOBIN_ENVIRONMENT_VARIBALE = "GOBIN"

// GOOS environment variable name
const GOOS_ENVIRONMENT_VARIBALE = "GOOS"

// GOARCH environment variable name
const GOARCH_ENVIRONMENT_VARIBALE = "GOARCH"

var (

	// command line flags
	verboseFlagIsSet                = flag.Bool("v", false, "Verbose mode")
	fmtFlagIsSet                    = flag.Bool("fmt", false, "Format the source files")
	testFlagIsSet                   = flag.Bool("test", false, "Execute all tests (go test")
	installFlagIsSet                = flag.Bool("install", false, "Force rebuild of everything (go install -a)")
	crossCompileFlagIsSet           = flag.Bool("crosscompile", false, "Cross-compile everything (Won't work if you don't have your golang installation prepared for this)")
	crossCompileWithDockerFlagIsSet = flag.Bool("crosscompile-with-docker", false, "Cross-compile everything using docker (Won't work if you don't have docker installed)")
	listDependenciesFlagIsSet       = flag.Bool("list-dependencies", false, "List all third-party dependencies")
	updateDependenciesFlagIsSet     = flag.Bool("update-dependencies", false, "Update all third-party dependencies")
	versionFlagIsSet                = flag.Bool("version", false, "Get the current version number of the repository")

	// The GOPATH for the current project
	goPath = getWorkingDirectory()

	// The GOBIN for the current project
	goBin = filepath.Join(goPath, "bin")

	// packages to build
	buildPackages = []string{
		fmt.Sprintf("%s/cmd/allmark", ProjectNamespace),
	}

	// a regular expression matching all non-standard go library packages (e.g. github.com/..., allmark.io/... )
	nonStandardPackagePattern = regexp.MustCompile(`^\w+[\.-].+/`)

	// The git version pattern.
	gitVersionPattern = regexp.MustCompile(`\b\d\d\d\d-\d\d-\d\d-[0-9a-f]{7,7}\b`)

	// A list of all supported compilation targets (e.g. "windows/386")
	compilationTargets = []compilationTarget{
		compilationTarget{"darwin", "386"},
		compilationTarget{"darwin", "amd64"},
		compilationTarget{"dragonfly", "386"},
		compilationTarget{"dragonfly", "amd64"},
		compilationTarget{"freebsd", "386"},
		compilationTarget{"freebsd", "amd64"},
		compilationTarget{"freebsd", "arm"},
		compilationTarget{"linux", "386"},
		compilationTarget{"linux", "amd64"},
		compilationTarget{"linux", "arm"},
		compilationTarget{"nacl", "386"},
		compilationTarget{"nacl", "amd64p32"},
		compilationTarget{"nacl", "arm"},
		compilationTarget{"netbsd", "386"},
		compilationTarget{"netbsd", "amd64"},
		compilationTarget{"netbsd", "arm"},
		compilationTarget{"openbsd", "386"},
		compilationTarget{"openbsd", "amd64"},
		compilationTarget{"solaris", "amd64"},
		compilationTarget{"windows", "386"},
		compilationTarget{"windows", "amd64"},
	}

	// The current version number (e.g. "2015-01-11-284c030+")
	version = gitVersion()
)

// Compilation Target Definition
type compilationTarget struct {
	OperatingSystem string
	Architecture    string
}

func (target *compilationTarget) String() string {
	return fmt.Sprintf("%s/%s", target.OperatingSystem, target.Architecture)
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	if *verboseFlagIsSet {
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("%s: %s\n", GOPATH_ENVIRONMENT_VARIABLE, goPath)
		fmt.Printf("%s: %s\n", GOBIN_ENVIRONMENT_VARIBALE, goBin)
	}

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

	if *crossCompileFlagIsSet {
		crossCompile()
		return
	}

	if *crossCompileWithDockerFlagIsSet {
		crossCompileWithDocker()
		return
	}

	if *listDependenciesFlagIsSet {
		listDependencies()
		return
	}

	if *updateDependenciesFlagIsSet {
		updateDependencies()
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

	// prepare the environment variables
	environmentVariables := cleanGoEnv()
	environmentVariables = setEnv(environmentVariables, GOPATH_ENVIRONMENT_VARIABLE, goPath)
	environmentVariables = setEnv(environmentVariables, GOBIN_ENVIRONMENT_VARIBALE, goBin)

	for _, packageName := range buildPackages {
		runCommand(os.Stdout, os.Stderr, goPath, environmentVariables, "go", "install", getBuildVersionFlag(), packageName)
	}

}

// Cross-compile all parts of allmark for all supported platforms.
func crossCompile() {

	// prepare the environment variables
	environmentVariables := cleanGoEnv()
	environmentVariables = setEnv(environmentVariables, GOPATH_ENVIRONMENT_VARIABLE, goPath)
	environmentVariables = setEnv(environmentVariables, GOBIN_ENVIRONMENT_VARIBALE, goBin)

	// iterate over all supported compilation targets
	for _, target := range compilationTargets {

		// iterate over all build packages
		for _, packageName := range buildPackages {

			// prepare environment variables for cross-compilation
			crossCompileEnvironemntVariables := environmentVariables
			crossCompileEnvironemntVariables = setEnv(crossCompileEnvironemntVariables, GOOS_ENVIRONMENT_VARIBALE, target.OperatingSystem)
			crossCompileEnvironemntVariables = setEnv(crossCompileEnvironemntVariables, GOARCH_ENVIRONMENT_VARIBALE, target.Architecture)

			// build the package for the specified os and arch
			fmt.Printf("Compiling %s for %s\n", packageName, target.String())
			runCommand(os.Stdout, os.Stderr, goPath, crossCompileEnvironemntVariables, "go", "install", getBuildVersionFlag(), packageName)
		}
	}

}

// Cross-compile all parts of allmark for all supported platforms using docker.
func crossCompileWithDocker() {
	dockerImageName := "golang:1.4.2-cross"
	projectPathInDocker := "/usr/src/allmark"
	binPath := filepath.Join(projectPathInDocker, "bin")
	command := `docker`

	// prepare the environment variables
	environmentVariables := cleanGoEnv()

	// iterate over all supported compilation targets
	for _, target := range compilationTargets {

		// iterate over all build packages
		for _, packageName := range buildPackages {

			// assemble the build command
			args := []string{
				"run",
				"--rm",
				"-v=" + fmt.Sprintf(`%s:%s`, goPath, projectPathInDocker),
				`-w=` + projectPathInDocker,
				"-e=" + envPair(GOOS_ENVIRONMENT_VARIBALE, target.OperatingSystem),
				"-e=" + envPair(GOARCH_ENVIRONMENT_VARIBALE, target.Architecture),
				"-e=" + envPair(GOPATH_ENVIRONMENT_VARIABLE, projectPathInDocker),
				"-e=" + envPair(GOBIN_ENVIRONMENT_VARIBALE, binPath),
				dockerImageName,
				"go",
				"install",
				getBuildVersionFlag(),
				packageName,
			}

			// build the package for the specified os and arch
			fmt.Printf("Compiling %s for %s\n", packageName, target.String())
			runCommand(os.Stdout, os.Stderr, goPath, environmentVariables, command, args...)
		}
	}

}

// Execute all test in the current project.
func executeTests() {
	packages := getPackagesWithTests()

	// prepare the environment variables
	environmentVariables := cleanGoEnv()
	environmentVariables = setEnv(environmentVariables, GOPATH_ENVIRONMENT_VARIABLE, goPath)
	environmentVariables = setEnv(environmentVariables, GOBIN_ENVIRONMENT_VARIBALE, goBin)

	for index, packageName := range packages {

		fmt.Printf("Testing package %02d of %v: %s\n", index+1, len(packages), packageName)
		runCommand(os.Stdout, os.Stderr, goPath, environmentVariables, "go", "test", packageName)
	}
}

// Format all packages of allmark using go fmt.
func format() {
	packages := getInternalPackages()

	// prepare the environment variables
	environmentVariables := cleanGoEnv()
	environmentVariables = setEnv(environmentVariables, GOPATH_ENVIRONMENT_VARIABLE, goPath)
	environmentVariables = setEnv(environmentVariables, GOBIN_ENVIRONMENT_VARIBALE, goBin)

	for index, packageName := range packages {

		fmt.Printf("Formatting package %02d of %v: %s\n", index+1, len(packages), packageName)
		runCommand(os.Stdout, os.Stderr, goPath, environmentVariables, "go", "fmt", packageName)

	}
}

// List all third-party packages that allmark depends on.
func listDependencies() {
	thirdPartyPackages := getThirdPartyPackages()

	for _, dependency := range thirdPartyPackages {
		fmt.Println(dependency)
	}
}

// Update all third-party packages that allmark depends on.
func updateDependencies() {

	// prepare the environment variables
	environmentVariables := cleanGoEnv()
	environmentVariables = setEnv(environmentVariables, GOPATH_ENVIRONMENT_VARIABLE, goPath)
	environmentVariables = setEnv(environmentVariables, GOBIN_ENVIRONMENT_VARIBALE, goBin)

	// Remove all existing third party packages
	for _, namespace := range getTopNamespacesOfExternalDependencies() {

		packagePath := getPackagePathByName(namespace)
		if err := os.RemoveAll(packagePath); err != nil {
			panic(err)
		}

	}

	// Download the packages
	thirdPartyPackages := getThirdPartyPackages()
	for index, dependency := range thirdPartyPackages {

		fmt.Printf("Updating package %02d of %v: %s\n", index+1, len(thirdPartyPackages), dependency)
		runCommand(os.Stdout, os.Stderr, goPath, environmentVariables, "go", "get", dependency)

	}

	// Remove any source control folders
	for _, namespace := range getTopNamespacesOfExternalDependencies() {

		externalPackagePath := getPackagePathByName(namespace)
		versionControlFolders := getSourceControlFolderPaths(externalPackagePath)

		for _, gitFolder := range versionControlFolders {
			if err := os.RemoveAll(gitFolder); err != nil && !os.IsExist(err) {
				panic(err)
			}
		}

	}
}

// Print the current version number of the project.
func printProjectVersionNumber() {
	fmt.Println(gitVersion())
}

// Get all top-namespaces of the used third-party packages (e.g. "golang.org", "github.com")
func getTopNamespacesOfExternalDependencies() []string {
	topNamespaces := make([]string, 0)
	lookupMap := make(map[string]int)
	thirdPartyPackages := getThirdPartyPackages()
	for _, depenency := range thirdPartyPackages {

		components := strings.Split(depenency, "/")
		if len(components) == 0 {
			continue
		}

		topComponent := components[0]
		if _, exists := lookupMap[topComponent]; exists {
			continue
		}

		topNamespaces = append(topNamespaces, topComponent)
		lookupMap[topComponent] = 1
	}

	return topNamespaces
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

	// prepare the environment variables
	environmentVariables := cleanGoEnv()
	environmentVariables = setEnv(environmentVariables, GOPATH_ENVIRONMENT_VARIABLE, goPath)
	environmentVariables = setEnv(environmentVariables, GOBIN_ENVIRONMENT_VARIBALE, goBin)

	// get all dependent packages (will include duplicates and standard library packages)
	allDependentPackages := make([]string, 0)
	for _, buildPackage := range buildPackages {
		output := new(bytes.Buffer)
		errors := new(bytes.Buffer)
		runCommand(output, errors, goPath, environmentVariables, "go", "list", "-f", `'{{ join .Deps ","}}'`, buildPackage)

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
	packagePath := strings.Join(strings.Split(packageName, "/"), string(os.PathSeparator))
	return filepath.Join(goPath, "src", packagePath)
}

// getWorkingDirectory returns the current working directory path or fails.
func getWorkingDirectory() string {
	goPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	return goPath
}

// Execute go in the specified go path with the supplied command arguments.
func runCommand(stdout, stderr io.Writer, workingDirectory string, environmentVariables []string, command string, args ...string) {

	// Create the command
	cmdName := fmt.Sprintf("%s %s", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)

	// Set the working directory
	cmd.Dir = workingDirectory

	// set environment variables
	cmd.Env = environmentVariables

	// Capture the output
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if *verboseFlagIsSet {
		log.Printf("Running %s", cmdName)
	}

	// execute the command
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error running %s: %v", cmdName, err)
	}
}

// cleanGoEnv returns a copy of the current environment with GOPATH_ENVIRONMENT_VARIABLE and GOBIN_ENVIRONMENT_VARIBALE removed.
func cleanGoEnv() (clean []string) {
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, GOPATH_ENVIRONMENT_VARIABLE+"=") || strings.HasPrefix(env, GOBIN_ENVIRONMENT_VARIBALE+"=") {
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

// gitVersion returns the git version of the git repo at goPath as a
// string of the form "yyyy-mm-dd-xxxxxxx", with an optional trailing
// '+' if there are any local uncomitted modifications to the tree.
func gitVersion() string {
	cmd := exec.Command("git", "rev-list", "--max-count=1", "--pretty=format:'%ad-%h'", "--date=short", "HEAD")
	cmd.Dir = goPath

	commandOutput, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error running git rev-list in %s: %v", goPath, err)
	}

	versionNumber := strings.TrimSpace(string(commandOutput))
	if m := gitVersionPattern.FindStringSubmatch(versionNumber); m != nil {
		versionNumber = m[0]
	} else {
		log.Fatalf("Failed to find git version in string %q", versionNumber)
	}

	cmd = exec.Command("git", "diff", "--exit-code")
	cmd.Dir = goPath
	if err := cmd.Run(); err != nil {
		versionNumber += "+"
	}

	return versionNumber
}

// Get the build version flag for the go linker (e.g. -X allmark.io/cmd/allmark 2015-01-11-284c030+).
func getBuildVersionFlag() string {
	return fmt.Sprintf(`--ldflags=-X %s.GitInfo %s`, "allmark.io/modules/common/buildinfo", version)
}

// Get the paths of all source control folders in the specified base path.
func getSourceControlFolderPaths(basePath string) []string {
	paths := make([]string, 0)

	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if info == nil || !info.IsDir() {
			return nil
		}

		if info.Name() == ".git" || info.Name() == ".hg" {
			paths = append(paths, path)
		}

		return nil

	})

	return paths
}
