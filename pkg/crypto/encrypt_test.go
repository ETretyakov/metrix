package crypto

import (
	"crypto/rsa"
	"reflect"
	"testing"
)

func TestEncryption_Encrypt(t *testing.T) {
	encryption, err := NewEncryption("assets/test_public.key")
	if err != nil {
		t.Errorf("failed to build encryption")
	}

	decryption, err := NewDecryption("assets/test_private.key")
	if err != nil {
		t.Errorf("failed to build decryption")
	}

	type fields struct {
		PublicKey *rsa.PublicKey
		Nonce     [12]byte
	}
	type args struct {
		payload []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:   "Test #1 Success",
			fields: fields(*encryption),
			args: args{
				payload: []byte("some test message"),
			},
			want:    []byte("some test message"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Encryption{
				PublicKey: tt.fields.PublicKey,
				Nonce:     tt.fields.Nonce,
			}
			encrypted, err := e.Encrypt(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encryption.Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := decryption.Decrypt(encrypted)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encryption.Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encryption.Encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
