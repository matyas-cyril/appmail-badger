package login

type LOGIN struct {
	Cmd      string
	Login    string
	Password string
	Nbf      int64
	Exp      int64
	Enable   bool
	File     string // Fichier de configuration
}
