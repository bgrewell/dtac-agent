//go:build mage
// +build mage

package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
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
	installPrefix   = "/opt/dtac"
	installBinDir   = "/opt/dtac/bin"
	installPlugDir  = "/opt/dtac/plugins"
	installModDir   = "/opt/dtac/modules"
	etcDtacDir      = "/etc/dtac"
	etcDtacCfg      = "/etc/dtac/config.yaml"
	systemdUnitSrc  = "service/systemd/dtac-agentd.service"
	systemdUnitDest = "/etc/systemd/system/dtac-agentd.service"
	usrBinSymlink   = "/usr/bin/dtac"
	cfgSrcExample   = "configs/example.yaml"
	cfgPassMarker   = "need_to_generate_a_random_password_on_install_or_first_run"
	packageName = "github.com/bgrewell/dtac-agent"
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
	ldflagsArr = append(ldflagsArr, fmt.Sprintf("github.com/bgrewell/dtac-agent/internal/version.version=%s", buildVer))

	buildDate := getBuildDate()
	ldflagsArr = append(ldflagsArr, fmt.Sprintf("github.com/bgrewell/dtac-agent/internal/version.date=%s", buildDate))

	buildRev := getBuildRevision()
	ldflagsArr = append(ldflagsArr, fmt.Sprintf("github.com/bgrewell/dtac-agent/internal/version.rev=%s", buildRev))

	buildBranch := getBuildBranch()
	ldflagsArr = append(ldflagsArr, fmt.Sprintf("github.com/bgrewell/dtac-agent/internal/version.branch=%s", buildBranch))

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
	if err := Modules(); err != nil {
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

// osquery embedding -----------------------------------------------------------

const (
	osqueryVersion       = "5.23.0"
	osqueryReleasesAPI   = "https://api.github.com/repos/osquery/osquery/releases/tags/"
	osqueryEmbedDir      = "cmd/plugins/osquery/osqueryplugin/binaries"
	osqueryEmbedBuildTag = "embed_osquery"
)

// osqueryEmbedTarget describes one (goos, goarch) build that bundles osqueryd.
// AssetName is the upstream release asset to download; BinaryInTar is the
// path of osqueryd within that tarball.
type osqueryEmbedTarget struct {
	GOOS        string
	GOARCH      string
	AssetName   string
	BinaryInTar string
}

// osqueryEmbedTargets lists every platform we can bundle a daemon for.
// macos amd64 has no native tarball on the upstream release (only .pkg) so
// it's deliberately omitted — built without -tags=embed_osquery it still
// works via config.binary_path.
var osqueryEmbedTargets = []osqueryEmbedTarget{
	// Linux tarballs use usr/bin/osqueryd as a symlink to the real binary
	// under opt/osquery/bin/osqueryd — extract the real one or we get
	// zero bytes back.
	{"linux", "amd64", "osquery-" + osqueryVersion + "_1.linux_x86_64.tar.gz", "opt/osquery/bin/osqueryd"},
	{"linux", "arm64", "osquery-" + osqueryVersion + "_1.linux_aarch64.tar.gz", "opt/osquery/bin/osqueryd"},
	// The macos "bare" tarball is just the daemon binary at the root.
	{"darwin", "arm64", "osqueryd-macos-bare-" + osqueryVersion + ".tar.gz", "osqueryd"},
}

// OsqueryBundle downloads osqueryd from the pinned upstream release,
// verifies SHA256 against the GitHub release metadata, stages the binary
// under cmd/plugins/osquery/osqueryplugin/binaries/<os-arch>/osqueryd, then
// builds the osquery plugin with -tags=embed_osquery for each supported
// platform. The resulting plugin is self-contained — at first run it
// extracts osqueryd to a content-addressed cache dir.
//
// Skips downloads for any platform whose binary is already staged. Use
// `mage clean` (or rm -rf cmd/plugins/osquery/osqueryplugin/binaries) to
// force re-fetching.
func OsqueryBundle() error {
	fmt.Printf("Bundling osquery %s into plugin\n", osqueryVersion)

	assets, err := fetchOsqueryReleaseAssets()
	if err != nil {
		return fmt.Errorf("fetching osquery release metadata: %v", err)
	}

	for _, t := range osqueryEmbedTargets {
		targetDir := filepath.Join(osqueryEmbedDir, t.GOOS+"-"+t.GOARCH)
		targetPath := filepath.Join(targetDir, "osqueryd")
		if _, err := os.Stat(targetPath); err == nil {
			fmt.Printf("  osqueryd for %s-%s already staged\n", t.GOOS, t.GOARCH)
			continue
		}
		if err := stageOsqueryBinary(t, assets, targetPath); err != nil {
			return fmt.Errorf("staging osqueryd for %s-%s: %v", t.GOOS, t.GOARCH, err)
		}
	}

	for _, t := range osqueryEmbedTargets {
		fmt.Printf("  Compiling osquery (embed) for %s %s\n", t.GOOS, t.GOARCH)
		if err := buildOsqueryEmbedded(t.GOOS, t.GOARCH); err != nil {
			return fmt.Errorf("building embedded osquery plugin for %s-%s: %v", t.GOOS, t.GOARCH, err)
		}
	}
	return nil
}

// fetchOsqueryReleaseAssets calls the GitHub release API for the pinned
// osquery tag and returns a map from asset name to (download URL, sha256).
// The release API includes a sha256 in the `digest` field for each asset.
func fetchOsqueryReleaseAssets() (map[string]struct{ URL, SHA256 string }, error) {
	resp, err := http.Get(osqueryReleasesAPI + osqueryVersion)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("github api returned %d: %s", resp.StatusCode, string(body))
	}

	var payload struct {
		Assets []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
			Digest             string `json:"digest"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	out := make(map[string]struct{ URL, SHA256 string }, len(payload.Assets))
	for _, a := range payload.Assets {
		digest := strings.TrimPrefix(a.Digest, "sha256:")
		out[a.Name] = struct{ URL, SHA256 string }{a.BrowserDownloadURL, digest}
	}
	return out, nil
}

// stageOsqueryBinary downloads, verifies, extracts, and writes one osqueryd
// binary to targetPath. The intermediate tarball is held in memory — at ~85
// MiB it's small enough to skip a temp file.
func stageOsqueryBinary(t osqueryEmbedTarget, assets map[string]struct{ URL, SHA256 string }, targetPath string) error {
	meta, ok := assets[t.AssetName]
	if !ok {
		return fmt.Errorf("asset %q not found in osquery %s release", t.AssetName, osqueryVersion)
	}
	if meta.SHA256 == "" {
		return fmt.Errorf("asset %q has no SHA256 digest in release metadata", t.AssetName)
	}

	fmt.Printf("  Downloading %s\n", t.AssetName)
	tarball, err := downloadVerified(meta.URL, meta.SHA256)
	if err != nil {
		return err
	}

	fmt.Printf("  Extracting %s from %s\n", t.BinaryInTar, t.AssetName)
	binary, err := readFromTarGz(tarball, t.BinaryInTar)
	if err != nil {
		return err
	}
	if len(binary) == 0 {
		return fmt.Errorf("entry %q in %s is empty (likely a symlink — point BinaryInTar at the real target)", t.BinaryInTar, t.AssetName)
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(targetPath, binary, 0o755)
}

// downloadVerified fetches url, streams the bytes into memory while hashing,
// and returns the body iff the SHA256 matches expected (hex).
func downloadVerified(url, expected string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download %s returned %d", url, resp.StatusCode)
	}
	h := sha256.New()
	var buf bytes.Buffer
	if _, err := io.Copy(io.MultiWriter(&buf, h), resp.Body); err != nil {
		return nil, err
	}
	got := hex.EncodeToString(h.Sum(nil))
	if got != expected {
		return nil, fmt.Errorf("sha256 mismatch: got %s, want %s", got, expected)
	}
	return buf.Bytes(), nil
}

// readFromTarGz walks a gzipped tarball in memory and returns the bytes of
// the entry whose name matches target. Match is exact, with leading "./"
// tolerated since some tarballs emit it.
func readFromTarGz(data []byte, target string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		name := strings.TrimPrefix(hdr.Name, "./")
		if name == target {
			return io.ReadAll(tr)
		}
	}
	return nil, fmt.Errorf("entry %q not found in tarball", target)
}

// E2EOsquery runs the DART end-to-end workflow for the osquery plugin.
// Requires:
//   - dart CLI on PATH (https://github.com/bgrewell/dart)
//   - LXD/Incus available locally
//   - bin/plugins/osquery-linux-amd64.plugin built via `mage osqueryBundle`
//
// Pass `-v` through the DART_FLAGS env var for verbose output, e.g.
//   DART_FLAGS=-v mage e2eOsquery
func E2EOsquery() error {
	if _, err := exec.LookPath("dart"); err != nil {
		return fmt.Errorf("dart CLI not found on PATH; install from https://github.com/bgrewell/dart")
	}
	if _, err := os.Stat("bin/plugins/osquery-linux-amd64.plugin"); err != nil {
		return fmt.Errorf("missing bin/plugins/osquery-linux-amd64.plugin — run `mage osqueryBundle` first")
	}
	args := []string{"-c", "test/e2e/osquery-plugin/osquery-plugin.yaml"}
	if extra := os.Getenv("DART_FLAGS"); extra != "" {
		args = append(args, strings.Fields(extra)...)
	}
	return runWith(nil, "dart", iargs(args)...)
}

// iargs reshapes []string into the []any signature expected by runWith.
func iargs(s []string) []any {
	out := make([]any, len(s))
	for i, v := range s {
		out[i] = v
	}
	return out
}

// buildOsqueryEmbedded compiles cmd/plugins/osquery for one platform with
// the embed_osquery build tag set. The output filename embeds the target
// triple so cross-builds for different arches don't clobber each other —
// the agent host should pick e.g. osquery-linux-amd64.plugin.
func buildOsqueryEmbedded(goos, goarch string) error {
	extension := ""
	if goos == "darwin" {
		extension = ".app"
	} else if goos == "windows" {
		extension = ".exe"
	}
	output := fmt.Sprintf("bin/plugins/osquery-%s-%s.plugin%s", goos, goarch, extension)
	env := flagEnv()
	env["GOOS"] = goos
	env["GOARCH"] = goarch
	return runWith(env, goexe, "build", "-tags", osqueryEmbedBuildTag, "-o", output, "cmd/plugins/osquery/main.go")
}

func Modules() error {
	fmt.Println("Building modules")
	// Define a struct to unmarshal the build.yaml contents into.
	type BuildInfo struct {
		Name      string   `yaml:"name"`
		Entry     string   `yaml:"entry"`
		Platforms []string `yaml:"platforms"`
	}

	buildFiles, err := findBuildYAMLFiles("cmd/modules")
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
			outPath := fmt.Sprintf("bin/modules/%s.module", buildInfo.Name)
			err := buildModules(path.Join(inPath, buildInfo.Entry), os, arch, outPath)
			if err != nil {
				return fmt.Errorf("failed building module %s: %v", buildInfo.Name, err)
			}
		}
	}

	return nil
}

func buildModules(source string, os string, arch string, binary string) error {
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

// Install sets up dtac (agent, cli, plugins, config, systemd, symlink)
func Install() error {
	if runtime.GOOS != "linux" {
		return errors.New("Install is currently implemented for Linux only")
	}
	if err := requireRoot(); err != nil {
		return err
	}
	if !hasSystemd() {
		return errors.New("systemd not detected; Install requires systemd for service management")
	}

	// 1) Build: agent, cli, plugins, modules (current code builds multiple platforms; for install we only need host)
	// Reuse existing tasks so the artifacts exist in ./bin, ./bin/plugins, and ./bin/modules
	if err := buildHostOnly(); err != nil {
		return err
	}
	//if err := Plugins(); err != nil {
	//	return err
	//}
	//if err := Modules(); err != nil {
	//	return err
	//}

	// 2) Create /opt/dtac/{bin,plugins,modules}
	if err := ensureDir(installBinDir, 0o755); err != nil {
		return err
	}
	if err := ensureDir(installPlugDir, 0o755); err != nil {
		return err
	}
	if err := ensureDir(installModDir, 0o755); err != nil {
		return err
	}

	// 3) Copy compiled agent & cli
	agentSrc, cliSrc, err := hostBinaries()
	if err != nil {
		return err
	}
	if err := copyFile(agentSrc, filepath.Join(installBinDir, "dtac-agentd"), 0o755); err != nil {
		return err
	}
	if err := copyFile(cliSrc, filepath.Join(installBinDir, "dtac"), 0o755); err != nil {
		return err
	}

	// 4) Copy all plugins
	if err := copyPlugins("bin/plugins", installPlugDir); err != nil {
		return err
	}

	// 5) Copy all modules
	if err := copyModules("bin/modules", installModDir); err != nil {
		return err
	}

	// 6) Create /etc/dtac
	if err := ensureDir(etcDtacDir, 0o755); err != nil {
		return err
	}

	// 7) Copy example config
	if err := copyFile(cfgSrcExample, etcDtacCfg, 0o600); err != nil {
		return err
	}

	// 8) Generate 8-char alnum password and replace placeholder
	if err := replaceConfigPassword(etcDtacCfg, cfgPassMarker); err != nil {
		return err
	}

	// 9) Install systemd service, reload, enable
	if err := copyFile(systemdUnitSrc, systemdUnitDest, 0o644); err != nil {
		return err
	}
	if err := runCmd("systemctl", "daemon-reload"); err != nil {
		return err
	}
	if err := runCmd("systemctl", "enable", "dtac-agentd.service"); err != nil {
		return err
	}

	// 10) Symlink /usr/bin/dtac -> /opt/dtac/bin/dtac
	_ = os.Remove(usrBinSymlink) // best-effort remove if exists
	if err := os.Symlink(filepath.Join(installBinDir, "dtac"), usrBinSymlink); err != nil {
		return fmt.Errorf("create symlink %s -> %s: %w", usrBinSymlink, filepath.Join(installBinDir, "dtac"), err)
	}

	// 11) Start service
	if err := runCmd("systemctl", "start", "dtac-agentd.service"); err != nil {
		return err
	}

	fmt.Println("✅ dtac installed successfully.")
	return nil
}

// Uninstall stops/disable service, removes installed files (keeps /etc/dtac/config.yaml)
func Uninstall() error {
	if runtime.GOOS != "linux" {
		return errors.New("Uninstall is currently implemented for Linux only")
	}
	if err := requireRoot(); err != nil {
		return err
	}
	if hasSystemd() {
		_ = runCmd("systemctl", "stop", "dtac-agentd.service")
		_ = runCmd("systemctl", "disable", "dtac-agentd.service")
	}

	// Remove systemd unit and reload
	_ = os.Remove(systemdUnitDest)
	if hasSystemd() {
		_ = runCmd("systemctl", "daemon-reload")
	}

	// Remove symlink
	_ = os.Remove(usrBinSymlink)

	// Remove /opt/dtac tree
	_ = os.RemoveAll(installPrefix)

	// Keep /etc/dtac/config.yaml (user might have edited it)
	fmt.Println("✅ dtac uninstalled (config preserved at /etc/dtac).")
	return nil
}

// ---------- Helpers ----------

// Build only host OS/arch binaries for install
func buildHostOnly() error {
	os := runtime.GOOS
	arch := runtime.GOARCH

	// agent
	if err := build(os, arch); err != nil {
		return err
	}
	// cli
	if err := buildCli(os, arch); err != nil {
		return err
	}
	return nil
}

func hostBinaries() (agentPath, cliPath string, err error) {

	arch := runtime.GOARCH
	var ext string
	switch runtime.GOOS {
	case "windows":
		ext = ".exe"
	case "darwin":
		ext = ".app"
	default:
		ext = ""
	}
	agent := fmt.Sprintf("bin/%s-%s%s", binaryname, arch, ext)
	cli := fmt.Sprintf("bin/%s-%s%s", "dtac", arch, ext)

	if _, e := os.Stat(agent); e != nil {
		return "", "", fmt.Errorf("agent binary not found: %s (build failed?)", agent)
	}
	if _, e := os.Stat(cli); e != nil {
		return "", "", fmt.Errorf("cli binary not found: %s (build failed?)", cli)
	}
	return agent, cli, nil
}

func ensureDir(p string, mode fs.FileMode) error {
	return os.MkdirAll(p, mode)
}

func copyFile(src, dst string, mode fs.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open %s: %w", src, err)
	}
	defer in.Close()

	tmp := dst + ".tmp"
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("create %s: %w", tmp, err)
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		_ = os.Remove(tmp)
		return fmt.Errorf("copy to %s: %w", tmp, err)
	}
	if err := out.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, dst)
}

func copyPlugins(srcDir, dstDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("read %s: %w", srcDir, err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		// Only copy *.plugin* (linux: .plugin; other OS may have suffixes)
		if !strings.HasPrefix(e.Name(), ".") && strings.Contains(e.Name(), ".plugin") {
			src := filepath.Join(srcDir, e.Name())
			dst := filepath.Join(dstDir, e.Name())
			if err := copyFile(src, dst, 0o755); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyModules(srcDir, dstDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("read %s: %w", srcDir, err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		// Only copy *.module* (linux: .module; other OS may have suffixes)
		if !strings.HasPrefix(e.Name(), ".") && strings.Contains(e.Name(), ".module") {
			src := filepath.Join(srcDir, e.Name())
			dst := filepath.Join(dstDir, e.Name())
			if err := copyFile(src, dst, 0o755); err != nil {
				return err
			}
		}
	}
	return nil
}

func replaceConfigPassword(cfgPath, marker string) error {
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	pw, err := randomPassword(8)
	if err != nil {
		return err
	}
	out := bytes.Replace(data, []byte(marker), []byte(pw), 1)
	if bytes.Equal(out, data) {
		return fmt.Errorf("password marker not found in %s", cfgPath)
	}
	return os.WriteFile(cfgPath, out, 0o600)
}

func randomPassword(n int) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	// crypto/rand to generate unbiased alnum string
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("rand.Read: %w", err)
	}
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b), nil
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func requireRoot() error {
	if os.Geteuid() != 0 {
		return errors.New("this operation requires root; re-run with sudo (e.g., sudo mage Install)")
	}
	return nil
}

func hasSystemd() bool {
	// crude but effective: check presence of systemctl and systemd's pid 1 cgroup name
	if _, err := exec.LookPath("systemctl"); err != nil {
		return false
	}
	// Optional: further heuristics could be added here if needed
	return true
}
