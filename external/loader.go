package external

import (
	"devops-lab/internal/middleware"
	"devops-lab/internal/model"
	"encoding/json"
	"os"
	"path/filepath"
)

type helmStackFile struct {
	Releases []model.HelmRelease `json:"releases"`
}

func ResolveRelativeToExecutable(path string) (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	base := filepath.Dir(exe)
	return filepath.Join(base, path), nil
}

func LoadHelmReleasesFromJSON(path string, namespace string) ([]model.HelmRelease, error) {

	path, err := ResolveRelativeToExecutable(path)
	if err != nil {
		middleware.Fail(err.Error())
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var stack helmStackFile
	if err := json.Unmarshal(data, &stack); err != nil {
		return nil, err
	}

	// 注入 namespace（运行时上下文）
	for i := range stack.Releases {
		stack.Releases[i].Namespace = namespace
	}

	return stack.Releases, nil
}
