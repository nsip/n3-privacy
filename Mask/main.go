package main

import (
	"flag"
	"os"
	"path/filepath"

	p "github.com/nsip/n3-privacy/Mask/process"
)

func main() {
	exe := filepath.Base(os.Args[0])
	if len(os.Args) < 3 {
		fPf("Usage: %s [-o='output'] <inputdata.json> <mask.json>\n", exe)
		return
	}

	inFilePath, maskFilePath := os.Args[1], os.Args[2]
	if len(os.Args) == 4 {
		inFilePath, maskFilePath = os.Args[2], os.Args[3]
	}

	outputPtr := flag.String("o", "out.json", "a string")
	flag.Parse()
	output := *outputPtr
	if !sHasSuffix(output, ".json") {
		output = output + ".json"
	}

	p.FileMask(inFilePath, maskFilePath, output)
}
