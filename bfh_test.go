package bytes4humans

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Encode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tests := []struct {
			Name           string
			Bytes          []byte
			ExpectedResult string
		}{
			{
				Name:           "empty",
				Bytes:          []byte{},
				ExpectedResult: "",
			},
			{
				Name:           "0x7e",
				Bytes:          []byte{126},
				ExpectedResult: "4-fr00-0000",
			},
			{
				Name:           "0xff",
				Bytes:          []byte{255},
				ExpectedResult: "4-zw00-0000",
			},
			{
				Name:           "0xff",
				Bytes:          []byte{255, 0, 0, 0, 0},
				ExpectedResult: "0-zw00-0000",
			},
			{
				Name:           "zeros",
				Bytes:          []byte{0, 0, 0, 0, 0},
				ExpectedResult: "0-0000-0000",
			},
			{
				Name:           "mod5 slice length of 0xff",
				Bytes:          []byte{255, 255, 255, 255, 255},
				ExpectedResult: "0-zzzz-zzzz",
			},
			{
				Name:           "non-mod5 slice length of 0xff",
				Bytes:          []byte{255, 255, 255, 255, 255, 255},
				ExpectedResult: "4-zzzz-zzzz-zw00-0000",
			},
			{
				Name:           "non-mod5 slice length of 0xff",
				Bytes:          []byte{255, 255, 255, 255},
				ExpectedResult: "1-zzzz-zzr0",
			},
			{
				Name:           "somewhat random numbers",
				Bytes:          []byte{255, 32, 167, 0, 253, 17},
				ExpectedResult: "4-zwga-e07x-2400-0000",
			},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				actualResult, err := Encode(tt.Bytes)

				assert.NoError(t, err)
				assert.Equal(t, tt.ExpectedResult, actualResult)
			})
		}
	})

	t.Run("random success", func(t *testing.T) {
		tests := []struct {
			Name        string
			Length      int
			PacketCount int
		}{
			{
				Name:        "37",
				Length:      37,
				PacketCount: 16,
			},
			{
				Name:        "69",
				Length:      69,
				PacketCount: 28,
			},
			{
				Name:        "120",
				Length:      120,
				PacketCount: 48,
			},
			{
				Name:        "141",
				Length:      141,
				PacketCount: 58,
			},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				b := make([]byte, tt.Length)

				_, err := rand.Read(b)
				require.NoError(t, err)

				str, err := Encode(b)

				assert.NoError(t, err)
				assert.Regexp(t, fmt.Sprintf("^[0-4]-([a-z0-9]{4}\\-){%d}$", tt.PacketCount), str+"-")
			})
		}
	})
}

func Test_Decode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tests := []struct {
			Name           string
			String         string
			ExpectedResult []byte
		}{
			{
				Name:           "empty",
				String:         "",
				ExpectedResult: []byte{},
			},
			{
				Name:           "alternative empty",
				String:         "0-",
				ExpectedResult: []byte{},
			},
			{
				Name:           "0x7e",
				String:         "4-fr00-0000",
				ExpectedResult: []byte{126},
			},
			{
				Name:           "0xff",
				String:         "4-zw00-0000",
				ExpectedResult: []byte{255},
			},
			{
				Name:           "0xff",
				String:         "4-zw00-0000",
				ExpectedResult: []byte{255},
			},
			{
				Name:           "zeros",
				String:         "0-0000-0000",
				ExpectedResult: []byte{0, 0, 0, 0, 0},
			},
			{
				Name:           "mod5 slice length of 0xff",
				String:         "0-zzzz-zzzz",
				ExpectedResult: []byte{255, 255, 255, 255, 255},
			},
			{
				Name:           "non-mod5 slice length of 0xff",
				String:         "4-zzzz-zzzz-zw00-0000",
				ExpectedResult: []byte{255, 255, 255, 255, 255, 255},
			},
			{
				Name:           "non-mod5 slice length of 0xff",
				String:         "1-zzzz-zzr0",
				ExpectedResult: []byte{255, 255, 255, 255},
			},
			{
				Name:           "somewhat random numbers",
				String:         "4-zwga-e07x-2400-0000",
				ExpectedResult: []byte{255, 32, 167, 0, 253, 17},
			},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				actualResult, err := Decode(tt.String)

				require.NoError(t, err)
				assert.Equal(t, tt.ExpectedResult, actualResult)
			})
		}
	})

	t.Run("random success", func(t *testing.T) {
		tests := []struct {
			Name   string
			Length int
		}{
			{
				Name:   "37",
				Length: 37,
			},
			{
				Name:   "69",
				Length: 69,
			},
			{
				Name:   "120",
				Length: 120,
			},
			{
				Name:   "141",
				Length: 141,
			},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				b := make([]byte, tt.Length)

				_, err := rand.Read(b)
				require.NoError(t, err)

				encoded, err := Encode(b)
				require.NoError(t, err)

				decoded, err := Decode(encoded)
				require.NoError(t, err)

				assert.Equal(t, b, decoded)
			})
		}
	})
}

func Test_IsWellFormatted(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tests := []struct {
			Name           string
			String         string
			ExpectedResult bool
		}{
			{
				Name:           "empty",
				String:         "",
				ExpectedResult: true,
			},
			// alternative empty is not well formatted but acceptable
			{
				Name:           "alternative empty",
				String:         "0-",
				ExpectedResult: false,
			},
			{
				Name:           "0x7e",
				String:         "4-fr00-0000",
				ExpectedResult: true,
			},
			{
				Name:           "0xff",
				String:         "4-zw00-0000",
				ExpectedResult: true,
			},
			{
				Name:           "0xff",
				String:         "4-zw00-0000",
				ExpectedResult: true,
			},
			{
				Name:           "zeros",
				String:         "0-0000-0000",
				ExpectedResult: true,
			},
			{
				Name:           "mod5 slice length of 0xff",
				String:         "0-zzzz-zzzz",
				ExpectedResult: true,
			},
			{
				Name:           "non-mod5 slice length of 0xff",
				String:         "4-zzzz-zzzz-zw00-0000",
				ExpectedResult: true,
			},
			{
				Name:           "non-mod5 slice length of 0xff",
				String:         "1-zzzz-zzr0",
				ExpectedResult: true,
			},
			{
				Name:           "somewhat random numbers",
				String:         "4-zwga-e07x-2400-0000",
				ExpectedResult: true,
			},
			{
				Name:           "invalid empty",
				String:         "a-",
				ExpectedResult: false,
			},
			{
				Name:           "wrong length",
				String:         "4-zwg",
				ExpectedResult: false,
			},
			{
				Name:           "extra dash at the end",
				String:         "4-zwga-e07x-2400-0000-",
				ExpectedResult: false,
			},
			{
				Name:           "different dash positions",
				String:         "4-zwg-ae0-7x2-400-00-00",
				ExpectedResult: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				actualResult := IsWellFormattedBfh(tt.String)

				assert.Equal(t, tt.ExpectedResult, actualResult)
			})
		}
	})

	t.Run("random success", func(t *testing.T) {
		tests := []struct {
			Name   string
			Length int
		}{
			{
				Name:   "37",
				Length: 37,
			},
			{
				Name:   "69",
				Length: 69,
			},
			{
				Name:   "120",
				Length: 120,
			},
			{
				Name:   "141",
				Length: 141,
			},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				b := make([]byte, tt.Length)

				_, err := rand.Read(b)
				require.NoError(t, err)

				encoded, err := Encode(b)
				require.NoError(t, err)

				isValid := IsWellFormattedBfh(encoded)

				assert.True(t, isValid)
			})
		}
	})
}

func Test_IsAcceptable(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tests := []struct {
			Name           string
			String         string
			ExpectedResult bool
		}{
			{
				Name:           "empty",
				String:         "",
				ExpectedResult: true,
			},
			{
				Name:           "alternative empty",
				String:         "0-",
				ExpectedResult: true,
			},
			{
				Name:           "0x7e",
				String:         "4-fr00-0000",
				ExpectedResult: true,
			},
			{
				Name:           "0xff",
				String:         "4-zw00-0000",
				ExpectedResult: true,
			},
			{
				Name:           "0xff",
				String:         "4-zw00-0000",
				ExpectedResult: true,
			},
			{
				Name:           "zeros",
				String:         "0-0000-0000",
				ExpectedResult: true,
			},
			{
				Name:           "mod5 slice length of 0xff",
				String:         "0-zzzz-zzzz",
				ExpectedResult: true,
			},
			{
				Name:           "non-mod5 slice length of 0xff",
				String:         "4-zzzz-zzzz-zw00-0000",
				ExpectedResult: true,
			},
			{
				Name:           "non-mod5 slice length of 0xff",
				String:         "1-zzzz-zzr0",
				ExpectedResult: true,
			},
			{
				Name:           "somewhat random numbers",
				String:         "4-zwga-e07x-2400-0000",
				ExpectedResult: true,
			},
			{
				Name:           "invalid empty",
				String:         "a-",
				ExpectedResult: false,
			},
			{
				Name:           "wrong length",
				String:         "4-zwg",
				ExpectedResult: false,
			},
			{
				Name:           "extra dash at the end",
				String:         "4-zwga-e07x-2400-0000-",
				ExpectedResult: true,
			},
			{
				Name:           "different dash positions",
				String:         "4-zwg-ae0-7x2-400-00-00",
				ExpectedResult: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				actualResult := IsAcceptableBfh(tt.String)

				assert.Equal(t, tt.ExpectedResult, actualResult)
			})
		}
	})

	t.Run("random success", func(t *testing.T) {
		tests := []struct {
			Name   string
			Length int
		}{
			{
				Name:   "37",
				Length: 37,
			},
			{
				Name:   "69",
				Length: 69,
			},
			{
				Name:   "120",
				Length: 120,
			},
			{
				Name:   "141",
				Length: 141,
			},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				b := make([]byte, tt.Length)

				_, err := rand.Read(b)
				require.NoError(t, err)

				encoded, err := Encode(b)
				require.NoError(t, err)

				isValid := IsAcceptableBfh(encoded)

				assert.True(t, isValid)
			})
		}
	})
}
