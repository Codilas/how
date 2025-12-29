# Formatter Edge Cases and Stress Tests

This document describes the comprehensive edge case and stress tests added to ensure the formatter's robustness.

## Test Files

### 1. `formatter_edge_cases_test.go`
Core edge case tests covering basic boundary conditions and malformed inputs.

#### Tests Included:
- **TestFormatEmptyString**: Tests handling of completely empty input with different configs
- **TestFormatWhitespaceOnly**: Tests input containing only whitespace (spaces, tabs, newlines)
- **TestFormatVeryLongInput**: Tests formatter with extremely large inputs (1KB to 1MB)
- **TestFormatVeryLongSingleLine**: Tests handling of single lines exceeding normal width
- **TestFormatDeeplyNestedMarkdown**: Tests nested markdown structures (lists, headers, inline formatting)
- **TestFormatMalformedTables**: Tests malformed table syntax (missing pipes, uneven columns, empty cells)
- **TestFormatMixedLineEndings**: Tests various line ending combinations (LF, CRLF, CR, mixed)
- **TestFormatUnicodeCharacters**: Tests handling of Unicode including emoji, CJK, Arabic, Hebrew, Cyrillic
- **TestFormatSpecialRegexCharacters**: Tests input containing regex metacharacters and special patterns
- **TestFormatExtremeConfigValues**: Tests extreme configuration values (very small/large widths, indents)
- **TestFormatComplexCodeBlocks**: Tests code block edge cases (unclosed, nested backticks, many lines)
- **TestFormatStructuredCommands**: Tests removal of structured command tags
- **TestFormatMalformedMarkdown**: Tests malformed markdown syntax
- **TestFormatAllConfigCombinations**: Tests various combinations of formatter options
- **TestFormatStressLongContent**: Tests complex documents with mixed content types
- **TestFormatBoundaryValues**: Tests exact boundary conditions (line width boundaries)
- **TestFormatRobustnessWithNoOptions**: Tests formatter with empty config using defaults

### 2. `formatter_stress_test.go`
Stress tests covering performance, concurrency, and complex scenarios.

#### Tests Included:
- **TestFormatMemoryIntensive**: Tests with memory-intensive inputs (10,000+ lines, large items)
- **TestFormatConcurrentStress**: Tests concurrent formatter usage (100+ goroutines)
- **TestFormatEdgeCaseMarkdownCombinations**: Tests complex markdown patterns
- **TestFormatLineWrappingEdgeCases**: Tests line wrapping with various edge cases
- **TestFormatTableEdgeCases**: Tests advanced table scenarios (wide columns, many columns, unicode)
- **TestFormatIndentationEdgeCases**: Tests indentation handling (deep nesting, mixed styles)
- **TestFormatColoringEdgeCases**: Tests coloring with various edge cases
- **TestFormatCompactModeEdgeCases**: Tests compact mode behavior
- **TestFormatRegexPerformance**: Tests regex performance with many patterns
- **TestFormatCommentPrefixEdgeCases**: Tests different comment prefix configurations
- **TestFormatNumberedLists**: Tests numbered list variations
- **TestFormatShowLineNumbersEdgeCases**: Tests line numbering in code blocks
- **TestFormatWhitespaceHandling**: Tests various whitespace scenarios

### 3. `formatter_regression_test.go`
Regression tests focusing on potential bugs and off-by-one errors.

#### Tests Included:
- **TestFormatNilInputHandling**: Tests nil/empty input handling
- **TestFormatRegexCatastrophicBacktracking**: Tests prevention of regex DoS attacks
- **TestFormatRegressionOffByOne**: Tests potential off-by-one errors
- **TestFormatRegressionSliceIndexing**: Tests slice indexing boundary conditions
- **TestFormatRegressionStringMethods**: Tests string operation edge cases
- **TestFormatRegressionMapAccess**: Tests safe map access
- **TestFormatRegressionBufferOperations**: Tests buffer handling
- **TestFormatRegressionTypeConversion**: Tests type conversion safety
- **TestFormatRegressionLoopConditions**: Tests loop boundary conditions
- **TestFormatRegressionRegexMatching**: Tests regex matching at document boundaries
- **TestFormatRegressionContentPreservation**: Tests that content is correctly preserved
- **TestFormatRegressionEmptyContentBlocks**: Tests empty block handling
- **TestFormatRegressionComplexFormatCombinations**: Tests complex format combinations
- **TestFormatRegressionIncrementalChanges**: Tests incremental content changes

## Test Coverage Summary

### Edge Cases Covered:
1. **Empty/Whitespace Input**: Empty strings, spaces, tabs, newlines, mixed whitespace
2. **Size Extremes**: 1KB to 1MB+ inputs, very long single lines
3. **Encoding**: Unicode (emoji, CJK, Arabic, Hebrew, Cyrillic), combining chars, RTL text
4. **Markdown**: Nested structures, malformed syntax, unclosed tags, mixed formatting
5. **Tables**: Malformed tables, missing pipes, uneven columns, empty cells, wide content
6. **Code Blocks**: Unclosed blocks, nested backticks, many lines, special regex chars
7. **Line Endings**: Unix (LF), Windows (CRLF), Classic Mac (CR), mixed
8. **Configuration**: Extreme values, zero/negative-like values, all feature combinations
9. **Performance**: Regex patterns, concurrent usage, memory-intensive content
10. **Regex Safety**: Catastrophic backtracking prevention, DoS protection

### Configuration Variations Tested:
- All combinations of boolean flags (UseColors, UseBoxes, UseBullets, etc.)
- Extreme LineWidth values (1 to 2,147,483,647)
- Extreme IndentSize values (0 to 999+)
- Various CommentPrefix options (empty, long, special characters, unicode)
- All three preset configurations (Default, Colored, Compact)

### Boundary Conditions:
- Exactly at line width boundaries
- Off-by-one scenarios (N-1, N, N+1)
- Empty content blocks
- Single character inputs
- Maximum/minimum config values

## Running the Tests

Run all tests:
```bash
make test
```

Run with race detection:
```bash
make test-race
```

Run specific test file:
```bash
go test -v ./pkg/text -run TestFormat
```

Run specific test:
```bash
go test -v ./pkg/text -run TestFormatEmptyString
```

## Expected Behavior

All tests expect the formatter to:
1. **Never panic**: All edge cases should be handled gracefully
2. **Produce output**: For non-empty input, should generally produce output
3. **Preserve content**: Important content should be preserved in output
4. **Handle gracefully**: Malformed input should not cause errors
5. **Complete promptly**: No regex DoS or infinite loops
6. **Be safe**: No buffer overflows, nil pointer dereferences, or index out of bounds

## Notes

- Tests use table-driven approach consistent with Go testing best practices
- Tests focus on robustness rather than specific output formats
- Some tests verify non-panic behavior rather than exact output
- Concurrent stress tests verify thread safety
- Performance tests use goroutines to detect hanging/slowness

## Future Improvements

Potential areas for additional testing:
- Benchmark tests for performance regression detection
- Fuzz testing for discovering unknown edge cases
- Visual regression tests for formatting consistency
- Property-based testing for invariant verification
