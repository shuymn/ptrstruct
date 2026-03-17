package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/shuymn/ptrstruct"
)

func main() {
	singlechecker.Main(ptrstruct.Analyzer)
}
