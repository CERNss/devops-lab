package platform

import (
	"devops-lab/internal/runtime"
)

type NamespaceManager struct {
	rt runtime.Runtime
}

func NewNamespaceManager(rt runtime.Runtime) *NamespaceManager {
	return &NamespaceManager{rt: rt}
}

// Exists checks whether the namespace exists
func (n *NamespaceManager) Exists(name string) bool {
	return n.rt.Run("kubectl", "get", "ns", name) == nil
}

// Create creates the namespace (no existence check)
func (n *NamespaceManager) Create(name string) error {
	return n.rt.Run("kubectl", "create", "ns", name)
}

// Delete deletes the namespace
func (n *NamespaceManager) Delete(name string) error {
	return n.rt.Run("kubectl", "delete", "ns", name)
}
