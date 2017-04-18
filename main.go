package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	goflags "github.com/jessevdk/go-flags"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

type options struct {
	ConfigFile string `short:"c" long:"config" description:"Configuration file"`
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %+v\n", os.Args[0], err)
		os.Exit(1)
	}
}

func run() error {
	opts, err := processArgs()
	if err != nil {
		return err
	}

	cfg, err := loadConfig(opts.ConfigFile)
	if err != nil {
		return err
	}

	dns.HandleFunc(".", makeHandler(cfg))
	if cfg.Ports.UDP4 != 0 {
		startUDPServer(cfg.Ports.UDP4)
	}
	if cfg.Ports.TCP4 != 0 {
		startTCPServer(cfg.Ports.TCP4)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Fatalf("Signal (%v) received, stopping\n", s)

	return nil
}

func processArgs() (*options, error) {
	var opts options
	_, err := goflags.ParseArgs(&opts, os.Args)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse arguments")
	}

	if opts.ConfigFile == "" {
		return nil, errors.New("mandatory argument `config` is not specified")
	}

	return &opts, nil
}

func makeHandler(cfg *Config) func(dns.ResponseWriter, *dns.Msg) {
	return func(w dns.ResponseWriter, r *dns.Msg) {
		handleDNSRequest(cfg, w, r)
	}
}

func startUDPServer(port uint16) {
	go func() {
		srv := &dns.Server{Addr: ":" + strconv.Itoa(int(port)), Net: "udp"}
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set udp listener %s\n", err.Error())
		}
	}()
}

func startTCPServer(port uint16) {
	go func() {
		srv := &dns.Server{Addr: ":" + strconv.Itoa(int(port)), Net: "tcp"}
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set tcp listener %s\n", err.Error())
		}
	}()
}

func handleDNSRequest(cfg *Config, w dns.ResponseWriter, r *dns.Msg) {
	for _, q := range r.Question {
		a := getAnswers(cfg, q)
		m := new(dns.Msg)
		m.SetReply(r)
		m.Authoritative = true
		m.Answer = a
		w.WriteMsg(m)
	}
}

func getAnswers(cfg *Config, q dns.Question) []dns.RR {
	qname := strings.TrimSuffix(q.Name, ".")
	r, found := cfg.Records[qname]
	if !found {
		fmt.Fprintf(os.Stderr, "no entry found for %s\n", qname)
		return nil
	}

	result := make([]dns.RR, 0, len(r))
	for _, e := range r {
		if q.Qtype != e.Type {
			fmt.Fprintf(os.Stderr, "q type and e type are different: q=%d, e=%d\n", q.Qtype, e.Type)
			continue
		}

		switch e.Type {
		case dns.TypeA:
			result = append(result, &dns.A{
				Hdr: dns.RR_Header{
					Name:     q.Name,
					Rrtype:   e.Type,
					Class:    e.Class,
					Ttl:      e.TTL,
					Rdlength: 4,
				},
				A: e.A,
			})
		default:
			fmt.Fprintf(os.Stderr, "type %d is not supported yet.\n", e.Type)
		}
	}
	return result
}
