package network

/*
On Linux you can only shape traffic as it exits a network interface. So for example if you wanted to shape traffic going
to/from the local host over eth0 to a server at 1.2.3.4 you could easily create a shaping rule for your outbound traffic
but you couldn't shape your input traffic. To get around this I need to look at doing the following.

1. Create a virtual network pair veth0 and veth1
2. Bridge the veth0 side of it to eth0
3. Migrate network configuration to veth1
4. Apply outbound shaping to the eth0 interface
5. Apply inbound shaping to the veth0 interface
*/

type ShapingRequest struct {
	UplinkInterface *string //UplinkInterface - the network interface to apply the shaping
	DownlinkInterface *string //DownlinkInterface - the network interface to apply the shaping
	SourceAddr *string //SourceAddr - either a single address like 1.2.3.4 or a range like 1.2.3.0/24
	DestinationAddr *string //DestinationAddr - either a single address like 1.2.3.4 or a range like 1.2.3.0/24
	SourcePort *string //SourcePort - either a single port like '80' or a range like '80:90' or pairs like '22,80,443'
	DestinationPort *string //DestinationPort - either a single port like '80' or a range like '80:90' or pairs like '22,80,443'
	Protocol *string //Protocol - the network protocol like tcp, udp etc...
	UplinkRate *string //UplinkRate - the rate to shape uplink traffic at (uplink is defined as from the source to the destination)
	DownlinkRate *string //DownlinkRate - the rate to shape downlink traffic at (downlink is defined as from the destination to the source)
}

// TODO: When called we need to check and see if the interfaces have been properly prepped
//tc qdisc add dev eth0 root handle 1: htb default 10
//tc class add dev eth0 parent 1: classid 1:1 htb rate 1000mbps ceil 1000mbps
//tc class add dev eth0 parent 1:1 classid 1:10 htb rate 1000mbps ceil 1000mbps

// TODO: This is the per rule specific portion but it relies on the above being set already
//tc class add dev eth0 parent 1:1 classid 1:11 htb rate {RATE} ceil {RATE}
//tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip src ${CLIENT_IP} flowid 1:11
