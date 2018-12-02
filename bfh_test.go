package bfh

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
				ExpectedResult: "0-",
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

func Test_EncodeStrict(t *testing.T) {
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
				Name:           "0xff",
				Bytes:          []byte{255, 0, 0, 0, 0},
				ExpectedResult: "zw00-0000",
			},
			{
				Name:           "zeros",
				Bytes:          []byte{0, 0, 0, 0, 0},
				ExpectedResult: "0000-0000",
			},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				actualResult, err := EncodeStrict(tt.Bytes)

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
				Name:        "15",
				Length:      15,
				PacketCount: 6,
			},
			{
				Name:        "25",
				Length:      25,
				PacketCount: 10,
			},
			{
				Name:        "80",
				Length:      80,
				PacketCount: 32,
			},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				b := make([]byte, tt.Length)

				_, err := rand.Read(b)
				require.NoError(t, err)

				str, err := EncodeStrict(b)

				assert.NoError(t, err)
				assert.Regexp(t, fmt.Sprintf("^([a-z0-9]{4}\\-){%d}$", tt.PacketCount), str+"-")
			})
		}
	})

	t.Run("fail", func(t *testing.T) {
		b := make([]byte, 14)

		_, err := rand.Read(b)
		require.NoError(t, err)

		_, err = EncodeStrict(b)

		assert.Error(t, err)
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

	t.Run("fail wrong padding", func(t *testing.T) {
		_, err := DecodeStrict("6-zwga-e07x-2400-0000")

		assert.Error(t, err)
	})

	t.Run("fail wrong length", func(t *testing.T) {
		_, err := DecodeStrict("0-zwga-e0")

		assert.Error(t, err)
	})
}

func Test_IsWellFormatted(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tests := []struct {
			Name           string
			String         string
			ExpectedResult bool
		}{
			// alternative empty is not well formatted but acceptable
			{
				Name:           "empty",
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
				Name:           "invalid empty 2",
				String:         "",
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
				Name:           "very empty",
				String:         "",
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

func Test_IsStrictBfh(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tests := []struct {
			Name           string
			String         string
			ExpectedResult bool
		}{
			{
				Name:           "empty",
				String:         "",
				ExpectedResult: false,
			},
			{
				Name:           "0x7e",
				String:         "fr00-0000",
				ExpectedResult: true,
			},
			{
				Name:           "0xff",
				String:         "zw00-0000",
				ExpectedResult: true,
			},
			{
				Name:           "0xff",
				String:         "zw00-0000",
				ExpectedResult: true,
			},
			{
				Name:           "zeros",
				String:         "0000-0000",
				ExpectedResult: true,
			},
			{
				Name:           "mod5 slice length of 0xff",
				String:         "zzzz-zzzz",
				ExpectedResult: true,
			},
			{
				Name:           "non-mod5 slice length of 0xff",
				String:         "zzzz-zzzz-zw00-0000",
				ExpectedResult: true,
			},
			{
				Name:           "non-mod5 slice length of 0xff",
				String:         "zzzz-zzr0",
				ExpectedResult: true,
			},
			{
				Name:           "somewhat random numbers",
				String:         "zwga-e07x-2400-0000",
				ExpectedResult: true,
			},
			{
				Name:           "wrong length",
				String:         "zwg",
				ExpectedResult: false,
			},
			{
				Name:           "extra dash at the end",
				String:         "zwga-e07x-2400-0000-",
				ExpectedResult: false,
			},
			{
				Name:           "different dash positions",
				String:         "zwg-ae0-7x2-400-00-00",
				ExpectedResult: false,
			},
			{
				Name:           "different dash positions",
				String:         "zwgae07x24000000",
				ExpectedResult: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				actualResult := IsStrictBfh(tt.String)

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
				Name:   "15",
				Length: 15,
			},
			{
				Name:   "25",
				Length: 25,
			},
			{
				Name:   "80",
				Length: 80,
			},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				b := make([]byte, tt.Length)

				_, err := rand.Read(b)
				require.NoError(t, err)

				encoded, err := EncodeStrict(b)
				require.NoError(t, err)

				isValid := IsStrictBfh(encoded)

				assert.True(t, isValid)
			})
		}
	})
}
