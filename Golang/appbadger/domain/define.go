package domain

import (
	"fmt"
	"regexp"
	"strings"
)

var domRegexp = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`)

// Vérifier la validité syntaxisque d'un domaine et retourne la mise en forme
func validDomain(domain string) (string, error) {

	dom := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(domain), "."))

	if len(dom) == 0 || len(dom) > 253 {
		return "", fmt.Errorf("domain '%s' length invalid", domain)
	}

	for _, label := range strings.Split(dom, ".") {
		if !domRegexp.MatchString(label) {
			return "", fmt.Errorf("domain '%s' syntax invalid", domain)
		}
	}
	return dom, nil
}
