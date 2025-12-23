package platform

import (
	"devops-lab/internal/model"
	"devops-lab/internal/runtime"
)

type Helm struct {
	rt runtime.Runtime
}

func NewHelm(rt runtime.Runtime) *Helm {
	return &Helm{rt: rt}
}

func (h *Helm) Exists(name, namespace string) bool {
	return h.rt.Run("helm", "status", name, "-n", namespace) == nil
}

func (h *Helm) Install(r model.HelmRelease) error {
	if r.Repo != nil && r.Repo.Name != "" && r.Repo.URL != "" {
		if err := h.ensureRepo(*r.Repo); err != nil {
			return err
		}
	}

	args := []string{
		"upgrade", "--install", r.Name, r.Chart,
		"-n", r.Namespace,
	}

	if r.Values != "" {
		args = append(args, "-f", r.Values)
	}
	if r.Timeout != "" {
		args = append(args, "--timeout", r.Timeout)
	}
	if r.Atomic {
		args = append(args, "--atomic")
	}

	return h.rt.Run("helm", args...)
}

func (h *Helm) ensureRepo(repo model.HelmRepo) error {
	// Use --force-update to make this idempotent.
	if err := h.rt.Run("helm", "repo", "add", repo.Name, repo.URL, "--force-update"); err != nil {
		return err
	}
	return h.rt.Run("helm", "repo", "update", repo.Name)
}

func (h *Helm) Delete(name, namespace string) error {
	return h.rt.Run("helm", "uninstall", name, "-n", namespace)
}
