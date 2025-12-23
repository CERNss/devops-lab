#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="${CLUSTER_NAME:-devops}"
NAMESPACE="${NAMESPACE:-devops}"
DELETE_CLUSTER="${DELETE_CLUSTER:-true}"

echo "==============================================="
echo "⚠️  EMERGENCY CLEANUP SCRIPT"
echo "-----------------------------------------------"
echo "This script is for LAST RESORT cleanup only."
echo "Prefer: devops-lab clean"
echo "-----------------------------------------------"
echo "Cluster   : ${CLUSTER_NAME}"
echo "Namespace : ${NAMESPACE}"
echo "Delete k3d cluster: ${DELETE_CLUSTER}"
echo "==============================================="
echo

read -rp "❗ Type 'YES' to continue: " CONFIRM
if [[ "${CONFIRM}" != "YES" ]]; then
  echo "Aborted."
  exit 0
fi

echo
echo "[1/3] Force uninstall Helm releases"

RELEASES=(gitlab jenkins harbor nexus monitoring loki)

for r in "${RELEASES[@]}"; do
  echo "  - Trying to uninstall $r"
  helm uninstall "$r" -n "${NAMESPACE}" --ignore-not-found || true
done

echo
echo "[2/3] Delete namespace (force if needed)"

kubectl delete ns "${NAMESPACE}" --wait=false || true

echo
echo "[3/3] Delete k3d cluster"

if [[ "${DELETE_CLUSTER}" == "true" ]]; then
  k3d cluster delete "${CLUSTER_NAME}" || true
else
  echo "  - Keep cluster"
fi

echo
echo "⚠️ Emergency cleanup finished."