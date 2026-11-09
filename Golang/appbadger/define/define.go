package define

import (
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

// Convertir un type string vers un type sqlc -> pgtype.UUID
func StringToPgUUID(data string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(data)
	return uuid, err
}

func TrimSpace(str ...*string) {
	for _, s := range str {
		if s != nil {
			*s = strings.TrimSpace(*s)
		}
	}
}

// Evite le doublons , mais sensible à la casse
func TrimSpaceArray(array []string) []string {

	var tab []string

	seen := make(map[string]bool, len(array))

	// Pré-allouer le tableau
	if len(array) > 0 {
		tab = make([]string, 0, len(array))
	}

	for _, value := range array {
		value = strings.TrimSpace(value)
		if value != "" {

			if seen[value] { // Existe déjà
				continue
			}

			seen[value] = true
			tab = append(tab, value)
		}
	}
	return tab
}

/*
Comptabiliser le nbr d'arg != de "" lors de l'utilisation des param d'app
*/
func CountDefinedStrArg(args ...string) uint64 {
	var cpt uint64 = 0
	for _, str := range args {
		if strings.TrimSpace(str) != "" {
			cpt++
		}
	}
	return cpt
}
