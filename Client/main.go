package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	eg "github.com/cdutwhu/n3-util/n3errs"
	clt "github.com/nsip/n3-privacy/Server/go-client"
)

func main() {
	route := clt.Config{}.Route
	fns, err := structFields(&route)
	failOnErr("%v", err)
	failOnErrWhen(len(os.Args) < 3, "%v: need [config.toml] %v", eg.CLI_SUBCMD_ERR, fns)

	cltcfg, fn := os.Args[1], os.Args[2]

	ok, err := xin(fn, []string{"HELP", "LsID", "LsContext", "LsUser", "LsObject"})
	failOnErr("%v", err)
	failOnErrWhen(!ok && len(os.Args) < 4, "%v: need %v [-id= -u= -c= -o= -rw= -p= -d= -w=]", eg.PARAM_INVALID, fns)

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
		clt.Args{
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

	ok, err = xin(fn, []string{"HELP", "LsID", "LsContext", "LsUser", "LsObject"})
	if failOnErr("%v", err); ok {
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
