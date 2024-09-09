package main

import (
	"errors"

	"fmt"
	"os"
	"path/filepath"

	"github.com/flexstack/envtpl"
	"github.com/joho/godotenv"
	"github.com/pterm/pterm"
	flag "github.com/spf13/pflag"
)

func main() {
	// The first argument is the path to the template file
	var file string
	flag.StringVarP(&file, "output", "o", ".env", "The path to output the .env file")
	flag.Parse()

	path := flag.Arg(0)
	if path == "" {
		defaultPaths := []string{".env.template", ".env.example"}
		for _, p := range defaultPaths {
			if _, err := os.Stat(p); err == nil {
				path = p
				break
			}
		}
	}

	// Check if the path is a directory
	if fi, err := os.Stat(path); err != nil || fi.IsDir() {
		defaultPaths := []string{".env.template", ".env.example"}
		found := false

		for _, p := range defaultPaths {
			maybePath := filepath.Join(path, p)

			if _, err := os.Stat(maybePath); err == nil {
				path = maybePath
				found = true
				break
			}
		}

		if !found {
			fmt.Println("No template file found")
			os.Exit(1)
		}
	}

	if path == "" {
		fmt.Println("No template file found")
		os.Exit(1)
	}

	// Parse the template file
	envVars, err := envtpl.Parse(path)
	if err != nil {
		fmt.Println("Error parsing template file:", path)
		os.Exit(1)
	}

	env := make(map[string]string)

	// If the .env file already exists, merge the contents
	if _, err = os.Stat(file); err == nil {
		existing, err := godotenv.Read(file)
		if err != nil {
			fmt.Println("Error reading existing .env file:", err)
			os.Exit(1)
		}

		for key, value := range existing {
			env[key] = value
		}
	}

	hasChanges := false
	for _, envVar := range envVars {
		if env[envVar.Key] != "" {
			continue
		}

		hasChanges = true
		value, err := envVar.Value.Generate()
		if err != nil {
			fmt.Println("Error generating value for:", envVar.Key)
			os.Exit(1)
		}
		switch envVar.Value.Type {
		case envtpl.Text, envtpl.Password:
			prompt := envVar.Key
			v, _ := envVar.Value.Generate()
			if p, ok := v.(string); ok && p != "" {
				prompt = p
			}
			var result string
			if envVar.Value.Type == envtpl.Password {
				result, _ = pterm.DefaultInteractiveTextInput.WithMask("*").Show(prompt)
			} else {
				result, _ = pterm.DefaultInteractiveTextInput.Show(prompt)
			}
			env[envVar.Key] = result

		case envtpl.Enum:
			rawEnum, err := envVar.Value.Generate()
			if errors.Is(err, envtpl.ErrInvalidArg) {
				fmt.Println("Error generating value for:", envVar.Key)
				os.Exit(1)
			}
			enum, ok := rawEnum.([]string)
			if !ok {
				fmt.Println("Error generating value for:", envVar.Key)
				os.Exit(1)
			}
			result, _ := pterm.DefaultInteractiveSelect.WithOptions(enum).WithFilter(false).Show(envVar.Key)
			env[envVar.Key] = result

		default:
			env[envVar.Key] = value.(string)
		}
	}

	if !hasChanges {
		fmt.Println("No changes to .env file")
		os.Exit(0)
	}

	// Generate the .env
	contents, err := godotenv.Marshal(env)
	if err != nil {
		fmt.Println("Error generating .env file:", err)
		os.Exit(1)
	}

	// Write the .env file
	fmt.Println("Writing to:", file)
	if err = os.WriteFile(file, []byte(contents), 0644); err != nil {
		fmt.Println("Error writing .env file:", err)
		os.Exit(1)
	}
}
