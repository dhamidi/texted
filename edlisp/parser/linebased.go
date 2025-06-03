// Package parser implements a line-based S-expression parser for texted scripts.
package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/dhamidi/texted/edlisp"
)

// ParseReader parses texted scripts from an io.Reader using the line-based format.
// Each line is parsed according to the rules:
// 1. Leading whitespace is stripped
// 2. If the next character is '(', read a regular S-expression list
// 3. Otherwise build a list by reading tokens until a newline is encountered
func ParseReader(r io.Reader) ([]edlisp.Value, error) {
	scanner := bufio.NewScanner(r)
	var expressions []edlisp.Value

	for scanner.Scan() {
		line := strings.TrimLeftFunc(scanner.Text(), unicode.IsSpace)
		if line == "" {
			continue // Skip empty lines
		}

		expr, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	return expressions, nil
}

// ParseString parses a single texted script from a string.
func ParseString(s string) ([]edlisp.Value, error) {
	return ParseReader(strings.NewReader(s))
}

// ParseSexp parses pure S-expression format where simple values are not wrapped in lists.
// For example, "5" returns [Number(5)], not [List(Number(5))].
func ParseSexp(s string) ([]edlisp.Value, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return []edlisp.Value{}, nil
	}

	// If it starts with '(', parse as regular S-expression
	if strings.HasPrefix(trimmed, "(") {
		expr, err := parseSExpression(trimmed)
		if err != nil {
			return nil, err
		}
		return []edlisp.Value{expr}, nil
	}

	// Otherwise, parse as a single token value
	value, err := parseToken(trimmed)
	if err != nil {
		return nil, err
	}

	return []edlisp.Value{value}, nil
}

// ParseFormat parses input according to the specified format.
// Supported formats: "sexp", "shell", "json"
// Defaults to "sexp" for unknown formats.
func ParseFormat(format, input string) ([]edlisp.Value, error) {
	switch format {
	case "sexp", "":
		return ParseSexp(input)
	case "shell":
		return ParseString(input)
	case "json":
		return ParseJSONString(input)
	default:
		// Default to sexp for unknown formats
		return ParseSexp(input)
	}
}

// parseLine parses a single line according to texted rules.
func parseLine(line string) (edlisp.Value, error) {
	if strings.HasPrefix(line, "(") {
		return parseSExpression(line)
	}
	return parseShellLike(line)
}

// parseSExpression parses a traditional S-expression starting with '('.
func parseSExpression(line string) (edlisp.Value, error) {
	tokens, err := tokenize(line)
	if err != nil {
		return nil, err
	}

	expr, _, err := parseTokens(tokens, 0)
	return expr, err
}

// parseShellLike parses a shell-like command line into a list.
func parseShellLike(line string) (edlisp.Value, error) {
	tokens, err := tokenize(line)
	if err != nil {
		return nil, err
	}

	var elements []edlisp.Value
	pos := 0

	for pos < len(tokens) {
		token := tokens[pos]

		// If we encounter an opening parenthesis, parse as S-expression
		if token == "(" {
			expr, newPos, err := parseTokens(tokens, pos)
			if err != nil {
				return nil, err
			}
			elements = append(elements, expr)
			pos = newPos
		} else if token == ")" {
			return nil, fmt.Errorf("unexpected closing parenthesis in shell-like syntax")
		} else {
			// Parse as regular token
			value, err := parseToken(token)
			if err != nil {
				return nil, err
			}
			elements = append(elements, value)
			pos++
		}
	}

	return &edlisp.List{Elements: elements}, nil
}

// tokenize splits a line into tokens, handling quoted strings and parentheses.
func tokenize(line string) ([]string, error) {
	var tokens []string
	var current strings.Builder
	inQuotes := false
	escapeNext := false

	for _, r := range line {
		if escapeNext {
			current.WriteRune(r)
			escapeNext = false
			continue
		}

		switch r {
		case '\\':
			current.WriteRune(r)
			if inQuotes {
				escapeNext = true
			}
		case '"':
			current.WriteRune(r)
			inQuotes = !inQuotes
		case '(', ')':
			if inQuotes {
				current.WriteRune(r)
			} else {
				if current.Len() > 0 {
					tokens = append(tokens, current.String())
					current.Reset()
				}
				tokens = append(tokens, string(r))
			}
		case ' ', '\t', '\n', '\r':
			if inQuotes {
				current.WriteRune(r)
			} else if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if inQuotes {
		return nil, fmt.Errorf("unterminated string literal")
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens, nil
}

// parseTokens recursively parses tokens into S-expressions.
func parseTokens(tokens []string, start int) (edlisp.Value, int, error) {
	if start >= len(tokens) {
		return nil, start, fmt.Errorf("unexpected end of input")
	}

	token := tokens[start]

	if token == "(" {
		return parseList(tokens, start+1)
	}

	if token == ")" {
		return nil, start, fmt.Errorf("unexpected closing parenthesis")
	}

	value, err := parseToken(token)
	return value, start + 1, err
}

// parseList parses a list starting after the opening parenthesis.
func parseList(tokens []string, start int) (edlisp.Value, int, error) {
	var elements []edlisp.Value
	pos := start

	for pos < len(tokens) && tokens[pos] != ")" {
		element, newPos, err := parseTokens(tokens, pos)
		if err != nil {
			return nil, pos, err
		}
		elements = append(elements, element)
		pos = newPos
	}

	if pos >= len(tokens) {
		return nil, pos, fmt.Errorf("unterminated list")
	}

	return &edlisp.List{Elements: elements}, pos + 1, nil
}

// parseToken converts a single token into a Value.
func parseToken(token string) (edlisp.Value, error) {
	// Try to parse as a string (quoted)
	if strings.HasPrefix(token, `"`) && strings.HasSuffix(token, `"`) {
		if len(token) < 2 {
			return nil, fmt.Errorf("invalid string literal: %s", token)
		}
		unquoted, err := strconv.Unquote(token)
		if err != nil {
			return nil, fmt.Errorf("invalid string literal %s: %w", token, err)
		}
		return &edlisp.String{Value: unquoted}, nil
	}

	// Try to parse as a number
	if num, err := strconv.ParseFloat(token, 64); err == nil {
		return &edlisp.Number{Value: num}, nil
	}

	// Otherwise, it's a symbol
	return &edlisp.Symbol{Name: token}, nil
}
