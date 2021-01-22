// This package provices functions to decode a CSGO share code.

package utils

import (
	"log"
	"math/big"
	"regexp"
	"strings"
)

// ShareCode holds the decoded match code
type ShareCode struct {
	MatchID   uint64
	OutcomeID uint64
	Token     uint32
}

// Dictionary is used for the share code decoding.
const Dictionary = "ABCDEFGHJKLMNOPQRSTUVWXYZabcdefhijkmnopqrstuvwxyz23456789"

// Used for share code decoding.
var bitmask64 uint64 = 18446744073709551615

// Decode decodes the share code. Taken from ValvePython/csgo
func Decode(code string) *ShareCode {
	var validateRe = regexp.MustCompile(`^CSGO(-?[\w]{5}){5}$`)
	found := validateRe.MatchString(code)

	if !found {
		log.Print("invalid share code")
		return nil
	}

	var re = regexp.MustCompile(`^CSGO|\-`)
	s := re.ReplaceAllString(code, "")
	s = reverse(s)

	bigNumber := big.NewInt(0)

	for _, c := range s {
		bigNumber = bigNumber.Mul(bigNumber, big.NewInt(int64(len(Dictionary))))
		bigNumber = bigNumber.Add(bigNumber, big.NewInt(int64(strings.Index(Dictionary, string(c)))))
	}

	a := swapEndianness(bigNumber)

	matchid := big.NewInt(0)
	outcomeid := big.NewInt(0)
	token := big.NewInt(0)

	matchid = matchid.And(a, big.NewInt(0).SetUint64(bitmask64))
	outcomeid = outcomeid.Rsh(a, 64)
	outcomeid = outcomeid.And(outcomeid, big.NewInt(0).SetUint64(bitmask64))
	token = token.Rsh(a, 128)
	token = token.And(token, big.NewInt(0xFFFF))

	sharecode := &ShareCode{MatchID: matchid.Uint64(), OutcomeID: outcomeid.Uint64(), Token: uint32(token.Uint64())}
	return sharecode
}

// reverse reverses a string
func reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

// swapEndianness changes the byte order
func swapEndianness(number *big.Int) *big.Int {
	result := big.NewInt(0)

	left := big.NewInt(0)
	rightTemp := big.NewInt(0)
	rightResult := big.NewInt(0)

	for i := 0; i < 144; i += 8 {
		left = left.Lsh(result, 8)
		rightTemp = rightTemp.Rsh(number, uint(i))
		rightResult = rightResult.And(rightTemp, big.NewInt(0xFF))
		result = left.Add(left, rightResult)
	}

	return result
}
