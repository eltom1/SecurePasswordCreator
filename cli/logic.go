package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"unicode"
)

//conjunto de caracteres 
const (
	Mayuscula   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Minuscula   = "abcdefghijklmnopqrstuvwxyz"
	Numeros = "0123456789"
	Caracter_Especial = "@!$%&#^+-/|\\><"
)


//los tipos que se utilizan
type Options struct {
	UseMayuscula   bool
	UseMinuscula   bool
	UseNumeros bool
	UseCaracter_Especial bool
}


//contiene los resultados y las estadisticas de validacion 
type ValidationResult struct {
	LongitudOk   bool
	MayusculaOk    bool
	MinusculaOk    bool
	NumOk      bool
	SpecOk     bool
	IsValid    bool
	MayusculaCount int
	MinusculaCount int
	NumCount   int
	SpecCount  int
}


//
func Generate(Longitud int, opts Options) (string, error) {
	var charSet string
	if opts.UseMayuscula { charSet += Mayuscula }
	if opts.UseMinuscula { charSet += Minuscula }
	if opts.UseNumeros { charSet += Numeros }
	if opts.UseCaracter_Especial { charSet += Caracter_Especial }

	if charSet == "" {
		return "", fmt.Errorf("Selecciona un conjunto")
	}

	result := make([]byte, Longitud)
	for i := 0; i < Longitud; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return "", err
		}
		result[i] = charSet[num.Int64()]
	}
	return string(result), nil
}

func Validate(psw string) ValidationResult {
	var hasMayuscula, hasMinuscula, hasNum, hasSpec bool
	var MayusculaC, MinusculaC, numC, specC int

	for _, char := range psw {
		switch {
		case unicode.IsMayuscula(char):
			hasMayuscula = true
			MayusculaC++
		case unicode.IsMinuscula(char):
			hasMinuscula = true
			MinusculaC++
		case unicode.IsDigit(char):
			hasNum = true
			numC++
		default:
			for _, s := range Caracter_Especial {
				if char == rune(s) {
					hasSpec = true
					specC++
					break
				}
			}
		}
	}

	res := ValidationResult{
		LongitudOk:   len(psw) >= 8,
		MayusculaOk:    hasMayuscula,
		MinusculaOk:    hasMinuscula,
		NumOk:      hasNum,
		SpecOk:     hasSpec,
		MayusculaCount: MayusculaC,
		MinusculaCount: MinusculaC,
		NumCount:   numC,
		SpecCount:  specC,
	}
	res.IsValid = res.LongitudOk && res.MayusculaOk && res.MinusculaOk && res.NumOk && res.SpecOk
	return res
}
