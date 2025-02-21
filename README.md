微服务 云原生 DevOps



# 环境配置

## habor

### 安装 docker

> harbor 依赖 Docker Engine V20.10.10-ce+ 或者更高。

更新 apt 包索引：

```bash
sudo apt-get update
```

安装必要的依赖包：

```bash
sudo apt-get install \
  ca-certificates \
  curl \
  gnupg \
  lsb-release
```

添加 Docker 的官方 GPG 密钥：

```bash
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo tee /etc/apt/keyrings/docker.asc
```

添加 docker 仓库：

```bash
echo "deb [arch=amd64 signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```

再次更新 apt 包索引：

```bash
sudo apt-get update
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



### 安装 docker-compose

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



### 安装 harbor

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
hostname: # 自己的域名或ip地址
http:
  port: 80
https:
  port: 443
  certificate:  /opt/cert/harbor.crt
  private_key:  /opt/cert/harbor.key
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
