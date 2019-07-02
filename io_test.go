package pbqp

import (
	"reflect"
	"testing"
)

func TestBStoI64(t *testing.T) {
	type args struct {
		bs []byte
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "normal",
			args: args{bs: []byte{1, 0, 0, 0, 0, 1, 0, 90}},
			want: 72057594037993562,
		},
		{
			name: "small",
			args: args{bs: []byte{0, 0, 0, 90}},
			want: 90,
		},
		{
			name: "large",
			args: args{bs: []byte{1, 0, 1, 0, 0, 0, 0, 1, 0, 90}},
			want: 72057594037993562,
		},
		{
			name: "-1",
			args: args{bs: []byte{255, 255, 255, 255, 255, 255, 255, 255}},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BStoI64(tt.args.bs); got != tt.want {
				t.Errorf("BStoI64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestI64toBS(t *testing.T) {
	type args struct {
		n int64
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "-1",
			args: args{n: -1},
			want: []byte{255, 255, 255, 255, 255, 255, 255, 255},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := I64toBS(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("I64toBS() = %v, want %v", got, tt.want)
			}
		})
	}
}
