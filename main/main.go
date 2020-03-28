package main

import (
	"fmt"
	"golang.org/x/tools/imports"
)

func main() {

	mock := &Mock{
		Name:      "Manager",
		FieldList: []*Field{},
		Constructor: &Constructor{
			Name: "New",
		},
		SetterList: []*Setter{},
		MethodList: []*Method{
			{
				Name: "Registration",
				ArgList: []*Field{
					{
						Name: "nickName",
						Type: "string",
					},
					{
						Name: "password",
						Type: "string",
					},
				},
				ResultList: []*Field{
					{
						Name: "accountId",
						Type: "int64",
					},
					{
						Name: "checkCode",
						Type: "string",
					},
				},
				ArgNameTypeList:    []string{},
				ResultNameTypeList: []string{},
			},
			{
				Name: "SignIn",
				ArgList: []*Field{
					{
						Name: "accountId",
						Type: "int64",
					},
					{
						Name: "password",
						Type: "string",
					},
				},
				ResultList: []*Field{
					{
						Name: "nickName",
						Type: "string",
					},
				},
				ArgNameTypeList:    []string{},
				ResultNameTypeList: []string{},
			},
		},
	}
	mock.ReceiverName = getReceiverName(mock.Name)
	mock.WantName = "want" + toPublic(mock.Name)
	mock.GotName = "got" + toPublic(mock.Name)

	argList := []*Field{}
	resultList := []*Field{}
	for _, method := range mock.MethodList {
		for _, arg := range method.ArgList {
			arg.WantName = "want" + toPublic(arg.Name)
			arg.GotName = "got" + toPublic(arg.Name)
			arg.NameType = arg.Name + " " + arg.Type

			switch {
			case arg.Type == "string":
				arg.ExampleValue = `"my` + toPublic(arg.Name) + `"`
			case arg.Type == "bool":
				arg.ExampleValue = "true"
			case arg.Type == "rune":
				arg.ExampleValue = `"X"`
			case arg.Type == "byte":
				arg.ExampleValue = `50`
			case len(arg.Type) >= 3 && arg.Type[0:3] == "int":
				arg.ExampleValue = "100"
			case len(arg.Type) >= 4 && arg.Type[0:4] == "uint":
				arg.ExampleValue = "200"
			case len(arg.Type) >= 5 && arg.Type[0:5] == "float":
				arg.ExampleValue = "3.14"
			}

			method.ArgNameTypeList = append(method.ArgNameTypeList, arg.NameType)
		}
		for _, result := range method.ResultList {
			result.WantName = "want" + toPublic(result.Name)
			result.GotName = "got" + toPublic(result.Name)
			result.NameType = result.Name + " " + result.Type

			switch {
			case result.Type == "string":
				result.ExampleValue = `"my` + toPublic(result.Name) + `"`
			case result.Type == "bool":
				result.ExampleValue = "true"
			case result.Type == "rune":
				result.ExampleValue = `"X"`
			case result.Type == "byte":
				result.ExampleValue = `50`
			case len(result.Type) >= 3 && result.Type[0:3] == "int":
				result.ExampleValue = "100"
			case len(result.Type) >= 4 && result.Type[0:4] == "uint":
				result.ExampleValue = "200"
			case len(result.Type) >= 5 && result.Type[0:5] == "float":
				result.ExampleValue = "3.14"
			}
			method.ResultNameTypeList = append(method.ResultNameTypeList, result.NameType)
		}

		argList = append(argList, method.ArgList...)
		resultList = append(resultList, method.ResultList...)

	}

	// TODO: check unique
ArgLoop:
	for _, arg := range argList {
		for _, mockField := range mock.FieldList {
			if mockField.Name == arg.Name {
				continue ArgLoop
			}
		}
		mock.FieldList = append(mock.FieldList, arg)
	}

ResultLoop:
	for _, result := range resultList {
		for _, mockField := range mock.FieldList {
			if mockField.Name == result.Name {
				continue ResultLoop
			}
		}
		mock.FieldList = append(mock.FieldList, result)
	}

	for _, field := range mock.FieldList {
		mock.SetterList = append(mock.SetterList, &Setter{
			Name:  "Set" + toPublic(field.Name),
			Field: field,
		})
	}

	result := "package main\n"
	result += PrintMock(mock)

	//fmt.Println(result)
	//return

	formattedResult, err := imports.Process("", []byte(result), &imports.Options{
		Fragment:   true,
		AllErrors:  true,
		Comments:   true,
		TabIndent:  true,
		TabWidth:   8,
		FormatOnly: false,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(formattedResult))

	return

	// TODO: arg get package path
	//packagePath := "domain"
	packagePath := "examples/ManagerInterface"

	rPkgs, err := getRPackages(packagePath)
	if err != nil {
		panic(err)
	}

	wPkgs, err := makeWPackages(rPkgs)
	if err != nil {
		panic(err)
	}

	_ = wPkgs

	err = savePackages(wPkgs)
	if err != nil {
		panic(err)
	}

	return
}
