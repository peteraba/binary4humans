package bfh

import (
	"errors"
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
	standardBfhRegex   = regexp.MustCompile(standardBfhRegexString)
	acceptableBfhRegex = regexp.MustCompile(acceptableBfhRegexString)
	strictBfhRegex     = regexp.MustCompile(strictBfhRegexString)
	padding1           = []string{"0", "8", "g", "r"}
	padding2           = []string{"0", "2", "4", "6", "8", "a", "c", "e", "g", "j", "m", "p", "r", "t", "w", "y"}
	padding3           = []string{"0", "g"}
	padding4           = []string{"0", "4", "8", "c", "g", "m", "r", "w"}
	paddingMap1        map[string]string
	paddingMap2        map[string]string
	paddingMap3        map[string]string
	paddingMap4        map[string]string
)

// nolint
func getDigit(r uint8) (byte, error) {
	switch r {
	case '0':
		return 0, nil
	case '1':
		return 1, nil
	case '2':
		return 2, nil
	case '3':
		return 3, nil
	case '4':
		return 4, nil
	case '5':
		return 5, nil
	case '6':
		return 6, nil
	case '7':
		return 7, nil
	case '8':
		return 8, nil
	case '9':
		return 9, nil
	case 'a':
		return 10, nil
	case 'b':
		return 11, nil
	case 'c':
		return 12, nil
	case 'd':
		return 13, nil
	case 'e':
		return 14, nil
	case 'f':
		return 15, nil
	case 'g':
		return 16, nil
	case 'h':
		return 17, nil
	case 'j':
		return 18, nil
	case 'k':
		return 19, nil
	case 'm':
		return 20, nil
	case 'n':
		return 21, nil
	case 'p':
		return 22, nil
	case 'q':
		return 23, nil
	case 'r':
		return 24, nil
	case 's':
		return 25, nil
	case 't':
		return 26, nil
	case 'v':
		return 27, nil
	case 'w':
		return 28, nil
	case 'x':
		return 29, nil
	case 'y':
		return 30, nil
	case 'z':
		return 31, nil
		// 0123456789abcdefghjkmnpqrstvwxyz
	}

	return 0, errors.New(errMsgContainsInvalidCharacter)
}

// Encode encodes binary data into a human readable string
func Encode(b []byte) (string, error) {
	if b == nil {
		return "", errors.New(errMsgBinaryDataMustNotBeNil)
	}

	result := newNormalResult(len(b))

	result = encode(b, result, 2)

	return string(result), nil
}

// EncodeStrict encodes binary data with a length dividable by 5 into a simplified human readable string
func EncodeStrict(b []byte) (string, error) {
	if b == nil {
		return "", errors.New(errMsgBinaryDataMustNotBeNil)
	}

	if len(b)%5 != 0 {
		return "", errors.New(errMsgStrictMustBeDividableBy5)
	}

	result := newStrictEncodeResult(len(b))

	result = encode(b, result, 0)

	return string(result), nil
}

// newNormalResult will create a byte slice and fill it with dashes and zeros as required, plus setting the padding byte
func newNormalResult(byteLength int) []byte {
	b32padding := (5 - byteLength%5) % 5
	b32Length := (byteLength+b32padding)*2 + 1

	if byteLength == 0 {
		b32Length++
	}

	s1 := make([]byte, b32Length)
	for i := 0; i < b32Length; i++ {
		s1[i] = byte('0')
	}
	for i := 6; i < b32Length; i += 5 {
		s1[i] = byte('-')
	}

	s1[0] = byte(b32padding + '0')
	s1[1] = byte('-')

	return s1
}

// newStrictEncodeResult will create a byte slice and fill it with dashes and zeros as required
func newStrictEncodeResult(byteLength int) []byte {
	b32Length := (byteLength)*2 - 1

	if byteLength == 0 {
		b32Length++
	}

	s1 := make([]byte, b32Length)
	for i := 0; i < b32Length; i++ {
		s1[i] = byte('0')
	}
	for i := 4; i < b32Length; i += 5 {
		s1[i] = byte('-')
	}

	return s1
}

func encode(b, result []byte, offset int) []byte {
	var (
		readCount = 0
		maxCount  = len(b)*8/5 + 1
		f         byte
		idx       int
	)

	for maxCount > readCount {
		f = readByte(b, readCount*5)

		idx = readCount + offset + (readCount / 4)
		if idx >= len(result) {
			break
		}

		result[idx] = digits[f]

		readCount++
	}

	return result
}

func readByte(b []byte, readBits int) byte {
	m := uint8(readBits % 8)
	l := readBits / 8

	// we need to mask out bits we've read before
	if l >= len(b) {
		return 0
	}
	f := b[l] & encodeMasks[m]

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
	if len(b) > l+1 {
		s = b[l+1] >> (11 - m)
	}

	return f | s
}

// Decode decodes a human readable string into a binary data
func Decode(str string) ([]byte, error) {
	// dashes are not needed, they only help readability
	str = strings.Replace(str, "-", "", -1)

	padding, err := getDigit(str[0])
	if err != nil {
		return nil, err
	}
	if padding > 4 {
		return nil, errors.New(errMsgPaddingNotBetween0and4)
	}

	str = str[1:]

	if len(str)%8 != 0 {
		return nil, errors.New(errMsgStrictMustBeDividableBy8)
	}

	data, err := decode(str)
	if err != nil {
		return nil, err
	}

	if padding > 0 {
		return data[:len(data)-int(padding)], nil
	}

	return data, nil
}

// DecodeStrict decodes a string into binary data without using any padding
func DecodeStrict(str string) ([]byte, error) {
	// dashes are not needed, they only help readability
	str = strings.Replace(str, "-", "", -1)

	if len(str)%4 != 0 {
		return nil, errors.New(errMsgStrictInvalid)
	}

	return decode(str)
}

func decode(str string) ([]byte, error) {
	// string length -> byte length:
	// - len(str)-1 as base since first byte represents the padding
	// - *5/8 as 1 byte represents 5 bits and 1 byte is 8 bits of course
	// - -padding to get the original length
	data := make([]byte, len(str)*5/8)

	for i := 0; i < len(str); i++ {
		charValue, err := getDigit(str[i])
		if err != nil {
			return nil, err
		}

		byteIndex := i * 5 / 8

		firstByte, secondByte := splitByte(charValue, i)

		data[byteIndex] |= firstByte

		if secondByte > 0 && len(data) > byteIndex+1 {
			data[byteIndex+1] |= secondByte
		}
	}

	return data, nil
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

	if !standardBfhRegex.MatchString(fixedStr) {
		return false
	}

	fixedStr = strings.Replace(fixedStr, "-", "", -1)

	return isPaddingCorrect(fixedStr)
}

// IsAcceptableBfh returns true if bfh can accept it for decoding
func IsAcceptableBfh(str string) bool {
	fixedStr := strings.Replace(str, "-", "", -1)

	if !acceptableBfhRegex.MatchString(fixedStr) {
		return false
	}

	return isPaddingCorrect(fixedStr)
}

func isPaddingCorrect(str string) bool {
	l := len(str)

	if l == 1 && str == "0" {
		return true
	}

	if l < 9 {
		return false
	}

	switch str[0:1] {
	case "1":
		return isPaddingCorrect1(str)
	case "2":
		return isPaddingCorrect2(str)
	case "3":
		return isPaddingCorrect3(str)
	case "4":
		return isPaddingCorrect4(str)
	}

	return true
}

func isPaddingCorrect1(str string) bool {
	l := len(str)

	if str[l-1:] != "0" {
		return false
	}

	if paddingMap1 == nil {
		paddingMap1 = map[string]string{}

		for _, v := range padding1 {
			paddingMap1[v] = v
		}
	}

	_, ok := paddingMap1[str[l-2:l-1]]

	return ok
}

func isPaddingCorrect2(str string) bool {
	l := len(str)

	if str[l-2:] != "000" {
		return false
	}

	if paddingMap2 == nil {
		paddingMap2 = map[string]string{}

		for _, v := range padding2 {
			paddingMap2[v] = v
		}
	}

	_, ok := paddingMap2[str[l-4:l-3]]

	return ok
}

func isPaddingCorrect3(str string) bool {
	l := len(str)

	if str[l-4:] != "0000" {
		return false
	}

	if paddingMap3 == nil {
		paddingMap3 = map[string]string{}

		for _, v := range padding3 {
			paddingMap3[v] = v
		}
	}

	_, ok := paddingMap3[str[l-5:l-4]]

	return ok
}

func isPaddingCorrect4(str string) bool {
	l := len(str)

	if str[l-6:] != "000000" {
		return false
	}

	if paddingMap4 == nil {
		paddingMap4 = map[string]string{}

		for _, v := range padding4 {
			paddingMap4[v] = v
		}
	}

	_, ok := paddingMap4[str[l-7:l-6]]

	return ok
}

// IsStrictBfh returns true if the string is strict-compatible
func IsStrictBfh(str string) bool {
	fixedStr := str + "-"

	return strictBfhRegex.MatchString(fixedStr)
}
