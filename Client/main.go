package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	"github.com/cdutwhu/n3-util/n3err"
	clt "github.com/nsip/n3-privacy/Server/goclient"
)

func main() {
	route := clt.Config{}.Route
	fns := structFields(&route)
	failOnErrWhen(len(os.Args) < 3, "%v: need [config.toml] %v", n3err.CLI_SUBCMD_ERR, fns)

	cltcfg, fn := os.Args[1], os.Args[2]

	ok := exist(fn, "HELP", "LsID", "LsContext", "LsUser", "LsObject")
	failOnErrWhen(!ok && len(os.Args) < 4, "%v: need %v [-id= -u= -c= -o= -rw= -p= -d= -w=]", n3err.PARAM_INVALID, fns)

	cmd := flag.NewFlagSet(fn, flag.ExitOnError)
	id := cmd.String("id", "", "policy ID")
	user := cmd.String("u", "", "user")
	ctx := cmd.String("c", "", "context")
	object := cmd.String("o", "", "object")
	rw := cmd.String("rw", "", "read/write")
	policyPtr := cmd.String("p", "", "the path of policy to be uploaded")
	dataPtr := cmd.String("d", "", "the path of json to be uploaded")
	wholeDump := cmd.Bool("w", false, "output all attributes content from response")
	cmd.Parse(os.Args[3:])

	policy, err := ioutil.ReadFile(*policyPtr)
	failOnErrWhen(fn == "Update", "%v: %s", err, *policyPtr)
	data, err := ioutil.ReadFile(*dataPtr)
	failOnErrWhen(fn == "Enforce", "%v: %s", err, *dataPtr)

	str, err := clt.DO(
		cltcfg,
		fn,
		&clt.Args{
			ID:     *id,
			Policy: policy,
			User:   *user,
			Ctx:    *ctx,
			RW:     *rw,
			Object: *object,
			Data:   data,
		},
	)
	failOnErr("%v", err)

	if exist(fn, "HELP", "LsID", "LsContext", "LsUser", "LsObject") {
		fPln(str)
		return
	}

	m := make(map[string]interface{})
	failOnErr("json.Unmarshal ... %v", json.Unmarshal([]byte(str), &m))
	if *wholeDump {
		fPf("Empty? %v\n%s\n", m["empty"], "-----------------------------")
		fPf("ERROR: %v\n%s\n", m["error"], "-----------------------------")
	}
	fPf("%s\n", m["data"])
}
