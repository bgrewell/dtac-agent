# DTAC AGENT

The Distributed Telemetry and Advanced Control (DTAC) framework is a collection of projects designed to reduce the time
to completion of software projects and testbeds by providing a highly reusable and extensible framework for the
collection of monitoring and manipulation of a wide variety of systems. 

This project, the DTAC Agent, is focused on the endpoints. It is designed to run on various operating systems including
Windows, Linux and MacOS (Darwin). It provides access to a wide variety of telemetry on these systems and also provides
the ability to control many of the system parameters out of the box. It has been designed to be highly extensible 
through a multitude of methodologies described in more detail in the [extensibility](#extensibility) section below.

## Installation

Release packages for DTAC-Agent contain scripts to help users install the agent on the operating systems that are 
supported out of the box. This includes Windows (x86 and amd64), Linux (x86 and amd64) and MacOS (amd64). For additional
operating systems or hardware architectures you will need to compile the agent by following the instructions in the 
[compilation](#compilation) section below. In addition to release packages found in this repository if you are on the
Intel network you can install the latest version using the commands below.

### Windows

Click `start` then type `powershell` and `right click` on `Windows Powershell` and select `Run as Administrator` and 
then in the powershell window type in the following

```powershell
powershell -exec bypass -c "(New-Object Net.WebClient).Proxy.Credentials=[Net.CredentialCache]::DefaultNetworkCredentials;iwr('https://software.labs.intel.com/dtac/agent/windows/install.ps1')|iex"
```

### Linux

### MacOS

## Usage

### 

## Compilation











## Basic Info

`/`

**Support:** `Windows` and `Linux`

The home page provides some basic information about the system. The type of information to expect is shown below. While much of this is available in both `Windows` and `Linux` there are some diffrences in the type of information that can be found so when leveraging this information on both types of systems ensure that the fields you are using are available and populated on both systems.

### Host

- hostname
- uptime
- process_count
- os
- platform
- family
- version
- kernel_version
- kernel_arch
- virtualization_system
- host_id

### CPU

- cpu_num
- vendor_id
- family
- model
- stepping
- physical_id
- core_id
- core_count
- model_name
- speed
- cache
- flags
- microcode
 
 ### Memory
 
- total
- available
- used
- used_percent
- free
- active
- inactive
- wired
- laundry
- buffers
- cached
- writeback
- dirty
- shared
- slab
- huge_pages_total
- huge_pages_free
- huge_pages_size

### Network

- name
- interface_index
- hardware_addr
- mtu
- flags
- addresses

### HTTP Routes

Routes section of the general information page `/` lists all of the supported endpoints active on the system. This output can be used to verify that a system supports the endpoint you wish to call. The nature of `dtac-agent` is the flexibility to support multiple operating systems, custom endpoints and plugins and as such what is available system to system can be diffrent so it is always good to programmatically consume this section to verify availaility to endpoints before calling them.

## Network

### Interfaces List

`/network/interfaces`

**Support:** `Windows` and `Linux`

This endpoint provides information about the network interfaces found on the system. This information includes the following

- name
- index
- hardware address
- flags
- addresses
- multicast addresses
- statistics

### Interface Names

`/network/interfaces/names`

**Support:** `Windows` and `Linux`

This endpoint provides a list of names of the interfaces on the system.

- interface name list

### Interface Details

`/network/interface/<name>`

**Support:** `Windows` and `Linux`

This endpoint provides details about the named network interface.

- name
- index
- hardware address
- flags
- addresses
- multicast addresses
- statistics

### QOS Policies

`/network/qos/policies`

**Support:** `Windows`

**Method:** `GET`

This endpoint is used to get Windows based network QoS policies also known as `NetQoSPolicies`
