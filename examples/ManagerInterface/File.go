package ManagerInterface

import (
	"reflect"
	"testing"
)

type File struct {
	Name   string
	Path   string
	isOpen bool
}

func NewFile() (f *File) {
	f = new(File)
	return
}

func (f *File) SetName(name string) *File {
	f.Name = name
	return f
}

func TestNewFile(t *testing.T) {
	tests := []struct {
		name  string
		wantF *File
	}{
		{
			name:  "Success",
			wantF: &File{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotF := NewFile(); !reflect.DeepEqual(gotF, tt.wantF) {
				t.Errorf("NewFile() = %v, wantF %v", gotF, tt.wantF)
			}
		})
	}
}

func TestFile_SetName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *File
	}{
		{
			name: "Success",
			args: args{
				name: "myName",
			},
			want: &File{
				Name: "myName",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{}
			if got := f.SetName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetName() = %v, want %v", got, tt.want)
			}
		})
	}
}
