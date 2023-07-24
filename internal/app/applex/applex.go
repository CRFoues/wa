// 版权 @2023 凹语言 作者。保留所有权利。

package applex

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/scanner"
	"wa-lang.org/wa/internal/token"
)

var CmdLex = &cli.Command{
	Hidden: true,
	Name:   "lex",
	Usage:  "lex Wa source code and print token list",
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "no input file")
			os.Exit(1)
		}

		err := Lex(c.Args().First())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}

func Lex(filename string) error {
	src, err := readSource(filename, nil)
	if err != nil {
		return err
	}

	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile(filename, fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
	}

	return nil
}

func readSource(filename string, src interface{}) ([]byte, error) {
	if src != nil {
		switch s := src.(type) {
		case string:
			return []byte(s), nil
		case []byte:
			return s, nil
		case *bytes.Buffer:
			if s != nil {
				return s.Bytes(), nil
			}
		case io.Reader:
			d, err := io.ReadAll(s)
			return d, err
		}
		return nil, errors.New("invalid source")
	}

	d, err := os.ReadFile(filename)
	return d, err
}
