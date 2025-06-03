package edlisp

import (
	"fmt"
	"regexp"
	"strings"
)

// BuiltinInsert implements the insert function that inserts text at the current point.
func BuiltinInsert(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("insert expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("insert expects a string argument")
	}

	str := args[0].(*String)
	buffer.Insert(str.Value)
	return NewString(""), nil
}

// BuiltinGotoChar moves the point to the specified character position.
func BuiltinGotoChar(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("goto-char expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheNumberKind) {
		return nil, fmt.Errorf("goto-char expects a number argument")
	}

	num := args[0].(*Number)
	pos := int(num.Value)
	
	content := buffer.String()
	if pos < 1 {
		pos = 1
	} else if pos > len(content)+1 {
		pos = len(content) + 1
	}
	
	buffer.SetPoint(pos)
	return NewString(""), nil
}

// BuiltinGotoLine moves the point to the beginning of the specified line.
func BuiltinGotoLine(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("goto-line expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheNumberKind) {
		return nil, fmt.Errorf("goto-line expects a number argument")
	}

	num := args[0].(*Number)
	lineNum := int(num.Value)
	
	content := buffer.String()
	lines := strings.Split(content, "\n")
	
	if lineNum < 1 {
		lineNum = 1
	} else if lineNum > len(lines) {
		lineNum = len(lines)
	}
	
	// Calculate position at beginning of target line
	pos := 1
	for i := 0; i < lineNum-1; i++ {
		pos += len(lines[i]) + 1 // +1 for newline
	}
	
	buffer.SetPoint(pos)
	return NewString(""), nil
}

// BuiltinSearchForward searches for a string forward from the current point.
func BuiltinSearchForward(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("search-forward expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("search-forward expects a string argument")
	}

	str := args[0].(*String)
	content := buffer.String()
	startPos := buffer.Point() - 1 // Convert to 0-based
	
	if startPos < 0 {
		startPos = 0
	}
	if startPos >= len(content) {
		return nil, fmt.Errorf("search failed")
	}
	
	index := strings.Index(content[startPos:], str.Value)
	if index == -1 {
		return nil, fmt.Errorf("search failed")
	}
	
	// Set point to end of found text
	matchStart := startPos + index + 1 // Convert back to 1-based
	matchEnd := matchStart + len(str.Value)
	buffer.SetPoint(matchEnd)
	buffer.lastSearchMatch = str.Value
	buffer.lastSearchStart = matchStart
	buffer.lastSearchEnd = matchEnd
	
	return NewString(""), nil
}

// BuiltinSearchBackward searches for a string backward from the current point.
func BuiltinSearchBackward(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("search-backward expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("search-backward expects a string argument")
	}

	str := args[0].(*String)
	content := buffer.String()
	endPos := buffer.Point() - 1 // Convert to 0-based
	
	if endPos > len(content) {
		endPos = len(content)
	}
	if endPos < 0 {
		return nil, fmt.Errorf("search failed")
	}
	
	searchArea := content[:endPos]
	index := strings.LastIndex(searchArea, str.Value)
	if index == -1 {
		return nil, fmt.Errorf("search failed")
	}
	
	// Set point to end of found text
	matchStart := index + 1 // Convert back to 1-based
	matchEnd := matchStart + len(str.Value)
	buffer.SetPoint(matchEnd)
	buffer.lastSearchMatch = str.Value
	buffer.lastSearchStart = matchStart
	buffer.lastSearchEnd = matchEnd
	
	return NewString(""), nil
}

// BuiltinReSearchForward searches for a regular expression forward from the current point.
func BuiltinReSearchForward(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("re-search-forward expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("re-search-forward expects a string argument")
	}

	str := args[0].(*String)
	content := buffer.String()
	startPos := buffer.Point() - 1 // Convert to 0-based
	
	if startPos < 0 {
		startPos = 0
	}
	if startPos >= len(content) {
		return nil, fmt.Errorf("search failed")
	}
	
	re, err := regexp.Compile(str.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid regexp: %v", err)
	}
	
	match := re.FindStringIndex(content[startPos:])
	if match == nil {
		return nil, fmt.Errorf("search failed")
	}
	
	// Set point to end of found text
	matchStart := startPos + match[0] + 1 // Convert back to 1-based
	matchEnd := startPos + match[1] + 1
	buffer.SetPoint(matchEnd)
	buffer.lastSearchMatch = content[startPos+match[0] : startPos+match[1]]
	buffer.lastSearchStart = matchStart
	buffer.lastSearchEnd = matchEnd
	
	return NewString(""), nil
}

// BuiltinEndOfBuffer moves the point to the end of the buffer.
func BuiltinEndOfBuffer(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("end-of-buffer expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	buffer.SetPoint(len(content) + 1)
	return NewString(""), nil
}

// BuiltinSetMark sets the mark at the current point position.
func BuiltinSetMark(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("set-mark expects 0 arguments, got %d", len(args))
	}
	
	buffer.SetMark(buffer.Point())
	return NewString(""), nil
}

// BuiltinMark returns the current mark position.
func BuiltinMark(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("mark expects 0 arguments, got %d", len(args))
	}
	
	return NewNumber(float64(buffer.Mark())), nil
}

// BuiltinPoint returns the current point position.
func BuiltinPoint(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("point expects 0 arguments, got %d", len(args))
	}
	
	return NewNumber(float64(buffer.Point())), nil
}

// BuiltinLineNumberAtPos returns the line number at the current point.
func BuiltinLineNumberAtPos(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("line-number-at-pos expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	if pos < 0 {
		pos = 0
	}
	if pos > len(content) {
		pos = len(content)
	}
	
	lineNum := 1
	for i := 0; i < pos; i++ {
		if content[i] == '\n' {
			lineNum++
		}
	}
	
	return NewNumber(float64(lineNum)), nil
}

// BuiltinReplaceMatch replaces the last search match with new text.
func BuiltinReplaceMatch(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("replace-match expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("replace-match expects a string argument")
	}

	str := args[0].(*String)
	
	if buffer.lastSearchMatch == "" {
		return nil, fmt.Errorf("no previous search")
	}
	
	content := buffer.String()
	start := buffer.lastSearchStart - 1 // Convert to 0-based
	end := buffer.lastSearchEnd - 1     // Convert to 0-based
	
	if start < 0 || end > len(content) || start >= end {
		return nil, fmt.Errorf("invalid search match positions")
	}
	
	newContent := content[:start] + str.Value + content[end:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	// Update point to end of replacement
	buffer.SetPoint(start + len(str.Value) + 1)
	
	return NewString(""), nil
}

// BuiltinDeleteRegion deletes the text between mark and point.
func BuiltinDeleteRegion(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("delete-region expects 0 arguments, got %d", len(args))
	}
	
	start := buffer.Mark()
	end := buffer.Point()
	
	if start > end {
		start, end = end, start
	}
	
	content := buffer.String()
	start-- // Convert to 0-based
	end--   // Convert to 0-based
	
	if start < 0 {
		start = 0
	}
	if end > len(content) {
		end = len(content)
	}
	if start >= end {
		return NewString(""), nil
	}
	
	newContent := content[:start] + content[end:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	// Set point to start of deleted region
	buffer.SetPoint(start + 1)
	
	return NewString(""), nil
}

// BuiltinBufferSubstring returns the substring between two positions.
func BuiltinBufferSubstring(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("buffer-substring expects 2 arguments, got %d", len(args))
	}

	if !IsA(args[0], TheNumberKind) || !IsA(args[1], TheNumberKind) {
		return nil, fmt.Errorf("buffer-substring expects number arguments")
	}

	start := int(args[0].(*Number).Value)
	end := int(args[1].(*Number).Value)
	
	content := buffer.String()
	
	// Handle special case: -1 means end of buffer
	if end == -1 {
		end = len(content) + 1
	}
	
	start-- // Convert to 0-based
	end--   // Convert to 0-based
	
	if start < 0 {
		start = 0
	}
	if end > len(content) {
		end = len(content)
	}
	if start >= end {
		return NewString(""), nil
	}
	
	return NewString(content[start:end]), nil
}

// BuiltinMarkWord marks the word at the current point.
func BuiltinMarkWord(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("mark-word expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	if pos < 0 || pos >= len(content) {
		return NewString(""), nil
	}
	
	// Find the start of the word (move backward to find non-letter)
	start := pos
	for start > 0 && isLetter(content[start-1]) {
		start--
	}
	
	// Find the end of the word (move forward to find non-letter)
	end := pos
	for end < len(content) && isLetter(content[end]) {
		end++
	}
	
	// Set mark at beginning of word, point at end
	buffer.SetMark(start + 1)     // Convert back to 1-based
	buffer.SetPoint(end + 1)      // Convert back to 1-based
	
	return NewString(""), nil
}

// isLetter checks if a character is a letter or digit (word character)
func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')
}

// BuiltinReplaceRegion replaces the text in the region with new text.
func BuiltinReplaceRegion(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("replace-region expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("replace-region expects a string argument")
	}

	str := args[0].(*String)
	
	start := buffer.Mark()
	end := buffer.Point()
	
	if start > end {
		start, end = end, start
	}
	
	content := buffer.String()
	start-- // Convert to 0-based
	end--   // Convert to 0-based
	
	if start < 0 {
		start = 0
	}
	if end > len(content) {
		end = len(content)
	}
	if start >= end {
		return NewString(""), nil
	}
	
	newContent := content[:start] + str.Value + content[end:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	// Set point to end of replacement
	buffer.SetPoint(start + len(str.Value) + 1)
	
	return NewString(""), nil
}

// BuiltinRegionBeginning returns the beginning position of the region.
func BuiltinRegionBeginning(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("region-beginning expects 0 arguments, got %d", len(args))
	}
	
	start := buffer.Mark()
	end := buffer.Point()
	
	if start > end {
		start = end
	}
	
	return NewNumber(float64(start)), nil
}

// BuiltinRegionEnd returns the end position of the region.
func BuiltinRegionEnd(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("region-end expects 0 arguments, got %d", len(args))
	}
	
	start := buffer.Mark()
	end := buffer.Point()
	
	if start > end {
		end = start
	}
	
	return NewNumber(float64(end)), nil
}

// BuiltinForwardChar moves the point forward by the specified number of characters.
func BuiltinForwardChar(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("forward-char expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("forward-char expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	newPos := buffer.Point() + count
	
	if newPos < 1 {
		newPos = 1
	} else if newPos > len(content)+1 {
		newPos = len(content) + 1
	}
	
	buffer.SetPoint(newPos)
	return NewString(""), nil
}

// BuiltinLength returns the length of a string.
func BuiltinLength(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("length expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("length expects a string argument")
	}

	str := args[0].(*String)
	return NewNumber(float64(len(str.Value))), nil
}

// BuiltinSubstring returns a substring of the given string.
func BuiltinSubstring(args []Value, buffer *Buffer) (Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("substring expects 2 or 3 arguments, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("substring expects a string as first argument")
	}

	if !IsA(args[1], TheNumberKind) {
		return nil, fmt.Errorf("substring expects a number as second argument")
	}

	str := args[0].(*String)
	start := int(args[1].(*Number).Value)
	end := len(str.Value)

	if len(args) == 3 {
		if !IsA(args[2], TheNumberKind) {
			return nil, fmt.Errorf("substring expects a number as third argument")
		}
		end = int(args[2].(*Number).Value)
	}

	// Convert from 1-based to 0-based indexing
	start--

	// For two-argument form, end should be to the end of string
	if len(args) == 2 {
		end = len(str.Value)
	} else {
		// For three-argument form, end is 1-based and exclusive
		// Convert to 0-based exclusive by decrementing
		end--
	}

	// Bounds checking
	if start < 0 {
		start = 0
	}
	if end > len(str.Value) {
		end = len(str.Value)
	}
	if start > end {
		return NewString(""), nil
	}

	result := str.Value[start:end]
	return NewString(result), nil
}

// BuiltinStringMatch searches for a pattern in a string and returns the match position.
func BuiltinStringMatch(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-match expects 2 arguments, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) || !IsA(args[1], TheStringKind) {
		return nil, fmt.Errorf("string-match expects string arguments")
	}

	pattern := args[0].(*String)
	str := args[1].(*String)

	// Try to compile as regular expression
	re, err := regexp.Compile(pattern.Value)
	if err != nil {
		// If not a valid regex, treat as literal string
		index := strings.Index(str.Value, pattern.Value)
		if index == -1 {
			return NewSymbol("nil"), nil
		}
		return NewNumber(float64(index)), nil
	}

	// Use regular expression matching
	match := re.FindStringIndex(str.Value)
	if match == nil {
		return NewSymbol("nil"), nil
	}

	return NewNumber(float64(match[0])), nil
}

// BuiltinBackwardChar moves the point backward by the specified number of characters.
func BuiltinBackwardChar(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("backward-char expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("backward-char expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	newPos := buffer.Point() - count
	
	if newPos < 1 {
		newPos = 1
	}
	
	buffer.SetPoint(newPos)
	return NewString(""), nil
}

// BuiltinBeginningOfBuffer moves the point to the beginning of the buffer.
func BuiltinBeginningOfBuffer(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("beginning-of-buffer expects 0 arguments, got %d", len(args))
	}
	
	buffer.SetPoint(1)
	return NewString(""), nil
}

// BuiltinBeginningOfLine moves the point to the beginning of the current line.
func BuiltinBeginningOfLine(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("beginning-of-line expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	if pos < 0 {
		buffer.SetPoint(1)
		return NewString(""), nil
	}
	if pos >= len(content) {
		pos = len(content) - 1
	}
	
	// Move backward to find beginning of line
	for pos > 0 && content[pos-1] != '\n' {
		pos--
	}
	
	buffer.SetPoint(pos + 1) // Convert back to 1-based
	return NewString(""), nil
}

// BuiltinEndOfLine moves the point to the end of the current line.
func BuiltinEndOfLine(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("end-of-line expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	if pos < 0 {
		pos = 0
	}
	if pos >= len(content) {
		buffer.SetPoint(len(content) + 1)
		return NewString(""), nil
	}
	
	// Move forward to find end of line (newline character or end of content)
	for pos < len(content) && content[pos] != '\n' {
		pos++
	}
	
	// pos now points to newline or beyond end of content
	// We want to be at the last character of the line, not the newline
	buffer.SetPoint(pos + 1) // Convert back to 1-based
	return NewString(""), nil
}

// BuiltinBufferSize returns the size of the buffer.
func BuiltinBufferSize(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("buffer-size expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	return NewNumber(float64(len(content))), nil
}

// BuiltinPointMax returns the maximum valid point position.
func BuiltinPointMax(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("point-max expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	return NewNumber(float64(len(content) + 1)), nil
}

// BuiltinPointMin returns the minimum valid point position.
func BuiltinPointMin(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("point-min expects 0 arguments, got %d", len(args))
	}
	
	return NewNumber(1), nil
}

// BuiltinCurrentColumn returns the current column position.
func BuiltinCurrentColumn(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("current-column expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	if pos < 0 {
		return NewNumber(0), nil
	}
	if pos >= len(content) {
		pos = len(content) - 1
	}
	
	column := 0
	// Count backward to find beginning of line
	for i := pos; i >= 0 && content[i] != '\n'; i-- {
		column++
	}
	
	return NewNumber(float64(column)), nil
}

// BuiltinUpcase converts a string to uppercase.
func BuiltinUpcase(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("upcase expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("upcase expects a string argument")
	}

	str := args[0].(*String)
	return NewString(strings.ToUpper(str.Value)), nil
}

// BuiltinDowncase converts a string to lowercase.
func BuiltinDowncase(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("downcase expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("downcase expects a string argument")
	}

	str := args[0].(*String)
	return NewString(strings.ToLower(str.Value)), nil
}

// BuiltinCapitalize capitalizes the first letter of a string.
func BuiltinCapitalize(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("capitalize expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("capitalize expects a string argument")
	}

	str := args[0].(*String)
	if len(str.Value) == 0 {
		return NewString(""), nil
	}
	
	result := strings.ToUpper(string(str.Value[0])) + strings.ToLower(str.Value[1:])
	return NewString(result), nil
}

// BuiltinConcat concatenates multiple strings.
func BuiltinConcat(args []Value, buffer *Buffer) (Value, error) {
	var result strings.Builder
	
	for i, arg := range args {
		if !IsA(arg, TheStringKind) {
			return nil, fmt.Errorf("concat expects string arguments, got non-string at position %d", i+1)
		}
		str := arg.(*String)
		result.WriteString(str.Value)
	}
	
	return NewString(result.String()), nil
}

// BuiltinSetMarkCommand sets the mark at the current point and shows a message.
func BuiltinSetMarkCommand(args []Value, buffer *Buffer) (Value, error) {
	var pos int
	
	if len(args) > 1 {
		return nil, fmt.Errorf("set-mark-command expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("set-mark-command expects a number argument")
		}
		pos = int(args[0].(*Number).Value)
	} else {
		pos = buffer.Point()
	}
	
	buffer.SetMark(pos)
	return NewString(""), nil
}

// BuiltinExchangePointAndMark swaps the point and mark positions.
func BuiltinExchangePointAndMark(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("exchange-point-and-mark expects 0 arguments, got %d", len(args))
	}
	
	point := buffer.Point()
	mark := buffer.Mark()
	
	buffer.SetPoint(mark)
	buffer.SetMark(point)
	
	return NewString(""), nil
}

// BuiltinMarkWholeBuffer marks the entire buffer.
func BuiltinMarkWholeBuffer(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("mark-whole-buffer expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	buffer.SetMark(1)
	buffer.SetPoint(len(content) + 1)
	
	return NewString(""), nil
}

// BuiltinDeleteChar deletes characters at the current point.
func BuiltinDeleteChar(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("delete-char expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("delete-char expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	pos := buffer.Point() // 1-based position
	
	if pos < 1 || pos > len(content) {
		return NewString(""), nil
	}
	
	// Point position is used as 0-based index into string
	startPos := pos // Start at the character AT the point position (0-based)
	endPos := startPos + count
	if endPos > len(content) {
		endPos = len(content)
	}
	
	// Delete characters starting at current position
	newContent := content[:startPos] + content[endPos:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	return NewString(""), nil
}

// BuiltinDeleteBackwardChar deletes characters backward from the current point.
func BuiltinDeleteBackwardChar(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("delete-backward-char expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("delete-backward-char expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	startPos := pos - count
	if startPos < 0 {
		startPos = 0
	}
	
	if startPos >= pos {
		return NewString(""), nil
	}
	
	newContent := content[:startPos] + content[pos:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	// Update point to the new position
	buffer.SetPoint(startPos + 1)
	
	return NewString(""), nil
}

// BuiltinForwardWord moves the point forward by the specified number of words.
func BuiltinForwardWord(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("forward-word expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("forward-word expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	for i := 0; i < count && pos < len(content); i++ {
		// Skip non-word characters to get to a word
		for pos < len(content) && !isLetter(content[pos]) {
			pos++
		}
		// Skip current word characters to get to end of word
		for pos < len(content) && isLetter(content[pos]) {
			pos++
		}
	}
	
	buffer.SetPoint(pos + 1) // Convert back to 1-based
	return NewString(""), nil
}

// BuiltinBackwardWord moves the point backward by the specified number of words.
func BuiltinBackwardWord(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("backward-word expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("backward-word expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	for i := 0; i < count && pos > 0; i++ {
		// Skip current non-word characters
		for pos > 0 && !isLetter(content[pos-1]) {
			pos--
		}
		// Skip word characters to get to beginning of word
		for pos > 0 && isLetter(content[pos-1]) {
			pos--
		}
	}
	
	buffer.SetPoint(pos + 1) // Convert back to 1-based
	return NewString(""), nil
}

// BuiltinLookingAt checks if text at point matches a pattern.
func BuiltinLookingAt(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("looking-at expects 1 argument, got %d", len(args))
	}
	
	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("looking-at expects a string argument")
	}
	
	pattern := args[0].(*String)
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	if pos < 0 || pos >= len(content) {
		return NewSymbol("nil"), nil
	}
	
	// Try to compile as regular expression
	re, err := regexp.Compile(pattern.Value)
	if err != nil {
		// If not a valid regex, treat as literal string
		if strings.HasPrefix(content[pos:], pattern.Value) {
			return NewSymbol("t"), nil
		}
		return NewSymbol("nil"), nil
	}
	
	// Use regular expression matching
	match := re.FindStringIndex(content[pos:])
	if match != nil && match[0] == 0 {
		return NewSymbol("t"), nil
	}
	
	return NewSymbol("nil"), nil
}

// BuiltinLookingBack checks if text before point matches a pattern.
func BuiltinLookingBack(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("looking-back expects 1 argument, got %d", len(args))
	}
	
	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("looking-back expects a string argument")
	}
	
	pattern := args[0].(*String)
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	if pos <= 0 {
		return NewSymbol("nil"), nil
	}
	
	// Try to compile as regular expression
	re, err := regexp.Compile(pattern.Value)
	if err != nil {
		// If not a valid regex, treat as literal string
		if pos >= len(pattern.Value) && strings.HasSuffix(content[:pos], pattern.Value) {
			return NewSymbol("t"), nil
		}
		return NewSymbol("nil"), nil
	}
	
	// Use regular expression matching on text before point
	beforeText := content[:pos]
	match := re.FindStringIndex(beforeText)
	if match != nil && match[1] == len(beforeText) {
		return NewSymbol("t"), nil
	}
	
	return NewSymbol("nil"), nil
}

// BuiltinMarkLine marks the entire current line.
func BuiltinMarkLine(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("mark-line expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("mark-line expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	// Find beginning of current line
	lineStart := pos
	for lineStart > 0 && content[lineStart-1] != '\n' {
		lineStart--
	}
	
	// Find end of line(s) based on count
	lineEnd := pos
	for i := 0; i < count; i++ {
		for lineEnd < len(content) && content[lineEnd] != '\n' {
			lineEnd++
		}
		if lineEnd < len(content) && content[lineEnd] == '\n' {
			lineEnd++ // Include the newline
		}
	}
	
	buffer.SetMark(lineStart + 1) // Convert back to 1-based
	buffer.SetPoint(lineEnd + 1)  // Convert back to 1-based
	
	return NewString(""), nil
}

// BuiltinDeleteLine deletes the current line.
func BuiltinDeleteLine(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("delete-line expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("delete-line expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	// Find beginning of current line
	lineStart := pos
	for lineStart > 0 && content[lineStart-1] != '\n' {
		lineStart--
	}
	
	// Find end of line(s) based on count
	lineEnd := pos
	for i := 0; i < count; i++ {
		for lineEnd < len(content) && content[lineEnd] != '\n' {
			lineEnd++
		}
		if lineEnd < len(content) && content[lineEnd] == '\n' {
			lineEnd++ // Include the newline
		}
	}
	
	newContent := content[:lineStart] + content[lineEnd:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	buffer.SetPoint(lineStart + 1) // Convert back to 1-based
	
	return NewString(""), nil
}

// BuiltinKillLine kills (deletes) text from point to end of line.
func BuiltinKillLine(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("kill-line expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("kill-line expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	if pos < 0 || pos >= len(content) {
		return NewString(""), nil
	}
	
	// Kill multiple lines if count > 1
	lineEnd := pos
	for i := 0; i < count; i++ {
		// Find end of current line
		for lineEnd < len(content) && content[lineEnd] != '\n' {
			lineEnd++
		}
		// Include the newline character if present and if we're killing multiple lines
		if lineEnd < len(content) && content[lineEnd] == '\n' && (i < count-1 || count > 1) {
			lineEnd++
		}
	}
	
	newContent := content[:pos] + content[lineEnd:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	return NewString(""), nil
}

// BuiltinKillWord kills (deletes) word forward from point.
func BuiltinKillWord(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("kill-word expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("kill-word expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	startPos := buffer.Point() - 1 // Convert to 0-based
	pos := startPos
	
	for i := 0; i < count && pos < len(content); i++ {
		// Skip non-word characters to get to a word
		for pos < len(content) && !isLetter(content[pos]) {
			pos++
		}
		// Skip current word characters to get to end of word
		for pos < len(content) && isLetter(content[pos]) {
			pos++
		}
	}
	
	newContent := content[:startPos] + content[pos:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	return NewString(""), nil
}

// BuiltinBackwardKillWord kills (deletes) word backward from point.
func BuiltinBackwardKillWord(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("backward-kill-word expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("backward-kill-word expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	endPos := buffer.Point() - 1 // Convert to 0-based
	pos := endPos
	
	for i := 0; i < count && pos > 0; i++ {
		// Move backward to start of previous word
		// First skip any non-word characters before current position
		for pos > 0 && !isLetter(content[pos-1]) {
			pos--
		}
		// Then skip the word characters to get to beginning of word
		for pos > 0 && isLetter(content[pos-1]) {
			pos--
		}
	}
	
	newContent := content[:pos] + content[endPos:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	buffer.SetPoint(pos + 1) // Convert back to 1-based
	
	return NewString(""), nil
}

// BuiltinReSearchBackward searches for a regular expression backward from the current point.
func BuiltinReSearchBackward(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("re-search-backward expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("re-search-backward expects a string argument")
	}

	str := args[0].(*String)
	content := buffer.String()
	endPos := buffer.Point() - 1 // Convert to 0-based
	
	if endPos > len(content) {
		endPos = len(content)
	}
	if endPos < 0 {
		return nil, fmt.Errorf("search failed")
	}
	
	re, err := regexp.Compile(str.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid regexp: %v", err)
	}
	
	searchArea := content[:endPos]
	matches := re.FindAllStringIndex(searchArea, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("search failed")
	}
	
	// Get the last match (rightmost before point)
	match := matches[len(matches)-1]
	
	// Set point to end of found text
	matchStart := match[0] + 1 // Convert back to 1-based
	matchEnd := match[1] + 1
	buffer.SetPoint(matchEnd)
	buffer.lastSearchMatch = content[match[0]:match[1]]
	buffer.lastSearchStart = matchStart
	buffer.lastSearchEnd = matchEnd
	
	return NewString(""), nil
}

// BuiltinReplaceRegexpInString replaces regexp matches in a string.
func BuiltinReplaceRegexpInString(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("replace-regexp-in-string expects 3 arguments, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) || !IsA(args[1], TheStringKind) || !IsA(args[2], TheStringKind) {
		return nil, fmt.Errorf("replace-regexp-in-string expects string arguments")
	}

	pattern := args[0].(*String)
	replacement := args[1].(*String)
	str := args[2].(*String)
	
	re, err := regexp.Compile(pattern.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid regexp: %v", err)
	}
	
	result := re.ReplaceAllString(str.Value, replacement.Value)
	return NewString(result), nil
}