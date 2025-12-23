#!/usr/bin/env bash
set -euo pipefail

# =========================
# Inputs (from env / Go)
# =========================
NAMESPACE="${NAMESPACE:-devops}"
DOMAIN_SUFFIX="${DOMAIN_SUFFIX:-local.test}"
INGRESS_CLASS="${INGRESS_CLASS:-traefik}"
INSTALL_TRAEFIK="${INSTALL_TRAEFIK:-false}"

echo "Generating Helm values (domain: ${DOMAIN_SUFFIX})"

mkdir -p values

# ---------------- Traefik ----------------
cat > values/traefik.yaml <<EOF
service:
  type: LoadBalancer
  ports:
    web:
      port: 80
      protocol: TCP
    websecure:
      port: 443
      protocol: TCP
ingressClass:
  enabled: true
  isDefaultClass: false
  name: ${INGRESS_CLASS}
EOF

# ---------------- GitLab ----------------
cat > values/gitlab.yaml <<EOF
global:
  edition: ce
  hosts:
    domain: ${DOMAIN_SUFFIX}
    externalIP: 127.0.0.1
  ingress:
      enabled: true
      class: ${INGRESS_CLASS}
      provider: traefik
      configureCertmanager: false
  minio:
    enabled: true
  appConfig:
    object_store:
      enabled: true

gitlab:
  webservice:
    minReplicas: 1
    maxReplicas: 1
  sidekiq:
    minReplicas: 1
    maxReplicas: 1

nginx-ingress:
  enabled: false

postgresql:
  install: true
redis:
  install: true

registry:
  enabled: false
prometheus:
  install: false
gitlab-runner:
  install: false
EOF

# ---------------- Jenkins ----------------
cat > values/jenkins.yaml <<EOF
controller:
  admin:
    username: admin
    password: admin123
  serviceType: ClusterIP
  ingress:
    enabled: true
    ingressClassName: ${INGRESS_CLASS}
    hostName: jenkins.${DOMAIN_SUFFIX}
    annotations:
      traefik.ingress.kubernetes.io/router.entrypoints: web
  installPlugins:
    - kubernetes
    - workflow-aggregator
    - git
    - configuration-as-code
    - docker-workflow
    - blueocean
persistence:
  enabled: true
  storageClass: local-path
  size: 20Gi
EOF

# ---------------- Harbor ----------------
cat > values/harbor.yaml <<EOF
expose:
  type: ingress
  ingress:
    className: ${INGRESS_CLASS}
    hosts:
      core: harbor.${DOMAIN_SUFFIX}
    annotations:
      traefik.ingress.kubernetes.io/router.entrypoints: web
  tls:
    enabled: false

externalURL: http://harbor.${DOMAIN_SUFFIX}
harborAdminPassword: "Harbor12345"

persistence:
  enabled: true
  persistentVolumeClaim:
    registry:
      storageClass: local-path
      size: 20Gi
    database:
      storageClass: local-path
      size: 10Gi
    redis:
      storageClass: local-path
      size: 5Gi
EOF

# ---------------- Nexus ----------------
cat > values/nexus.yaml <<EOF
ingress:
  enabled: true
  ingressClassName: ${INGRESS_CLASS}
  hostRepo: nexus.${DOMAIN_SUFFIX}
  annotations:
    traefik.ingress.kubernetes.io/router.entrypoints: web

persistence:
  enabled: true
  storageClass: local-path
  size: 20Gi
EOF

# ---------------- Observability ----------------
cat > values/kube-prometheus-stack.yaml <<EOF
grafana:
  enabled: true
  adminPassword: "Grafana12345"
  ingress:
    enabled: true
    ingressClassName: ${INGRESS_CLASS}
    hosts:
      - grafana.${DOMAIN_SUFFIX}
EOF

cat > values/loki-stack.yaml <<EOF
loki:
  enabled: true
  persistence:
    enabled: true
    storageClassName: local-path
    size: 10Gi
promtail:
  enabled: true
grafana:
  enabled: false
EOF

echo "Helm values generated under ./values/"
