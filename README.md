![NanoBus Logo](docs/images/nanobus-logo.svg)

NanoBus is a lightweight application runtime that reduces developer responsibility so that teams can **focus on core logic**.

Fore detailed information see the [overview](./docs/overview.md) and [architecture](./docs/architecture.md) pages.

## Install

**Linux** - Install from Terminal to `/usr/local/bin`:

```shell
wget -q https://raw.githubusercontent.com/nanobus/nanobus/main/install/install.sh -O - | /bin/bash
```

**MacOS** - Install from Terminal to `/usr/local/bin`:

```shell
curl -fsSL https://raw.githubusercontent.com/nanobus/nanobus/main/install/install.sh | /bin/bash
```

**Windows** - Install from Command Prompt:

```shell
powershell -Command "iwr -useb https://raw.githubusercontent.com/nanobus/nanobus/main/install/install.ps1 | iex"
```

**Note**: Updates to PATH might not be visible until you restart your terminal application.

## Create a NanoBus Application

## Developer Setup

### Dependencies

To setup a local development environment

| Dependency | Check            | Description                                   |
|:---------- |:---------------- |:--------------------------------------------- |
| [go]       | $ go version     | Go compiler.  Ensure $HOME/go/bin is in PATH. |
| [just]     | $ just --version | Like Makefile [just] runs the needed commands |

### Install from source

```shell
git clone https://github.com/nanobus/nanobus.git
cd nanobus
just install
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on the code of conduct and the process for submitting pull requests.

## License

This project is licensed under the [Mozilla Public License Version 2.0](https://mozilla.org/MPL/2.0/).

[apex]: https://apexlang.io/docs/getting-started
[apexlang.io]: https://apexlang.io
[docker]: https://docs.docker.com/engine/install/
[docker-compose]: https://docs.docker.com/compose/install/
[go]: https://go.dev/doc/install
[iota]: https://github.com/nanobus/iota
[iotas]: https://github.com/nanobus/iota
[just]: https://github.com/casey/just#Installation
[nanobus]: https://github.com/nanobus/nanobus#Install
[npm]: https://docs.npmjs.com/downloading-and-installing-node-js-and-npm
[npx]: https://www.npmjs.com/package/npx#Install
[postgres]: https://www.postgresql.org/download/
[postgresql database]: https://www.postgresql.org/
[rust]: https://rustup.rs/
