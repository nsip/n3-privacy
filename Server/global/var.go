package global

import (
	"os"

	cfg "github.com/nsip/n3-privacy/Server/config"
)

var (
	// Cfg : global variable
	Cfg *cfg.Config

	// EncPwd :
	EncPwd = "password"

	// WD : original work directory
	WD, _ = os.Getwd()
)

// Init : initialize the global variables
func Init(configs ...string) bool {
	configs = append(configs, "./config.toml", "../config.toml", "../../config.toml")
	Cfg = cfg.NewCfg(configs...)
	return Cfg != nil
}

// WDCheck :
// func WDCheck() {
// 	done := make(chan string)
// 	go func() {
// 	AGAIN:
// 		if path, _ := os.Getwd(); path != WD {
// 			time.Sleep(10 * time.Millisecond)
// 			goto AGAIN
// 		}
// 		done <- "done"
// 		// fmt.Println("done")
// 	}()
// 	<-done
// }
