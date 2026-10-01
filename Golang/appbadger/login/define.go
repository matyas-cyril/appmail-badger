package login

import (
	"regexp"
	"unicode/utf8"
)

type LOGIN struct {
	Cmd      string
	Login    string
	Password string
	Nbf      int64
	Exp      int64
	Enable   bool
	File     string // Fichier de configuration
}

const (

	// Définir la taille MIN et MAX du LOGIN
	LOGIN_MIN_SIZE uint8 = 1
	LOGIN_MAX_SIZE uint8 = 64

	// Définir la taille MIN et MAX du mot de passe
	PASSWD_MIN_SIZE uint8 = 8
	PASSWD_MAX_SIZE uint8 = 255
)

var (

	// Go ne supporte pas les lookaheads comme en perl, donc 2 Regexp
	loginRegex       = regexp.MustCompile(`^[A-Za-z0-9_](?:[A-Za-z0-9._-]*[A-Za-z0-9_])?$`)
	loginLetterRegex = regexp.MustCompile(`[A-Za-z]`)

	passwordRegex = regexp.MustCompile(`[^\x21-\x7E\x{A1}-\x{A8}\x{B5}\x{BF}-\x{FF}]+`)
)

// IsValidLogin vérifie la validité syntaxique d'un login
func IsValidLogin(login string) error {

	// Valider la longueur en nbr de carac
	length := utf8.RuneCountInString(login)
	if length < int(LOGIN_MIN_SIZE) || length > int(LOGIN_MAX_SIZE) {
		return ErrLoginInvalidLength
	}

	// Valider le format
	if loginRegex.MatchString(login) && loginLetterRegex.MatchString(login) {
		return nil
	}

	return ErrLoginInvalidFormat
}

// IsValidLogin vérifie la validité syntaxique d'un mot de passe
func IsValidPassword(passwd string) error {

	// Valider la longueur en nbr de carac
	length := utf8.RuneCountInString(passwd)
	if length < int(PASSWD_MIN_SIZE) || length > int(PASSWD_MAX_SIZE) {
		return ErrPasswordInvalidLength
	}

	// Valider le format
	if passwordRegex.MatchString(passwd) {
		return ErrPasswordInvalidFormat
	}

	return nil

}

func IsValidPasswordRune(r rune) bool {
	return !passwordRegex.MatchString(string(r))
}
