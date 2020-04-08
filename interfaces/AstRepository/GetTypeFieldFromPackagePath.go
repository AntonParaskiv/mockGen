package AstRepository

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
)

func (r *Repository) GetTypeFieldFromPackagePath(packagePath, typeName string) (field *domain.Field, err error) {
	typeSpec, err := r.getTypeDeclarationFromPackagePath(packagePath, typeName)
	if err != nil {
		err = fmt.Errorf("get type declaration from package path failed: %w", err)
		return
	}
	if typeSpec == nil {
		return
	}

	fieldType, err := getFieldType(typeSpec.Type)
	if err != nil {
		err = fmt.Errorf("get field type from failed: %w", err)
		return
	}

	field = &domain.Field{
		Name: getNodeName(typeSpec),
		Type: fieldType,
	}
	return
}

func (r *Repository) getTypeDeclarationFromPackagePath(packagePath, typeName string) (typeSpec *ast.TypeSpec, err error) {
	packageGoPath := createPackageGoPath(packagePath)

	astPackage, err := r.CodeStorage.GetAstPackage(packageGoPath)
	if err != nil {
		err = fmt.Errorf("get ast package %s failed: %w", packageGoPath, err)
		return
	}

	typeSpec = getTypeDeclarationFromAstPackage(astPackage, typeName)
	return
}

func createPackageGoPath(packagePath string) (packageGoPath string) {
	goPathSrc := filepath.Join(os.Getenv("GOPATH"), "src")

	if len(packagePath) < len(goPathSrc) {
		packageGoPath = packagePath
		return
	}

	if packagePath[:len(goPathSrc)] == goPathSrc {
		packageGoPath = packagePath
		return
	}

	packageGoPath = filepath.Join(goPathSrc, packagePath)
	return
}

func getTypeDeclarationFromAstPackage(astPackage *ast.Package, typeName string) (typeSpec *ast.TypeSpec) {
	for _, astFile := range astPackage.Files {
		for _, decl := range astFile.Decls {
			switch decl := decl.(type) {
			case *ast.GenDecl:
				switch decl.Tok {
				case token.TYPE:
					spec := decl.Specs[0].(*ast.TypeSpec) // TODO: check array
					if getNodeName(spec) == typeName {
						typeSpec = spec
						return
					}
				}
			}
		}
	}
	return
}
