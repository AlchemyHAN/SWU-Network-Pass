# SWU-Network-Pass

**SWU-Network-Pass** is a robust automatic login client for Southwest University (Chongqing, China) campus network. Until 18th, September, 2024, Southwest University campus network uses RuiJie web-based authentication system. Designed to ensure continuous internet access without interruptions, this Go based tool automates network verification and re-login procedures. It's especially useful for managing network connectivity issues seamlessly with built-in network status checks and automatic reconnection.

## Features

- **Automatic Network Detection and Recovery**: Ensures your device maintains uninterrupted network access by detecting disruptions and performing automatic logins.
- **Support for Encrypted Passwords**: Enhances security by encrypting passwords before transmission. Easy toggling of encryption based on environment settings.
- **Cross-Platform Compatibility**: Runs efficiently across a variety of operating systems and architectures.
- **Lightweight and Efficient**: Optimized for minimal resource consumption, perfect for devices with limited hardware capabilities.
- **User-Friendly Configuration**: Simple setup with a straightforward accounts file and environment variables for customization.

## Table of Contents

1. [Installation](#installation)
2. [Configuration](#configuration)
3. [Usage](#usage)
4. [Supported Platforms](#supported-platforms)
5. [Contributing](#contributing)
6. [License](#license)

## Installation

Download the latest release suitable for your system from the [Releases](https://github.com/AlchemyHAN/SWU-Network-Pass/releases) page.

**Installation Steps:**
1. Extract the downloaded file to a desired location.
(For Unix-like systems, you can use the following command to extract the tarball: `tar -xzvf SWU-Network-Pass_Linux_arm64.tar.gz`)
2. Ensure the binary is executable (you might need to change permissions on Unix-like systems).

## Configuration

**Accounts File Setup:**
- Create a file named `accounts.txt` in the same directory as the executable.
- Add your login credentials in the format: `username password`

**Example:**
```
lilei abc123456@
hanmeimei xyz789012#
```

**Environment Variables:**
- `SWU_NEED_ENCRYPTION`: Set to `true` (default) to enable password encryption or `false` to disable it.

## Usage

Change working directory to the location of the binary before running the client, then execute the binary simply to start the client.
The client will:
- Monitor the network status continuously.
- Automatically attempt to login using the credentials from `accounts.txt` when it detects network disruptions.

**Running the Client:**
```bash
cd /path/to/binary
./swu-network-pass
```

## Supported Platforms

This client is precompiled for a range of operating systems and architectures:

- **Operating Systems**: Linux, Windows, macOS (Darwin)
- **Architectures**: amd64, arm64, armv6, armv7, i386, ppc64, ppc64le, s390x, loongarch64, riscv64

Refer to the [Releases](https://github.com/AlchemyHAN/SWU-Network-Pass/releases) page for detailed information about each supported platform.

## Contributing

Contributions are welcome and greatly appreciated. To contribute:
1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/myNewFeature`).
3. Commit your changes (`git commit -am 'Add some great feature'`).
4. Push to the branch (`git push origin feature/myNewFeature`).
5. Open a Pull Request.

## License

This project is licensed under the GPL-3.0 License. For more information, refer to the [LICENSE](LICENSE) file.

---

**Disclaimer:** This software is intended for educational and research purposes in Southwest University only. The developer is not responsible for any misuse or damage caused by this tool. Use at your own risk.
```