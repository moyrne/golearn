package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
)

func main() {
	const dir = "files"
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, fName := range dirs {
		filename := fName.Name()
		if len(filename) > 3 && filename[len(filename)-3:] == ".in" {
			log.Println(scanFile(dir, filename))
		}
	}
}

func scanFile(dir, filename string) error {
	src, err := ioutil.ReadFile(dir + "/" + filename)
	if err != nil {
		return errors.WithStack(err)
	}
	// 初始化 scanner
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile(filename, fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	out, err := os.Create(dir + "/" + filename + ".out")
	if err != nil {
		return errors.WithStack(err)
	}
	defer out.Close()
	// 扫描
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		_, _ = out.WriteString(fmt.Sprintf("%s\t%s\t%q\n", fset.Position(pos), tok, lit))
	}
	return nil
}
