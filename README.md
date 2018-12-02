BFH - Binary 4 Humans
=====================

[![Travis](https://img.shields.io/travis/com/peteraba/binary4humans.svg?style=flat-square&&branch=master)](https://travis-ci.com/peteraba/binary4humans)
[![GoReportCard](https://goreportcard.com/badge/github.com/peteraba/binary4humans?style=flat-square)](https://goreportcard.com/report/github.com/peteraba/binary4humans)
[![Releases](https://img.shields.io/github/release/peteraba/binary4humans.svg?style=flat-square)](https://github.com/peteraba/binary4humans/releases)
[![LICENSE](https://img.shields.io/github/license/peteraba/binary4humans.svg?style=flat-square)](https://github.com/peteraba/binary4humans/blob/master/LICENSE)

This library aims to help displaying binary data to human users of systems, primary goal was displaying user tokens.
In purpose it is very similar to the standard [base32](https://golang.org/pkg/encoding/base32/) library, in some details 
it is inspired by [Crockford's Base32 Encoding](https://www.crockford.com/wrmg/base32.html) definition.

**WARNING!** Using `encoding/base32` library is currently about 50 to 100 times faster than `bfh`, so keep that in mind when making decisions! (See benchmarks below)


Definition details
------------------

`bfh` uses 32 characters to encode binary data into a string representation. The symbols used are the same as defined by
[Crockford's Base32 Encoding](https://www.crockford.com/wrmg/base32.html), except that `bfh` uses lower case characters
and there are no check symbols in the current implementation.

Since the encoded characters will only hold 5 bits of data, `bfh` will create packets of 8 characters, each encoding 40
bits of useful data and each will be displayed in two 4-character long subpackets.

If the length of the binary data is not dividable by 5 then it will be padded by zeros so that it will become
dividable. The number of zero bytes added to the data will become the first character of the encoded string.

As an alternatively you can also rely on `encodeStrict` and `decodeStrict` which can only be used with binary data of
length dividable by 5 but then the padding is no longer used.

Dashes are generated automatically during encoding, but ignored completely during decoding.

### Empty byte arrays

Normally empty byte arrays should translate to `0-`, but was decided to be represented by an empty string instead. `0-` 
will still decode to an empty byte array.

### 

Example
-------

### Example 1

*Binary data:* `[167, 13]` = `[10100111, 00001101]`  

*Encoded string:* "3-mw6g-0000"

### Example 2

*Binary data:* `[167, 13, 0, 0, 0]` = `[10100111, 00001101, 00000000, 00000000, 00000000]`  

*Encoded string:* "0-mw6g-0000"

*Strictly encoded string:* "mw6g-0000"


Note that the padding here is important, because without that, we'd decode the encoded string to `[167, 13, 0, 0, 0]`

Installation
------------

```
go get github.com/peteraba/binary4humans
```

Usage
-----

### With random length binary data

```go
package main

import (
	"fmt"
	
	bfh "github.com/peteraba/binary4humans"
)

func main() {
    binaryData := []byte{255, 32, 167, 0, 253, 17, 215, 43}
    encoded, err := bfh.Encode(binaryData)
    if err != nil {
        // handle error...
    }
    fmt.Println(encoded)
    // 2-zwga-e07x-27bj-p000
    
    decoded, err := bfh.Decode(encoded)
    if err != nil {
        // handle error...
    }
    fmt.Printf("%v\n", decoded)
    // [255 32 167 0 253 17 215 43]
}
```

### In strict mode

```go
package main

import (
	"fmt"
	
	bfh "github.com/peteraba/binary4humans"
)

func main() {
    binaryData := []byte{255, 32, 167, 0, 253, 17, 215, 43, 0, 0}
    encoded, err := bfh.EncodeStrict(binaryData)
    if err != nil {
        // handle error...
    }
    fmt.Println(encoded)
    // zwga-e07x-27bj-p000
    
    decoded, err := bfh.DecodeStrict(encoded)
    if err != nil {
        // handle error...
    }
    fmt.Printf("%v\n", decoded)
    // [255 32 167 0 253 17 215 43 0 0]
}
```

Extra
-----

`bfh` ships with altogether three validators.

 - For checking strings encoding random length binary data there are two validators:
   1. a relaxed validator called `IsAcceptableBfh` which ignores dashes
   1. and a validator called `IsWellFormattedBfh` which expects the dashes to be properly placed
 - For checking strings encoding in `strict` mode there's a validator called `IsStrictBfh` which also expects the dashes
 to be placed properly

Benchmarks
----------

As you can see the performance leaves a lot to be desired at the moment. Improvements are planned but not promised. Pull requests are welcome!

```
➤ go test -bench=. -cpu=1                  🡑 ✹   master 965ee68 "Add benchmarks"
goos: linux
goarch: amd64
pkg: github.com/peteraba/binary4humans
Benchmark_BfhEncode_23        	  300000	      4655 ns/op
Benchmark_BfhEncodeStrict_25  	  300000	      4428 ns/op
Benchmark_BfhEncode_238       	   30000	     41659 ns/op
Benchmark_BfhEncodeStrict_240 	   30000	     42284 ns/op
Benchmark_Base32Encode_23     	10000000	       117 ns/op
Benchmark_Base32Encode_20     	20000000	       101 ns/op
Benchmark_Base32Encode_238    	 2000000	       605 ns/op
Benchmark_Base32Encode_240    	 2000000	       605 ns/op
Benchmark_BfhDecode_23        	 1000000	      1004 ns/op
Benchmark_BfhDecodeStrict_25  	 1000000	      2050 ns/op
Benchmark_BfhDecode_238       	  200000	      9861 ns/op
Benchmark_BfhDecodeStrict_240 	  100000	     20419 ns/op
Benchmark_Base32Decode_23     	 5000000	       279 ns/op
Benchmark_Base32Decode_20     	10000000	       203 ns/op
Benchmark_Base32Decode_238    	 1000000	      2042 ns/op
Benchmark_Base32Decode_240    	 1000000	      2056 ns/op
PASS
ok  	github.com/peteraba/binary4humans	28.796s
```

TODO
----

 - [ ] Improve performance
