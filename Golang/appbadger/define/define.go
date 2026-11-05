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
