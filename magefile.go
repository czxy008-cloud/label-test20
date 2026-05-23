//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var (
	appName  = "task-scheduler"
	buildDir = "dist"
	mainPath = "./cmd/server"
)

type Platform struct {
	OS   string
	Arch string
}

var platforms = []Platform{
	{"linux", "amd64"},
	{"linux", "arm64"},
	{"darwin", "amd64"},
	{"darwin", "arm64"},
	{"windows", "amd64"},
}

func init() {
	os.Setenv("GO111MODULE", "on")
	os.Setenv("CGO_ENABLED", "0")
}

func Build() error {
	fmt.Println("开始构建...")
	clean()

	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("创建构建目录失败: %w", err)
	}

	platform := Platform{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	return buildForPlatform(platform)
}

func BuildAll() error {
	fmt.Println("开始构建所有平台...")
	clean()

	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("创建构建目录失败: %w", err)
	}

	for _, platform := range platforms {
		if err := buildForPlatform(platform); err != nil {
			fmt.Printf("构建 %s/%s 失败: %v\n", platform.OS, platform.Arch, err)
		}
	}

	fmt.Println("所有平台构建完成")
	return nil
}

func buildForPlatform(p Platform) error {
	fmt.Printf("构建 %s/%s...\n", p.OS, p.Arch)

	ext := ""
	if p.OS == "windows" {
		ext = ".exe"
	}

	outputName := fmt.Sprintf("%s-%s-%s%s", appName, p.OS, p.Arch, ext)
	outputPath := fmt.Sprintf("%s/%s", buildDir, outputName)

	env := os.Environ()
	env = append(env, fmt.Sprintf("GOOS=%s", p.OS))
	env = append(env, fmt.Sprintf("GOARCH=%s", p.Arch))
	env = append(env, "CGO_ENABLED=0")

	version := getVersion()
	buildDate := getBuildDate()
	ldflags := fmt.Sprintf("-s -w -X main.version=%s -X main.buildDate=%s", version, buildDate)

	cmd := exec.Command("go", "build",
		"-o", outputPath,
		"-ldflags", ldflags,
		"-trimpath",
		mainPath,
	)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Run() error {
	if err := Build(); err != nil {
		return err
	}

	fmt.Println("运行服务器...")

	binPath := fmt.Sprintf("%s/%s-%s-%s", buildDir, appName, runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}

	cmd := exec.Command(binPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func Clean() error {
	return clean()
}

func clean() error {
	if _, err := os.Stat(buildDir); err == nil {
		fmt.Println("清理构建目录...")
		return os.RemoveAll(buildDir)
	}
	return nil
}

func Test() error {
	fmt.Println("运行测试...")
	cmd := exec.Command("go", "test", "./...", "-v", "-cover")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Lint() error {
	fmt.Println("运行代码检查...")
	cmd := exec.Command("go", "vet", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Fmt() error {
	fmt.Println("格式化代码...")
	cmd := exec.Command("go", "fmt", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ModTidy() error {
	fmt.Println("整理依赖...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getVersion() string {
	version := os.Getenv("VERSION")
	if version == "" {
		version = "0.0.1"
	}
	return version
}

func getBuildDate() string {
	return fmt.Sprintf("%d", 1700000000)
}
