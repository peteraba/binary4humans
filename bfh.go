package bfh

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	digits                         = "0123456789abcdefghjkmnpqrstvwxyz"
	errMsgBinaryDataMustNotBeNil   = "binary data must not be nil"
	errMsgStrictMustBeDividableBy5 = "length of binary data must be some multiple of 5 for strict encoding"
	errMsgStrictMustBeDividableBy8 = "length of encoded string must be some multiple of 8 for strict decoding"
	errMsgStrictInvalid            = "invalid encoded string for strict mode"
	errMsgPaddingNotBetween0and4   = "non empty string must start with 0, 1, 2, 3 or 4"
	errMsgContainsInvalidCharacter = "string contains invalid character: %s"

	// note that valid encoded strings will not end in a hyphen, it needs to be added when validating
	standardBfhRegexString   = "^[0-4]\\-([0123456789abcdefghjkmnpqrstvwxyz]{4}\\-)*$"
	acceptableBfhRegexString = "^[0-4]([0123456789abcdefghjkmnpqrstvwxyz]{4})*$"
	strictBfhRegexString     = "^([0123456789abcdefghjkmnpqrstvwxyz]{4}\\-)*$"
)

var (
	encodeMasks        = []uint8{0xff, 0x7f, 0x3f, 0x1f, 0x0f, 0x07, 0x03, 0x01}
	decodeMasks        = []uint8{0x1, 0x3, 0x7, 0xf}
	digitMap           map[string]uint8
	standardBfhRegex   = regexp.MustCompile(standardBfhRegexString)
	acceptableBfhRegex = regexp.MustCompile(acceptableBfhRegexString)
	strictBfhRegex     = regexp.MustCompile(strictBfhRegexString)
)

func init() {
	digitMap = map[string]uint8{}
	for i := 0; i < len(digits); i++ {
		digitMap[digits[i:i+1]] = uint8(i)
	}
}

// Encode encodes binary data into a human readable string
func Encode(b []byte) (string, error) {
	var (
		strBuilder strings.Builder
		err        error
	)

	if b == nil {
		return "", errors.New(errMsgBinaryDataMustNotBeNil)
	}

	b, err = padBytes(&strBuilder, b)
	if err != nil {
		return "", err
	}

	return encode(b, &strBuilder)
}

// EncodeStrict encodes binary data with a length dividable by 5 into a simplified human readable string
func EncodeStrict(b []byte) (string, error) {
	var (
		strBuilder strings.Builder
	)

	if b == nil {
		return "", errors.New(errMsgBinaryDataMustNotBeNil)
	}

	if len(b)%5 != 0 {
		return "", errors.New(errMsgStrictMustBeDividableBy5)
	}

	return encode(b, &strBuilder)
}

func encode(b []byte, strBuilder *strings.Builder) (string, error) {
	var (
		readBits = 0
		maxBits  = len(b) * 8
		err      error
	)

	for maxBits > readBits {
		f := readByte(b, readBits)

		_, err = fmt.Fprintf(strBuilder, "%s", digits[f:f+1])
		if err != nil {
			return "", err
		}

		readBits += 5

		if readBits%20 != 0 || readBits >= maxBits {
			continue
		}

		_, err = fmt.Fprint(strBuilder, "-")
		if err != nil {
			return "", err
		}
	}

	return strBuilder.String(), nil
}

func padBytes(strBuilder *strings.Builder, b []byte) ([]byte, error) {
	pad := (5 - len(b)%5) % 5
	_, err := fmt.Fprintf(strBuilder, "%d-", pad)
	if err != nil {
		return nil, err
	}

	if len(b)%5 != 0 {
		c := make([]byte, len(b)+pad)
		copy(c, b)
		return c, nil
	}

	return b, nil
}

func readByte(b []byte, readBits int) byte {
	m := uint8(readBits % 8)

	// we need to mask out bits we've read before
	f := b[readBits/8] & encodeMasks[m]

	// if we're reading from the first half of the byte than we have it easy...
	// otherwise we need to push some bits towards larger value to leave place
	// for bits from the next byte
	if m < 4 {
		return f >> (3 - m)
	}

	f <<= m - 3

	// next byte may not exist
	// otherwise we'll need push the first bits to the end
	var s byte
	if len(b) > readBits/8+1 {
		s = b[readBits/8+1] >> (11 - m)
	}

	return f | s
}

// Decode decodes a human readable string into a binary data
func Decode(str string) ([]byte, error) {
	// dashes are not needed, they only help readability
	str = strings.Replace(str, "-", "", -1)

	padding, ok := digitMap[str[0:1]]
	if !ok || padding > 4 {
		return nil, errors.New(errMsgPaddingNotBetween0and4)
	}

	str = str[1:]

	if len(str)%8 != 0 {
		return nil, errors.New(errMsgStrictMustBeDividableBy8)
	}

	bytes, err := decode(str)
	if err != nil {
		return nil, err
	}

	if padding > 0 {
		return bytes[:len(bytes)-int(padding)], nil
	}

	return bytes, nil
}

// DecodeStrict decodes a string into binary data without using any padding
func DecodeStrict(str string) ([]byte, error) {
	if !IsStrictBfh(str) {
		return nil, errors.New(errMsgStrictInvalid)
	}
	// dashes are not needed, they only help readability
	str = strings.Replace(str, "-", "", -1)

	return decode(str)
}

func decode(str string) ([]byte, error) {
	// string length -> byte length:
	// - len(str)-1 as base since first byte represents the padding
	// - *5/8 as 1 byte represents 5 bits and 1 byte is 8 bits of course
	// - -padding to get the original length
	bytes := make([]byte, len(str)*5/8)

	for i := 0; i < len(str); i++ {
		charValue, ok := digitMap[str[i:i+1]]
		if !ok {
			return nil, fmt.Errorf(errMsgContainsInvalidCharacter, str[i:i+1])
		}

		byteIndex := i * 5 / 8

		firstByte, secondByte := splitByte(charValue, i)

		bytes[byteIndex] |= firstByte

		if secondByte > 0 && len(bytes) > byteIndex+1 {
			bytes[byteIndex+1] |= secondByte
		}
	}

	return bytes, nil
}

func splitByte(charValue uint8, stringIndex int) (byte, byte) {
	mod := uint8((stringIndex * 5) % 8)

	if mod < 4 {
		return charValue << (3 - mod), 0
	}

	return charValue >> (mod - 3), charValue & decodeMasks[mod-4] << (11 - mod)
}

// IsWellFormattedBfh returns true if the string is a well-formatted string
func IsWellFormattedBfh(str string) bool {
	fixedStr := str + "-"

	return standardBfhRegex.MatchString(fixedStr)
}

// IsAcceptableBfh returns true if bfh can accept it for decoding
func IsAcceptableBfh(str string) bool {
	fixedStr := strings.Replace(str, "-", "", -1)

	return acceptableBfhRegex.MatchString(fixedStr)
}

// IsStrictBfh returns true if the string is strict-compatible
func IsStrictBfh(str string) bool {
	fixedStr := str + "-"

	return strictBfhRegex.MatchString(fixedStr)
}
