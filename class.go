package main

import (
	"github.com/miekg/dns"
)

var g_classMap map[string]uint16

func init() {
	m := map[string]uint16{
		"IN": dns.ClassINET,
	}
	g_classMap = m
}
