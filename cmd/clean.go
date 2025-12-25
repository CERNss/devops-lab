package cmd

import (
	"devops-lab/external"
	"devops-lab/internal/middleware"
	"devops-lab/internal/platform"
	"devops-lab/internal/runtime"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	deleteCluster bool
	force         bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove DevOps platform and optionally delete cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		middleware.Warn("Starting DevOps platform cleanup")

		rt := &runtime.ShellRuntime{}

		helm := platform.NewHelm(rt)
		ns := platform.NewNamespaceManager(rt)
		cluster := platform.NewCluster(rt)

		// ======================
		// 1. Uninstall Helm releases
		// ======================

		releases, err := external.LoadHelmReleasesFromJSON(helmStackFile, Namespace)
		if err != nil {
			middleware.Fail(err.Error())
			return err
		}

		for _, r := range releases {
			if helm.Exists(r.Name, r.Namespace) {
				err := helm.Delete(r.Name, r.Namespace)
				if err != nil {
					middleware.Fail(err.Error())
					return err
				}
			}
		}

		// ======================
		// 2. Delete namespace
		// ======================
		if ns.Exists(Namespace) {
			middleware.Warn(fmt.Sprintf("Deleting namespace: %s", Namespace))
			if err := ns.Delete(Namespace); err != nil && !force {
				middleware.Fail(err.Error())
				return err
			}
		} else {
			middleware.Info(fmt.Sprintf("Namespace %s not found, skip", Namespace))
		}

		// ======================
		// 3. Delete cluster (optional)
		// ======================
		if deleteCluster {
			if cluster.Exists(ClusterName) {
				middleware.Warn(fmt.Sprintf("Deleting cluster: %s", ClusterName))
				if err := cluster.Delete(ClusterName); err != nil && !force {
					middleware.Fail(err.Error())
					return err
				}
			} else {
				middleware.Info(fmt.Sprintf("Cluster %s not found, skip", ClusterName))
			}
		} else {
			middleware.Info("Skip cluster deletion (--delete-cluster=false)")
		}

		middleware.OK("Cleanup completed")
		middleware.Info("If resources are stuck, use scripts/emergency-clean.sh")

		return nil
	},
}

func init() {
	cleanCmd.Flags().BoolVar(
		&deleteCluster,
		"delete-cluster",
		true,
		"Delete k3d cluster",
	)

	cleanCmd.Flags().BoolVar(
		&force,
		"force",
		false,
		"Continue cleanup even if some steps fail",
	)

	rootCmd.AddCommand(cleanCmd)
}
