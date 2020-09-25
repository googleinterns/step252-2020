package io

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"html/template"

	"github.com/googleinterns/terraform-cost-estimation/io/web"
	"github.com/googleinterns/terraform-cost-estimation/resources"
)

// GetOutputWriter returns the output os.File (stdout/file) for a given output path or an error.
func GetOutputWriter(outputPath string) (*os.File, error) {
	if outputPath == "stdout" {
		fmt.Println("-----------------------------------------------------------------------------------------------------------------------------")
		return os.Stdout, nil
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// FinishOutput closes the file when the output is done and returns an error where is the case.
// If the output file is Stdout, thenit is not closed and 2 newlines are printed.
func FinishOutput(outputFile *os.File) error {
	if outputFile == nil {
		return nil
	}
	if outputFile != os.Stdout {
		return outputFile.Close()
	}
	fmt.Println("\n-----------------------------------------------------------------------------------------------------------------------------")
	fmt.Printf("\n\n\n")
	return nil
}

// GenerateWebPage generates a webpage file with the pricing information of the specified resources.
func GenerateWebPage(outputPath string, res []resources.ResourceState) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, callerFile, _, _ := runtime.Caller(0)
	t, err := template.ParseFiles(filepath.Dir(callerFile) + "/web/web_template.gohtml")
	if err != nil {
		return err
	}

	if err = t.Execute(f, mapToWebTables(res)); err != nil {
		return err
	}

	return nil
}

func mapToWebTables(res []resources.ResourceState) (t []*web.PricingTypeTables) {
	for i, r := range res {
		t = append(t, r.GetWebTables(i))
	}
	return
}
