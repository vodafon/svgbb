package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// The //go:embed directive tells the Go compiler to embed the contents
// of the templates directory into the templateFiles variable.
// This variable is an embed.FS, which is a read-only file system.
//
//go:embed templates/*.tmpl
var templateFiles embed.FS

func main() {
	// Define the command-line flags for the tool.
	// -dir: Specifies the output directory. Defaults to the current directory (".").
	// -s: The string to be injected into the templates. This is a required flag.
	// -prefix: An optional prefix for the generated filenames.
	targetDir := flag.String("dir", ".", "Target directory for generated SVG files")
	injectionString := flag.String("s", "", "String to inject into the templates (required)")
	filenamePrefix := flag.String("prefix", "", "Prefix for the output filenames")

	// Parse the command-line flags provided by the user.
	flag.Parse()

	// Validate that the required -s flag was provided.
	if *injectionString == "" {
		fmt.Println("Error: The -s flag is required to provide the payload/string for the templates.")
		flag.Usage() // Print the default usage message.
		os.Exit(1)
	}

	// The 'templates' directory is now the root within our embedded filesystem.
	templatesDir := "templates"

	// Read all files from the embedded 'templates' directory.
	files, err := templateFiles.ReadDir(templatesDir)
	if err != nil {
		// This error would typically happen if the 'templates' directory
		// didn't exist at compile time.
		fmt.Printf("Error: Could not read the embedded '%s' directory: %v\n", templatesDir, err)
		os.Exit(1)
	}

	// Ensure the target output directory exists, creating it if necessary.
	if err := os.MkdirAll(*targetDir, 0755); err != nil {
		fmt.Printf("Error: Could not create target directory '%s': %v\n", *targetDir, err)
		os.Exit(1)
	}

	fmt.Printf("Generating SVGs in directory: '%s'\n", *targetDir)
	generatedCount := 0

	// Iterate over each file found in the embedded 'templates' directory.
	for _, file := range files {
		// We only care about files ending in .tmpl, not directories.
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".tmpl") {
			continue
		}

		// Construct the path within the embedded filesystem.
		templatePath := filepath.Join(templatesDir, file.Name())

		// Read the content of the template file from the embedded filesystem.
		templateContent, err := templateFiles.ReadFile(templatePath)
		if err != nil {
			fmt.Printf("Warning: Skipping file. Could not read embedded file '%s': %v\n", templatePath, err)
			continue
		}

		// Parse the template content.
		tmpl, err := template.New(file.Name()).Parse(string(templateContent))
		if err != nil {
			fmt.Printf("Warning: Skipping file. Could not parse template '%s': %v\n", templatePath, err)
			continue
		}

		// Construct the final output filename.
		// It will be [prefix] + [template_name_without_.tmpl].
		outputFilename := *filenamePrefix + strings.TrimSuffix(file.Name(), ".tmpl")
		outputPath := filepath.Join(*targetDir, outputFilename)

		// Create the new SVG file.
		outputFile, err := os.Create(outputPath)
		if err != nil {
			fmt.Printf("Warning: Skipping file. Could not create '%s': %v\n", outputPath, err)
			continue
		}

		// Execute the template, injecting the string from the -s flag
		// and writing the result to the output file.
		err = tmpl.Execute(outputFile, *injectionString)
		outputFile.Close() // Ensure the file is closed.
		if err != nil {
			fmt.Printf("Warning: Failed to execute template for '%s': %v\n", outputPath, err)
			continue
		}

		fmt.Printf("  -> Created %s\n", outputPath)
		generatedCount++
	}

	if generatedCount == 0 {
		fmt.Println("No valid '.tmpl' files found in the embedded templates. Nothing was generated.")
	} else {
		fmt.Printf("\nDone. Successfully generated %d SVG file(s).\n", generatedCount)
	}
}
