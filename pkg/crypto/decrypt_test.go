package crypto

import (
	"bytes"
	"compress/gzip"
	"crypto/rsa"
	"fmt"
	"io"
	"testing"
)

func TestDecryption_Decrypt(t *testing.T) {
	encryption, err := NewEncryption("assets/test_public.key")
	if err != nil {
		t.Errorf("failed to build encryption")
	}

	decryption, err := NewDecryption("assets/test_private.key")
	if err != nil {
		t.Errorf("failed to build decryption")
	}

	payload := []byte(`here is my message`)

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	if _, err = writer.Write(payload); err != nil {
		t.Errorf("failed to compress data: %s", err)
	}

	if err := writer.Close(); err != nil {
		t.Errorf("failed to close writer: %s", err)
	}

	encrypted, err := encryption.Encrypt(buf.Bytes())
	if err != nil {
		t.Errorf("failed to encrypt: %s", err)
	}

	reader, err := gzip.NewReader(&buf)
	if err != nil {
		t.Errorf("failed to get reader: %s", err)
	}
	decompressed, err := io.ReadAll(reader)
	if err != nil {
		t.Errorf("failed to decompress: %s", err)
	}

	if !bytes.Equal(payload, decompressed) {
		t.Errorf("failed to compress: %s", err)
	}

	type fields struct {
		PrivateKey *rsa.PrivateKey
		Nonce      [12]byte
	}
	type args struct {
		encrypted []byte
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
			fields: fields(*decryption),
			args: args{
				encrypted: encrypted,
			},
			want:    decompressed,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Decryption{
				PrivateKey: tt.fields.PrivateKey,
				Nonce:      tt.fields.Nonce,
			}
			decrypted, err := d.Decrypt(tt.args.encrypted)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decryption.Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			reader, err := gzip.NewReader(bytes.NewBuffer(decrypted))
			if err != nil {
				t.Errorf("failed to get reader: %s", err)
			}
			got, err := io.ReadAll(reader)
			if err != nil {
				t.Errorf("failed to decompress: %s", err)
			}

			if !bytes.Equal(got, tt.want) {
				fmt.Println("got: ", string(got))
				fmt.Println("want: ", string(tt.want))
				t.Errorf("Decryption.Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
