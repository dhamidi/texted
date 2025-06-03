# Texted Function Reference

This document provides a comprehensive reference for all functions available in texted's simplified elisp implementation. Texted is designed for scriptable text editing without a user interface.

## Start Here

The following functions are essential for making changes to existing source code files:

1. **`search-forward`** _string_ - Find text in your code
2. **`search-backward`** _string_ - Search backwards for text  
3. **`re-search-forward`** _regexp_ - Find patterns using regular expressions
4. **`set-mark`** - Mark the start of a region to select
5. **`replace-region`** _string_ - Replace selected text with new content
6. **`replace-match`** _string_ - Replace the last search match
7. **`insert`** _string_ - Add new text at the current position
8. **`delete-region`** - Remove selected text
9. **`goto-line`** _line-number_ - Jump to a specific line number
10. **`mark-word`** - Select a complete word for modification

## Search and Navigation Functions

### `search-forward` _string_

Search for the first occurrence of _string_ forward from the current point. Sets point to the end of the match.

### `search-backward` _string_

Search for the first occurrence of _string_ backward from the current point. Sets point to the beginning of the match.

### `re-search-forward` _regexp_

Search for the first occurrence of regular expression _regexp_ forward from the current point. Sets point to the end of the match.

### `re-search-backward` _regexp_

Search for the first occurrence of regular expression _regexp_ backward from the current point. Sets point to the beginning of the match.

### `goto-char` _position_

Move point to the specified character position.

### `goto-line` _line-number_

Move point to the beginning of the specified line number.

### `beginning-of-line`

Move point to the beginning of the current line.

### `end-of-line`

Move point to the end of the current line.

### `beginning-of-buffer`

Move point to the beginning of the buffer.

### `end-of-buffer`

Move point to the end of the buffer.

### `forward-char` [_count_]

Move point forward by _count_ characters (default 1).

### `backward-char` [_count_]

Move point backward by _count_ characters (default 1).

### `forward-word` [_count_]

Move point forward by _count_ words (default 1).

### `backward-word` [_count_]

Move point backward by _count_ words (default 1).

## Mark and Region Functions

### `set-mark`

Set the mark at the current point position.

### `set-mark-command` [_position_]

Set the mark at the specified position, or at point if no position given.

### `exchange-point-and-mark`

Exchange the positions of point and mark.

### `mark-whole-buffer`

Set mark at beginning of buffer and point at end, selecting entire buffer.

### `mark-word` [_count_]

Set mark at point and move point forward by _count_ words (default 1).

### `mark-line` [_count_]

Set mark at beginning of current line and move point to beginning of next _count_ lines.

## Text Modification Functions

### `replace-region` _string_

Replace the current region (between point and mark) with _string_.

### `replace-match` _string_

Replace the text matched by the last search operation with _string_.

### `insert` _string_

Insert _string_ at the current point position.

### `delete-char` [_count_]

Delete _count_ characters forward from point (default 1).

### `delete-backward-char` [_count_]

Delete _count_ characters backward from point (default 1).

### `delete-region`

Delete the text in the current region (between point and mark).

### `delete-line` [_count_]

Delete _count_ lines starting from the current line (default 1).

### `kill-line` [_count_]

Delete from point to end of line, or _count_ lines if specified.

### `kill-word` [_count_]

Delete _count_ words forward from point (default 1).

### `backward-kill-word` [_count_]

Delete _count_ words backward from point (default 1).

## Query and Information Functions

### `point`

Return the current position of point.

### `mark`

Return the current position of mark.

### `point-min`

Return the minimum valid point position in the buffer.

### `point-max`

Return the maximum valid point position in the buffer.

### `buffer-size`

Return the size of the buffer in characters.

### `line-number-at-pos` [_position_]

Return the line number at the specified position, or at point if no position given.

### `current-column`

Return the column number of point on the current line.

### `region-beginning`

Return the beginning position of the current region.

### `region-end`

Return the end position of the current region.

### `buffer-substring` _start_ _end_

Return the text between positions _start_ and _end_.

### `looking-at` _regexp_

Return true if text at point matches regular expression _regexp_.

### `looking-back` _regexp_

Return true if text before point matches regular expression _regexp_.

## String and Text Analysis Functions

### `length` _string_

Return the length of _string_ in characters.

### `substring` _string_ _start_ [_end_]

Return substring of _string_ from _start_ to _end_ (or end of string).

### `concat` _string1_ _string2_

Concatenate multiple strings into one.

### `upcase` _string_

Return _string_ converted to uppercase.

### `downcase` _string_

Return _string_ converted to lowercase.

### `capitalize` _string_

Return _string_ with first character capitalized.

### `string-match` _regexp_ _string_

Return position of first match of _regexp_ in _string_, or nil if no match.

### `replace-regexp-in-string` _regexp_ _replacement_ _string_

Return _string_ with all matches of _regexp_ replaced by _replacement_.

## List and Data Functions

### `list` _item1_ _item2_

Create a list containing the specified items.

### `length` _list_

Return the number of elements in _list_.

### `nth` _index_ _list_

Return the element at _index_ in _list_ (0-based).

### `first` _list_

Return the first element of _list_.

### `rest` _list_

Return all elements of _list_ except the first.

### `append` _list1_ _list2_

Concatenate multiple lists into one.

## Numeric Functions

### `+` _number1_ _number2_

Add numbers together.

### `-` _number1_ _number2_

Subtract numbers (or negate if single argument).

### `*` _number1_ _number2_

Multiply numbers together.

### `/` _number1_ _number2_

Divide numbers.

### `=` _number1_ _number2_

Return true if numbers are equal.

### `<` _number1_ _number2_

Return true if _number1_ is less than _number2_.

### `>` _number1_ _number2_

Return true if _number1_ is greater than _number2_.

### `<=` _number1_ _number2_

Return true if _number1_ is less than or equal to _number2_.

### `>=` _number1_ _number2_

Return true if _number1_ is greater than or equal to _number2_.

## Control Flow Functions

### `progn` _expression1_ _expression2_

Evaluate expressions in sequence and return the value of the last one.

### `if` _condition_ _then-expression_ [_else-expression_]

Evaluate _then-expression_ if _condition_ is true, otherwise _else-expression_.

### `when` _condition_ _expression1_ _expression2_

Evaluate expressions if _condition_ is true.

### `unless` _condition_ _expression1_ _expression2_

Evaluate expressions if _condition_ is false.

### `and` _expression1_ _expression2_

Evaluate expressions until one returns false, or return the last value.

### `or` _expression1_ _expression2_

Evaluate expressions until one returns true, or return the last value.

### `not` _expression_

Return true if _expression_ is false or nil, otherwise false.

## Notes

- All position arguments are 1-based (first character is at position 1)
- String arguments should be quoted when used in shell-like syntax
- Functions with optional arguments use default values when arguments are omitted
- Regular expressions follow Go's regexp syntax
- This implementation focuses on text editing operations suitable for automated scripting

