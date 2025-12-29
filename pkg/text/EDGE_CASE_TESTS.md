# Edge Case and Stress Tests for Terminal Formatter

## Overview

This document provides a comprehensive reference for all edge case and stress tests added to the `pkg/text` formatter package. The tests are designed to ensure robustness, prevent regressions, and verify correct behavior with unusual or extreme inputs.

## Test Files Created

1. **formatter_edge_cases_test.go** - 17 test functions
2. **formatter_stress_test.go** - 13 test functions
3. **formatter_regression_test.go** - 14 test functions

**Total: 44 comprehensive test functions**

## Edge Cases Test File: formatter_edge_cases_test.go

### 1. TestFormatEmptyString
- **Purpose**: Verify handling of completely empty input
- **Test Cases**: 3 configurations (default, colored, compact)
- **Coverage**: Empty string with different formatter settings

### 2. TestFormatWhitespaceOnly
- **Purpose**: Handle input containing only whitespace
- **Test Cases**: 7 scenarios (spaces, tabs, newlines, mixed)
- **Coverage**: Various whitespace-only inputs

### 3. TestFormatVeryLongInput
- **Purpose**: Test with extremely large inputs
- **Test Cases**: 5 scenarios (1KB, 10KB, 100KB, 1MB inputs)
- **Coverage**: Memory and performance with large content

### 4. TestFormatVeryLongSingleLine
- **Purpose**: Test single lines exceeding normal width constraints
- **Test Cases**: 4 scenarios (200 to 5000 character lines)
- **Coverage**: Line wrapping with extremely long content

### 5. TestFormatDeeplyNestedMarkdown
- **Purpose**: Handle deeply nested markdown structures
- **Test Cases**: 4 scenarios (nested lists, heading levels, code, formatting)
- **Coverage**: Multi-level nesting of markdown elements

### 6. TestFormatMalformedTables
- **Purpose**: Handle malformed table syntax gracefully
- **Test Cases**: 6 scenarios (missing pipes, uneven columns, long content)
- **Coverage**: Table parsing robustness

### 7. TestFormatMixedLineEndings
- **Purpose**: Handle various line ending combinations
- **Test Cases**: 5 scenarios (LF, CRLF, CR, mixed, consecutive)
- **Coverage**: Cross-platform line ending normalization

### 8. TestFormatUnicodeCharacters
- **Purpose**: Verify Unicode character support
- **Test Cases**: 12 scenarios (emoji, CJK, Arabic, Hebrew, Cyrillic, combining chars, RTL)
- **Coverage**: Comprehensive Unicode handling including complex scripts

### 9. TestFormatSpecialRegexCharacters
- **Purpose**: Handle input with regex metacharacters safely
- **Test Cases**: 8 scenarios (., *, +, [], {}, |, ^, $, etc.)
- **Coverage**: Regex-safe input processing

### 10. TestFormatExtremeConfigValues
- **Purpose**: Test with extreme configuration values
- **Test Cases**: 7 scenarios (very small/large widths, zero/large indents, long prefixes)
- **Coverage**: Configuration boundary conditions

### 11. TestFormatComplexCodeBlocks
- **Purpose**: Handle various code block edge cases
- **Test Cases**: 10 scenarios (unclosed, empty, nested, many lines, special chars)
- **Coverage**: Code block parsing robustness

### 12. TestFormatStructuredCommands
- **Purpose**: Verify proper removal of structured commands
- **Test Cases**: 4 scenarios (simple, nested, multiple, with special chars)
- **Coverage**: Structured command tag removal

### 13. TestFormatMalformedMarkdown
- **Purpose**: Handle malformed markdown syntax
- **Test Cases**: 8 scenarios (unclosed tags, invalid headers, quotes)
- **Coverage**: Markdown parsing error handling

### 14. TestFormatAllConfigCombinations
- **Purpose**: Test various configuration combinations
- **Test Cases**: 4 major combinations (all enabled, all disabled, mixed 1, mixed 2)
- **Coverage**: Feature interaction verification

### 15. TestFormatStressLongContent
- **Purpose**: Stress test with complex mixed content
- **Test Cases**: 3 configurations with large documents
- **Coverage**: Complex document formatting under load

### 16. TestFormatBoundaryValues
- **Purpose**: Test exact boundary conditions
- **Test Cases**: 7 scenarios (width boundaries, character counts)
- **Coverage**: Off-by-one and boundary condition verification

### 17. TestFormatRobustnessWithNoOptions
- **Purpose**: Test with completely empty configuration
- **Test Cases**: 6 different inputs with empty config
- **Coverage**: Default value handling

## Stress Tests File: formatter_stress_test.go

### 1. TestFormatMemoryIntensive
- **Purpose**: Test with memory-intensive inputs
- **Test Cases**: 4 scenarios (10K+ lines, large items, many code blocks, many tables)
- **Coverage**: Memory efficiency and performance

### 2. TestFormatConcurrentStress
- **Purpose**: Verify thread-safe concurrent usage
- **Test Cases**: 100 concurrent formatter invocations
- **Coverage**: Concurrency safety

### 3. TestFormatEdgeCaseMarkdownCombinations
- **Purpose**: Test complex markdown pattern combinations
- **Test Cases**: 9 scenarios (nested formatting, links, code combinations)
- **Coverage**: Markdown interaction patterns

### 4. TestFormatLineWrappingEdgeCases
- **Purpose**: Test line wrapping edge cases
- **Test Cases**: 4 scenarios (long words, boundary conditions, disabled wrapping)
- **Coverage**: Wrapping algorithm robustness

### 5. TestFormatTableEdgeCases
- **Purpose**: Advanced table edge cases
- **Test Cases**: 6 scenarios (wide columns, many columns, unicode, adjacent tables)
- **Coverage**: Table rendering edge cases

### 6. TestFormatIndentationEdgeCases
- **Purpose**: Test indentation handling edge cases
- **Test Cases**: 4 scenarios (deep nesting, mixed styles, zero indent, nested quotes)
- **Coverage**: Indentation calculation correctness

### 7. TestFormatColoringEdgeCases
- **Purpose**: Test coloring with edge cases
- **Test Cases**: 5 scenarios (empty input, special chars, unicode, long lines)
- **Coverage**: Color application robustness

### 8. TestFormatCompactModeEdgeCases
- **Purpose**: Test compact mode behavior
- **Test Cases**: 4 scenarios (blank lines, code blocks, lists)
- **Coverage**: Compact mode spacing logic

### 9. TestFormatRegexPerformance
- **Purpose**: Test regex performance with many patterns
- **Test Cases**: 5 scenarios (bold, italic, links, code, mixed patterns)
- **Coverage**: Performance with high pattern density

### 10. TestFormatCommentPrefixEdgeCases
- **Purpose**: Test different comment prefix configurations
- **Test Cases**: 4 scenarios (empty, long, special, unicode prefixes)
- **Coverage**: Prefix handling variations

### 11. TestFormatNumberedLists
- **Purpose**: Test numbered list variations
- **Test Cases**: 5 scenarios (simple, large numbers, mixed types, non-sequential)
- **Coverage**: Numbered list parsing

### 12. TestFormatShowLineNumbersEdgeCases
- **Purpose**: Test line numbering in code blocks
- **Test Cases**: 4 scenarios (empty blocks, single line, many lines)
- **Coverage**: Line numbering logic

### 13. TestFormatWhitespaceHandling
- **Purpose**: Test various whitespace scenarios
- **Test Cases**: 5 scenarios (trailing, leading, tabs, mixed, non-breaking spaces)
- **Coverage**: Whitespace normalization

## Regression Tests File: formatter_regression_test.go

### 1. TestFormatNilInputHandling
- **Purpose**: Prevent nil pointer dereferences
- **Test Cases**: 3 configurations with empty input
- **Coverage**: Nil safety

### 2. TestFormatRegexCatastrophicBacktracking
- **Purpose**: Prevent regex DoS attacks
- **Test Cases**: 5 scenarios with potentially problematic patterns
- **Coverage**: Regex security

### 3. TestFormatRegressionOffByOne
- **Purpose**: Catch off-by-one errors
- **Test Cases**: 5 scenarios at line width boundaries
- **Coverage**: Boundary calculation accuracy

### 4. TestFormatRegressionSliceIndexing
- **Purpose**: Prevent slice index out of bounds
- **Test Cases**: 8 scenarios (single char, special chars, combinations)
- **Coverage**: Slice indexing safety

### 5. TestFormatRegressionStringMethods
- **Purpose**: Safe string method usage
- **Test Cases**: 4 scenarios (long strings, trim operations, split)
- **Coverage**: String operation safety

### 6. TestFormatRegressionMapAccess
- **Purpose**: Safe map access
- **Test Cases**: 3 scenarios (empty cells, many columns)
- **Coverage**: Map access safety

### 7. TestFormatRegressionBufferOperations
- **Purpose**: Correct buffer handling
- **Test Cases**: 5 scenarios (consecutive newlines, alternating content)
- **Coverage**: Buffer operation correctness

### 8. TestFormatRegressionTypeConversion
- **Purpose**: Safe type conversions
- **Test Cases**: 3 configurations with extreme values
- **Coverage**: Type conversion safety

### 9. TestFormatRegressionLoopConditions
- **Purpose**: Correct loop boundary conditions
- **Test Cases**: 5 scenarios (N-1, N, N+1 iterations)
- **Coverage**: Loop correctness

### 10. TestFormatRegressionRegexMatching
- **Purpose**: Correct regex matching at boundaries
- **Test Cases**: 5 scenarios (elements at document boundaries)
- **Coverage**: Boundary case regex matching

### 11. TestFormatRegressionContentPreservation
- **Purpose**: Verify content is preserved correctly
- **Test Cases**: 5 scenarios (text, special chars, unicode, code, numbers)
- **Coverage**: Content integrity

### 12. TestFormatRegressionEmptyContentBlocks
- **Purpose**: Handle empty blocks correctly
- **Test Cases**: 6 scenarios (empty code, table, list, quote, header)
- **Coverage**: Empty block handling

### 13. TestFormatRegressionComplexFormatCombinations
- **Purpose**: Complex format interaction safety
- **Test Cases**: 5 scenarios (code with formatting, nested features)
- **Coverage**: Format combination correctness

### 14. TestFormatRegressionIncrementalChanges
- **Purpose**: Handle content variations correctly
- **Test Cases**: 4 scenarios (original, with space, with newline, with content)
- **Coverage**: Incremental change handling

## Test Coverage Summary

### Input Types Tested:
- ✅ Empty strings
- ✅ Whitespace-only strings
- ✅ Very long inputs (1KB-1MB)
- ✅ Very long single lines (200-5000 chars)
- ✅ Unicode characters (emoji, CJK, RTL, combining)
- ✅ Special regex characters
- ✅ Malformed markup
- ✅ Mixed line endings
- ✅ Special characters and escape sequences

### Features Tested:
- ✅ All formatter configuration combinations
- ✅ Headers and nested headers
- ✅ Bold, italic, and mixed formatting
- ✅ Code blocks and inline code
- ✅ Lists (ordered, unordered, nested)
- ✅ Block quotes
- ✅ Tables (valid and malformed)
- ✅ Line wrapping
- ✅ Line numbers in code
- ✅ Color support
- ✅ Compact mode
- ✅ Custom prefixes

### Safety Verified:
- ✅ No panics on any input
- ✅ No nil pointer dereferences
- ✅ No index out of bounds errors
- ✅ No buffer overflows
- ✅ No regex DoS vulnerabilities
- ✅ No infinite loops
- ✅ Thread-safe concurrent usage
- ✅ Content preservation

### Performance:
- ✅ Handles memory-intensive inputs (10K+ lines)
- ✅ Handles stress with 100+ concurrent calls
- ✅ Handles high regex pattern density (1000+ patterns)
- ✅ Completes promptly (no hanging on complex inputs)

## Running Tests

### Run all formatter tests:
```bash
go test -v ./pkg/text
```

### Run only edge case tests:
```bash
go test -v ./pkg/text -run Edge
```

### Run only stress tests:
```bash
go test -v ./pkg/text -run Stress
```

### Run only regression tests:
```bash
go test -v ./pkg/text -run Regression
```

### Run with race detection:
```bash
go test -v -race ./pkg/text
```

### Run with coverage:
```bash
go test -v -cover ./pkg/text
```

## Test Execution Time

Expected execution time for all 44 tests: <5 seconds on modern hardware

## Maintenance

When modifying the formatter:
1. Ensure all 44 tests continue to pass
2. Add new edge case tests for any new features
3. Update TEST_COVERAGE.md if tests are added/removed
4. Run with race detection to check for concurrency issues
