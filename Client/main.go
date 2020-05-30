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
	fns := structFields(clt.Config{}.Route)
	failOnErrWhen(len(os.Args) < 2, "%v: need %v", eg.PARAM_INVALID, fns)

	fn, args := os.Args[1], ""
	if !xin(fn, []string{"HELP", "LsID", "LsContext", "LsUser", "LsObject"}) {
		failOnErrWhen(len(os.Args) < 3, "%v: need %v [-id= -u= -c= -o= -rw= -p= -d= -w=]", eg.PARAM_INVALID, fns)
		args = os.Args[2]
	}

	cmd := flag.NewFlagSet(args, flag.ExitOnError)
	id := cmd.String("id", "", "policy ID")
	user := cmd.String("u", "", "user")
	ctx := cmd.String("c", "", "context")
	object := cmd.String("o", "", "object")
	rw := cmd.String("rw", "", "read/write")
	policyPtr := cmd.String("p", "", "the path of policy to be uploaded")
	dataPtr := cmd.String("d", "", "the path of json to be uploaded")
	wholeDump := cmd.Bool("w", false, "output all attributes content from response")
	cmd.Parse(os.Args[2:])

	policy, err := ioutil.ReadFile(*policyPtr)
	failOnErrWhen(fn == "Update", "%v: %s", err, *policyPtr)
	data, err := ioutil.ReadFile(*dataPtr)
	failOnErrWhen(fn == "Enforce", "%v: %s", err, *dataPtr)

	str, err := clt.DO(
		"cfg-clt-privacy.toml",
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

	if xin(fn, []string{"HELP", "LsID", "LsContext", "LsUser", "LsObject"}) {
		fPln(str)
		return
	}

	m := make(map[string]interface{})
	failOnErr("json.Unmarshal ... %v", json.Unmarshal([]byte(str), &m))
	if *wholeDump {
		fPf("Empty? %v\n%s\n", m["empty"], "-----------------------------")
		fPf("ERROR: %v\n%s\n", m["error"], "-----------------------------")
	}
	if m["data"] != nil && m["data"] != "" {
		fPf("%s\n", m["data"])
	}
}
