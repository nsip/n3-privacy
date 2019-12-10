package preprocess

import (
	"os"
	"os/exec"
	"strings"

	cmn "github.com/nsip/n3-privacy/common"
)

func prepareJQ(jqDirs ...string) (jqWD, oriWD string, err error) {
	fn := "prepareJQ"
	oriWD, err = os.Getwd()
	cmn.FailOnErr("Getwd() 1 fatal @ %v: %w", fn, err)
	jqDirs = append(jqDirs, "./", "../", "../../")
	for _, jqWD = range jqDirs {
		if !strings.HasSuffix(jqWD, "/") {
			jqWD += "/"
		}
		if _, err = os.Stat(jqWD + jq); err == nil {
			cmn.FailOnErr("Chdir() fatal @ %v: %w", fn, os.Chdir(jqWD))
			jqWD, err = os.Getwd()
			cmn.FailOnErr("Getwd() 2 fatal @ %v: %w", fn, err)
			return jqWD, oriWD, nil
		}
	}
	cmn.FailOnErr("[%s] is not found @ %v", jq, fEf(fn))
	return "", "", nil
}

// FmtJSONStr : <json string> must not have single quote <'>
func FmtJSONStr(json string, jqDirs ...string) string {
	_, oriWD, _ := prepareJQ(jqDirs...)
	defer func() { os.Chdir(oriWD) }()

	json = "'" + strings.ReplaceAll(json, "'", "\\'") + "'" // *** deal with <single quote> in "echo" ***
	cmdstr := "echo " + json + ` | ./` + jq + " ."
	cmd := exec.Command(execCmdName, execCmdP0, cmdstr)

	if output, err := cmd.Output(); err == nil {
		return string(output)
	}
	cmn.FailOnErr("cmd.Output() error @ %v", fEf("FmtJSONStr"))
	return ""
}

// FmtJSONFile : <file> is the <relative path> to <jq>
func FmtJSONFile(file string, jqDirs ...string) string {
	if _, oriWD, err := prepareJQ(jqDirs...); err == nil {
		defer func() { os.Chdir(oriWD) }()
		cmdstr := "cat " + file + ` | ./` + jq + " ."
		// cmdstr := "cat " + file
		cmd := exec.Command(execCmdName, execCmdP0, cmdstr)
		if output, err := cmd.Output(); err == nil {
			return string(output)
		}
		cmn.FailOnErr("cmd.Output() error @ %v", fEf("FmtJSONFile"))
		return ""
	}
	cmn.FailOnErr("prepareJQ error @ %v", fEf("FmtJSONFile"))
	return ""
}
