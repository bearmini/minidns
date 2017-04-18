package main

import "github.com/miekg/dns"

var g_rrTypeMap map[string]uint16

func init() {
	m := map[string]uint16{
		"A": dns.TypeA,
	}

	g_rrTypeMap = m
}
