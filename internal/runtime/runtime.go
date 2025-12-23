package runtime

type Runtime interface {
	Run(cmd string, args ...string) error
	RunQuiet(cmd string, args ...string) error
}
