package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/speps/go-hashids"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"usage:\n"+
				"\t%s [options] <intlist> [<intlist>...]\n"+
				"\t%s [options] -d <hashid> [<hashid>...]\n\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	var decode bool
	var params hashids.HashIDData
	var separator string
	var useHex bool
	flag.StringVar(&params.Salt, `salt`, "", `salt`)
	flag.StringVar(&params.Alphabet, `alphabet`, hashids.DefaultAlphabet, `minimum 16 characters`)
	flag.IntVar(&params.MinLength, `min`, 0, `minimum length (for encoding)`)
	flag.BoolVar(&decode, `d`, false, `decode (instead of encoding)`)
	flag.StringVar(&separator, `sep`, ",", `separator for integers`)
	flag.BoolVar(&useHex, `hex`, false, `use hex for encoding/decoding (strings)`)
	flag.Parse()

	codec, err := hashids.NewWithData(&params)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(decode, useHex)

	args := os.Args[len(os.Args)-flag.NArg():]
	if decode {
		for _, arg := range args {
			var output []byte
			var err error
			if useHex {
				result, err := codec.DecodeHex(arg)
				if err == nil {
					output, err = hex.DecodeString(result)
				}
			} else {
				result, err := codec.DecodeInt64WithError(arg)
				if err == nil {
					for _, x := range result {
						if len(output) != 0 {
							output = append(output, separator...)
						}
						output = strconv.AppendInt(output, x, 10)
					}
				}
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", arg, err)
			} else {
				fmt.Printf("%s: %s\n", arg, output)
			}
		}
	} else {

	ARGS:
		for _, arg := range args {
			var err error
			var result string
			if useHex {
				result, err = codec.EncodeHex(hex.EncodeToString([]byte(arg)))
			} else {
				spl := strings.Split(arg, separator)
				ints := make([]int64, len(spl))
				for i, s := range spl {
					ints[i], err = strconv.ParseInt(s, 0, 64)
					if err != nil {
						fmt.Fprintf(os.Stderr, "%s: %s\n", arg, err)
						continue ARGS
					}
				}
				result, err = codec.EncodeInt64(ints)
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", arg, err)
			} else {
				fmt.Printf("%s: %v\n", arg, result)
			}
		}
	}
}
