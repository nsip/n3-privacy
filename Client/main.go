package main

import (
	"fmt"

	glb "github.com/nsip/n3-privacy/Client/global"
	// cfg "github.com/nsip/n3-privacy/Client/config"
)

func main() {

	glb.Init()
	fmt.Println(glb.Cfg.Path)

	// if resp, err := http.Get( glb.Cfg.Path,    "http://192.168.92.133:1323/policy-service/0.1.0/list/object"); err == nil {
	// 	data, _ := ioutil.ReadAll(resp.Body)
	// 	fmt.Println(string(data))
	// }
}
