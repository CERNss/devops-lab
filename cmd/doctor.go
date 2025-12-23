package cmd

import (
	"devops-lab/external"
	"fmt"

	"devops-lab/internal/middleware"
	"devops-lab/internal/platform"
	"devops-lab/internal/runtime"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check DevOps Lab environment health",
	RunE: func(cmd *cobra.Command, args []string) error {
		middleware.Info("Running environment diagnostics")
		fmt.Println()

		checkBinaries()
		fmt.Println()

		rt := &runtime.ShellRuntime{}
		checkCluster(rt)
		fmt.Println()

		checkNamespace(rt)
		fmt.Println()

		checkHelmComponents(rt)

		return nil
	},
}

func checkBinaries() {
	middleware.Info("Checking required binaries")

	checkBinary("docker")
	checkBinary("k3d")
	checkBinary("kubectl")
	checkBinary("helm")
}

func checkBinary(name string) {
	if runtime.CommandExists(name) {
		middleware.OK(fmt.Sprintf("%s found", name))
	} else {
		middleware.Fail(fmt.Sprintf("%s not found in PATH", name))
	}
}

func checkCluster(rt runtime.Runtime) {
	middleware.Info("Checking cluster")

	cluster := platform.NewCluster(rt)

	if cluster.Exists(ClusterName) {
		middleware.OK(fmt.Sprintf("cluster '%s' exists", ClusterName))
	} else {
		middleware.Fail(fmt.Sprintf("cluster '%s' not found", ClusterName))
	}
}

func checkNamespace(rt runtime.Runtime) {
	middleware.Info("Checking namespace")

	cluster := platform.NewCluster(rt)
	ns := platform.NewNamespaceManager(rt)

	if !cluster.IsReachable() {
		middleware.Warn("cluster not reachable, skip namespace check")
		return
	}

	if ns.Exists(Namespace) {
		middleware.OK(fmt.Sprintf("namespace '%s' exists", Namespace))
	} else {
		middleware.Warn(fmt.Sprintf("namespace '%s' not found", Namespace))
	}
}

func checkHelmComponents(rt runtime.Runtime) {
	middleware.Info("Checking Helm components")

	helm := platform.NewHelm(rt)
	releases, err := external.LoadHelmReleasesFromJSON(helmStackFile, Namespace)
	if err != nil {
		middleware.Fail(err.Error())
		return
	}

	for _, r := range releases {
		if helm.Exists(r.Name, r.Namespace) {
			middleware.OK("helm release " + r.Name + " installed")
		} else {
			middleware.Warn("helm release " + r.Name + " not installed")
		}
	}

}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
