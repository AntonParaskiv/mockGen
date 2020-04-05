package old

import "testing"

func Test_getMapKeyValueTypes(t *testing.T) {
	type args struct {
		fieldType string
	}
	tests := []struct {
		name          string
		args          args
		wantKeyType   string
		wantValueType string
	}{
		{
			name: "Easy",
			args: args{
				fieldType: "map[int]string",
			},
			wantKeyType:   "int",
			wantValueType: "string",
		},
		{
			name: "Medium",
			args: args{
				fieldType: "map[int][]string",
			},
			wantKeyType:   "int",
			wantValueType: "[]string",
		},
		{
			name: "Hard",
			args: args{
				fieldType: "map[string]map[string]int",
			},
			wantKeyType:   "string",
			wantValueType: "map[string]int",
		},
		{
			name: "Godlike",
			args: args{
				fieldType: "map[int][]map[string][]int",
			},
			wantKeyType:   "map[[]int][]string",
			wantValueType: "map[[]string][]int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKeyType, gotValueType := getMapKeyValueTypes(tt.args.fieldType)
			if gotKeyType != tt.wantKeyType {
				t.Errorf("getMapKeyValueTypes() gotKeyType = %v, want %v", gotKeyType, tt.wantKeyType)
			}
			if gotValueType != tt.wantValueType {
				t.Errorf("getMapKeyValueTypes() gotValueType = %v, want %v", gotValueType, tt.wantValueType)
			}
		})
	}
}
