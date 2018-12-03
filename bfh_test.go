package bfh

import (
	"crypto/rand"
	"encoding/base32"
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
				Name:           "0xff without padding",
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
				Name:           "6 bytes long 0xff",
				Bytes:          []byte{255, 255, 255, 255, 255, 255},
				ExpectedResult: "4-zzzz-zzzz-zw00-0000",
			},
			{
				Name:           "4 bytes long 0xff",
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

	t.Run("fail on nil", func(t *testing.T) {
		_, err := Encode(nil)

		assert.Error(t, err)
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

	t.Run("fail on wrong length", func(t *testing.T) {
		b := make([]byte, 14)

		_, err := rand.Read(b)
		require.NoError(t, err)

		_, err = EncodeStrict(b)

		assert.Error(t, err)
	})

	t.Run("fail on nil", func(t *testing.T) {
		_, err := EncodeStrict(nil)

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
				Name:           "6 bytes long 0xff",
				String:         "4-zzzz-zzzz-zw00-0000",
				ExpectedResult: []byte{255, 255, 255, 255, 255, 255},
			},
			{
				Name:           "4 bytes long 0xff",
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
				Name:           "6 bytes long 0xff",
				String:         "4-zzzz-zzzz-zw00-0000",
				ExpectedResult: true,
			},
			{
				Name:           "4 bytes long 0xff",
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
			{
				Name:           "fishy ending",
				String:         "1-zwga-e07x-2400-00a0",
				ExpectedResult: false,
			},
			{
				Name:           "fishy ending #2",
				String:         "4-zwga-e07x-24a0-0000",
				ExpectedResult: false,
			},
			{
				Name:           "random failed before #1",
				String:         "3-1sfc-xrn8-dqek-dqxt-9he2-f3rf-5nj8-29ya-zyj1-9wfj-fppx-maa6-x75g-5jnb-1azg-0000",
				ExpectedResult: true,
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

				assert.True(t, isValid, fmt.Sprintf("Failing value: %s", encoded))
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
				Name:           "6 bytes long 0xff",
				String:         "4-zzzz-zzzz-zw00-0000",
				ExpectedResult: true,
			},
			{
				Name:           "4 bytes long 0xff",
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

				assert.True(t, isValid, fmt.Sprintf("Failing value: %s", encoded))
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
				ExpectedResult: true,
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
				Name:           "6 bytes long 0xff",
				String:         "zzzz-zzzz-zw00-0000",
				ExpectedResult: true,
			},
			{
				Name:           "4 bytes long 0xff",
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

				assert.True(t, isValid, fmt.Sprintf("Failing value: %s", encoded))
			})
		}
	})
}

var (
	data = []byte{
		133, 237, 158, 181, 220, 197, 196, 29, 137, 199,
		57, 60, 172, 6, 70, 72, 184, 186, 18, 169,
		127, 57, 94, 20, 222, 115, 22, 237, 37, 66,
		98, 242, 148, 29, 200, 192, 166, 73, 26, 153,

		142, 217, 24, 147, 47, 184, 130, 55, 122, 177,
		195, 40, 6, 102, 228, 221, 252, 97, 64, 118,
		208, 235, 117, 219, 120, 86, 119, 121, 41, 164,
		249, 39, 91, 72, 133, 143, 157, 132, 0, 99,

		109, 183, 234, 164, 20, 6, 63, 156, 15, 52,
		117, 213, 115, 208, 106, 39, 18, 248, 157, 181,
		93, 65, 149, 25, 140, 110, 139, 189, 105, 64,
		196, 62, 6, 63, 194, 122, 168, 63, 123, 87,

		207, 204, 2, 143, 143, 156, 23, 197, 216, 75,
		231, 61, 104, 203, 24, 199, 71, 13, 243, 97,
		62, 162, 99, 218, 240, 114, 207, 58, 198, 216,
		213, 58, 96, 12, 67, 78, 109, 13, 101, 222,

		216, 83, 37, 19, 77, 119, 214, 95, 158, 127,
		125, 67, 38, 106, 254, 9, 38, 108, 186, 125,
		59, 187, 96, 203, 54, 107, 197, 250, 135, 90,
		175, 159, 41, 78, 65, 71, 190, 244, 154, 198,

		210, 107, 17, 6, 91, 107, 68, 195, 18, 222,
		212, 220, 4, 55, 196, 157, 193, 58, 59, 182,
		153, 101, 7, 13, 233, 124, 31, 11, 161, 191,
		190, 236, 128, 176, 165, 83, 93, 170, 195, 236,
	}
	encodedData240  = "gqps-xdew-rq21-v2e7-74ya-r1j6-92wb-m4n9-fwwn-w56y-ecbe-t9a2-cbs9-87e8-r2k4-j6ms-hvch-h4sf-q213-eynh-rcm0-csq4-vqy6-2g3p-t3nq-bpvr-asvq-jad4-z4kn-pj45-hyer-8033-dpvy-n90m-0rzs-r3sm-eqaq-7m3a-4w9f-h7dn-bn0s-a6cc-dt5v-tta0-rgz0-cfy2-fam3-yytq-sz60-53wf-kgbw-bp2b-wwyp-hjrr-rx3g-vwv1-7th6-7pqg-eb7k-nhpr-tmx6-0323-9spg-tsey-v19j-a4td-ezb5-z7kz-fn1j-ctqy-14k6-sekx-7exp-1jsp-df2z-n1tt-nyfj-jkj1-8yzf-96p6-t9nh-21jv-dd2c-64py-tke0-8dy4-kq0k-mexp-k5jg-e3f9-fgfg-q8dz-qvp8-1c55-adet-ngzc"
	encodedData238  = "2-gqps-xdew-rq21-v2e7-74ya-r1j6-92wb-m4n9-fwwn-w56y-ecbe-t9a2-cbs9-87e8-r2k4-j6ms-hvch-h4sf-q213-eynh-rcm0-csq4-vqy6-2g3p-t3nq-bpvr-asvq-jad4-z4kn-pj45-hyer-8033-dpvy-n90m-0rzs-r3sm-eqaq-7m3a-4w9f-h7dn-bn0s-a6cc-dt5v-tta0-rgz0-cfy2-fam3-yytq-sz60-53wf-kgbw-bp2b-wwyp-hjrr-rx3g-vwv1-7th6-7pqg-eb7k-nhpr-tmx6-0323-9spg-tsey-v19j-a4td-ezb5-z7kz-fn1j-ctqy-14k6-sekx-7exp-1jsp-df2z-n1tt-nyfj-jkj1-8yzf-96p6-t9nh-21jv-dd2c-64py-tke0-8dy4-kq0k-mexp-k5jg-e3f9-fgfg-q8dz-qvp8-1c55-adet-m000"
	encodedData25   = "gqps-xdew-rq21-v2e7-74ya-r1j6-92wb-m4n9-fwwn-w56y"
	encodedData23   = "2-gqps-xdew-rq21-v2e7-74ya-r1j6-92wb-m4n9-fwwn-w000"
	decodedResult   []byte
	encodeResult    string
	validatedResult bool
)

func Benchmark_Encode_23(b *testing.B) {
	var str string

	for n := 0; n < b.N; n++ {
		str, _ = Encode(data[:23])
	}

	encodeResult = str
}

func Benchmark_EncodeStrict_25(b *testing.B) {
	var str string

	for n := 0; n < b.N; n++ {
		str, _ = EncodeStrict(data[:25])
	}

	encodeResult = str
}

func Benchmark_Encode_238(b *testing.B) {
	var str string

	for n := 0; n < b.N; n++ {
		str, _ = Encode(data[:238])
	}

	encodeResult = str
}

func Benchmark_EncodeStrict_240(b *testing.B) {
	var str string

	for n := 0; n < b.N; n++ {
		str, _ = EncodeStrict(data)
	}

	encodeResult = str
}

func Benchmark_Decode_23(b *testing.B) {
	var decodedData []byte

	for n := 0; n < b.N; n++ {
		decodedData, _ = Decode(encodedData23)
	}

	decodedResult = decodedData
}

func Benchmark_DecodeStrict_25(b *testing.B) {
	var decodedData []byte

	for n := 0; n < b.N; n++ {
		decodedData, _ = DecodeStrict(encodedData25)
	}

	decodedResult = decodedData
}

func Benchmark_Decode_238(b *testing.B) {
	var decodedData []byte

	for n := 0; n < b.N; n++ {
		decodedData, _ = Decode(encodedData238)
	}

	decodedResult = decodedData
}

func Benchmark_DecodeStrict_240(b *testing.B) {
	var decodedData []byte

	for n := 0; n < b.N; n++ {
		decodedData, _ = DecodeStrict(encodedData240)
	}

	decodedResult = decodedData
}

func Benchmark_IsWellFormattedBfh_238(b *testing.B) {
	var validatedData bool

	for n := 0; n < b.N; n++ {
		validatedData = IsWellFormattedBfh(encodedData238)
	}

	validatedResult = validatedData
}

func Benchmark_IsAcceptableBfh_238(b *testing.B) {
	var validatedData bool

	for n := 0; n < b.N; n++ {
		validatedData = IsAcceptableBfh(encodedData238)
	}

	validatedResult = validatedData
}

func Benchmark_IsStrictBfh_240(b *testing.B) {
	var validatedData bool

	for n := 0; n < b.N; n++ {
		validatedData = IsStrictBfh(encodedData240)
	}

	validatedResult = validatedData
}

func Benchmark_Base32Encode_23(b *testing.B) {
	var str string

	for n := 0; n < b.N; n++ {
		str = base32.StdEncoding.EncodeToString(data[:23])
	}

	encodeResult = str
}

func Benchmark_Base32Encode_20(b *testing.B) {
	var str string

	for n := 0; n < b.N; n++ {
		str = base32.StdEncoding.EncodeToString(data[:20])
	}

	encodeResult = str
}

func Benchmark_Base32Encode_238(b *testing.B) {
	var str string

	for n := 0; n < b.N; n++ {
		str = base32.StdEncoding.EncodeToString(data[:238])
	}

	encodeResult = str
}

func Benchmark_Base32Encode_240(b *testing.B) {
	var str string

	for n := 0; n < b.N; n++ {
		str = base32.StdEncoding.EncodeToString(data[:240])
	}

	encodeResult = str
}

func Benchmark_Base32Decode_23(b *testing.B) {
	var decodedData []byte
	var str = base32.StdEncoding.EncodeToString(data[:23])

	for n := 0; n < b.N; n++ {
		decodedData, _ = base32.StdEncoding.DecodeString(str)
	}

	decodedResult = decodedData
}

func Benchmark_Base32Decode_20(b *testing.B) {
	var decodedData []byte
	var str = base32.StdEncoding.EncodeToString(data[:20])

	for n := 0; n < b.N; n++ {
		decodedData, _ = base32.StdEncoding.DecodeString(str)
	}

	decodedResult = decodedData
}

func Benchmark_Base32Decode_238(b *testing.B) {
	var decodedData []byte
	var str = base32.StdEncoding.EncodeToString(data[:238])

	for n := 0; n < b.N; n++ {
		decodedData, _ = base32.StdEncoding.DecodeString(str)
	}

	decodedResult = decodedData
}

func Benchmark_Base32Decode_240(b *testing.B) {
	var decodedData []byte
	var str = base32.StdEncoding.EncodeToString(data[:240])

	for n := 0; n < b.N; n++ {
		decodedData, _ = base32.StdEncoding.DecodeString(str)
	}

	decodedResult = decodedData
}
