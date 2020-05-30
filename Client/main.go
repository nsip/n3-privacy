package main

import (
	"flag"
	"io/ioutil"
	"os"

	clt "github.com/nsip/n3-privacy/Server/go-client"
)

func main() {

	fn := os.Args[1]
	args := os.Args[2]

	cmd := flag.NewFlagSet(args, flag.ExitOnError)
	id := cmd.String("id", "", "policy ID")
	user := cmd.String("u", "", "user")
	ctx := cmd.String("c", "", "context")
	object := cmd.String("o", "", "object")
	rw := cmd.String("rw", "", "read/write")
	policyPtr := cmd.String("p", "", "the path of policy to be uploaded")
	// wholeDump := cmd.Bool("w", false, "output all attributes content from response")
	dataPtr := cmd.String("d", "", "the path of json to be uploaded")
	cmd.Parse(os.Args[3:])

	policy, err := ioutil.ReadFile(*policyPtr)
	failOnErr("%v", err)

	data, err := ioutil.ReadFile(*dataPtr)
	failOnErr("%v", err)

	// mngMode := false
	// // case "LsID", "LsContext", "LsUser", "LsObject": mngMode = true

	clt.DO(
		"cfg-clt-privacy.toml",
		fn,
		clt.Args{
			ID:     *id, // 1615307cc4bf38ffcad912dea96fec4024700fd9r
			Policy: policy,
			User:   *user,
			Ctx:    *ctx,
			RW:     *rw,
			Object: *object,
			Data:   data,
		},
	)

	// data, err = ioutil.ReadAll(resp.Body)
	// failOnErr("%v", err)

	// const SepLn = "-----------------------------"

	// if *fullDump {
	// 	fPf("accessing... %s\n%s\n", url, SepLn)
	// }

	// if data != nil {
	// 	if os.Args[1] == "HELP" {
	// 		fPt(string(data))
	// 	} else {
	// 		m := make(map[string]interface{})
	// 		failOnErr("json.Unmarshal ... %v", json.Unmarshal(data, &m))
	// 		if !mngMode {
	// 			if *fullDump {
	// 				if m["empty"] != nil && m["empty"] != "" {
	// 					fPf("Empty? %v\n%s\n", m["empty"], SepLn)
	// 				}
	// 				if m["error"] != nil && m["error"] != "" {
	// 					fPf("ERROR: %v\n%s\n", m["error"], SepLn)
	// 				}
	// 			}
	// 			if m["data"] != nil && m["data"] != "" {
	// 				fPf("%s\n", m["data"])
	// 			}
	// 		} else {
	// 			key := ""
	// 			switch {
	// 			case *user != "" && *ctx != "":
	// 				key = fSf("%s@%s", *user, *ctx)
	// 			case *user != "":
	// 				key = *user
	// 			case *ctx != "":
	// 				key = *ctx
	// 			default:
	// 				key = "all"
	// 			}
	// 			fPf("%s\n", m[key])
	// 		}
	// 	}
	// }
}
