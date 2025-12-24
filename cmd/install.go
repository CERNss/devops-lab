package cmd

import (
	"devops-lab/internal/middleware"
	"devops-lab/internal/runtime"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Initialize cluster and install DevOps platform",
	RunE: func(cmd *cobra.Command, args []string) error {
		middleware.Info("Starting DevOps platform installation")

		rt := &runtime.ShellRuntime{}
		if err := setupCluster(cmd, rt); err != nil {
			return err
		}
		if err := ensureNamespace(rt); err != nil {
			return err
		}
		if err := generateValues(); err != nil {
			return err
		}
		if err := installHelmReleases(rt); err != nil {
			return err
		}

		middleware.OK("Installation completed 🎉")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
