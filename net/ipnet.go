package net

import "net"

type IPNetEx struct {
	*net.IPNet
	net.IP
}

func ParseCIDR(cidr string) (*IPNetEx, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	return &IPNetEx{
		IPNet: ipnet,
		IP:    ip,
	}, nil
}

func (net *IPNetEx) Range() (net.IP, net.IP) {

}

func (net *IPNetEx) AvailableRange() (net.IP, net.IP) {

}
