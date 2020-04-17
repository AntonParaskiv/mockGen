package Interactor

import (
	"github.com/AntonParaskiv/mockGen/domain"
	"github.com/AntonParaskiv/mockGen/interfaces/AstRepository"
	"reflect"
	"testing"
)

func TestInteractor_createField(t *testing.T) {
	type fields struct {
		AstRepository       AstRepository.Repository
		mockFile            *domain.GoCodeFile
		interfacePackage    *domain.GoCodePackage
		mockPackage         *domain.GoCodePackage
		CreateFieldExamples bool
	}
	type args struct {
		iFaceField *domain.Field
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantMockField *domain.Field
	}{
		{
			name:   "String",
			fields: fields{},
			args: args{
				iFaceField: &domain.Field{
					Name: "myName",
					Type: "string",
				},
			},
			wantMockField: &domain.Field{
				Name: "myName",
				Type: "string",
			},
		},
		{
			name:   "Pointer string",
			fields: fields{},
			args: args{
				iFaceField: &domain.Field{
					Name: "myName",
					Type: "*string",
				},
			},
			wantMockField: &domain.Field{
				Name: "myName",
				Type: "*string",
				BaseType: &domain.Field{
					Name: "myName",
					Type: "string",
				},
			},
		},
		{
			name:   "Array string",
			fields: fields{},
			args: args{
				iFaceField: &domain.Field{
					Name: "myName",
					Type: "[]string",
				},
			},
			wantMockField: &domain.Field{
				Name: "myName",
				Type: "[]string",
				BaseType: &domain.Field{
					Name: "myName",
					Type: "string",
				},
			},
		},
		{
			name:   "Array pointer string",
			fields: fields{},
			args: args{
				iFaceField: &domain.Field{
					Name: "myName",
					Type: "[]*string",
				},
			},
			wantMockField: &domain.Field{
				Name: "myName",
				Type: "[]*string",
				BaseType: &domain.Field{
					Name: "myName",
					Type: "*string",
					BaseType: &domain.Field{
						Name: "myName",
						Type: "string",
					},
				},
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interactor{
				AstRepository:       tt.fields.AstRepository,
				mockFile:            tt.fields.mockFile,
				interfacePackage:    tt.fields.interfacePackage,
				mockPackage:         tt.fields.mockPackage,
				CreateFieldExamples: tt.fields.CreateFieldExamples,
			}
			if gotMockField := i.createField(tt.args.iFaceField); !reflect.DeepEqual(gotMockField, tt.wantMockField) {
				t.Errorf("createField() = %v, want %v", gotMockField, tt.wantMockField)
			}
		})
	}
}
