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

func TestManager_Registration(t *testing.T) {
	type fields struct {
		accountId int64 // returnList only
		checkCode string
	}
	type args struct {
		nickName string
		password string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantAccountId int64
		wantCheckCode string
		wantManager   *Manager
	}{
		{
			name: "Success",
			fields: fields{
				accountId: 100,
				checkCode: "myCheckCode",
			},
			args: args{
				nickName: "myNickName",
				password: "myPassword",
			},
			wantAccountId: 100,
			wantCheckCode: "myCheckCode",
			wantManager: &Manager{
				accountId: 100,
				password:  "myPassword",
				nickName:  "myNickName",
				checkCode: "myCheckCode",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				accountId: tt.fields.accountId,
				checkCode: tt.fields.checkCode,
			}
			gotAccountId, gotCheckCode := m.Registration(tt.args.nickName, tt.args.password)
			if gotAccountId != tt.wantAccountId {
				t.Errorf("gotAccountId = %v, want %v", gotAccountId, tt.wantAccountId)
			}
			if gotCheckCode != tt.wantCheckCode {
				t.Errorf("gotCheckCode = %v, want %v", gotCheckCode, tt.wantCheckCode)
			}
			if !reflect.DeepEqual(m, tt.wantManager) {
				t.Errorf("Manager = %v, want %v", m, tt.wantManager)
			}
		})
	}
}
