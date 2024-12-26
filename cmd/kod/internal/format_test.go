package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImportsCode(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		want    string
		wantErr bool
	}{
		{
			name: "simple format",
			code: `package main
import "fmt"
func main(){
fmt.Println("hello")}`,
			want: `package main

import "fmt"

func main() {
	fmt.Println("hello")
}
`,
			wantErr: false,
		},
		{
			name:    "invalid code",
			code:    "invalid{",
			wantErr: true,
		},
		{
			name: "already formatted",
			code: `package main

import "fmt"

func main() {
	fmt.Println("hello")
}`,
			want: `package main

import "fmt"

func main() {
	fmt.Println("hello")
}
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ImportsCode(tt.code)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, string(got))
		})
	}
}
