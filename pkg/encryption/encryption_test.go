package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/go-playground/assert/v2"
	"reflect"
	"testing"
)

func TestEncryption_MapReflectField(t *testing.T) {
	type s struct {
		Int          int    `encryption:""`
		String       string `encryption:""`
		String2      string
		StringArray  []string `encryption:""`
		StringArray2 []string ``
		IntArray     []int    `encryption:""`
		CustomStruct struct {
			String string `encryption:""`
		} `encryption:""`
	}

	type args struct {
		value s
	}
	tests := []struct {
		name string
		args args
		want s
	}{
		{
			name: "Valid test",
			args: args{
				value: s{
					Int:          10,
					String:       "a",
					String2:      "a",
					StringArray:  []string{"1", "2", "3"},
					StringArray2: []string{"1", "2", "3"},
					IntArray:     []int{1, 2, 3},
					CustomStruct: struct {
						String string `encryption:""`
					}{"a"},
				},
			},
			want: s{
				Int:          10,
				String:       "encrypted",
				String2:      "a",
				StringArray:  []string{"encrypted", "encrypted", "encrypted"},
				StringArray2: []string{"1", "2", "3"},
				IntArray:     []int{1, 2, 3},
				CustomStruct: struct {
					String string `encryption:""`
				}{"encrypted"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &encryption{}
			e.mapReflectField(reflect.ValueOf(&tt.args.value).Elem(), false, func(str string, salt bool) (string, error) {
				return "encrypted", nil
			})

			assert.Equal(t, tt.args.value, tt.want)
		})
	}
}

func Test_encryption_EncryptString(t *testing.T) {
	type fields struct {
		cipher cipher.Block
	}
	type args struct {
		s1 string
		s2 string
	}

	cipherImpl, _ := aes.NewCipher([]byte("z$B&E)H@McQfTjWnZr4u7x!A%D*F-JaN"))

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantEqual bool
	}{
		{
			name: "without salt",
			fields: fields{
				cipher: cipherImpl,
			},
			args: args{
				s1: "string",
				s2: "string",
			},
			wantEqual: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &encryption{
				cipher: tt.fields.cipher,
			}

			e1, _ := e.EncryptString(tt.args.s1, false)
			e2, _ := e.EncryptString(tt.args.s2, false)

			d1, _ := e.DecryptString(e1, false)
			d2, _ := e.DecryptString(e2, false)

			if (e1 == e2) != tt.wantEqual {
				t.Errorf("EncryptString() = %v, EncryptString() = %v, wantEqual %v", e1, e2, tt.wantEqual)
			}

			if d1 != tt.args.s1 {
				t.Errorf("DecryptString() = %v, want %v", d1, tt.args.s1)
			}

			if d2 != tt.args.s2 {
				t.Errorf("DecryptString() = %v, want %v", d2, tt.args.s1)
			}
		})
	}
}
