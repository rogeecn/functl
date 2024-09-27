package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	mod "golang.org/x/mod/modfile"
)

func Test_parseModFile(t *testing.T) {
	// Read go.mod file
	data, err := os.ReadFile("go.mod")
	if err != nil {
		log.Fatalf("Error reading go.mod file: %v", err)
	}

	// Parse go.mod file
	modFile, err := mod.Parse("go.mod", data, nil)
	if err != nil {
		log.Fatalf("Error parsing go.mod file: %v", err)
	}

	// Print module path
	fmt.Println("Module Path:", modFile)
}
