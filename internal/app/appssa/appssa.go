// 版权 @2023 凹语言 作者。保留所有权利。

package appssa

import (
	"fmt"
	"os"
	"sort"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/ssa"
)

var CmdSsa = &cli.Command{
	Hidden: true,
	Name:   "ssa",
	Usage:  "print Wa ssa code",
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "no input file")
			os.Exit(1)
		}

		opt := appbase.BuildOptions(c)
		err := SSARun(opt, c.Args().First())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}

func SSARun(opt *appbase.Option, filename string) error {
	cfg := opt.Config()
	prog, err := loader.LoadProgram(cfg, filename)
	if err != nil {
		return err
	}

	prog.SSAMainPkg.WriteTo(os.Stdout)

	var funcNames []string
	for name, x := range prog.SSAMainPkg.Members {
		if _, ok := x.(*ssa.Function); ok {
			funcNames = append(funcNames, name)
		}
	}
	sort.Strings(funcNames)
	for _, s := range funcNames {
		prog.SSAMainPkg.Func(s).WriteTo(os.Stdout)
	}

	return nil
}
