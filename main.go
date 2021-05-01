package main

import (
	"flag"
	"fmt"

	. "github.com/dave/jennifer/jen"
	"github.com/ui-kreinhard/go-params-to-interfaces/ast-parser"
)

func main() {
	file := flag.String("file", "", "The file to be analyzed")
	function := flag.String("function", "", "Name of the function to be 'builderized'")
	structName := flag.String("struct", "", "receiver of the struct. can be empty")

	flag.Parse()

	if *file == "" || *function == "" {
		flag.PrintDefaults()
		return
	}

	f := NewFile("main")
	ast_parser.ExtractModel(*file, *function, *structName).GetStruct(f)
	ast_parser.ExtractModel(*file, *function, *structName).GetInterfaces(f)
	ast_parser.ExtractModel(*file, *function, *structName).GetImplementations(f)
	ast_parser.ExtractModel(*file, *function, *structName).GetEntryMethod(f)
	ast_parser.ExtractModel(*file, *function, *structName).GetInterfaceContract(f)
	fmt.Println(f.GoString())
}
