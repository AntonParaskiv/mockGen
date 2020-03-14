package main

import (
	"fmt"
	"go/ast"
	"strings"
)

func getNodeName(node ast.Node) (name string) {
	switch nodeItem := node.(type) {
	case *ast.Package:
		name = nodeItem.Name
	case *ast.File:
		name = nodeItem.Name.Name
	case *ast.TypeSpec:
		name = nodeItem.Name.Name
	case *ast.Field:
		name = nodeItem.Names[0].Name
	default:
		panic(fmt.Sprintf("no getting name case for type %T", node))
	}
	return
}

func toPublic(name string) (publicName string) {
	firstLetterUpper := strings.ToUpper(getFirstLetter(name))
	publicName = firstLetterUpper + getFollowingLetters(name)
	return
}

func toPrivate(name string) (privateName string) {
	firstLetterLower := strings.ToLower(getFirstLetter(name))
	privateName = firstLetterLower + getFollowingLetters(name)
	return
}

func getFirstLetter(text string) (firstLetter string) {
	firstLetter = text[0:1]
	return
}

func getFollowingLetters(text string) (followingLetters string) {
	followingLetters = text[1:]
	return
}

func createName(name string) (names *ast.Ident) {
	names = &ast.Ident{
		Name: name,
	}
	return
}

func createNames(name string) (names []*ast.Ident) {
	names = []*ast.Ident{
		{
			Name: name,
		},
	}
	return
}
