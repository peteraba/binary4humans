BFH - Binary 4 Humans
=====================

[![Travis](https://img.shields.io/travis/peteraba/binary4humans.svg?style=flat-square&&branch=master)](https://travis-ci.com/peteraba/binary4humans)
[![GoReportCard](https://goreportcard.com/badge/github.com/peteraba/binary4humans?style=flat-square)](https://goreportcard.com/report/github.com/peteraba/binary4humans)
[![Releases](https://img.shields.io/github/release/peteraba/binary4humans.svg?style=flat-square)](https://github.com/peteraba/binary4humans/releases)
[![LICENSE](https://img.shields.io/github/license/peteraba/binary4humans.svg?style=flat-square)](https://github.com/peteraba/binary4humans/blob/master/LICENSE)

This library aims to help displaying binary data to human users of systems, primary goal was displaying user tokens. In purpose it is very similar to the standard [base32](https://golang.org/pkg/encoding/base32/) library, in some details it is inspired by [Crockford's Base32 Encoding](https://www.crockford.com/wrmg/base32.html) definition.


Definition details
------------------

`bfh` uses 32 characters to encode binary data into a string representation. The symbols used are the same as defined by [Crockford's Base32 Encoding](https://www.crockford.com/wrmg/base32.html), except that `bfh` uses lower case characters and there are no check symbols in the current implementation.

Since the encoded characters will only hold 5 bits of data, `bfh` will create packets of 8 characters, 40 bits, which will be displayed in two 4-character long subpackets.
If the length of the binary data is not dividable by 5, then it will be padded by zeros so that it will become dividable. The number of zero bytes added to the data will become the first character of the encoded string.

Dashes are generated automatically during encoding, but ignored completely during decoding.

There's is a special case for empty bytes, which would normally translate to `0-`, but was decided to be represented by an empty string instead. `0-` will still decode to an empty byte array.

Example
-------

*Binary data:* `[167, 13]` = `[10100111, 00001101]`
*Encoded string:* "3-mw6g-0000"

Note that the padding here is important, because without that, we'd decode the encoded string to `[167, 13, 0, 0, 0]`

Installation
------------

```
go get github.com/peteraba/binary4humans
```

Usage
-----

```go
package main

import (
	"fmt"
	
	bfh "github.com/peteraba/binary4humans"
)

func main() {
    bytes1 := []byte{255, 32, 167, 0, 253, 17}
    encoded, err := bfh.Encode(bytes1)
    if err != nil {
        // handle error...
    }
    fmt.Println(encoded)
    // 4-zwga-e07x-2400-0000
    
    decoded, err := bfh.Decode(encoded)
    if err != nil {
        // handle error...
    }
    fmt.Printf("%v\n", decoded)
    // [255 32 167 0 253 17]
}
```

Extra
-----

`bfh` ships with two validators. There's a rather strict validator called `IsWellFormattedBfh` and a more relaxed `IsWellFormattedBfh`.

TODO
----

 - [ ] Improve the validators so that the padding number is checked against the last characters
