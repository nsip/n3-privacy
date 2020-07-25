package goclient

// Config : AUTO Created From /home/qmiao/Desktop/n3/n3-privacy/Server/goclient/config.toml
type Config struct {
    Path string
    Service string
    Route struct {
        LsID string
        LsUser string
        GetHash string
        Delete string
        Get string
        LsContext string
        Update string
        GetID string
        Enforce string
        HELP string
        LsObject string
    }
    Server struct {
        Port int
        Protocol string
        IP string
    }
    Access struct {
        Timeout int
    }
}
