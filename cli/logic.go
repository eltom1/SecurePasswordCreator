package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"unicode"
)

const (
	Upper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Lower   = "abcdefghijklmnopqrstuvwxyz"
	Numbers = "0123456789"
	Special = "@!$%&#^+-/|\\><"
)

type Options struct {
	UseUpper   bool
	UseLower   bool
	UseNumbers bool
	UseSpecial bool
}

type ValidationResult struct {
	LengthOk   bool
	UpperOk    bool
	LowerOk    bool
	NumOk      bool
	SpecOk     bool
	IsValid    bool
	UpperCount int
	LowerCount int
	NumCount   int
	SpecCount  int
}

func Generate(length int, opts Options) (string, error) {
	var charSet string
	if opts.UseUpper { charSet += Upper }
	if opts.UseLower { charSet += Lower }
	if opts.UseNumbers { charSet += Numbers }
	if opts.UseSpecial { charSet += Special }

	if charSet == "" {
		return "", fmt.Errorf("select at least one charset")
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return "", err
		}
		result[i] = charSet[num.Int64()]
	}
	return string(result), nil
}

func Validate(psw string) ValidationResult {
	var hasUpper, hasLower, hasNum, hasSpec bool
	var upperC, lowerC, numC, specC int

	for _, char := range psw {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
			upperC++
		case unicode.IsLower(char):
			hasLower = true
			lowerC++
		case unicode.IsDigit(char):
			hasNum = true
			numC++
		default:
			for _, s := range Special {
				if char == rune(s) {
					hasSpec = true
					specC++
					break
				}
			}
		}
	}

	res := ValidationResult{
		LengthOk:   len(psw) >= 8,
		UpperOk:    hasUpper,
		LowerOk:    hasLower,
		NumOk:      hasNum,
		SpecOk:     hasSpec,
		UpperCount: upperC,
		LowerCount: lowerC,
		NumCount:   numC,
		SpecCount:  specC,
	}
	res.IsValid = res.LengthOk && res.UpperOk && res.LowerOk && res.NumOk && res.SpecOk
	return res
}
