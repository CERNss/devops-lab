package cmd

import (
	"devops-lab/external"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	ClusterName  string
	Namespace    string
	DomainSuffix string

	HttpPort  string
	HttpsPort string

	helmStackFile string
	// disable k3d bundled Traefik so we can install our own controller
	DisableDefaultTraefik bool

	ClusterEnv     string
	IngressClass   string
	InstallTraefik bool
)

var allowedClusterEnvList = []string{"k8s", "k3d", "k3s"}

var allowedClusterEnvSet = map[string]struct{}{
	"k8s": {},
	"k3d": {},
	"k3s": {},
}

var rootCmd = &cobra.Command{
	Use:   "devops-lab",
	Short: "DevOps Lab Platform CLI",
	Long:  "DevOps Lab - local DevOps platform based on k3d + Helm",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&ClusterName,
		"cluster",
		"devops",
		"k3d cluster name",
	)

	rootCmd.PersistentFlags().StringVar(
		&Namespace,
		"namespace",
		"devops",
		"k8s namespace",
	)

	rootCmd.PersistentFlags().StringVar(
		&DomainSuffix,
		"domain",
		"local.test",
		"ingress domain suffix",
	)

	rootCmd.PersistentFlags().StringVar(
		&HttpPort,
		"http-port",
		"8080",
		"host HTTP port mapped to ingress (k3d)",
	)

	rootCmd.PersistentFlags().StringVar(
		&HttpsPort,
		"https-port",
		"8443",
		"host HTTPS port mapped to ingress (k3d)",
	)

	rootCmd.PersistentFlags().StringVar(
		&helmStackFile,
		"helm-stack",
		external.DefaultHelmStackPath,
		"Helm stack definition file (json)",
	)

	rootCmd.PersistentFlags().BoolVar(
		&DisableDefaultTraefik,
		"disable-built-in-traefik",
		true,
		"disable bundled Traefik for k3d/k3s so we can provision our own controller",
	)

	ClusterEnv = "k8s"
	rootCmd.PersistentFlags().Var(
		newEnumValue(&ClusterEnv, allowedClusterEnvSet),
		"env",
		"cluster environment: k8s, k3d, or k3s",
	)

	rootCmd.PersistentFlags().BoolVar(
		&InstallTraefik,
		"install-traefik",
		false,
		"install traefik release; otherwise use existing ingress controller",
	)

	rootCmd.PersistentFlags().StringVar(
		&IngressClass,
		"ingress-class",
		"traefik",
		"ingress class name for app ingresses",
	)
}

type enumValue struct {
	value   *string
	allowed map[string]struct{}
}

func newEnumValue(target *string, allowed map[string]struct{}) *enumValue {
	return &enumValue{
		value:   target,
		allowed: allowed,
	}
}

func (e *enumValue) String() string {
	if e == nil || e.value == nil {
		return ""
	}
	return *e.value
}

func (e *enumValue) Set(val string) error {
	if _, ok := e.allowed[val]; !ok {
		return fmt.Errorf("invalid value %q, allowed: %s", val, allowedClusterEnvHint())
	}
	*e.value = val
	return nil
}

func (e *enumValue) Type() string {
	return "string"
}

func allowedClusterEnvHint() string {
	return strings.Join(allowedClusterEnvList, ", ")
}
