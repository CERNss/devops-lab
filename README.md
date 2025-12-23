## DevOps Lab 部署脚本

该项目通过 Helm 自动化安装 GitLab、Jenkins、Harbor 等 DevOps 组件。支持三种环境模式：
- `k8s`：使用当前 kubeconfig 上下文，不创建集群
- `k3d`：本地创建 k3d 集群并安装组件
- `k3s`：本地 k3s（当前仅使用 kubeconfig，上层保留扩展接口）

### 前置条件
- 安装 `kubectl`, `helm`, `docker`（使用 `k3d` 时还需要安装 `k3d`）
- 添加需要的 Helm 仓库，例如：

```sh
helm repo add gitlab https://charts.gitlab.io
helm repo add jenkins https://charts.jenkins.io
helm repo add harbor https://helm.goharbor.io
helm repo add traefik https://traefik.github.io/charts
helm repo update
```

### 使用方式
默认环境为 `k8s`，直接使用当前 kubeconfig 上下文：

```sh
go run . --env k8s install
```

创建本地 k3d 并安装：

```sh
go run . --env k3d --cluster devops --http-port 8080 --https-port 8443 install
```

使用本地 k3s（当前行为：使用 kubeconfig，不创建集群）：

```sh
go run . --env k3s install
```

### Ingress / Traefik 说明
- `--install-traefik`：安装自定义 Traefik（否则使用集群已有的 IngressController）
- `--ingress-class`：业务组件的 IngressClass 名称（默认 `traefik`）
- `--disable-built-in-traefik`：仅对 `k3d/k3s` 有意义，禁用内置 Traefik

示例：

```sh
# 使用已有 IngressController（例如 nginx）
go run . --env k8s --install-traefik=false --ingress-class=nginx install

# k3d 中禁用内置 Traefik，改用自定义 Traefik
go run . --env k3d --disable-built-in-traefik=true --install-traefik=true install
```

### 命名空间
所有组件默认安装到同一个 namespace（`--namespace` 控制）。
