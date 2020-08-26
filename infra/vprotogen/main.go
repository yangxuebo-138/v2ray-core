package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"v2ray.com/core/common"
)

var protocMap = map[string]string{
	"windows": filepath.Join(os.Getenv("GOPATH"), "src", "v2ray.com", "core", ".dev", "protoc", "windows", "protoc.exe"),
	"darwin":  filepath.Join(os.Getenv("GOPATH"), "src", "v2ray.com", "core", ".dev", "protoc", "macos", "protoc"),
	"linux":   filepath.Join(os.Getenv("GOPATH"), "src", "v2ray.com", "core", ".dev", "protoc", "linux", "protoc"),
}

var (
	repo = flag.String("repo", "", "Repo for protobuf generation, such as v2ray.com/core")
)

func main() {
	flag.Parse()

	protofiles := make(map[string][]string)
	protoc := protocMap[runtime.GOOS]
	fmt.Printf("%s\n", protoc)
	protoc = "G:\\Golang\\code\\src\\v2ray.com\\core\\.dev\\protoc\\windows\\protoc.exe"
	gosrc := filepath.Join(os.Getenv("GOPATH"), "src")
	fmt.Printf("%s\n", gosrc)
	gosrc = "G:\\Golang\\code\\src"
	reporoot := filepath.Join(os.Getenv("GOPATH"), "src", *repo)
	fmt.Printf("%s\n", reporoot)
	reporoot = "G:\\Golang\\code\\src"

	filepath.Walk(reporoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		dir := filepath.Dir(path)
		filename := filepath.Base(path)
		if strings.HasSuffix(filename, ".proto") {
			protofiles[dir] = append(protofiles[dir], path)
		}

		return nil
	})

	for _, files := range protofiles {
		args := []string{"--proto_path", gosrc, "--go_out", "plugins=grpc:" + gosrc}
		args = append(args, files...)
		cmd := exec.Command(protoc, args...)
		cmd.Env = append(cmd.Env, os.Environ()...)
		fmt.Printf("%v\n", cmd.Env)
		output, err := cmd.CombinedOutput()
		if len(output) > 0 {
			//fmt.Println(string(output))
			fmt.Printf("%s\n",output)
		}
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}

	common.Must(filepath.Walk(reporoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".pb.go") {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		content = bytes.Replace(content, []byte("\"golang.org/x/net/context\""), []byte("\"context\""), 1)

		pos := bytes.Index(content, []byte("\npackage"))
		if pos > 0 {
			content = content[pos+1:]
		}

		return ioutil.WriteFile(path, content, info.Mode())
	}))
}
