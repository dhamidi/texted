# Texted Command Line Interface

## Synopsis

```
texted [OPTIONS] [SCRIPT] [FILE...]
texted parse [OPTIONS]
```

## Description

Texted is a scriptable, headless text editor for automated file editing. It processes scripts written in shell-like syntax, S-expressions, or JSON format to perform text transformations.

## Options

### Script Input Options

- `-s, --script SCRIPT`     Execute SCRIPT directly (can be used multiple times)
- `-f, --file SCRIPT_FILE`  Read script from SCRIPT_FILE
- `-e, --expression EXPR`   Execute single expression EXPR (can be used multiple times)

### Input/Output Options

- `-i, --in-place`          Edit files in place (modify original files)
- `-o, --output FILE`       Write output to FILE (single file mode only)
- `--backup SUFFIX`         Create backup files with SUFFIX when using --in-place

### Script Format Options

- `--format FORMAT`         Specify script format: shell, sexp, json (default: auto-detect)
- `--shell`                 Force shell-like syntax parsing
- `--sexp`                  Force S-expression syntax parsing  
- `--json`                  Force JSON syntax parsing

### Behavior Options

- `-v, --verbose`           Enable verbose output
- `-q, --quiet`             Suppress all output except errors
- `-n, --dry-run`           Show what would be done without making changes
- `--encoding ENCODING`     Specify text encoding (default: utf8)

### Help and Information

- `-h, --help`              Show this help message
- `--version`               Show version information
- `--list-functions`        List all available functions

## Parse Subcommand

The `texted parse` subcommand converts scripts between different formats without executing them.

### Synopsis

```
texted parse [--input-format FORMAT] [--output-format FORMAT]
```

### Description

Reads a script from stdin in the specified input format and writes the parsed version to stdout in the specified output format. This is useful for converting between shell-like syntax, S-expressions, and JSON formats.

### Options

- `--input-format FORMAT`   Input script format: shell, sexp, json (default: shell)
- `--output-format FORMAT` Output script format: shell, sexp, json (default: sexp)

### Examples

Convert shell-like syntax to S-expressions:
```bash
echo 'search-forward "text"' | texted parse
```

Convert S-expressions to JSON:
```bash
echo '(search-forward "text")' | texted parse --input-format sexp --output-format json
```

Convert JSON to shell-like syntax:
```bash
echo '["search-forward", "text"]' | texted parse --input-format json --output-format shell
```

## Usage Examples

### Basic Usage

Execute a script on a single file:

```bash
texted -s 'search-forward "old"' -s 'replace-match "new"' file.txt
```

Edit multiple files in place:

```bash
texted -i -s 'search-forward "TODO"' -s 'insert "[DONE] "' *.js
```

### Script from File

```bash
# Create script file
cat > script.txt << 'EOF'
search-forward "function doIt"
set-mark
search-forward "("
replace-region "helloWorld"
EOF

# Apply to file
texted -f script.txt -i example.js
```

### Different Script Formats

Shell-like syntax:

```bash
texted -s 'search-forward "text"' -s 'replace-match "replacement"' file.txt
```

S-expression format:

```bash
texted --sexp -s '(search-forward "text")' -s '(replace-match "replacement")' file.txt
```

JSON format:

```bash
texted --json -s '["search-forward", "text"]' -s '["replace-match", "replacement"]' file.txt
```

### Output Options

Save to new file:

```bash
texted -s 'upcase-region' -o output.txt input.txt
```

Edit in place with backup:

```bash
texted -i --backup .bak -s 'replace-region "new content"' file.txt
```

### Advanced Usage

Dry run to preview changes:

```bash
texted -n -s 'search-forward "old"' -s 'replace-match "new"' *.txt
```

Verbose execution:

```bash
texted -v -f complex_script.txt file1.txt file2.txt
```

Process stdin:

```bash
echo "hello world" | texted -s 'search-forward "world"' -s 'replace-match "universe"'
```

## Exit Codes

- `0` - Success
- `1` - General error (script execution failed, file not found, etc.)
- `2` - Invalid command line arguments
- `3` - Script parsing error
- `4` - File I/O error

## Script Execution

When multiple files are specified:

- The script is executed once for each file
- Each file is processed independently
- Exit code reflects overall success/failure

When no files are specified:

- Input is read from stdin
- Output is written to stdout (unless `-o` specified)

## Notes

- File encoding is auto-detected when possible, or specify with `--encoding`
- Script format is auto-detected based on content (leading `(` for sexp, `[` for JSON)
- Backup files are only created when using `--in-place --backup`
- Regular expressions use Go's regexp syntax
- Position numbering is 1-based (first character is position 1)

