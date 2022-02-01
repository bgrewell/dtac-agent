# SYSTEM API

## TODO BEFORE RELEASE

 - Auto-Restart on update
 - Deploy to all systems
 - Package as Kubernetes Sidecar

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

Routes section of the general information page `/` lists all of the supported endpoints active on the system. This output can be used to verify that a system supports the endpoint you wish to call. The nature of `system-api` is the flexibility to support multiple operating systems, custom endpoints and plugins and as such what is available system to system can be diffrent so it is always good to programmatically consume this section to verify availaility to endpoints before calling them.

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
