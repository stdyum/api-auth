package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"reflect"

	"github.com/pkg/errors"
)

type Encryption interface {
	EncryptString(s string, salt bool) (string, error)
	DecryptString(s string, salt bool) (string, error)

	Encrypt(value any) error
	Decrypt(value any) error
}

type encryption struct {
	cipher cipher.Block
}

func NewEncryption(secret string) (Encryption, error) {
	cipher_, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &encryption{cipher: cipher_}, nil
}

func (e *encryption) Encrypt(value interface{}) error {
	return e.mapReflectField(reflect.ValueOf(value).Elem(), false, func(str string, salt bool) (string, error) {
		return e.EncryptString(str, salt)
	})
}

func (e *encryption) Decrypt(value interface{}) error {
	return e.mapReflectField(reflect.ValueOf(value).Elem(), false, func(str string, salt bool) (string, error) {
		return e.DecryptString(str, salt)
	})
}

func (e *encryption) EncryptString(s string, salt bool) (string, error) {
	plaintext := []byte(s)

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if salt {
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return "", err
		}
	}

	stream := cipher.NewCFBEncrypter(e.cipher, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	if !salt {
		ciphertext = ciphertext[aes.BlockSize:]
	}

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (e *encryption) DecryptString(s string, salt bool) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(s)
	if !salt {
		ciphertext = append(make([]byte, aes.BlockSize), ciphertext...)
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("chipper less than block size")
	}

	if len(ciphertext) == aes.BlockSize {
		return "", nil
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(e.cipher, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil
}

func (e *encryption) mapReflectSlice(encrypt bool, salt bool, field reflect.Value, mapFunc func(str string, salt bool) (string, error)) error {
	for j := 0; j < field.Len(); j++ {
		switch field.Index(j).Kind() {
		case reflect.String:
			if !encrypt {
				break
			}

			if err := e.mapReflectField(field.Index(j), salt, mapFunc); err != nil {
				return err
			}
		case reflect.Struct:
			if err := e.mapReflectField(field.Index(j), salt, mapFunc); err != nil {
				return err
			}
		default:
		}
	}

	return nil
}

func (e *encryption) mapReflectField(value reflect.Value, salt bool, mapFunc func(str string, salt bool) (string, error)) error {
	switch value.Kind() {
	case reflect.String:
		valueToSet, err := mapFunc(value.String(), salt)
		if err != nil {
			return err
		}
		value.SetString(valueToSet)
		return nil
	case reflect.Array, reflect.Slice:
		return e.mapReflectSlice(true, salt, value, mapFunc)
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			tField := value.Type().Field(i)
			field := value.Field(i)

			switch field.Kind() {
			case reflect.Struct:
				if err := e.mapReflectField(field, salt, mapFunc); err != nil {
					return err
				}
			case reflect.Slice, reflect.Array:
				tag, ok := tField.Tag.Lookup("encryption")
				salt = tag != "-salt"
				if err := e.mapReflectSlice(ok, salt, field, mapFunc); err != nil {
					return err
				}
			case reflect.String:
				tag, ok := tField.Tag.Lookup("encryption")
				if !ok {
					continue
				}
				salt = tag != "-salt"

				fieldValue := field.String()
				encrypted, err := mapFunc(fieldValue, salt)
				if err != nil {
					return err
				}

				field.SetString(encrypted)
			default:
			}
		}
	default:
	}

	return nil
}
