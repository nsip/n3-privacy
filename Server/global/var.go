package global

import (
	"os"
	"time"

	cfg "github.com/nsip/n3-privacy/Server/config"
)

var (
	// Cfg : global variable
	Cfg *cfg.Config
	// WD : original work directory
	WD, _ = os.Getwd()
)

// Init : initialize the global variables
func Init() bool {
	Cfg = cfg.NewCfg("./config.toml", "../config.toml")
	return Cfg != nil
}

// WDCheck :
func WDCheck() {
	done := make(chan string)
	go func() {
	AGAIN:
		if path, _ := os.Getwd(); path != WD {
			time.Sleep(10 * time.Millisecond)
			goto AGAIN
		}
		done <- "done"
		// fmt.Println("done")
	}()
	<-done
}
