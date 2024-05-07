package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"flag"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

var APP_VERSION = "unreleased"

// cannot continue, exit immediately without a stacktrace.
func fatal() {
	fmt.Printf("cannot continue, ") // "cannot continue, exit status 1"
	os.Exit(1)
}

func configure_validator(schema_file string) *jsonschema.Schema {
	label := "release.json"

	compiler := jsonschema.NewCompiler()
	//compiler.Draft = jsonschema.Draft4 // todo: either drop schema version or raise this one

	file_bytes, err := os.ReadFile(schema_file)
	if err != nil {
		slog.Error("failed to read the json schema", "path", schema_file)
		fatal()
	}

	err = compiler.AddResource(label, bytes.NewReader(file_bytes))
	if err != nil {
		slog.Error("failed to add schema to compiler", "error", err)
		fatal()
	}
	schema, err := compiler.Compile(label)
	if err != nil {
		slog.Error("failed to compile schema", "error", err)
		fatal()
	}

	return schema
}

func path_exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func read_cli_args(arg_list []string) (string, string) {
	input_file_ptr := flag.String("in", "", "path to release.json file")
	schema_file_ptr := flag.String("schema", "", "path to release.json schema file")
	phelp := flag.Bool("help", false, "print this help and exit")
	pversion := flag.Bool("version", false, "print program version and exit")

	flag.Parse()

	if phelp != nil && *phelp {
		flag.Usage()
		os.Exit(0)
	}

	if pversion != nil && *pversion {
		fmt.Println(APP_VERSION)
		os.Exit(0)
	}

	// ---

	if input_file_ptr == nil || *input_file_ptr == "" {
		fmt.Println("--in is required")
		fatal()
	}

	input_file := *input_file_ptr

	if !strings.HasSuffix(input_file, ".json") && !strings.HasSuffix(input_file, ".jsonl") {
		fmt.Printf("input file has unsupported file extension. supported extensions: .json, .jsonl\n")
		fatal()
	}

	if !path_exists(input_file) {
		fmt.Printf("input file does not exist: %s\n", input_file)
		fatal()
	}

	// ---

	if schema_file_ptr == nil || *schema_file_ptr == "" {
		fmt.Println("--schema is required")
		fatal()
	}

	schema_file := *schema_file_ptr

	if !strings.HasSuffix(schema_file, ".json") {
		fmt.Printf("schema file has unsupported file extension. supported extensions: .json\n")
		fatal()
	}

	if !path_exists(schema_file) {
		fmt.Printf("schema file does not exist: %s\n", schema_file)
		fatal()
	}

	return schema_file, input_file
}

func validate(schema *jsonschema.Schema, release_dot_json_bytes []byte) error {
	var raw interface{}
	err := json.Unmarshal(release_dot_json_bytes, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal release.json bytes into a generic struct for validation: %w", err)
	}

	return schema.Validate(raw)
}

func read_input_file(input_file string) ([][]byte, error) {
	empty_response := [][]byte{}

	switch filepath.Ext(input_file) {
	case ".json":
		bl, err := os.ReadFile(input_file)
		if err != nil {
			return empty_response, err
		}
		return [][]byte{bl}, nil

	case ".jsonl":
		bl := [][]byte{}
		fh, err := os.Open(input_file)
		if err != nil {
			return empty_response, fmt.Errorf("failed to open input file for reading: %w", err)
		}
		scanner := bufio.NewScanner(fh)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			blt := []byte(scanner.Text()) // this seems to work when scanner.Bytes() does not ...
			bl = append(bl, blt)
		}
		if err := scanner.Err(); err != nil {
			return empty_response, fmt.Errorf("failed to read contents of input file: %w", err)
		}
		return bl, nil

	default:
		return empty_response, nil
	}
}

func main() {
	schema_file, input_file := read_cli_args(os.Args)
	schema := configure_validator(schema_file)
	release_dot_json_bytes_list, err := read_input_file(input_file)
	if err != nil {
		slog.Error("failed to read input file", "error", err)
		fatal()
	}

	for i, release_dot_json_bytes := range release_dot_json_bytes_list {
		err = validate(schema, release_dot_json_bytes)
		if err != nil {
			fmt.Printf("%#v\n", err)
			break
		}
		fmt.Printf("%d valid\n", i+1)
	}
}
