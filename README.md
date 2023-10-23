# DTAC AGENT

[![Audit](https://github.com/intel-innersource/frameworks.automation.dtac.agent/actions/workflows/audit.yml/badge.svg)](https://github.com/intel-innersource/frameworks.automation.dtac.agent/actions/workflows/audit.yml)
[![Goreleaser](https://github.com/intel-innersource/frameworks.automation.dtac.agent/actions/workflows/release.yml/badge.svg)](https://github.com/intel-innersource/frameworks.automation.dtac.agent/actions/workflows/release.yml)

The Distributed Telemetry and Advanced Control (DTAC) framework is a collection of projects designed to reduce the time
to completion of software projects and testbeds by providing a highly reusable and extensible framework for the
collection of monitoring and manipulation of a wide variety of systems. 

This project, the DTAC Agent, is focused on the endpoints. It is designed to run on various operating systems including
Windows, Linux and MacOS (Darwin). It provides access to a wide variety of telemetry on these systems and also provides
the ability to control many of the system parameters out of the box. It has been designed to be highly extensible 
through a multitude of methodologies described in more detail in the [extensibility](#extensibility) section below.

## Installation

### Install from package

after the install is complete you will need to get the api password for the administrative user that has been generated. 
You can do this by running the following command:

```bash
sudo dtac config view authn.pass
````

### Install from source



### Windows

### Linux

### MacOS (Darwin)

## Usage

### Configuration

### Authentication

### Basic Requests

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
