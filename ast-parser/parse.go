package ast_parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
)

func getType(src []byte, param ast.Node, offset token.Pos) string {
	return string(src[param.Pos()-offset:param.End()-offset])	
}

func getReceiverType(src []byte, recv *ast.FieldList, offset token.Pos) string {
	if recv == nil || recv.List == nil || len(recv.List) <= 0 {
		return ""
	}
	return getType(src, recv.List[0].Type, offset)
}

func ExtractModel(filename string, methodName string, receiver string) *Method {
	ret := Method{
		"",
		[]Param{},
		"",
	}
	src, _ := ioutil.ReadFile(filename)
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	offset := node.Pos()

	ast.Inspect(node, func(n ast.Node) bool {
		// Find Return Statements
		methodFunction, ok := n.(*ast.FuncDecl)
		
		if ok {
			recvType := getReceiverType(src, methodFunction.Recv, offset)
			if methodFunction.Name.Name == methodName && (recvType == receiver) {
				ret.Name = methodName
				ret.Receiver = recvType
				for _, param := range methodFunction.Type.Params.List {
					for _, name := range param.Names {
						parameter := Param{
							getType(src, param.Type, offset),
							name.Name,
						}
						ret.Params = append(ret.Params, parameter)
					}
				}	
			}
			
			return true
		}
		return true
	})
	return &ret
}