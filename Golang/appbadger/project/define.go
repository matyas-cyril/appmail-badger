package project

import (
	"net/netip"
	"strings"
)

type parseSrv struct {
	name []string
	ip   []netip.Addr
}

// Permet de convertir une saisie de string en struct
// pour répartir entre nom de serveur et ip (juste la syntaxe)
// Retourne nil si vide
func parseServers(args []string) *parseSrv {

	var srv parseSrv

	for _, arg := range args {

		value := strings.TrimSpace(arg)
		if value == "" {
			continue
		}

		addr, err := netip.ParseAddr(value)
		if err != nil {
			// On va dire que c'est une chaine
			srv.name = append(srv.name, value)
			continue
		}

		srv.ip = append(srv.ip, addr)

	}

	if len(srv.ip) == 0 && len(srv.name) == 0 {
		return nil
	}

	return &srv
}
