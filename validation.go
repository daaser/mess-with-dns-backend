package main

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

var disallowedDomains = map[string]bool{
	"ns1":    true,
	"ns2":    true,
	"orange": true,
	"purple": true,
	"www":    true,
}

func validateDomainName(domain string, username string) error {
	if !strings.HasSuffix(domain, ".") {
		return fmt.Errorf("domain must end with a period")
	}
	if _, ok := dns.IsDomainName(domain); !ok {
		return fmt.Errorf("invalid domain name: %s", domain)
	}
	if !strings.HasSuffix(domain, ".flatbo.at.") {
		return fmt.Errorf("subdomain must end with .flatbo.at")
	}
	// get last component of domain
	name := strings.TrimSuffix(domain, ".flatbo.at.")
	subdomain := ExtractSubdomain(domain)
	if subdomain != username {
		return fmt.Errorf("subdomain must be '%s'", username)
	}
	if _, ok := disallowedDomains[subdomain]; ok {
		return fmt.Errorf("sorry, you're not allowed to make changes to '%s' :)", subdomain)
	}
	if strings.Contains(name, "flatbo.at") {
		return fmt.Errorf(
			"you tried to create a record for %s, you probably didn't want that",
			domain,
		)
	}
	return nil
}
