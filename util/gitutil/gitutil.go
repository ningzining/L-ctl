package gitutil

import (
	"os/exec"
)

const (
	templateGitUrl = "https://github.com/ningzining/L-ctl-template.git"
)

func CloneTemplate(localDir string) ([]byte, error) {
	return Clone(templateGitUrl, localDir)
}

func Clone(targetUrl string, localDir string) ([]byte, error) {
	cmd := exec.Command("git", "clone", targetUrl, localDir)
	return cmd.Output()
}

func Pull(dir string) ([]byte, error) {
	cmd := exec.Command("git", "pull")
	cmd.Dir = dir
	return cmd.Output()
}
