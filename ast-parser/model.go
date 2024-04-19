package ast_parser

import (
	"strings"

	. "github.com/dave/jennifer/jen"
)

type Param struct {
	Type string
	Name string
}

type ReturnValue struct {
	Type string
}

type Method struct {
	Name         string
	Params       []Param
	Receiver     string
	ReturnValues []ReturnValue
}

func (m *Method) getStructName() string {
	return strings.Title(m.Name) + "Struct"
}

func (m *Method) GetStruct(file *File) {
	members := []Code{}
	for _, param := range m.Params {
		members = append(members, Id(param.Name).Id(param.Type))
	}

	if m.hasReceiver() {
		members = append(members, Id("t").Id(m.getInterfaceContractName()))
	}
	file.Type().Id(m.getStructName()).Struct(members...)
}

func (m *Method) GetNextInterface(offset int) string {
	if offset == len(m.Params)-1 {
		return strings.Title(m.Name)
	} else {
		return strings.Title(m.Params[offset+1].Name)
	}
}

func (m *Method) getReturnValues() []Code {
	retValues := []Code{}
	for _, retValue := range m.ReturnValues {
		retValues = append(retValues, Id(retValue.Type))
	}
	return retValues
}

func (m *Method) genFinalInterface(file *File) {
	file.Type().Id(strings.Title(m.Name)).Interface(Id(strings.Title(m.Name)).Params().Parens(List(m.getReturnValues()...)))
}

func (m *Method) hasReceiver() bool {
	return m.Receiver != ""
}

func (m *Method) GetInterfaces(file *File) {
	for i, param := range m.Params {
		file.Type().Id(strings.Title(param.Name)).Interface(
			Id(strings.Title(param.Name)).Params(Id(param.Name).Id(param.Type)).Id(m.GetNextInterface(i)),
		)
	}
	m.genFinalInterface(file)
}

func (m *Method) getFinalImplementation(file *File) {
	params := []Code{}
	for _, paramOfMethod := range m.Params {
		params = append(params, Id("t").Dot(paramOfMethod.Name))
	}

	if m.hasReceiver() {
		file.Func().Params(Id("t").Id("*" + m.getStructName())).Id(strings.Title(m.Name)).
			Params().
			Parens(List(m.getReturnValues()...)).
			Block(
				Return(Id("t").Dot("t").Dot(m.Name).Call(params...)),
			)
	} else {
		file.Func().Params(Id("t").Id("*" + m.getStructName())).Id(strings.Title(m.Name)).
			Params().
			Parens(List(m.getReturnValues()...)).
			Block(
				Return(Id(m.Name).Call(params...)),
			)
	}
}

func (m *Method) GetImplementations(file *File) {
	for i, param := range m.Params {
		file.Func().Params(Id("t").Id("*"+m.getStructName())).Id(strings.Title(param.Name)).
			Params(Id(param.Name).Id(param.Type)).
			Id(m.GetNextInterface(i)).Block(
			Id("t."+param.Name).Op("=").Id(param.Name),
			Return(Id("t")),
		)
	}
	m.getFinalImplementation(file)
}

func (m *Method) GetEntryMethod(file *File) {
	if m.hasReceiver() {
		receiverName := strings.ReplaceAll(m.Receiver, "*", "")
		file.Func().Id("With"+strings.Title(receiverName)).Params(Id("t").Id(m.getInterfaceContractName())).Id(strings.Title(m.Params[0].Name)).
			Block(
				Id("ret").Op(":=").Id(m.getStructName()).Block(),
				Id("ret").Dot("t").Op("=").Id("t"),
				Return().Id("&ret"),
			)
	} else {
		file.Func().Id(strings.Title(m.Name)).Params().Id(strings.Title(m.Params[0].Name)).
			Block(
				Id("ret").Op(":=").Id(m.getStructName()).Block(),
				Return().Id("&ret"),
			)
	}
}

func (m *Method) getInterfaceContractName() string {
	return strings.Title(m.Name) + "Contract"
}

func (m *Method) getParameters() []Code {
	parameters := []Code{}
	for _, param := range m.Params {
		parameters = append(parameters, Id(param.Name).Id(param.Type))
	}
	return parameters
}

func (m *Method) GetInterfaceContract(file *File) {
	if m.hasReceiver() {
		parameters := m.getParameters()
		file.Type().Id(m.getInterfaceContractName()).Interface(Id(m.Name).Params(parameters...).Parens(List(m.getReturnValues()...)))
	}
}
