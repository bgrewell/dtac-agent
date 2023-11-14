# DTAC AGENT

[![Audit](https://github.com/intel-innersource/frameworks.automation.dtac.agent/actions/workflows/audit.yml/badge.svg)](https://github.com/intel-innersource/frameworks.automation.dtac.agent/actions/workflows/audit.yml)
[![Goreleaser](https://github.com/intel-innersource/frameworks.automation.dtac.agent/actions/workflows/release.yml/badge.svg)](https://github.com/intel-innersource/frameworks.automation.dtac.agent/actions/workflows/release.yml)

![dtac logo](https://github.com/intel-innersource/frameworks.automation.dtac.agent/blob/main/assets/logo/DTAC.png?raw=true)

The Distributed Telemetry and Advanced Control (DTAC) framework is a collection of projects designed to reduce the time
to completion of software projects and testbeds by providing a highly reusable and extensible framework for the
collection of monitoring and manipulation of a wide variety of systems. 

This project, the DTAC Agent, is focused on the endpoints. It is designed to run on various operating systems including
Windows, Linux and MacOS (Darwin). It provides access to a wide variety of telemetry on these systems and also provides
the ability to control many operating system and application parameters out of the box. The goal of the DTAC agent and
broader DTAC ecosystem is to replace legacy insecure APIs, custom tooling and shell scripts that do things like `ssh user@host <command>`
with a more feature complete, flexible and secure framework for automation and telemetry. To this end it has been designed 
to be highly extensible through a multitude of methodologies described in more detail in the [extensibility](#extensibility) 
section below.

It is important to note that this project is under active development and as such it may have incomplete or missing documentation as well as incomplete implementations of some functionality. While we will do our best to not change the public APIs it is possible that they will have breaking changes prior to the first production release.

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

#### Windows systems (.msi)

*These instructions are not yet available*

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

## Usage

The DTAC Agent is implemented as a REST API and can be consumed by any client that can query REST APIs, including using tools such as curl, Postman, or your browser.

By default, most of the API endpoints are secured and require authentication. To access these endpoints, you will need to provide an authentication token in the Authorization header of your HTTP requests.

You can obtain an authentication token by sending a POST request to the /api/authn/login endpoint with a valid username and password. The response will contain an access_token field that you can use to authenticate subsequent requests.

To retrieve the administrative password you can use the following command:

```bash
sudo dtac config view auth.pass
```

Once this is done you can use a tool like curl to request an access token with a request like shown below:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"username": "<username>", "password": "<password>"}' https://localhost:8180/api/authn/login
```

*Replace `<username>` and `<password>` with your actual username and password.*


Alternatively you can get an access token using the dtac command which will get the credentials from the configuration file directly. An example of this command is shown below

```bash
TOKEN=`sudo dtac token`
```

Once you have a token you can query the API using curl with a command similar to the one below which shows how to retreive the systems uuid:

```bash
curl -ks -H "Authorization: Bearer $TOKEN" https://localhost:8180/system/uuid
```

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
sudo dtac config view auth.pass
```

The format for viewing configuration elements is `dtac config view <path>` where `<path>` is the path to the element you
wish to view. The path is a dot separated list of keys to the element you wish to view. For example if you wanted to
view the entire authentication section of the configuration you could do so with the following command:

```bash
sudo dtac config view auth
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

Writing plugins in Go is fairly simple and straighforward. Below is an example of a struct that implements the plugin interface. It has comments
explaining each section of code and why it is needed. This sample can also be found in `cmd/plugins/hello`

```go
package helloplugin

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins/utility"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/types/endpoint"
	"reflect"
)

// HelloMessage is just a simple helper struct to encapsulate the hello world message
type HelloMessage struct {
	Message string `json:"message"`
}

// This sets a non-existent variable to the interface type of plugin then attempts to assign
// a pointer to HelloPlugin to it. This isn't needed, but it's a good way to ensure that the
// HelloPlugin struct implements the Plugin interface. If there are missing functions, this
// will fail to compile.
var _ plugins.Plugin = &HelloPlugin{}

// NewHelloPlugin is a constructor that returns a new instance of the HelloPlugin
func NewHelloPlugin() *HelloPlugin {
	// Create a new instance of the plugin
	hp := &HelloPlugin{
		PluginBase: plugins.PluginBase{
			Methods: make(map[string]endpoint.Func),
		},
		message: HelloMessage{
			Message: "this is an example of how to create a plugin. See the source at https://github.com/intel-innersource/frameworks.automation.dtac.agent/tree/main/plugin/examples/hello",
		},
	}
	// Ensure we set our root path which will be appended to all of our methods to help namespace them
	hp.SetRootPath("hello")

	// Return the new instance
	return hp
}

// HelloPlugin is the plugin struct that implements the Plugin interface
type HelloPlugin struct {
	// PluginBase provides some helper functions
	plugins.PluginBase
	message HelloMessage
}

// Name returns the name of the plugin type
// NOTE: this is intentionally not a pointer receiver otherwise it wouldn't work. This must be set at your plugin struct
// level. otherwise it will return the type of the PluginBase struct instead.
func (h HelloPlugin) Name() string {
	t := reflect.TypeOf(h)
	return t.Name()
}

// Register registers the plugin with the plugin manager
func (h *HelloPlugin) Register(args plugins.RegisterArgs, reply *plugins.RegisterReply) error {
	*reply = plugins.RegisterReply{Endpoints: make([]*plugins.PluginEndpoint, 0)}

	// Check if the configuration has the message set
	if message, ok := args.Config["message"]; ok {
		h.message = HelloMessage{
			Message: message.(string),
		}
	}

	// Declare our endpoint(s)
	endpoints := []*endpoint.Endpoint{
		{
			Path:           "hello",
			Action:         endpoint.ActionRead,
			UsesAuth:       args.DefaultSecure,
			ExpectedArgs:   nil,
			ExpectedBody:   nil,
			ExpectedOutput: &HelloMessage{},
			Function:       h.Hello,
		},
	}

	// Register them with the plugin
	h.RegisterMethods(endpoints)

	// Convert to plugin endpoints and return
	for _, ep := range endpoints {
		reply.Endpoints = append(reply.Endpoints, plugins.ToAPIEndpoint(ep))
	}

	// Return no error
	return nil
}

// Hello is the handler for the hello world route
func (h *HelloPlugin) Hello(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	// Here we use the utility wrapper to help us add some additional context to the call and simplify the
	// code by having a helper function build the ReturnVal object for us.
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, interface{}, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {h.Name()},
		}

		return headers, h.message, nil
	}, "hello plugin output message")
}

```

Once the HelloPlugin type has been defined you simply need to create a new instance of it and serve it. 

```go
package main

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/hello/helloplugin"
	"log"

	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins"
)

func main() {

	p := helloplugin.NewHelloPlugin()

	h, err := plugins.NewPluginHost(p)
	if err != nil {
		log.Fatal(err)
	}

	err = h.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
```

#### Python

Below is an example code snippet of a functional example plugin written in Python. The code is heavily commented to explain the process. 

```python
#!/usr/bin/env python3
# ^^ the shebang is required so that the script can be executed directly without being called by python ^^

# IMPORTANT NOTES: ################################################################################################################################################################
#   NOTE-1: The file must have the extension .plugin such as example.plugin as this is how the plugins are found by the DTAC Agent
#   Note-2: The plugin must be placed inside the DTAC Agent's plugin directory; default is /opt/dtac/plugins on Linux or C:\Program Files\Intel\dtac-agent\plugins on Windows
###################################################################################################################################################################################

# Import PluginBase and BaseConfig to assist in building your plugin.
# - PluginBase: A base class for creating plugins.
# - BaseConfig: A base configuration class for plugin settings.
from dtac_tools.plugins.base import PluginBase, BaseConfig

# Import various types used by the plugin system:
# - PluginEndpoint: Defines the endpoints your plugin will serve, detailing each endpoint.
# - InputArgs and ReturnVal: Types for input arguments and return values in the plugin system.
# - LoggingLevel: Specifies log levels for the centralized logging facilities.
# - EndpointAction: Enumerates possible actions an endpoint can perform.
from dtac_tools.plugins.types import PluginEndpoint, InputArgs, ReturnVal, LoggingLevel, EndpointAction

# Import DefaultPluginHost for managing the hosting aspects of your plugin.
from dtac_tools.host.default_host import DefaultPluginHost


# Define the ExamplePluginConfig class to handle configuration options for your plugin.
# This class is used for deserializing the configuration during the plugin's registration phase.
# It should consist of serializable types such as str, int, bool, etc.
# In this example, we define a configuration with a single 'message' field of type str.
# This is optional and not needed if you do not have any configuration options.
class ExamplePluginConfig(BaseConfig):
    message: str


# ExamplePlugin here defines a plugin object and it inherits PluginBase
class ExamplePlugin(PluginBase):

    def __init__(self):
        # All plugins get a configuration object that is passed to them by the DTAC Agent. This configuration blob is defined in the main DTAC Agent configuration file under the
        # plugins: section and can be any valid yaml. Whatever values are found below the config: entry in the DTAC Agent main configuration for the specified plugin will be passed
        # without processing or modification to the plugin to be handled internally as you the developer sees fit. Configuration is not required and can be ignored and left out in
        # cases where there are no configurable parameters.
        #
        # In the case of our example here we will only have one configuration element with a key called message and a value of the message we want returned when this endpoint is
        # called.
        #
        # So inside of the main DTAC Agent configuration we may have a section like this which sets our message. In this case that message will be used below to overwrite the default
        # value of self.message
        #
        # plugins:
        #   ...
        #   entries:
        #     example:
        #       enabled: true
        #       hash: <sha256 hash>
        #       config:
        #         message: this is a message that has been overridden

        # This is where we hold our message string and we give it a defualt value (which will be overwritten if we get a confguration from the DTAC Agent at load time with the message value)
        self.message = "this is the default message which can be overridden by dtac"

    # The `root_path` helper function is optional and allows you to customize the path to your plugin.
    # By default, the path is derived from the plugin's filename, excluding its extension.
    # For this plugin, the default path would be 'example', leading to endpoint paths like 'example/<method_name>', e.g., 'example/message'.
    # If you wish to use a different root path, you can override it by implementing this method.
    # For instance, using the REST API protocol, a custom path would be exposed as 'http(s)://<host>:<port>/<plugins_root>/example/message'.
    #
    # Uncomment and modify the following method to customize the plugin's root path:
    # def root_path(self) -> str:
    #     return "custom_path"  # Replace 'custom_path' with your desired root path


    # The `register` function is essential in all plugins and should adhere to this signature.
    # It's responsible for setting up the handlers that the plugin will support.
    # In this example, there's a single handler named 'message' that returns the 'message' attribute of this class.
    def register(self, params) -> list[PluginEndpoint]:
        # The `params` object includes configuration data from the DTAC Agent's main configuration file.
        # Here, we check if the configuration is provided and, if so, load it to set the internal 'message' attribute
        # according to the 'message' value from the configuration.
        if params["config"] is not None:
            config = self.load_config(params["config"], ExamplePluginConfig)
            self.message = config.message

        # `default_secure` is a flag provided by the DTAC Agent indicating whether to use authenticated endpoints by default.
        # The general recommendation is to adhere to this default unless there's a specific need for unauthenticated endpoints.
        # Although you can force authentication by setting `uses_auth` to True for each endpoint, doing so might conflict
        # with the DTAC Agent administrator's settings if `default_secure` is set to false.
        default_secure = True
        if params["default_secure"] is not None:
            default_secure = params["default_secure"]

        # `endpoints` is a list of supported endpoints in the plugin.
        # This example includes a single 'read' endpoint that returns a message.
        # The actions 'read', 'write', 'create', and 'delete' correspond to the HTTP methods GET, POST, PUT, and DELETE, respectively.
        endpoints = []


        # Create an endpoint for a 'read' request to example/message
        message_endpoint = PluginEndpoint(
            function=self.print_message,        # function is the function that will be called when this endpoint is hit
            path="message",                     # path is the path that will be appended to the root_path to create the full path for this endpoint
            action=EndpointAction.ActionRead,   # action is the action that this endpoint will be registered as. This is used to determine what type of request this endpoint will handle
            uses_auth=default_secure,           # uses_auth is a boolean that determines if this endpoint requires authentication or not. generally it should use default_secure which is a setting passed from the DTAC Agent
            expected_args=None,                 # expected args should be an object that represents the arguments you expect the client to pass into the endpoint
            expected_body=None,                 # expected body should represent the body of the request that you expect the client to pass into the endpoint
            expected_output=None,               # expected output should represent the output that you expect the endpoint to return to the client
        )

        # Add the endpoint to the endpoints collection
        endpoints.append(message_endpoint)

        # The following line demonstrate how to use centralized logging to output log messages via the agent.
        # Note that this is only an example and not a requirement in this context. Centralized logging is accessible
        # throughout your plugin.
        # It's highly recommended to use the LevelFatal log level when your plugin encounters a critical error and
        # needs to terminate. This ensures that the error is properly logged by the Agent, and it's clear that the
        # error was fatal, indicating the plugin is about to shut down.
        self.log(LoggingLevel.LevelInfo, f"{self.name()} plugin registered", {"endpoint_count": str(len(endpoints))})

        # Return the endpoints collection
        return endpoints

    # The `print_message` function serves as an example of a handler function.
    # Handler functions accept an InputArgs object and return a ReturnVal object.
    def print_message(self, args: InputArgs) -> ReturnVal:
        return ReturnVal(
            # Specify any response headers to be added to the client's request.
            # The final location of these headers depends on the API protocol used by the DTAC Agent.
            # For instance, in REST, these headers will be part of the response headers.
            headers={},

            # Define any parameters to be returned in the response.
            # These parameters are passed back to the client.
            # Similar to headers, their location in the response depends on the DTAC Agent's API protocol.
            params={},

            # The 'value' field contains the data returned to the caller.
            # Here, it returns the 'message' attribute, which could be a preset value or one from the DTAC Agent configuration.
            # This field can hold any serializable data type, like a string, int, list, dict, etc.
            value={"message": self.message},
        )


def main():
    # This section initializes and sets up the plugin for execution.
    # When placed in the plugins directory with the appropriate extension and shebang,
    # the DTAC Agent will detect and run it, triggering this main() method.
    # The process involves creating an instance of the plugin, registering it with the plugin host,
    # and then waiting for a connection. The plugin host communicates the connection details back to the DTAC Agent.

    plugin = ExamplePlugin()  # Instantiate your plugin
    host = DefaultPluginHost(plugin)  # Create a host to manage your plugin
    host.serve()  # Start the plugin host


if __name__ == '__main__':
     # Running this script directly will typically result in a message indicating that plugins cannot be run standalone,
     # followed by an exit. This is the expected behavior in normal circumstances.
     # However, if you wish to run the plugin for debugging purposes without it exiting, set the environment variable
     # DTAC_PLUGINS=true. This allows the plugin to run, and you should see an output similar to:
     # CONNECT{{ExamplePlugin:example:grpc:tcp:127.0.0.1:45985:plug_api_1.0:[enc=aNECin4i6qwaDMDVNRRs2n9wbJ603SdBjlwUUJ2fGNo%3D]}}
     # This output is a positive indication, as it describes how the plugin loader communicates with your plugin.

     # IMPORTANT: Avoid using standard output (stdout) or standard error (stderr) for any purpose within your plugin.
     # The plugin's output is essential for negotiating connections with the DTAC Agent plugin framework, and any
     # additional output could interfere with this process.
     main()

```
