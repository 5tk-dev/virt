package virt

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

/*
Network options:

-netdev l2tpv3,id=str,src=srcaddr,dst=dstaddr[,srcport=srcport][,dstport=dstport]
         [,rxsession=rxsession],txsession=txsession[,ipv6=on|off][,udp=on|off]
         [,cookie64=on|off][,counter][,pincounter][,txcookie=txcookie]
         [,rxcookie=rxcookie][,offset=offset]
                configure a network backend with ID 'str' connected to
                an Ethernet over L2TPv3 pseudowire.
                Linux kernel 3.3+ as well as most routers can talk
                L2TPv3. This transport allows connecting a VM to a VM,
                VM to a router and even VM to Host. It is a nearly-universal
                standard (RFC3931). Note - this implementation uses static
                pre-configured tunnels (same as the Linux kernel).
                use 'src=' to specify source address
                use 'dst=' to specify destination address
                use 'udp=on' to specify udp encapsulation
                use 'srcport=' to specify source udp port
                use 'dstport=' to specify destination udp port
                use 'ipv6=on' to force v6
                L2TPv3 uses cookies to prevent misconfiguration as
                well as a weak security measure
                use 'rxcookie=0x012345678' to specify a rxcookie
                use 'txcookie=0x012345678' to specify a txcookie
                use 'cookie64=on' to set cookie size to 64 bit, otherwise 32
                use 'counter=off' to force a 'cut-down' L2TPv3 with no counter
                use 'pincounter=on' to work around broken counter handling in peer
                use 'offset=X' to add an extra offset between header and data
-netdev socket,id=str[,fd=h][,listen=[host]:port][,connect=host:port]
                configure a network backend to connect to another network
                using a socket connection
-netdev socket,id=str[,fd=h][,mcast=maddr:port[,localaddr=addr]]
                configure a network backend to connect to a multicast maddr and port
                use 'localaddr=addr' to specify the host address to send packets from
-netdev socket,id=str[,fd=h][,udp=host:port][,localaddr=host:port]
                configure a network backend to connect to another network
                using an UDP tunnel
-netdev stream,id=str[,server=on|off],addr.type=inet,addr.host=host,addr.port=port[,to=maxport][,numeric=on|off][,keep-alive=on|off][,mptcp=on|off][,addr.ipv4=on|off][,addr.ipv6=on|off][,reconnect-ms=milliseconds]
-netdev stream,id=str[,server=on|off],addr.type=unix,addr.path=path[,abstract=on|off][,tight=on|off][,reconnect-ms=milliseconds]
-netdev stream,id=str[,server=on|off],addr.type=fd,addr.str=file-descriptor[,reconnect-ms=milliseconds]
                configure a network backend to connect to another network
                using a socket connection in stream mode.
-netdev dgram,id=str,remote.type=inet,remote.host=maddr,remote.port=port[,local.type=inet,local.host=addr]
-netdev dgram,id=str,remote.type=inet,remote.host=maddr,remote.port=port[,local.type=fd,local.str=file-descriptor]
                configure a network backend to connect to a multicast maddr and port
                use ``local.host=addr`` to specify the host address to send packets from
-netdev dgram,id=str,local.type=inet,local.host=addr,local.port=port[,remote.type=inet,remote.host=addr,remote.port=port]
-netdev dgram,id=str,local.type=unix,local.path=path[,remote.type=unix,remote.path=path]
-netdev dgram,id=str,local.type=fd,local.str=file-descriptor
                configure a network backend to connect to another network
                using an UDP tunnel

-netdev af-xdp,id=str,ifname=name[,mode=native|skb][,force-copy=on|off]
         [,queues=n][,start-queue=m][,inhibit=on|off][,sock-fds=x:y:...:z]
         [,map-path=/path/to/socket/map][,map-start-index=i]
                attach to the existing network interface 'name' with AF_XDP socket
                use 'mode=MODE' to specify an XDP program attach mode
                use 'force-copy=on|off' to force XDP copy mode even if device supports zero-copy (default: off)
                use 'inhibit=on|off' to inhibit loading of a default XDP program (default: off)
                with inhibit=on,
                  use 'sock-fds' to provide file descriptors for already open AF_XDP sockets
                  added to a socket map in XDP program.  One socket per queue.
                  use 'map-path' to provide the socket map location to populate AF_XDP sockets with,
                  and use 'map-start-index' to specify the starting index for the map (default: 0) (Since 10.1)
                use 'queues=n' to specify how many queues of a multiqueue interface should be used
                use 'start-queue=m' to specify the first queue that should be used


*/

/*
original command

	-nic [tap|bridge|passt|user|l2tpv3|vde|af-xdp|vhost-user|socket][,option][,...][mac=macaddr]
					initialize an on-board / default host NIC (using MAC address
					macaddr) and connect it to the given host network backend
	-nic none       use it alone to have zero network devices (the default is to
					provided a 'user' network connection)
*/
type NicType int

const (
	Tap       NicType = iota // tap
	Bridge                   // bridge
	Passt                    // passt
	User                     // user
	L2tpv3                   // l2tpv3
	Vde                      // vde
	AfXdp                    // af-xdp
	VhostUser                // vhost-user
	Socket                   // socket
)

func (n NicType) String() string {
	switch n {
	case Tap:
		return "tap"
	case Bridge:
		return "bridge"
	case Passt:
		return "passt"
	case User:
		return "user"
	case L2tpv3:
		return "l2tpv3"
	case Vde:
		return "vde"
	case AfXdp:
		return "af-xdp"
	case VhostUser:
		return "vhost-user"
	case Socket:
		return "socket"
	default:
		return "none"
	}
}

func (n *NicType) MarshalYAML() (interface{}, error) {
	return n.String(), nil
}

func (n *NicType) UnmarshalYAML(value *yaml.Node) error {
	switch value.Value {
	default:
		return fmt.Errorf("status inválido: %s", value.Value)
	case "tap":
		*n = Tap
	case "bridge":
		*n = Bridge
	case "passt":
		*n = Passt
	case "user":
		*n = User
	case "l2tpv3":
		*n = L2tpv3
	case "vde":
		*n = Vde
	case "af-xdp":
		*n = AfXdp
	case "vhost-user":
		*n = VhostUser
	case "socket":
		*n = Socket
	}
	return nil
}

type NicOptions struct {
	Type   NicType `yaml:",omitempty"`
	Mac    string  `yaml:",omitempty"`
	Option string  `yaml:",omitempty"`
}

func (n *NicOptions) ToArgs() []string {
	args := []string{n.Type.String()}
	if n.Option != "" {
		args = append(args, n.Option)
	}
	if n.Mac != "" {
		args = append(args, n.Mac)
	}
	return []string{"-nic", strings.Join(args, ",")}
}

/*
original command:

-netdev user,id=str[,ipv4=on|off][,net=addr[/mask]][,host=addr]

	[,ipv6=on|off][,ipv6-net=addr[/int]][,ipv6-host=addr]
	[,restrict=on|off][,hostname=host][,dhcpstart=addr]
	[,dns=addr][,ipv6-dns=addr][,dnssearch=domain][,domainname=domain]
	[,tftp=dir][,tftp-server-name=name][,bootfile=f][,hostfwd=rule][,guestfwd=rule][,smb=dir[,smbserver=addr]]
	       configure a user mode network backend with ID 'str',
	       its DHCP server and optional services
*/
type Netdev_UserOptions struct {
	ID               string `yaml:",omitempty"` // id=str
	Ipv4             string `yaml:",omitempty"` // [,ipv4=on|off]
	Net              string `yaml:",omitempty"` // [,net=addr[/mask]]
	Host             string `yaml:",omitempty"` // [,host=addr]
	Ipv6             string `yaml:",omitempty"` // [,ipv6=on|off]
	Ipv6_net         string `yaml:",omitempty"` // [,ipv6-net=addr[/int]]
	Ipv6_host        string `yaml:",omitempty"` // [,ipv6-host=addr]
	Restrict         string `yaml:",omitempty"` // [,restrict=on|off]
	Hostname         string `yaml:",omitempty"` // [,hostname=host]
	Dhcpstart        string `yaml:",omitempty"` // [,dhcpstart=addr]
	Dns              string `yaml:",omitempty"` // [,dns=addr]
	Ipv6_dns         string `yaml:",omitempty"` // [,ipv6-dns=addr]
	Dnssearch        string `yaml:",omitempty"` // [,dnssearch=domain]
	Domainname       string `yaml:",omitempty"` // [,domainname=domain]
	Tftp             string `yaml:",omitempty"` // [,tftp=dir]
	Tftp_server_name string `yaml:",omitempty"` // [,tftp-server-name=name]
	Bootfile         string `yaml:",omitempty"` // [,bootfile=f]
	Hostfwd          string `yaml:",omitempty"` // [,hostfwd=rule]
	Guestfwd         string `yaml:",omitempty"` // [,guestfwd=rule]
	Smb              string `yaml:",omitempty"` // [,smb=dir[,smbserver=addr]]
}

func (n *Netdev_UserOptions) ToArgs() []string {
	args := []string{"user"}
	if n.ID != "" {
		args = append(args, fmt.Sprintf("id=%s", n.ID))
	}
	if n.Ipv4 != "" {
		args = append(args, fmt.Sprintf("ipv4=%s", n.Ipv4))
	}
	if n.Net != "" {
		args = append(args, fmt.Sprintf("net=%s", n.Net))
	}
	if n.Host != "" {
		args = append(args, fmt.Sprintf("host=%s", n.Host))
	}
	if n.Ipv6 != "" {
		args = append(args, fmt.Sprintf("ipv6=%s", n.Ipv6))
	}
	if n.Ipv6_net != "" {
		args = append(args, fmt.Sprintf("ipv6-net=%s", n.Ipv6_net))
	}
	if n.Ipv6_host != "" {
		args = append(args, fmt.Sprintf("ipv6-host=%s", n.Ipv6_host))
	}
	if n.Restrict != "" {
		args = append(args, fmt.Sprintf("restrict=%s", n.Restrict))
	}
	if n.Hostname != "" {
		args = append(args, fmt.Sprintf("hostname=%s", n.Hostname))
	}
	if n.Dhcpstart != "" {
		args = append(args, fmt.Sprintf("dhcpstart=%s", n.Dhcpstart))
	}
	if n.Dns != "" {
		args = append(args, fmt.Sprintf("dns=%s", n.Dns))
	}
	if n.Ipv6_dns != "" {
		args = append(args, fmt.Sprintf("ipv6-dns=%s", n.Ipv6_dns))
	}
	if n.Dnssearch != "" {
		args = append(args, fmt.Sprintf("dnssearch=%s", n.Dnssearch))
	}
	if n.Domainname != "" {
		args = append(args, fmt.Sprintf("domainname=%s", n.Domainname))
	}
	if n.Tftp != "" {
		args = append(args, fmt.Sprintf("tftp=%s", n.Tftp))
	}
	if n.Tftp_server_name != "" {
		args = append(args, fmt.Sprintf("tftp-server-name=%s", n.Tftp_server_name))
	}
	if n.Bootfile != "" {
		args = append(args, fmt.Sprintf("bootfile=%s", n.Bootfile))
	}
	if n.Hostfwd != "" {
		args = append(args, fmt.Sprintf("hostfwd=%s", n.Hostfwd))
	}
	if n.Guestfwd != "" {
		args = append(args, fmt.Sprintf("guestfwd=%s", n.Guestfwd))
	}
	if n.Smb != "" {
		args = append(args, fmt.Sprintf("smb=%s", n.Smb))
	}
	return []string{"-netdev", strings.Join(args, ",")}
}

/*
original command:

	-netdev tap,id=str[,fd=h][,fds=x:y:...:z][,ifname=name][,script=file][,downscript=dfile]
         [,br=bridge][,helper=helper][,sndbuf=nbytes][,vnet_hdr=on|off][,vhost=on|off]
         [,vhostfd=h][,vhostfds=x:y:...:z][,vhostforce=on|off][,queues=n]
         [,poll-us=n]
                configure a host TAP network backend with ID 'str'
                connected to a bridge (default=br0)
                use network scripts 'file' (default=/etc/qemu-ifup)
                to configure it and 'dfile' (default=/etc/qemu-ifdown)
                to deconfigure it
                use '[down]script=no' to disable script execution
                use network helper 'helper' (default=/usr/lib/qemu/qemu-bridge-helper) to
                configure it
                use 'fd=h' to connect to an already opened TAP interface
                use 'fds=x:y:...:z' to connect to already opened multiqueue capable TAP interfaces
                use 'sndbuf=nbytes' to limit the size of the send buffer (the
                default is disabled 'sndbuf=0' to enable flow control set 'sndbuf=1048576')
                use vnet_hdr=off to avoid enabling the IFF_VNET_HDR tap flag
                use vnet_hdr=on to make the lack of IFF_VNET_HDR support an error condition
                use vhost=on to enable experimental in kernel accelerator
                    (only has effect for virtio guests which use MSIX)
                use vhostforce=on to force vhost on for non-MSIX virtio guests
                use 'vhostfd=h' to connect to an already opened vhost net device
                use 'vhostfds=x:y:...:z to connect to multiple already opened vhost net devices
                use 'queues=n' to specify the number of queues to be created for multiqueue TAP
                use 'poll-us=n' to specify the maximum number of microseconds that could be
                spent on busy polling for vhost net
*/

type Netdev_TapOptions struct {
	ID         string `yaml:",omitempty"` // id=str
	Fd         string `yaml:",omitempty"` // [,fd=h]
	Fds        string `yaml:",omitempty"` // [,fds=x:y:...:z]
	Ifname     string `yaml:",omitempty"` // [,ifname=name]
	Script     string `yaml:",omitempty"` // [,script=file]
	Downscript string `yaml:",omitempty"` // [,downscript=dfile]
	Br         string `yaml:",omitempty"` // [,br=bridge]
	Helper     string `yaml:",omitempty"` // [,helper=helper]
	Sndbuf     string `yaml:",omitempty"` // [,sndbuf=nbytes]
	Vnet_hdr   string `yaml:",omitempty"` // [,vnet_hdr=on|off]
	Vhost      string `yaml:",omitempty"` // [,vhost=on|off]
	Vhostfd    string `yaml:",omitempty"` // [,vhostfd=h]
	Vhostfds   string `yaml:",omitempty"` // [,vhostfds=x:y:...:z]
	Vhostforce string `yaml:",omitempty"` // [,vhostforce=on|off]
	Queues     string `yaml:",omitempty"` // [,queues=n]
	Poll       string `yaml:",omitempty"` // [,poll-us=n]
}

func (n Netdev_TapOptions) ToArgs() []string {
	args := []string{"tap"}
	if n.ID != "" {
		args = append(args, fmt.Sprintf("id=%s", n.ID))
	}
	if n.Fd != "" {
		args = append(args, fmt.Sprintf("fd=%s", n.Fd))
	}
	if n.Fds != "" {
		args = append(args, fmt.Sprintf("fds=%s", n.Fds))
	}
	if n.Ifname != "" {
		args = append(args, fmt.Sprintf("ifname=%s", n.Ifname))
	}
	if n.Script != "" {
		args = append(args, fmt.Sprintf("script=%s", n.Script))
	}
	if n.Downscript != "" {
		args = append(args, fmt.Sprintf("downscript=%s", n.Downscript))
	}
	if n.Br != "" {
		args = append(args, fmt.Sprintf("br=%s", n.Br))
	}
	if n.Helper != "" {
		args = append(args, fmt.Sprintf("helper=%s", n.Helper))
	}
	if n.Sndbuf != "" {
		args = append(args, fmt.Sprintf("sndbuf=%s", n.Sndbuf))
	}
	if n.Vnet_hdr != "" {
		args = append(args, fmt.Sprintf("vnet_hdr=%s", n.Vnet_hdr))
	}
	if n.Vhost != "" {
		args = append(args, fmt.Sprintf("vhost=%s", n.Vhost))
	}
	if n.Vhostfd != "" {
		args = append(args, fmt.Sprintf("vhostfd=%s", n.Vhostfd))
	}
	if n.Vhostfds != "" {
		args = append(args, fmt.Sprintf("vhostfds=%s", n.Vhostfds))
	}
	if n.Vhostforce != "" {
		args = append(args, fmt.Sprintf("vhostforce=%s", n.Vhostforce))
	}
	if n.Queues != "" {
		args = append(args, fmt.Sprintf("queues=%s", n.Queues))
	}
	if n.Poll != "" {
		args = append(args, fmt.Sprintf("poll=%s", n.Poll))
	}
	return []string{"-netdev", strings.Join(args, ",")}
}

/*
original command:

	-netdev passt,id=str[,path=file][,quiet=on|off][,vhost-user=on|off]
	[,mtu=mtu][,address=addr][,netmask=mask][,mac=addr][,gateway=addr]

		[,interface=name][,outbound=address][,outbound-if4=name]
		[,outbound-if6=name][,dns=addr][,search=list][,fqdn=name]
		[,dhcp-dns=on|off][,dhcp-search=on|off][,map-host-loopback=addr]
		[,map-guest-addr=addr][,dns-forward=addr][,dns-host=addr]
		[,tcp=on|off][,udp=on|off][,icmp=on|off][,dhcp=on|off]
		[,ndp=on|off][,dhcpv6=on|off][,ra=on|off][,freebind=on|off]
		[,ipv4=on|off][,ipv6=on|off][,tcp-ports=spec][,udp-ports=spec]
		[,param=list]
			configure a passt network backend with ID 'str'
			if 'path' is not provided 'passt' will be started according to PATH
			by default, informational message of passt are not displayed (quiet=on)
			to display this message, use 'quiet=off'
			by default, passt will be started in socket-based mode, to enable vhost-mode,
			use 'vhost-user=on'
			for details on other options, refer to passt(1)
			'param' allows to pass any option defined by passt(1)
*/
type Netdev_PasstOptions struct {
	ID                string `yaml:",omitempty"` // id=str
	Path              string `yaml:",omitempty"` // [,path=file]
	Quiet             string `yaml:",omitempty"` // [,quiet=on|off]
	Vhost             string `yaml:",omitempty"` // [,vhost-user=on|off]
	Mtu               string `yaml:",omitempty"` // [,mtu=mtu]
	Address           string `yaml:",omitempty"` // [,address=addr]
	Netmask           string `yaml:",omitempty"` // [,netmask=mask]
	Mac               string `yaml:",omitempty"` // [,mac=addr]
	Gateway           string `yaml:",omitempty"` // [,gateway=addr]
	Interface         string `yaml:",omitempty"` // [,interface=name]
	Outbound          string `yaml:",omitempty"` // [,outbound=address]
	Outbound_if4      string `yaml:",omitempty"` // [,outbound-if4=name]
	Outbound_if6      string `yaml:",omitempty"` // [,outbound-if6=name]
	Dns               string `yaml:",omitempty"` // [,dns=addr]
	Search            string `yaml:",omitempty"` // [,search=list]
	Fqdn              string `yaml:",omitempty"` // [,fqdn=name]
	Dhcp_dns          string `yaml:",omitempty"` // [,dhcp-dns=on|off]
	Dhcp_search       string `yaml:",omitempty"` // [,dhcp-search=on|off]
	Map_host_loopback string `yaml:",omitempty"` // [,map-host-loopback=addr]
	Map_guest_addr    string `yaml:",omitempty"` // [,map-guest-addr=addr]
	Dns_forward       string `yaml:",omitempty"` // [,dns-forward=addr]
	Dns_host          string `yaml:",omitempty"` // [,dns-host=addr]
	Tcp               string `yaml:",omitempty"` // [,tcp=on|off]
	Udp               string `yaml:",omitempty"` // [,udp=on|off]
	Icmp              string `yaml:",omitempty"` // [,icmp=on|off]
	Dhcp              string `yaml:",omitempty"` // [,dhcp=on|off]
	Ndp               string `yaml:",omitempty"` // [,ndp=on|off]
	Dhcpv6            string `yaml:",omitempty"` // [,dhcpv6=on|off]
	RA                string `yaml:",omitempty"` // [,ra=on|off]
	Freebind          string `yaml:",omitempty"` // [,freebind=on|off]
	Ipv4              string `yaml:",omitempty"` // [,ipv4=on|off]
	Ipv6              string `yaml:",omitempty"` // [,ipv6=on|off]
	Tcp_ports         string `yaml:",omitempty"` // [,tcp-ports=spec]
	Udp_ports         string `yaml:",omitempty"` // [,udp-ports=spec]
	Param             string `yaml:",omitempty"` // [,param=list]
}

func (n *Netdev_PasstOptions) ToArgs() []string {
	args := []string{"passt"}
	if n.ID != "" {
		args = append(args, fmt.Sprintf("id=%s", n.ID))
	}
	if n.Path != "" {
		args = append(args, fmt.Sprintf("path=%s", n.Path))
	}
	if n.Quiet != "" {
		args = append(args, fmt.Sprintf("quiet=%s", n.Quiet))
	}
	if n.Vhost != "" {
		args = append(args, fmt.Sprintf("vhost=%s", n.Vhost))
	}
	if n.Mtu != "" {
		args = append(args, fmt.Sprintf("mtu=%s", n.Mtu))
	}
	if n.Address != "" {
		args = append(args, fmt.Sprintf("address=%s", n.Address))
	}
	if n.Netmask != "" {
		args = append(args, fmt.Sprintf("nettmask=%s", n.Netmask))
	}
	if n.Mac != "" {
		args = append(args, fmt.Sprintf("mac=%s", n.Mac))
	}
	if n.Gateway != "" {
		args = append(args, fmt.Sprintf("gateway=%s", n.Gateway))
	}
	if n.Interface != "" {
		args = append(args, fmt.Sprintf("interface=%s", n.Interface))
	}
	if n.Outbound != "" {
		args = append(args, fmt.Sprintf("outbound=%s", n.Outbound))
	}
	if n.Outbound_if4 != "" {
		args = append(args, fmt.Sprintf("outbound-if4=%s", n.Outbound_if4))
	}
	if n.Outbound_if6 != "" {
		args = append(args, fmt.Sprintf("outbound-if6=%s", n.Outbound_if6))
	}
	if n.Dns != "" {
		args = append(args, fmt.Sprintf("dns=%s", n.Dns))
	}
	if n.Search != "" {
		args = append(args, fmt.Sprintf("search=%s", n.Search))
	}
	if n.Fqdn != "" {
		args = append(args, fmt.Sprintf("fqdn=%s", n.Fqdn))
	}
	if n.Dhcp_dns != "" {
		args = append(args, fmt.Sprintf("dhcp-dns=%s", n.Dhcp_dns))
	}
	if n.Dhcp_search != "" {
		args = append(args, fmt.Sprintf("dhcp-search=%s", n.Dhcp_search))
	}
	if n.Map_host_loopback != "" {
		args = append(args, fmt.Sprintf("map-host-loopback=%s", n.Map_host_loopback))
	}
	if n.Map_guest_addr != "" {
		args = append(args, fmt.Sprintf("map-guest-addr=%s", n.Map_guest_addr))
	}
	if n.Dns_forward != "" {
		args = append(args, fmt.Sprintf("dns-forward=%s", n.Dns_forward))
	}
	if n.Dns_host != "" {
		args = append(args, fmt.Sprintf("dns-host=%s", n.Dns_host))
	}
	if n.Tcp != "" {
		args = append(args, fmt.Sprintf("tcp=%s", n.Tcp))
	}
	if n.Udp != "" {
		args = append(args, fmt.Sprintf("udp=%s", n.Udp))
	}
	if n.Icmp != "" {
		args = append(args, fmt.Sprintf("icmp=%s", n.Icmp))
	}
	if n.Dhcp != "" {
		args = append(args, fmt.Sprintf("dhcp=%s", n.Dhcp))
	}
	if n.Ndp != "" {
		args = append(args, fmt.Sprintf("ndp=%s", n.Ndp))
	}
	if n.Dhcpv6 != "" {
		args = append(args, fmt.Sprintf("dhcpv6=%s", n.Dhcpv6))
	}
	if n.RA != "" {
		args = append(args, fmt.Sprintf("ra=%s", n.RA))
	}
	if n.Freebind != "" {
		args = append(args, fmt.Sprintf("freebind=%s", n.Freebind))
	}
	if n.Ipv4 != "" {
		args = append(args, fmt.Sprintf("ipv4=%s", n.Ipv4))
	}
	if n.Ipv6 != "" {
		args = append(args, fmt.Sprintf("ipv6=%s", n.Ipv6))
	}
	if n.Tcp_ports != "" {
		args = append(args, fmt.Sprintf("tcp-ports=%s", n.Tcp_ports))
	}
	if n.Udp_ports != "" {
		args = append(args, fmt.Sprintf("udp-ports=%s", n.Udp_ports))
	}
	if n.Param != "" {
		args = append(args, fmt.Sprintf("param=%s", n.Param))
	}
	return []string{"-netdev", strings.Join(args, "")}
}

/*
original command

-netdev bridge,id=str[,br=bridge][,helper=helper]

	configure a host TAP network backend with ID 'str' that is
	connected to a bridge (default=br0)
	using the program 'helper (default=/usr/lib/qemu/qemu-bridge-helper)
*/
type NetDev_BridgeOptions struct {
	ID     string `yaml:",omitempty"` // id=str
	Br     string `yaml:",omitempty"` // [,br=bridge]
	Helper string `yaml:",omitempty"` // [,helper=helper]
}

func (n *NetDev_BridgeOptions) ToArgs() []string {
	args := []string{"bridge"}
	if n.ID != "" {
		args = append(args, fmt.Sprintf("id=%s", n.ID))
	}
	if n.Br != "" {
		args = append(args, fmt.Sprintf("br=%s", n.Br))
	}
	if n.Helper != "" {
		args = append(args, fmt.Sprintf("helper=%s", n.Helper))
	}
	return []string{"-netdev", strings.Join(args, ",")}
}

/*
original command:

	-netdev vhost-user,id=str,chardev=dev[,vhostforce=on|off]
			configure a vhost-user network, backed by a chardev 'dev'
*/
type Netdev_Vhost_userOptions struct {
	ID         string `yaml:",omitempty"` // id=str
	Chardev    string `yaml:",omitempty"` // chardev=dev
	Vhostforce string `yaml:",omitempty"` // [,vhostforce=on|off]
}

func (n *Netdev_Vhost_userOptions) ToArgs() []string {
	args := []string{"vhost-user"}

	if n.ID != "" {
		args = append(args, fmt.Sprintf("id=%s", n.ID))
	}
	if n.Chardev != "" {
		args = append(args, fmt.Sprintf("chardev=%s", n.Chardev))
	}
	if n.Vhostforce != "" {
		args = append(args, fmt.Sprintf("vhostforce=%s", n.Vhostforce))
	}

	return []string{"-netdev", strings.Join(args, ",")}
}

/*
original command:

			-netdev vhost-vdpa,id=str[,vhostdev=/path/to/dev][,vhostfd=h]
	                configure a vhost-vdpa network,Establish a vhost-vdpa netdev
	                use 'vhostdev=/path/to/dev' to open a vhost vdpa device
	                use 'vhostfd=h' to connect to an already opened vhost vdpa device
*/
type Netdev_Vhost_vdpaOptions struct {
	ID       string `yaml:",omitempty"` // id=str
	Vhostdev string `yaml:",omitempty"` // [,vhostdev=/path/to/dev]
	Vhostfd  string `yaml:",omitempty"` // [,vhostfd=h]
}

func (n *Netdev_Vhost_vdpaOptions) ToArgs() []string {
	args := []string{"vhost-vdpa"}
	if n.Vhostdev != "" {
		args = append(args, fmt.Sprintf("vhostdev=%s", n.Vhostdev))
	}
	if n.Vhostfd != "" {
		args = append(args, fmt.Sprintf("vhostfd=%s", n.Vhostfd))
	}
	return []string{"-netdev", strings.Join(args, ",")}
}

/*
original command

	-netdev hubport,id=str,hubid=n[,netdev=nd]
		configure a hub port on the hub with ID 'n'
*/
type Netdev_HubportOptions struct {
	ID     string `yaml:",omitempty"` // id=str
	Hubid  string `yaml:",omitempty"` // hubid=n
	Netdev string `yaml:",omitempty"` // [,netdev=nd]
}

func (n *Netdev_HubportOptions) ToArgs() []string {
	args := []string{"hubport"}
	if n.ID != "" {
		args = append(args, fmt.Sprintf("id=%s", n.ID))
	}
	if n.Hubid != "" {
		args = append(args, fmt.Sprintf("hubid=%s", n.Hubid))
	}
	if n.Netdev != "" {
		args = append(args, fmt.Sprintf("netdev=%s", n.Netdev))
	}
	return []string{"-netdev", strings.Join(args, ",")}
}

/*
original command:

-netdev vde,id=str[,sock=socketpath][,port=n][,group=groupname][,mode=octalmode]

	configure a network backend to connect to port 'n' of a vde switch
	running on host and listening for incoming connections on 'socketpath'.
	Use group 'groupname' and mode 'octalmode' to change default
	ownership and permissions for communication port.
*/
type Netdev_VdeOptions struct {
	ID    string `yaml:",omitempty"` //  id=str
	Sock  string `yaml:",omitempty"` // [,sock=socketpath]
	Port  string `yaml:",omitempty"` // [,port=n]
	Group string `yaml:",omitempty"` // [,group=groupname]
	Mode  string `yaml:",omitempty"` // [,mode=octalmode]
}

func (n *Netdev_VdeOptions) ToArgs() []string {
	args := []string{"vde"}
	if n.ID != "" {
		args = append(args, fmt.Sprintf("id=%s", n.ID))
	}
	if n.Sock != "" {
		args = append(args, fmt.Sprintf("sock=%s", n.Sock))
	}
	if n.Port != "" {
		args = append(args, fmt.Sprintf("port=%s", n.Port))
	}
	if n.Group != "" {
		args = append(args, fmt.Sprintf("group=%s", n.Group))
	}
	if n.Mode != "" {
		args = append(args, fmt.Sprintf("mode=%s", n.Mode))
	}
	return []string{"-netdev", strings.Join(args, ",")}
}
