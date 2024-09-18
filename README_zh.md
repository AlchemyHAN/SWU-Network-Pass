# SWU-Network-Pass

**SWU-Network-Pass** 是一款为西南大学设计的校园网络自动登录客户端。截至2024年9月18日，西大校园网络使用锐捷的基于Web的认证系统。该工具基于Go语言开发，旨在保障持续的网络访问不受中断，自动化网络验证和重新登录流程，特别适用于想要在校园网中保持互联网连接的设备。该工具还具有网络状态心跳检查和自动重连功能。

- [简体中文](README_zh.md)
- [English](README.md)

## 功能特点

- **自动网络检测与恢复**：通过检测网络中断并执行自动登录，确保设备持续连接网络。
- **支持加密密码**：传输前加密密码以增强安全性，根据环境设置轻松切换加密功能。
- **跨平台兼容性**：在多种操作系统和架构上高效运行。
- **轻量与高效**：优化资源消耗，非常适合硬件能力有限的设备。
- **用户友好的配置**：简单的设置流程，通过账户文件和环境变量轻松定制。

## 目录

1. [安装](#安装)
2. [配置](#配置)
3. [使用](#使用)
4. [支持的平台](#支持的平台)
5. [如何贡献](#贡献)
6. [许可证](#许可证)

## 安装

从[发布页面](https://github.com/AlchemyHAN/SWU-Network-Pass/releases)下载适合您系统的最新版本。

**安装步骤：**

1. 将下载的文件解压到所需位置。（对于类Unix系统，可以使用以下命令解压tar包：`tar -xzvf SWU-Network-Pass_Linux_arm64.tar.gz`）
2. 确保二进制文件可执行（可能需要在类Unix系统上更改权限）。

## 配置

**账户文件设置：**

- 在可执行文件同一目录下创建名为`accounts.txt`的文件。
- 添加登录凭据，格式为：`username password`

**示例：**

``` plaintext
lilei abc123456@
hanmeimei xyz789012#
```

**环境变量：**

- `SWU_NEED_ENCRYPTION`：设置为`true`（默认）启用密码加密，或设置为`false`禁用。

## 使用

在运行客户端之前，切换到二进制文件所在的目录，然后执行该二进制文件即可启动客户端。
从此，客户端将：

- 持续监控网络状态。
- 在检测到网络中断时，自动使用`accounts.txt`中的凭据尝试登录。

**运行客户端：**

```bash
cd /path/to/binary
./swu-network-pass
```

## 支持的平台

已经为一系列操作系统和架构预编译了此客户端：

- **操作系统和架构：**
  - Windows (x86_64, x86, arm64, armv7)
  - Linux (x86_64, x86, arm64, armv7, armv6, mips64, ppc64, ppc64le, s390x, riscv64, loong64)
  - macOS (Darwin) (x86_64, arm64 (Apple Silicon))

请参阅[发布页面](https://github.com/AlchemyHAN/SWU-Network-Pass/releases)，获取每个支持平台的详细信息。

## 从源码构建

**前提条件：**

- Go 1.21或更高版本

按以下步骤构建客户端：

1. 克隆仓库：`git clone https://github.com/AlchemyHAN/SWU-Network-Pass.git`
2. 更改工作目录：`cd SWU-Network-Pass`
3. 构建客户端：`go build`

## 贡献

非常欢迎并感谢您贡献本项目。要参与贡献，请执行以下操作：

1. Fork 该仓库。
2. 创建您的分支（`git checkout -b feature/myNewFeature`）。
3. 提交您的更改（`git commit -am 'Add some great feature'`）。
4. 推送到分支（`git push origin feature/myNewFeature`）。
5. 创建一个 Pull Request。

## 许可证

该项目根据GPL-3.0许可证授权。有关更多信息，请参阅[LICENSE](LICENSE)文件。

---

**免责声明：** 此软件仅用于西南大学内部的教育和研究目的。本开发者不对该工具造成的任何误用或损害负责。请自行承担使用风险。
