package config

// Config : AUTO Created From /home/qmiao/Desktop/n3/n3-privacy/Server/config/config.toml
type Config struct {
	Path    string
	Service string
	Version string
	Log     string
	Storage struct {
		DB      string
		DBPath  string
		Tracing bool
	}
	Loggly struct {
		Token string
	}
	WebService struct {
		Port int
	}
	Route struct {
		GetHash   string
		Delete    string
		HELP      string
		LsID      string
		Enforce   string
		LsUser    string
		Update    string
		GetID     string
		Get       string
		LsContext string
		LsObject  string
	}
	File struct {
		EnforcerMac     string
		EnforcerWin64   string
		ClientConfig    string
		ClientLinux64   string
		ClientMac       string
		ClientWin64     string
		EnforcerLinux64 string
	}
	Server struct {
		IP       interface{}
		Port     interface{}
		Protocol string
	}
	Access struct {
		Timeout int
	}
}
