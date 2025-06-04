# Texted Command Line Interface

## Synopsis

```
texted edit [OPTIONS] [SCRIPT] [FILE...]
```

## Description

Texted is a scriptable, headless text editor for automated file editing. It processes scripts written in shell-like syntax, S-expressions, or JSON format to perform text transformations.

## Options

### Script Input Options

- `-s, --script SCRIPT`     Execute SCRIPT directly (can be used multiple times)
- `-f, --file SCRIPT_FILE`  Read script from SCRIPT_FILE
- `-e, --expression EXPR`   Execute single expression EXPR (can be used multiple times), print result of evaluation

### Input/Output Options

- `-i, --in-place`          Edit files in place (modify original files)
- `-o, --output FILE`       Write output to FILE (single file mode only), defaults to
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

### Help and Information

- `-h, --help`              Show this help message
- `--version`               Show version information
