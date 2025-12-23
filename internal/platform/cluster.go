package platform

import (
	"devops-lab/internal/runtime"
	"strconv"
)

type Cluster struct {
	rt runtime.Runtime
}

func NewCluster(rt runtime.Runtime) *Cluster {
	return &Cluster{rt: rt}
}

type CreateClusterOptions struct {
	HTTPPort              string
	HTTPSPort             string
	Servers               int
	Agents                int
	DisableDefaultTraefik bool
}

func DefaultCreateClusterOptions() CreateClusterOptions {
	return CreateClusterOptions{
		Servers:               1,
		Agents:                0,
		DisableDefaultTraefik: true,
	}
}

// IsReachable checks whether kubectl can talk to kube-apiserver
func (c *Cluster) IsReachable() bool {
	// kubectl cluster-info 是最轻量、最稳定的健康探测
	// 不依赖具体资源，不触发 discovery cache 风暴
	return c.rt.RunQuiet(
		"kubectl",
		"cluster-info",
		"--request-timeout=5s",
	) == nil
}

// Exists check whether a cluster exists
func (c *Cluster) Exists(name string) bool {
	return c.rt.Run("k3d", "cluster", "get", name) == nil
}

// Create create a cluster
func (c *Cluster) Create(name string, opts CreateClusterOptions) error {
	args := []string{
		"cluster", "create", name,
		"--servers", strconv.Itoa(opts.Servers),
		"--agents", strconv.Itoa(opts.Agents),
		"-p", opts.HTTPPort + ":80@loadbalancer",
		"-p", opts.HTTPSPort + ":443@loadbalancer",
	}

	if opts.DisableDefaultTraefik {
		args = append(args, "--k3s-arg", "--disable=traefik@server:*")
	}

	return c.rt.Run("k3d", args...)
}

// Delete delete a cluster
func (c *Cluster) Delete(name string) error {
	return c.rt.Run("k3d", "cluster", "delete", name)
}
