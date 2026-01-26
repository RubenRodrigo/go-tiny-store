package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	gormadapter "github.com/RubenRodrigo/go-tiny-store/internal/adapters/persistence/gorm"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		gormadapter.AllModels()...,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	io.WriteString(os.Stdout, stmts)

}
