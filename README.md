# DTAC AGENT

[![Audit](https://github.com/intel-innersource/frameworks.automation.dtac.agent/actions/workflows/audit.yml/badge.svg)](https://github.com/intel-innersource/frameworks.automation.dtac.agent/actions/workflows/audit.yml)
[![Goreleaser](https://github.com/intel-innersource/frameworks.automation.dtac.agent/actions/workflows/release.yml/badge.svg)](https://github.com/intel-innersource/frameworks.automation.dtac.agent/actions/workflows/release.yml)

![dtac logo](https://github.com/intel-innersource/frameworks.automation.dtac.agent/blob/main/assets/logo/DTAC.png?raw=true)

The Distributed Telemetry and Advanced Control (DTAC) framework is a collection of projects designed to reduce the time
to completion of software projects and testbeds by providing a highly reusable and extensible framework for the
collection of monitoring and manipulation of a wide variety of systems. 

This project, the DTAC Agent, is focused on the endpoints. It is designed to run on various operating systems including
Windows, Linux and MacOS (Darwin). It provides access to a wide variety of telemetry on these systems and also provides
the ability to control many of the system parameters out of the box. It has been designed to be highly extensible 
through a multitude of methodologies described in more detail in the [extensibility](#extensibility) section below.

## Installation

### Install from package

#### Debian-based systems (.deb)

1. Download the latest .deb package from the releases section of the GitHub repo.
2. Open a terminal and navigate to the directory where the .deb package was downloaded.
3. Run the following command to install the package:
   ```bash
   sudo dpkg -i <package-name>.deb  
   ```
   *Replace `<package-name>` with the actual name of the package you downloaded.*

4. If there are any missing dependencies, run the following command to install them:
   ```bash
   sudo apt-get install -f  
   ```

#### Red Hat-based systems (.rpm)

1. Download the latest .rpm package from the releases section of the GitHub repo.
2. Open a terminal and navigate to the directory where the .rpm package was downloaded.
3. Run the following command to install the package:
   ```bash
   sudo rpm -i <package-name>.rpm  
   ```
   *Replace `<package-name>` with the actual name of the package you downloaded.*

4. If there are any missing dependencies, run the following command to install them:
   ```bash
   sudo yum install -y <dependency-name>  
   ```
   *Replace `<dependency-name>`` with the actual name of the missing dependency.*

### Install from source

#### Prerequisites
 
Before proceeding with the installation, make sure you have the following prerequisites:
- Git
- Go
- Mage
- yq (optional)

#### Clone the repository

1. Open a terminal and clone the repository by running the following command:
   ```bash
   git clone https://github.com/intel-innersource/frameworks.automation.dtac.agent.git 
   ```
2. Navigate to the cloned repository:
   ```bash
   cd frameworks.automation.dtac.agent  
   ```

3. Build the agent, plugins and tools
   ```bash
   go run tools/mage/mage.go build plugins buildCli
   ```

4. Create directories
   ```bash
   sudo mkdir -p /opt/dtac/{bin,plugins} /etc/dtac
   ```

5. Copy files
   
   1. Copy bin/dtac-agentd-<architecture> to /opt/dtac/bin/    dtac-agentd:
      ```bash
      sudo cp bin/dtac-agentd-<architecture> /opt/dtac/bin/dtac-agentd
      ```
      *Replace <architecture> with the actual architecture of your system (e.g., amd64, arm64).*
   
   2. Copy bin/dtac-<architecture> to /opt/dtac/bin/dtac:
      ```bash
      sudo cp bin/dtac-<architecture> /opt/dtac/bin/dtac  
      ```
      *Replace <architecture> with the actual architecture of your system (e.g., amd64, arm64).*
   
   3. Copy any files from bin/plugins/ to /opt/dtac/plugins:
      ```bash
      sudo cp bin/plugins/* /opt/dtac/plugins/
      ```

   4. Create a symlink for the dtac commandline tool:
      ```bash
      sudo ln -s /opt/dtac/bin/dtac /usr/bin/dtac
      ```

6. Create service

   1. Copy the dtac-agentd.service file to /lib/systemd/system/dtac-agentd.service:
      ```bash
      sudo cp service/systemd/dtac-agentd.service /lib/systemd/system/dtac-agentd.service  
      ```

   2. Enable the service:
      ```bash
      sudo systemctl enable dtac-agentd.service  
      ```

   3. Start the service:
      ```bash
      sudo systemctl start dtac-agentd.service  
      ```

7. Update admin password (optional, highly recommended)

   1. If you have yq installed, run the following command to update the password in the /etc/dtac/config.yaml file:
      ```bash
      password=$(openssl rand -base64 32)
      yq eval -i '.authn.pass = "'"$password"'"' /etc/dtac/config.yaml  
      ```

      This will generate a random password and update the authn.pass field in the config.yaml file.

   2. Restart the dtac-agentd service:
      ```bash
      sudo systemctl restart dtac-agentd.service  
      ```

### Windows

### Linux

### MacOS (Darwin)

## Usage

after the install is complete you will need to get the api password for the administrative user that has been generated. 
You can do this by running the following command:

```bash
sudo dtac config view authn.pass
````

### Configuration

The configuration file for the `dtac-agent` is found by default in `/etc/dtac/config.yaml` and contains all of the
settings for the agent. The configuration file is in YAML format and can be edited with any text editor. It is important
to note that currently edits made to the configuration are not automatically reloaded by the agent. To reload the agent
after a change you need to run `sudo systemctl restart dtac-agentd`. Due to the sensitive nature of this configuration
file the default permissions only allow read/write access to the root user. This means that you will need to use `sudo`
when editing the file and be careful not to inadvertently change the permissions on the file which could result in the
credentials being exposed to other users on the system.

The configuration can also be viewed and edited using the `dtac` cli tool located on the system after install. This tool
is designed to make it easier to perform common tasks such as viewing and editing the configuration file. For example
if you want to view the default administrative credentials after install you can do so with the following command:

```bash
sudo dtac config view authn.pass
```

The format for viewing configuration elements is `dtac config view <path>` where `<path>` is the path to the element you
wish to view. The path is a dot separated list of keys to the element you wish to view. For example if you wanted to
view the entire authentication section of the configuration you could do so with the following command:

```bash
sudo dtac config view authn
```

### Authentication

The `dtac` command line configuration tool has been designed to assist with authentication which can be helpful for when
you wish to perform manual operations from the command line. Below is an example of how to get an authentication token
using the `dtac` command line tool which can then be used to perform queries against the API.

```bash
TOKEN=`sudo dtac token`
```

The `dtac` tool will get the certificate and credential information from the configuration file and use that to request
a token from the API. The token will be stored in the `TOKEN` environment variable and can be used in subsequent API 
calls as shown in `Basic Requests` below.

Calls made to secure API endpoints without proper authentication will result in a `401 Unauthorized` response such as 
the one shown below.

```bash
curl -ks -w " | status_code: %{http_code}"  https://localhost:8180/system/uuid
{"time":"2023-10-24T06:39:46.251154486-07:00","error":"invalid authorization header"} | status_code: 401
```

### Basic Requests

The following is an example request to show how to use an API access token to perform a request against the API. The
request will return a uuid identifying the system. The `-ks` options make curl ignore the self signed certificate used
and not show the progress meter or error messages in the output.

```bash
curl -ks -H "Authorization: Bearer $TOKEN" https://localhost:8180/system/uuid
```

## Development

This guide will walk you through preparing your system for DTAC Agent development. We'll cover:

- How to ensure commit messages are consistently formatted. 
- Maintaining standardized versioning for your contributions.
- Installing the compiler and essential tools required to build the project.

Follow the steps below to ensure a smooth and consistent development experience with DTAC Agent.

### System Setup

Development Prerequisites for the DTAC Agent

To efficiently develop and build the DTAC Agent, certain tools are indispensable. These tools facilitate a streamlined 
development process, ensuring high code quality, adherence to coding conventions, and efficient release management. 
The necessary tools are:

- go: The Go programming language, which forms the backbone of the DTAC Agent.
- goreleaser: Simplifies the release process by automating the creation of binaries, packaging, and distribution.
- golint: A linter for Go source code that flags potential style issues to maintain code consistency.
- staticcheck: A state-of-the-art linter for Go that checks for bugs, performance issues, and more.
- commitizen: Assists in creating consistent commit messages, following the "conventional commits" format.
- mage: Mage is a make/rake-like build tool using Go.

For developers using a Linux-based system, specifically Ubuntu 22.04, the aforementioned tools can be installed using 
the commands provided below. It's important to note that while the Go-related tools will be installed using the Go 
package manager, commitizen will be installed using Python's package manager, pip.

```bash
# Install Go
sudo apt update
sudo apt install -y golang-go

# Set GOPATH (assuming you're using bash)
echo "export GOPATH=\$HOME/go" >> ~/.bashrc
echo "export PATH=\$PATH:\$GOPATH/bin" >> ~/.bashrc
source ~/.bashrc

# Install Go-related tools
go get -u golang.org/x/lint/golint
go get -u honnef.co/go/tools/cmd/staticcheck
curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh

# Install pip and commitizen
sudo apt install -y python3-pip
pip3 install --user commitizen
```

Developers are encouraged to familiarize themselves with the functionalities of each tool to leverage them effectively 
during the DTAC Agent development process.

### Using Commitizen for Commits

`commitizen` streamlines the commit message creation process, ensuring messages are consistent and in line with 
conventional commit formats. Instead of using the regular `git commit -m "..." ` approach, follow these steps to 
leverage `commitizen`:

1. **Stage Your Changes**: Before you can commit anything, you need to stage the changes you wish to commit.
   ```bash
   git add .
   ```

2. **Craft the Commit Message**: Instead of manually writing a commit message, simply run:
   ```bash
   cz commit
   ```

   This command will initiate `commitizen`. You will be prompted to:

    - Select the type of change you're committing (e.g., feat, fix, chore, docs, style, refactor, perf, test).
    - Provide a short (under 100 characters) description of the change.
    - (Optional) Provide a longer description.
    - (Optional) List any breaking changes or issues affected by the change.

3. **Finalize the Commit**: Once you've provided the necessary information, `commitizen` will craft a commit message based on your inputs and commit the changes.

4. **Pushing to Remote Repository**: As with any other commit, if you wish to push your changes to a remote repository:
   ```bash
   git push
   ```

Using `commitizen` helps maintain a coherent commit history, which can be particularly beneficial when generating 
changelogs, understanding project history, or navigating through changes.

### Using Mage for Building

`mage` is a build tool that simplifies the process of building the DTAC Agent. It is a make/rake-like tool that uses Go
code to define tasks. For convinence a zero-install option has been included in this repository which allows you to run
mage without installation. To build the DTAC Agent, simply run:

```bash
go run tools/mage/mage.go <target>
```

Where `<target>` is the name of the target you wish to execute. For example, to build the DTAC Agent, run:

```bash
go run tools/mage/mage.go build
```

## Extensibility

### Plugin Development

#### Go

#### Python
