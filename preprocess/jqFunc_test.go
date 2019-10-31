package preprocess

import (
	"fmt"
	"io/ioutil"
	"testing"
)

// func MakeRTJQ() {
// 	switch runtime.GOOS {
// 	case "windows":
// 	case "linux":
// 	case "darwin":
// 	}
// }

func TestJQ(t *testing.T) {

	// fmt.Println(prepareJQ("../", "./", "./utils/"))
	// fmt.Println(os.Getwd())

	// fmt.Println(FmtJSONStr("{\"abc\": 123}", "../", "./", "./utils/"))

	// if data, err := ioutil.ReadFile("../data/sample.json"); err == nil {
	// 	// fmt.Println(string(data))
	// 	fmt.Println(FmtJSONStr(string(data), "../", "./", "./utils/"))
	// } else {
	// 	fmt.Println(err.Error())
	// }

	// formatted := FmtJSONFile("../../data/xapi.json", "../", "./", "../build/Linux64/")
	// ioutil.WriteFile("fmt.json", []byte(formatted), 0666)

	formatted := FmtJSONFile("../../Server/config/meta.json", "./utils/")
	ioutil.WriteFile("../Server/config/meta.json", []byte(formatted), 0666)

	fmt.Println()
	// FmtJSONFile("../data/xapi1.json", "../", "./", "./utils/")
}
