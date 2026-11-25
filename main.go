package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

type DomainProfile struct {
	A     []net.IP
	AAAA  []net.IP
	MX    []*net.MX
	TXT   []string
	NS    []*net.NS
	CNAME string
}

type PrintOptions struct {
	aFlag     bool
	aaaaFlag  bool
	mxFlag    bool
	txtFlag   bool
	nsFlag    bool
	cnameFlag bool
}

func main() {
	aFlag := flag.Bool("a", false, "lookup A records")
	aaaaFlag := flag.Bool("aaaa", false, "lookup AAAA records")
	mxFlag := flag.Bool("mx", false, "lookup MX records")
	txtFlag := flag.Bool("txt", false, "lookup A records")
	nsFlag := flag.Bool("ns", false, "lookup A records")
	cnameFlag := flag.Bool("cname", false, "lookup A records")

	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("No Domain Found")
	}
	domain := flag.Arg(0)
	dp, _ := BuildProfile(domain)

	opts := PrintOptions{
		aFlag:     *aFlag,
		aaaaFlag:  *aaaaFlag,
		mxFlag:    *mxFlag,
		txtFlag:   *txtFlag,
		nsFlag:    *nsFlag,
		cnameFlag: *cnameFlag,
	}

	PrintProfile(dp, domain, opts)
}

func BuildProfile(domain string) (profile DomainProfile, err error) {
	d := strings.ToLower(domain)
	profile.A, profile.AAAA, _ = DomainARecords(d)
	profile.MX, _ = net.LookupMX(d)
	profile.TXT, _ = net.LookupTXT(d)
	profile.NS, _ = net.LookupNS(d)
	profile.CNAME, _ = net.LookupCNAME(d)
	return profile, nil
}

func DomainARecords(domain string) (a []net.IP, aaaa []net.IP, err error) {
	p := DomainProfile{}
	d := strings.ToLower(domain)

	arecords, err := net.LookupIP(d)
	if err != nil {
		return nil, nil, err
	}

	for _, r := range arecords {
		if r.To4() != nil {
			p.A = append(p.A, r)
		} else {
			p.AAAA = append(p.AAAA, r)
		}
	}
	return p.A, p.AAAA, nil
}

func PrintProfile(profile DomainProfile, domain string, opts PrintOptions) {
	all := !opts.aFlag && !opts.aaaaFlag && !opts.mxFlag && !opts.txtFlag && !opts.nsFlag && !opts.cnameFlag
	if all || opts.aFlag {
		aStrings := ToString(profile.A, func(ip net.IP) string { return ip.String() })
		PrintLine("A Records:", aStrings)
	}
	if all || opts.aaaaFlag {
		aaaaStrings := ToString(profile.AAAA, func(ip net.IP) string { return ip.String() })
		PrintLine("AAAA Records:", aaaaStrings)
	}
	if all || opts.mxFlag {
		mxStrings := ToString(profile.MX, func(mx *net.MX) string { return fmt.Sprintf("Host: %s | Priority: %d", mx.Host, mx.Pref) })
		PrintLine("MX Records:", mxStrings)
	}
	if all || opts.txtFlag {
		fmt.Println("TXT Records:")
		for _, t := range profile.TXT {
			fmt.Println(t)
		}
	}
	if all || opts.nsFlag {
		nsStrings := ToString(profile.NS, func(ns *net.NS) string { return fmt.Sprintf("%s", ns.Host) })
		PrintLine("NS Records:", nsStrings)
	}
	if all || opts.cnameFlag {
		cname := strings.TrimSuffix(profile.CNAME, ".")
		if cname == domain {
			fmt.Println("CNAME Record: None Found")
		} else {
			fmt.Println("CNAME Record:", profile.CNAME)
		}
	}
}

func PrintLine(label string, items []string) {
	if len(items) == 0 {
		fmt.Println("None Found")
		return
	}
	fmt.Println(label, strings.Join(items, ", "))
}

func ToString[T any](items []T, format func(T) string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, format(item))
	}
	return out
}
