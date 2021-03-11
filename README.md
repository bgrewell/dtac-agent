# SYSTEM API

## Basic Info

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

### Routes

Routes section of the general information page `/` lists all of the supported endpoints active on the system. This output can be used to verify that a system supports the endpoint you wish to call. The nature of `system-api` is the flexibility to support multiple operating systems, custom endpoints and plugins and as such what is available system to system can be diffrent so it is always good to programmatically consume this section to verify availaility to endpoints before calling them.

## Network

### Interfaces List

`/network/interfaces`

This endpoint provides information about the network interfaces found on the system. This information includes the following

- Name
- Index
- Hardware Address
- Flags
- Addresses
- Multicast Addresses
- Statistics
