package main

import "strings"

func ExtractSubdomain(name string) string {
	if !strings.HasSuffix(name, ".flatbo.at.") {
		return ""
	}
	name = strings.TrimSuffix(name, ".flatbo.at.")
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}
