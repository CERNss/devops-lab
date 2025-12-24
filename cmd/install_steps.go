package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"devops-lab/external"
	"devops-lab/internal/middleware"
	"devops-lab/internal/model"
	"devops-lab/internal/platform"
	"devops-lab/internal/runtime"

	"github.com/spf13/cobra"
)

func setupCluster(cmd *cobra.Command, rt runtime.Runtime) error {
	middleware.Info("===== Cluster =====")
	cluster := platform.NewCluster(rt)
	if ClusterEnv == "k8s" && cmd.Flags().Changed("disable-built-in-traefik") {
		middleware.Info("disable-built-in-traefik ignored for env=k8s")
	}

	switch ClusterEnv {
	case "k8s":
		middleware.Info("Using existing cluster context (kubeconfig)")
		if !cluster.IsReachable() {
			err := errors.New("kubectl cannot reach the current cluster context")
			middleware.Fail(err.Error())
			return err
		}
	case "k3s":
		middleware.Info("k3s selected; using current kubeconfig context for now")
		if !cluster.IsReachable() {
			err := errors.New("kubectl cannot reach the current cluster context (k3s)")
			middleware.Fail(err.Error())
			return err
		}
	case "k3d":
		if cluster.Exists(ClusterName) {
			middleware.Info("Cluster already exists, skip creation")
		} else {
			middleware.Info("Creating k3d cluster")
			opts := platform.DefaultCreateClusterOptions()
			opts.HTTPPort = HttpPort
			opts.HTTPSPort = HttpsPort
			opts.DisableDefaultTraefik = DisableDefaultTraefik

			err := cluster.Create(ClusterName, opts)
			if err != nil {
				middleware.Fail(err.Error())
				return err
			}
		}
	default:
		err := fmt.Errorf("invalid env: %s (expected %s)", ClusterEnv, allowedClusterEnvHint())
		middleware.Fail(err.Error())
		return err
	}

	middleware.Info("Waiting for cluster to be ready...")
	time.Sleep(5 * time.Second)
	return nil
}

func ensureNamespace(rt runtime.Runtime) error {
	middleware.Info("===== Namespace =====")
	ns := platform.NewNamespaceManager(rt)
	if ns.Exists(Namespace) {
		middleware.Info("Namespace already exists, skip creation")
		return nil
	}

	middleware.Info("Creating namespace: " + Namespace)
	if err := ns.Create(Namespace); err != nil {
		middleware.Fail(err.Error())
		return err
	}
	return nil
}

func generateValues() error {
	middleware.Info("===== Generate Values =====")
	path, err := external.ResolveRelativeToExecutable("scripts/gen-values.sh")
	if err != nil {
		middleware.Fail(err.Error())
		return err
	}
	if err := runtime.RunShell(
		path,
		map[string]string{
			"NAMESPACE":       Namespace,
			"DOMAIN_SUFFIX":   DomainSuffix,
			"INGRESS_CLASS":   IngressClass,
			"INSTALL_TRAEFIK": strconv.FormatBool(InstallTraefik),
		},
	); err != nil {
		middleware.Fail(err.Error())
		return err
	}
	return nil
}

func installHelmReleases(rt runtime.Runtime) error {
	middleware.Info("===== Helm Install =====")
	helm := platform.NewHelm(rt)
	releases, err := external.LoadHelmReleasesFromJSON(
		helmStackFile, // flag or default
		Namespace,
	)
	if err != nil {
		middleware.Fail(err.Error())
		return err
	}

	if !InstallTraefik {
		middleware.Info("Skipping traefik release (install-traefik=false)")
		releases = filterTraefikRelease(releases)
	}

	for _, r := range releases {
		if helm.Exists(r.Name, r.Namespace) {
			middleware.Info(r.Name + " already installed, upgrading")
		}
		if err := helm.Install(r); err != nil {
			middleware.Fail(err.Error())
			return err
		}
	}
	return nil
}

func filterTraefikRelease(releases []model.HelmRelease) []model.HelmRelease {
	filtered := make([]model.HelmRelease, 0, len(releases))
	for _, r := range releases {
		if r.Name == "traefik" {
			continue
		}
		filtered = append(filtered, r)
	}
	return filtered
}
