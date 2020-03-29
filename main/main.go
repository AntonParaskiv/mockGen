package main

func main() {
	iFace := &Interface{
		Name: "Manager",
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
			},
		},
	}

	mock := CreateMock(iFace)
	GenCodeMock(mock)
	_ = mock

	//fmt.Println(result)
	//return

	//formattedResult, err := imports.Process("", []byte(result), &imports.Options{
	//	Fragment:   true,
	//	AllErrors:  true,
	//	Comments:   true,
	//	TabIndent:  true,
	//	TabWidth:   8,
	//	FormatOnly: false,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(string(formattedResult))

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
