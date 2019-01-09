// +build ignore

package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	versionRe = regexp.MustCompile(`-[0-9]{1,3}-g[0-9a-f]{5,10}`)
	goarch    string
	goos      string
	gocc      string
	gocxx     string
	cgo       string
	pkgArch   string
	version   string = "v1"
	// deb & rpm does not support semver so have to handle their version a little differently
	linuxPackageVersion   string = "v1"
	linuxPackageIteration string = ""
	race                  bool
	phjsToRelease         string
	workingDir            string
	includeBuildNumber    bool     = true
	buildNumber           int      = 0
	binaries              []string = []string{"ddcli"}
)

const minGoVersion = 1.8

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	ensureGoPath()

	flag.StringVar(&goarch, "goarch", runtime.GOARCH, "GOARCH")
	flag.StringVar(&goos, "goos", runtime.GOOS, "GOOS")
	flag.StringVar(&gocc, "cc", "", "CC")
	flag.StringVar(&gocxx, "cxx", "", "CXX")
	flag.StringVar(&cgo, "cgo-enabled", "", "CGO_ENABLED")
	flag.StringVar(&pkgArch, "pkg-arch", "", "PKG ARCH")
	flag.StringVar(&phjsToRelease, "phjs", "", "PhantomJS binary")
	flag.BoolVar(&race, "race", race, "Use race detector")
	flag.BoolVar(&includeBuildNumber, "includeBuildNumber", includeBuildNumber, "IncludeBuildNumber in package name")
	flag.IntVar(&buildNumber, "buildNumber", 0, "Build number from CI system")
	flag.Parse()

	readVersionFromLatestJson()

	log.Printf("Version: %s, Linux Version: %s, Package Iteration: %s\n", version, linuxPackageVersion, linuxPackageIteration)

	if flag.NArg() == 0 {
		log.Println("Usage: go run build.go build")
		return
	}

	workingDir, _ = os.Getwd()

	for _, cmd := range flag.Args() {
		switch cmd {
		case "setup":
			setup()

		case "build":
			clean()
			for _, binary := range binaries {
				build(binary, "./pkg/cmd", []string{})
			}

		case "test":
			test("./pkg/...")

		case "sha-dist":
			shaFilesInDist()

		case "latest":
			makeLatestDistCopies()

		case "clean":
			clean()

		default:
			log.Fatalf("Unknown command %q", cmd)
		}
	}
}

func makeLatestDistCopies() {
	files, err := ioutil.ReadDir("dist")
	if err != nil {
		log.Fatalf("failed to create latest copies. Cannot read from /dist")
	}

	latestMapping := map[string]string{
		".deb":    "dist/ddcli_latest_amd64.deb",
		".rpm":    "dist/ddcli-latest-1.x86_64.rpm",
		".tar.gz": "dist/ddcli-latest.linux-x64.tar.gz",
	}

	for _, file := range files {
		for extension, fullName := range latestMapping {
			if strings.HasSuffix(file.Name(), extension) {
				runError("cp", path.Join("dist", file.Name()), fullName)
			}
		}
	}
}

func readVersionFromLatestJson() {
	reader, err := os.Open("latest.json")
	if err != nil {
		log.Fatal("Failed to open latest.json")
		return
	}
	defer reader.Close()

	jsonObj := map[string]interface{}{}
	jsonParser := json.NewDecoder(reader)

	if err := jsonParser.Decode(&jsonObj); err != nil {
		log.Fatal("Failed to decode latest.json")
	}

	version = jsonObj["stable"].(string)
	linuxPackageVersion = version
	linuxPackageIteration = ""

	// handle pre version stuff (deb / rpm does not support semver)
	parts := strings.Split(version, "-")

	if len(parts) > 1 {
		linuxPackageVersion = parts[0]
		linuxPackageIteration = parts[1]
	}

	// add timestamp to iteration
	if includeBuildNumber {
		if buildNumber != 0 {
			linuxPackageIteration = fmt.Sprintf("%d%s", buildNumber, linuxPackageIteration)
		} else {
			linuxPackageIteration = fmt.Sprintf("%d%s", time.Now().Unix(), linuxPackageIteration)
		}
	}
}

func verifyGitRepoIsClean() {
	rs, err := runError("git", "ls-files", "--modified")
	if err != nil {
		log.Fatalf("Failed to check if git tree was clean, %v, %v\n", string(rs), err)
		return
	}
	count := len(string(rs))
	if count > 0 {
		log.Fatalf("Git repository has modified files, aborting")
	}

	log.Println("Git repository is clean")
}

func ensureGoPath() {
	if os.Getenv("GOPATH") == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		gopath := filepath.Clean(filepath.Join(cwd, "../../../../"))
		log.Println("GOPATH is", gopath)
		os.Setenv("GOPATH", gopath)
	}
}

func ChangeWorkingDir(dir string) {
	os.Chdir(dir)
}

func setup() {
	runPrint("go", "get", "-v", "github.com/kardianos/govendor")
	runPrint("go", "install", "-v", "./pkg/cmd")
}

func test(pkg string) {
	setBuildEnv()
	runPrint("go", "test", "-short", "-timeout", "60s", pkg)
}

func build(binaryName, pkg string, tags []string) {
	binary := "./bin/" + binaryName
	if goos == "windows" {
		binary += ".exe"
	}

	rmr(binary, binary+".md5")
	args := []string{"build", "-ldflags", ldflags()}
	if len(tags) > 0 {
		args = append(args, "-tags", strings.Join(tags, ","))
	}
	if race {
		args = append(args, "-race")
	}

	args = append(args, "-o", binary)
	args = append(args, pkg)
	setBuildEnv()

	runPrint("go", "version")
	runPrint("go", args...)

	// Create an md5 checksum of the binary, to be included in the archive for
	// automatic upgrades.
	err := md5File(binary)
	if err != nil {
		log.Fatal(err)
	}
}

func ldflags() string {
	var b bytes.Buffer
	b.WriteString("-w")
	b.WriteString(fmt.Sprintf(" -X main.version=%s", version))
	b.WriteString(fmt.Sprintf(" -X main.commit=%s", getGitSha()))
	b.WriteString(fmt.Sprintf(" -X main.buildstamp=%d", buildStamp()))
	return b.String()
}

func rmr(paths ...string) {
	for _, path := range paths {
		log.Println("rm -r", path)
		os.RemoveAll(path)
	}
}

func clean() {
	rmr("dist")
	rmr("tmp")
	rmr(filepath.Join(os.Getenv("GOPATH"), fmt.Sprintf("pkg/%s_%s/xiak/remote-command", goos, goarch)))
}

func setBuildEnv() {
	os.Setenv("GOOS", goos)
	if strings.HasPrefix(goarch, "armv") {
		os.Setenv("GOARCH", "arm")
		os.Setenv("GOARM", goarch[4:])
	} else {
		os.Setenv("GOARCH", goarch)
	}
	if goarch == "386" {
		os.Setenv("GO386", "387")
	}
	if cgo != "" {
		os.Setenv("CGO_ENABLED", cgo)
	}
	if gocc != "" {
		os.Setenv("CC", gocc)
	}
	if gocxx != "" {
		os.Setenv("CXX", gocxx)
	}
}

func getGitSha() string {
	v, err := runError("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return "unknown-dev"
	}
	return string(v)
}

func buildStamp() int64 {
	bs, err := runError("git", "show", "-s", "--format=%ct")
	if err != nil {
		return time.Now().Unix()
	}
	s, _ := strconv.ParseInt(string(bs), 10, 64)
	return s
}

func buildArch() string {
	os := goos
	if os == "darwin" {
		os = "macosx"
	}
	return fmt.Sprintf("%s-%s", os, goarch)
}

func run(cmd string, args ...string) []byte {
	bs, err := runError(cmd, args...)
	if err != nil {
		log.Println(cmd, strings.Join(args, " "))
		log.Println(string(bs))
		log.Fatal(err)
	}
	return bytes.TrimSpace(bs)
}

func runError(cmd string, args ...string) ([]byte, error) {
	ecmd := exec.Command(cmd, args...)
	bs, err := ecmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(bs), nil
}

func runPrint(cmd string, args ...string) {
	log.Println(cmd, strings.Join(args, " "))
	ecmd := exec.Command(cmd, args...)
	ecmd.Stdout = os.Stdout
	ecmd.Stderr = os.Stderr
	err := ecmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func md5File(file string) error {
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	h := md5.New()
	_, err = io.Copy(h, fd)
	if err != nil {
		return err
	}

	out, err := os.Create(file + ".md5")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, "%x\n", h.Sum(nil))
	if err != nil {
		return err
	}

	return out.Close()
}

func shaFilesInDist() {
	filepath.Walk("./dist", func(path string, f os.FileInfo, err error) error {
		if path == "./dist" {
			return nil
		}

		if strings.Contains(path, ".sha256") == false {
			err := shaFile(path)
			if err != nil {
				log.Printf("Failed to create sha file. error: %v\n", err)
			}
		}
		return nil
	})
}

func shaFile(file string) error {
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	h := sha256.New()
	_, err = io.Copy(h, fd)
	if err != nil {
		return err
	}

	out, err := os.Create(file + ".sha256")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, "%x\n", h.Sum(nil))
	if err != nil {
		return err
	}

	return out.Close()
}
