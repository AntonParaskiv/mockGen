package main

func main() {
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
