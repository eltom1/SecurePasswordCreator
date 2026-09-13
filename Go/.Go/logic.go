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
type Opciones struct {
	UsaMayus   bool
	UsaMinus   bool
	UsaNum bool
	UsaEsp bool
}


//contiene los resultados y las estadisticas de validacion
type ResultadoVal struct {
	LargoOk   bool
	MayusOk    bool
	MinusOk    bool
	NumOk      bool
	EspOk     bool
	EsValido    bool
	CantMayus int
	CantMinus int
	CantNum   int
	CantEsp  int
}


//
func Generar(Largo int, opts Opciones) (string, error) {
	var conjunto string
	if opts.UsaMayus { conjunto += Mayuscula }
	if opts.UsaMinus { conjunto += Minuscula }
	if opts.UsaNum { conjunto += Numeros }
	if opts.UsaEsp { conjunto += Caracter_Especial }

	if conjunto == "" {
		return "", fmt.Errorf("Selecciona un conjunto")
	}

	res := make([]byte, Largo)
	for i := 0; i < Largo; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(conjunto))))
		if err != nil {
			return "", err
		}
		res[i] = conjunto[num.Int64()]
	}
	return string(res), nil
}

func Validar(psw string) ResultadoVal {
	var hasMayus, hasMinus, hasNum, hasEsp bool
	var MayusC, MinusC, numC, espC int

	for _, char := range psw {
		switch {
		case unicode.IsUpper(char):
			hasMayus = true
			MayusC++
		case unicode.IsLower(char):
			hasMinus = true
			MinusC++
		case unicode.IsDigit(char):
			hasNum = true
			numC++
		default:
			for _, s := range Caracter_Especial {
				if char == rune(s) {
					hasEsp = true
					espC++
					break
				}
			}
		}
	}

	res := ResultadoVal{
		LargoOk:   len(psw) >= 8,
		MayusOk:    hasMayus,
		MinusOk:    hasMinus,
		NumOk:      hasNum,
		EspOk:     hasEsp,
		CantMayus: MayusC,
		CantMinus: MinusC,
		CantNum:   numC,
		CantEsp:  espC,
	}
	res.EsValido = res.LargoOk && res.MayusOk && res.MinusOk && res.NumOk && res.EspOk
	return res
}
