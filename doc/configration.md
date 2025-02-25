本项目基于 [完整版Kubernetes（K8S）全套入门+微服务实战项目，带你一站式深入掌握K8S核心能力]( https://www.bilibili.com/video/BV1MT411x7GH/?p=100&share_source=copy_web&vd_source=60c88f65925a3eefc5da1bf33dfcd433) 内容进行扩展，详细内容请观看原视频。

> 实验环境：
>
> |  节点  |      IP      |
> | :----: | :----------: |
> | Master | 192.168.64.7 |
> | Node1  | 192.168.64.8 |
> | Node2  | 192.168.64.9 |
>
> 本实验所有的资源均在 devops 命名空间内。
>
> 实验操作目录为 `./shop-flashsale`。



# Habor

## 安装 docker

> harbor 依赖 Docker Engine V20.10.10-ce+ 或者更高。

```bash
sudo apt-get update
```

```bash
sudo apt install docker.io
```

启动 Docker 并设置开机自启：

```bash
sudo systemctl start docker
sudo systemctl enable docker
```

验证 Docker 是否安装成功：

```bash
docker --version
```



## 安装 docker-compose

> harbor 依赖 docker-compose (v1.18.0+) 或 docker compose v2 (docker-compose-plugin)。

以 Docker Compose 1.29.2 版本为例，需要其它版本可访问 [Docker Compose Releases](https://github.com/docker/compose/releases) ：

```bash
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
```

> 或
>
> ```bash
> sudo curl -L "https://github.com/docker/compose/releases/download/$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep -Po '"tag_name": "\K.*\d')" -o /usr/local/bin/docker-compose
> ```

赋予 Docker Compose 执行权限：

```bash
sudo chmod +x /usr/local/bin/docker-compose
```

验证 Docker Compose 是否安装成功：

```bash
docker-compose --version
```



## 安装 harbor

> 在安装 harbor 前，确保您的设备/虚拟机满足 [harbor 安装条件](https://goharbor.io/docs/2.12.0/install-config/installation-prereqs/)。

安装 [harbor 发行包](https://github.com/goharbor/harbor/releases)，以 v2.12.2 offline 安装包为例，在 /opt/harbor 路径下安装：

```bash
sudo mkdir /opt/harbor && cd /opt/harbor
wegt https://github.com/goharbor/harbor/releases/download/v2.12.2/harbor-offline-installer-v2.12.2.tgz
```

> [!note]
>
> 注意：如果您使用 ARM 架构的 MacOS 系统，官方 harbor 发行版本可能会有平台架构 **不兼容** 情况（上述版本默认为 amd 64位操作系统）。
>
> 解决方法：使用以下安装包，仓库地址：https://github.com/wise2c-devops/build-harbor-aarch64
>
> ```bash
> wget https://github.com/wise2c-devops/build-harbor-aarch64/releases/download/v2.10.1/harbor-offline-installer-aarch64-v2.10.1.tgz
> ```

解压安装包：

```bash
tar xzvf harbor-offline-installer-v2.12.2.tgz
```

> arm64：
>
> ```bash
> tar xzvf harbor-offline-installer-aarch64-v2.10.1.tgz
> ```

进入 harbor 目录：

```bash
sudo cd harbor
```

复制 `harbor.yml.tmpl` 文件，作为我们自己的配置文件：

```bash
sudo cp harbor.yml.tmpl harbor.yml
```

> [!tip]
>
> 如果需要 https 访问：
>
> 创建一个目录 /opt/cert 并且进入目录：
>
> ```bash
> sudo mkdir /opt/cert && cd /opt/cert
> ```
>
> 创建一个配置文件 `openssl.cnf`，内容如下：
>
> ```bash
> [req]
> default_bits = 2048
> prompt = no
> default_md = sha256
> req_extensions = req_ext
> distinguished_name = dn
>  
> [dn]
> C = CN
> ST = State
> L = Locality
> O = Organization
> OU = Organizational Unit
> CN = harbor.igmwx.com
>  
> [req_ext]
> subjectAltName = @alt_names
>  
> [alt_names]
> DNS.1 = harbor.igmwx.com
> ```
>
> 生成私钥和证书签名请求 (CSR)：
>
> ```bash
> openssl req -new -sha256 -nodes -out harbor.csr -newkey rsa:2048 -keyout harbor.key -config openssl.cnf
> ```
>
> 使用 CSR 和配置文件生成自签名证书：
>
> ```bash
> openssl x509 -req -in harbor.csr -signkey harbor.key -out harbor.crt -days 365 -extfile openssl.cnf -extensions req_ext
> ```

修改 `harbor.yml` 配置：

```yaml
hostname: 192.168.64.7
http:
  port: 8858
https:
  port: 443
  certificate:  /opt/cert/harbor.crt
  private_key:  /opt/cert/harbor.key
# 可以修改登录密码，默认用户名 admin
harbor_admin_password: Harbor12345
```

> 如果不需要 https，可以自行注释 https 以下内容。

运行安装脚本 `install.sh`：

```bash
sudo ./install.sh
```

此时看到以下内容，则 harbor 安装成功：

```bash
----Harbor has been installed and started successfully.----
```

harbor 安装成功默认启动，可用以下命令手动启动或关闭：

```bash
docker-compose up -d
docker-compose down -v
```



## 配置 secret

创建命名空间：

```bash
kubectl create namespace devops 
```

创建 secret：

```bash
kubectl create secret docker-registry harbor-secret --docker-server=192.168.64.7:8858 --docker-userbname=admin --dokcer-password=Harbor12345 -n devops
```

添加 docker 认证（需要在每个节点添加）：

```bash
sudo vim /etc/docker/daemon.json
```

```json
{
  "insecure-registries": ["192.168.64.7:8858"]
}
```



## docker 登录 harbor

```bash
docker login -uadmin 192.168.64.7:8858
```

***



# NFS 配置

> 所有 NFS 配置都在 devops 命名空间内。

三台机器配置 nfs：

```bash
sudo apt update
sudo apt install nfs-kernel-server
```

```bash
sudo systemctl start nfs-kernel-server
sudo systemctl enable nfs-kernel-server
```

编辑 NFS 配置文件 `/etc/exports`，添加共享目录和访问权限：

```bash
sudo vim /etc/exports
```

```
/data/nfs 192.168.64.0/24(rw,sync,no_subtree_check,no_root_squash)
```

重新加载：

```bash
exportfs -f
systemctl reload nfs-server
```

创建 nfs provisioner、nfs storage class：

> [!NOTE]
>
> 如果是 arm 架构，需要在 [nfs provisioner 配置文件](https://github.com/StrayCa7/shop-flashsale/blob/main/devops/storage/nfs-provisioner-deployment.yaml) 中修改容器镜像地址：
>
> ```yaml
> containers:
>   - name: nfs-client-provisioner
>     image: k8s.gcr.io/sig-storage/nfs-subdir-external-provisioner:v4.0.0 # 支持 arm64版本
>     # image: registry.cn-beijing.aliyuncs.com/pylixm/nfs-subdir-external-provisioner:v4.0.0	# 默认为国内地址
>     # image: quay.io/external_storage/nfs-client-provisioner:latest
> ```

修改 nfs 服务地址：

```yaml
spec:
  serviceAccountName : nfs-client-provisioner
  containers:
    ...
      env:
        - name: PROVISIONER_NAME
          value: fuseim.pri/ifs
        - name: NFS_SERVER
          value: 192.168.64.9 	# 修改
        - name: NFS_PATH
          value: /data/nfs
  volumes:
    - name: nfs-client-root
      nfs:
        server: 192.168.64.9	# 修改
        path: /data/nfs
```

***



# SonarQube

创建 SonarQube：

```bash
kubectl apply -f shop-flashsale/devops/storage
```

```bash
kubectl apply -f shop-flashsale/devops/sonarqube
```

根据 NodePort 开放的端口进行访问。

> 用户名、密码均为 admin。

***



# Jenkins

## 采用 helm 的形式安装。

```bash
helm repo add jenkins https://charts.jenkins.io
helm repo update
```

安装 Jenkins：

```bash
helm install jenkins jenkins/jenkins -n devops
```

> 默认情况下，Jenkins 会使用 `LoadBalancer` 类型的 Service 来暴露它，可能会根据你的环境情况调整此配置。如果你使用的是 Minikube 或某些没有 LoadBalancer 支持的环境，可以使用 `NodePort` 或 `port-forward`。

配置 NodePort：

1. 采用 Helm Chart 中的 Jenkins 服务配置。

	```bash
	cd devops/jenkins
	```

	```bash
	helm show values jenkins/jenkins > jenkins-values.yaml
	```

	设置 Jenkins 服务的类型为 NodePort，并且为它分配一个端口 30000，你可以根据需要调整这个端口：

	```bash
	# For minikube, set this to NodePort, elsewhere uses LoadBalancer
	# Use ClusterIP if your setup includes ingress controller
	# -- k8s service type
	serviceType: NodePort	# 原本为  ClusterIP
	
	# -- k8s service clusterIP. Only used if serviceType is ClusterIP
	clusterIp:
	# -- k8s service port
	servicePort: 8080
	# -- k8s target port
	targetPort: 8080
	# -- k8s node port. Only used if serviceType is NodePort
	nodePort: 30000
	```

	修改完 `values.yaml` 后，使用以下命令更新 Helm 部署：

	```bash
	helm upgrade jenkins jenkins/jenkins -f jenkins-values.yaml -n devops
	```

	如果是第一次安装 Jenkins，也可以直接通过以下命令安装：

	```bash
	helm install jenkins jenkins/jenkins -f jenkins-values.yaml -n devops
	```

2. 使用 `kubectl expose` 命令修改服务：

	如果你已经通过 Helm 安装了 Jenkins，并且想在不修改 Helm 配置的情况下暴露 NodePort，可以使用 `kubectl expose` 命令：

	```bash
	kubectl expose svc jenkins --type=NodePort --name=jenkins-nodeport --port=8080 --target-port=8080 --node-port=30000 -n devops
	```

3. 手动配置 NodePort：

	可用以下命令查看 Jenkins chart 的标签：

	```bash
	kubectl describe sts jenkins -n devops | grep Labels
	```

	将 [jenkins-nodeport 配置文件](https://github.com/StrayCa7/shop-flashsale/blob/main/devops/jenkins/jenkins-nodeport.yaml) 中的 `selector` 改为所显示的 label。

获取 Jenkins 初始密码（初始用户为 admin）：

```bash
kubectl get secret -n devops jenkins -o jsonpath='{.data.jenkins-admin-password}' | base64 --decode; echo
```

> 或者使用 Jenkins 日志 ：
>
> ```bash
> kubectl logs -f pod/jenkins-0 -n devops
> ```

访问 `192.168.64.7:30000` 登录 Jenkins。



## 配置数据持久化（可选）

> 如果不配置持久化 jenkins 会默认使用 `local-path` 类型的 Storage Class 创建一个 PVC。
>
> ```bash
> root@master:shop-flashsale/devops/jenkins# kubectl get pvc -n devops
> NAME             STATUS   VOLUME    CAPACITY   ACCESS MODES   STORAGECLASS          VOLUMEATTRIBUTESCLASS   AGE
> jenkins          Bound    pvc-...   8Gi        RWO            local-path            <unset>                 1h
> ```

修改 `jenkins-values.yaml` 配置文件：

```yaml
# Enable persistence for Jenkins home directory
persistence:
  enabled: true
  storageClass: "your-storage-class"  # 可以选择特定的存储类
  size: 10Gi  # 配置 PVC 的大小，可以根据需求调整
  accessMode: ReadWriteOnce  # 配置访问模式
```

检查 PVC 是否生效：

```bash
kubectl get pvc -n devops
```

















