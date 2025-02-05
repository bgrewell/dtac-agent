//go:build mage
// +build mage

package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/magefile/mage/sh"
	log "github.com/sirupsen/logrus"
)

const (
	packageName = "github.com/intel-innersource/frameworks.automation.dtac.agent"
)

var (
	ldflagsArr []string
	ldflags    = "-X "
	goexe      = "go"
	binaryname = "dtac-agentd"
)

func init() {

	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	if name := os.Getenv("BINARY_NAME"); name != "" {
		binaryname = name
	}

	// Setup ldflags
	buildVer := getBuildVersion()
	ldflagsArr = append(ldflagsArr, fmt.Sprintf("github.com/intel-innersource/frameworks.automation.dtac.agent/internal/version.version=%s", buildVer))

	buildDate := getBuildDate()
	ldflagsArr = append(ldflagsArr, fmt.Sprintf("github.com/intel-innersource/frameworks.automation.dtac.agent/internal/version.date=%s", buildDate))

	buildRev := getBuildRevision()
	ldflagsArr = append(ldflagsArr, fmt.Sprintf("github.com/intel-innersource/frameworks.automation.dtac.agent/internal/version.rev=%s", buildRev))

	buildBranch := getBuildBranch()
	ldflagsArr = append(ldflagsArr, fmt.Sprintf("github.com/intel-innersource/frameworks.automation.dtac.agent/internal/version.branch=%s", buildBranch))

	ldflags += strings.Join(ldflagsArr, " -X ")
}

func outputWith(env map[string]string, cmd string, inArgs ...any) (string, error) {
	s := argsToStrings(inArgs...)
	return sh.OutputWith(env, cmd, s...)
}

func runWith(env map[string]string, cmd string, inArgs ...any) error {
	s := argsToStrings(inArgs...)
	return sh.RunWith(env, cmd, s...)
}

func argsToStrings(v ...any) []string {
	var args []string
	for _, arg := range v {
		switch v := arg.(type) {
		case string:
			if v != "" {
				args = append(args, v)
			}
		case []string:
			if v != nil {
				args = append(args, v...)
			}
		default:
			panic("invalid type")
		}
	}

	return args
}

func getBuildVersion() string {
	v, err := outputWith(nil, "git", "describe", "--tags")
	if err != nil {
		log.Fatal(err)
	}
	vr := strings.Split(v, "-g")
	return vr[0]
}

func getBuildDate() string {
	return time.Now().Format("2006.01.02_150405")
}

func getBuildRevision() string {
	r, err := outputWith(nil, "git", "rev-parse", "--short", "HEAD")
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func getBuildBranch() string {
	b, err := outputWith(nil, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		log.Fatal(err)
	}
	b = strings.TrimSpace(b)
	b = strings.ReplaceAll(b, "\040", "")
	b = strings.ReplaceAll(b, "\011", "")
	b = strings.ReplaceAll(b, "\012", "")
	b = strings.ReplaceAll(b, "\015", "")
	return b
}

func flagEnv() map[string]string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return map[string]string{
		"PACKAGE":     packageName,
		"COMMIT_HASH": hash,
		"BUILD_DATE":  time.Now().Format("2006-01-02T15:04:05Z0700"),
	}
}

func buildFlags() []string {
	if runtime.GOOS == "windows" {
		return []string{"-buildmode", "exe"}
	}
	return nil
}

func buildTags() string {
	// NOT USED CURRENTLY
	if envtags := os.Getenv("DTAC_BUILD_TAGS"); envtags != "" {
		return envtags
	}
	return "none"
}

func findBuildYAMLFiles(rootDir string) ([]string, error) {
	var paths []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is named build.yaml and it's not a directory.
		if !info.IsDir() && info.Name() == "build.yaml" {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			paths = append(paths, absPath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return paths, nil
}

func buildLin64() error {
	fmt.Println("  Compiling Linux amd64")
	return build("linux", "amd64")
}

func buildLinArm() error {
	fmt.Println("  Compiling Linux Arm")
	return build("linux", "arm")
}

func buildWin64() error {
	fmt.Println("  Compiling Windows amd64")
	return build("windows", "amd64")
}

func buildMac64() error {
	fmt.Println("  Compiling MacOS amd64")
	return build("darwin", "amd64")
}

func buildCliLin64() error {
	fmt.Println("  Compiling Linux amd64")
	return buildCli("linux", "amd64")
}

func buildCliLinArm() error {
	fmt.Println("  Compiling Linux Arm")
	return buildCli("linux", "arm")
}

func buildCliWin64() error {
	fmt.Println("  Compiling Windows amd64")
	return buildCli("windows", "amd64")
}

func buildCliMac64() error {
	fmt.Println("  Compiling MacOS amd64")
	return buildCli("darwin", "amd64")
}

func build(os string, arch string) error {
	extension := ""
	if os == "windows" {
		extension = ".exe"
	} else if os == "darwin" {
		extension = ".app"
	}
	env := flagEnv()
	env["GOOS"] = os
	env["GOARCH"] = arch
	output := fmt.Sprintf("bin/%s%s%s", binaryname, fmt.Sprintf("-%s", arch), extension)
	return runWith(env, goexe, "build", "-ldflags", ldflags, buildFlags(), "-tags", buildTags(), "-o", output, "cmd/agent/main.go")
}

func buildCli(os string, arch string) error {
	extension := ""
	if os == "windows" {
		extension = ".exe"
	} else if os == "darwin" {
		extension = ".app"
	}
	env := flagEnv()
	env["GOOS"] = os
	env["GOARCH"] = arch
	output := fmt.Sprintf("bin/%s%s%s", "dtac", fmt.Sprintf("-%s", arch), extension)
	return runWith(env, goexe, "build", "-ldflags", ldflags, buildFlags(), "-tags", buildTags(), "-o", output, "cmd/cli/main.go")
}

func Build() error {
	fmt.Println("Building agent")
	funcs := []func() error{buildLin64, buildLinArm, buildWin64, buildMac64}
	for _, f := range funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

func BuildCli() error {
	fmt.Println("Building cli")
	funcs := []func() error{buildCliLin64, buildCliLinArm, buildCliWin64, buildCliMac64}
	for _, f := range funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

func Container() error {
	if err := Build(); err != nil {
		return err
	}
	if err := Plugins(); err != nil {
		return err
	}
	if err := runWith(nil, "cp", "-r", "bin", "deployments/docker"); err != nil {
		return err
	}
	if err := runWith(nil, "docker", "build", "-t", fmt.Sprintf("dtac-agent:%s", getBuildVersion()), "deployments/docker/."); err != nil {
		return err
	}
	return runWith(nil, "docker", "build", "-t", fmt.Sprintf("dtac-agent:%s", "latest"), "deployments/docker/.")
}

func Debug() error {
	// Launch container with "tail -f /dev/null"
	// Execute command to install datc-agentd  "/tmp/dtac-agentd --install" piping to stdin/stdout/stderr
	return errors.New("this method has not been implemented")
}

func Deps() error {
	fmt.Println("Updating dependencies")
	env := make(map[string]string)
	env["GOPRIVATE"] = "github.com/bgrewell"
	env["GOPROXY"] = "direct"
	env["GO111MODULE"] = "on"
	env["GOSUMDB"] = "off"
	if err := runWith(env, goexe, "get", "-u", "./..."); err != nil {
		return err
	}
	if err := runWith(env, goexe, "mod", "tidy"); err != nil {
		return err
	}
	return runWith(nil, goexe, "install", "google.golang.org/protobuf/cmd/protoc-gen-go")
}

func Run() error {
	env := make(map[string]string)
	env["DTAC_CFG_LOCATION"] = "configs/example.yaml"
	//// TODO: Execute but pipe to stdin, stdout, stderr
	return runWith(nil, "sudo", "-E", "/usr/local/go/bin/go", "run", "cmd/agent/main.go")
}

func Plugins() error {
	fmt.Println("Building plugins")
	// Define a struct to unmarshal the build.yaml contents into.
	type BuildInfo struct {
		Name      string   `yaml:"name"`
		Entry     string   `yaml:"entry"`
		Platforms []string `yaml:"platforms"`
	}

	buildFiles, err := findBuildYAMLFiles("cmd/plugins")
	if err != nil {
		return err
	}

	for _, filename := range buildFiles {
		var buildInfo BuildInfo

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed reading %s: %v", filename, err)
		}

		err = yaml.Unmarshal(data, &buildInfo)
		if err != nil {
			return fmt.Errorf("failed unmarshaling %s: %v", filename, err)
		}

		// Run the go build command using the extracted name, entry, and platforms values.
		for _, platform := range buildInfo.Platforms {
			parts := strings.Split(platform, ":")
			os := parts[0]
			arch := "amd64"
			if len(parts) > 1 {
				arch = parts[1]
			}

			fmt.Printf("  Compiling %s for %s %s\n", buildInfo.Name, os, arch)
			inPath := filepath.Dir(filename)
			outPath := fmt.Sprintf("bin/plugins/%s.plugin", buildInfo.Name)
			err := buildPlugins(path.Join(inPath, buildInfo.Entry), os, arch, outPath)
			if err != nil {
				return fmt.Errorf("failed building plugin %s: %v", buildInfo.Name, err)
			}
		}
	}

	return nil
}

func buildPlugins(source string, os string, arch string, binary string) error {
	extension := ""
	if os == "windows" {
		extension = ".exe"
	} else if os == "darwin" {
		extension = ".app"
	}
	env := flagEnv()
	env["GOOS"] = os
	env["GOARCH"] = arch
	output := fmt.Sprintf("%s%s", binary, extension)
	return runWith(env, goexe, "build", "-tags", buildTags(), "-o", output, source)
}

func Clean() error {
	os.RemoveAll("dist")
	os.RemoveAll("bin")
	return nil
}

func Test() error {
	return runWith(nil, goexe, "test", "-v", "./...")
}

func Check() error {
	if err := runWith(nil, goexe, "install", "honnef.co/go/tools/cmd/staticcheck@latest"); err != nil {
		return err
	}
	if err := runWith(nil, goexe, "install", "golang.org/x/lint/golint@latest"); err != nil {
		return err
	}
	if err := runWith(nil, "staticcheck", "./..."); err != nil {
		return err
	}
	if err := runWith(nil, "golint", "./..."); err != nil {
		return err
	}
	return nil
}

func FindTODOs() error {
	// Run `git grep` to find all files that contain TODO comments
	grepCmd := exec.Command("git", "grep", "-l", "TODO")
	grepOutput, err := grepCmd.Output()
	if err != nil {
		return fmt.Errorf("error running `git grep`: %v", err)
	}

	// Split the output into separate file names
	fileNames := strings.Split(string(grepOutput), "\n")

	// Loop over the file names and run `git blame` on each file
	for _, fileName := range fileNames {
		if fileName == "" {
			continue
		}
		blameCmd := exec.Command("git", "blame", fileName)
		blameOutput, err := blameCmd.Output()
		if err != nil {
			return fmt.Errorf("error running `git blame` on %s: %v", fileName, err)
		}

		// Search the output of `git blame` for TODO comments
		for _, line := range strings.Split(string(blameOutput), "\n") {
			if strings.Contains(line, "TODO") {
				fmt.Printf("%s: %s\n", fileName, line)
			}
		}
	}

	return nil
}
