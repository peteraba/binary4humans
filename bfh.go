package bytes4humans

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const digits = "0123456789abcdefghjkmnpqrstvwxyz"

var (
	encodeMasks = []uint8{0xff, 0x7f, 0x3f, 0x1f, 0x0f, 0x07, 0x03, 0x01}
	decodeMasks = []uint8{0x1, 0x3, 0x7, 0xf}
	digitMap    map[string]uint8
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
		readBits   = 0
		err        error
	)

	if len(b) == 0 {
		return "", nil
	}

	b, err = padBytes(&strBuilder, b)
	maxBits := len(b) * 8
	if err != nil {
		return "", err
	}

	for maxBits > readBits {
		f := readByte(b, readBits)

		_, err = fmt.Fprintf(&strBuilder, "%s", digits[f:f+1])
		if err != nil {
			return "", err
		}

		readBits += 5

		if readBits%20 != 0 || readBits >= maxBits {
			continue
		}

		_, err = fmt.Fprint(&strBuilder, "-")
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
	if str == "" {
		return []byte{}, nil
	}

	// dashes are not needed, they only help readability
	str = strings.Replace(str, "-", "", -1)

	padding, ok := digitMap[str[0:1]]
	if !ok || padding > 4 {
		return nil, errors.New("non empty string must start with 0, 1, 2, 3 or 4")
	}

	str = str[1:]

	// string length -> byte length:
	// - len(str)-1 as base since first byte represents the padding
	// - *5/8 as 1 byte represents 5 bits and 1 byte is 8 bits of course
	// - -padding to get the original length
	bytes := make([]byte, len(str)*5/8)

	for i := 0; i < len(str); i++ {
		charValue, ok := digitMap[str[i:i+1]]
		if !ok {
			return nil, errors.New("string is invalid: " + str)
		}

		byteIndex := i * 5 / 8

		firstByte, secondByte := splitByte(charValue, i)

		bytes[byteIndex] |= firstByte

		if secondByte > 0 && len(bytes) > byteIndex+1 {
			bytes[byteIndex+1] |= secondByte
		}
	}

	if padding > 0 {
		return bytes[:len(bytes)-int(padding)], nil
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

const (
	// note that valid encoded strings will not end in a hyphen, it needs to be added when validating
	standardBfhRegexString   = "^[0-4]\\-([0123456789abcdefghjkmnpqrstvwxyz]{4}\\-)*$"
	acceptableBfhRegexString = "^[0-4]([0123456789abcdefghjkmnpqrstvwxyz]{4})*$"
)

var (
	standardBfhRegex   = regexp.MustCompile(standardBfhRegexString)
	acceptableBfhRegex = regexp.MustCompile(acceptableBfhRegexString)
)

// IsWellFormattedBfh returns true if the string is a well-formatted string
func IsWellFormattedBfh(str string) bool {
	if str == "" {
		return true
	}

	fixedStr := str + "-"

	return standardBfhRegex.MatchString(fixedStr)
}

// IsAcceptableBfh returns true if bfh can accept it for decoding
func IsAcceptableBfh(str string) bool {
	if str == "" {
		return true
	}

	fixedStr := strings.Replace(str, "-", "", -1)

	return acceptableBfhRegex.MatchString(fixedStr)
}
