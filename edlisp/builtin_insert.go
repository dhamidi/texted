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