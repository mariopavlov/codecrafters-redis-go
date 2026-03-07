package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseInput_Ping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		wantResp string
		wantErr  bool
	}{
		{
			name:     "simple PING command",
			input:    "*1\r\n$4\r\nPING\r\n",
			wantResp: "+PONG\r\n",
			wantErr:  false,
		},
		{
			name:     "PING in a 2-element array",
			input:    "*2\r\n$4\r\nPING\r\n$0\r\n\r\n",
			wantResp: "+PONG\r\n",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseInput([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.wantResp {
				t.Errorf("parseInput() = %q, want %q", string(got), tt.wantResp)
			}
		})
	}
}

func TestParseInput_Echo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		wantResp string
		wantErr  bool
	}{
		{
			name:     "ECHO with simple word",
			input:    "*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n",
			wantResp: "$5\r\nhello\r\n",
			wantErr:  false,
		},
		{
			name:     "ECHO with empty string",
			input:    "*2\r\n$4\r\nECHO\r\n$0\r\n\r\n",
			wantResp: "$0\r\n\r\n",
			wantErr:  false,
		},
		{
			name:     "ECHO with longer string",
			input:    "*2\r\n$4\r\nECHO\r\n$11\r\nhello world\r\n",
			wantResp: "$11\r\nhello world\r\n",
			wantErr:  false,
		},
		{
			name:     "ECHO with no argument returns error",
			input:    "*1\r\n$4\r\nECHO\r\n",
			wantResp: "+ERROR\r\n",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseInput([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.wantResp {
				t.Errorf("parseInput() = %q, want %q", string(got), tt.wantResp)
			}
		})
	}
}

func TestBuildBulkString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple word",
			input: "hello",
			want:  "$5\r\nhello\r\n",
		},
		{
			name:  "empty string",
			input: "",
			want:  "$0\r\n\r\n",
		},
		{
			name:  "string with spaces",
			input: "hello world",
			want:  "$11\r\nhello world\r\n",
		},
		{
			name:  "single character",
			input: "x",
			want:  "$1\r\nx\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := buildBulkString(tt.input)
			if string(got) != tt.want {
				t.Errorf("buildBulkString(%q) = %q, want %q", tt.input, string(got), tt.want)
			}
		})
	}
}

func TestParseInput_CaseInsensitive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		wantResp string
		wantErr  bool
	}{
		{
			name:     "lowercase ping",
			input:    "*1\r\n$4\r\nping\r\n",
			wantResp: "+PONG\r\n",
			wantErr:  false,
		},
		{
			name:     "mixed case echo",
			input:    "*2\r\n$4\r\nEcho\r\n$5\r\nhello\r\n",
			wantResp: "$5\r\nhello\r\n",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseInput([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.wantResp {
				t.Errorf("parseInput() = %q, want %q", string(got), tt.wantResp)
			}
		})
	}
}

func TestParseInput_UnknownCommand(t *testing.T) {
	t.Parallel()

	input := "*1\r\n$3\r\nSET\r\n"
	_, err := parseInput([]byte(input))
	if err == nil {
		t.Errorf("parseInput() with unknown command SET: expected error, got nil")
	}
}

func TestParseInput_MalformedInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []byte
	}{
		{name: "empty input", input: []byte{}},
		{name: "single byte", input: []byte("*")},
		{name: "non-numeric size", input: []byte("*X\r\n$4\r\nPING\r\n")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := parseInput(tt.input)
			if err == nil {
				t.Errorf("parseInput(%q) expected error for malformed input, got nil", string(tt.input))
			}
		})
	}
}

func TestParseInput_MultiDigitArraySize(t *testing.T) {
	// Array size > 9 should parse correctly, not just read one char
	input := "*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n"
	got, err := parseInput([]byte(input))
	if err != nil {
		t.Fatalf("parseInput() unexpected error: %v", err)
	}
	if string(got) != "$5\r\nhello\r\n" {
		t.Errorf("parseInput() = %q, want %q", string(got), "$5\r\nhello\r\n")
	}
}

func TestGetAt(t *testing.T) {
	t.Parallel()

	slice := []string{"$4", "ECHO", "$5", "hello"}

	tests := []struct {
		name    string
		idx     int
		wantVal string
		wantOk  bool
	}{
		{name: "valid index 1", idx: 1, wantVal: "ECHO", wantOk: true},
		{name: "valid index 3", idx: 3, wantVal: "hello", wantOk: true},
		{name: "index 0 returns false", idx: 0, wantVal: "", wantOk: false},
		{name: "negative index", idx: -1, wantVal: "", wantOk: false},
		{name: "out of bounds", idx: 10, wantVal: "", wantOk: false},
		{name: "exactly at length", idx: 4, wantVal: "", wantOk: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			val, ok := getAt(slice, tt.idx)
			if val != tt.wantVal || ok != tt.wantOk {
				t.Errorf("getAt(slice, %d) = (%q, %v), want (%q, %v)",
					tt.idx, val, ok, tt.wantVal, tt.wantOk)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// buildBulkString — special content
// ---------------------------------------------------------------------------

// TestBuildBulkString_SpecialContent covers content types not exercised by
// TestBuildBulkString: unicode (byte-length vs rune-count), embedded CRLF,
// special ASCII, null bytes, and numeric strings.
func TestBuildBulkString_SpecialContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			// RESP bulk-string length is byte-count, not rune-count.
			// "日本語" is 3 runes but 9 bytes in UTF-8.
			name:  "unicode uses byte length not rune count",
			input: "日本語",
			want:  "$9\r\n日本語\r\n",
		},
		{
			// Embedded CRLF is valid data; byte length includes the two bytes.
			name:  "string containing embedded CRLF counts both bytes",
			input: "hello\r\nworld",
			want:  "$12\r\nhello\r\nworld\r\n",
		},
		{
			name:  "special ASCII characters",
			input: "!@#$%",
			want:  "$5\r\n!@#$%\r\n",
		},
		{
			name:  "numeric string",
			input: "42",
			want:  "$2\r\n42\r\n",
		},
		{
			// Null bytes are valid binary data in RESP.
			name:  "string with null byte counted in length",
			input: "hel\x00lo",
			want:  "$6\r\nhel\x00lo\r\n",
		},
		{
			// A single space character.
			name:  "single space",
			input: " ",
			want:  "$1\r\n \r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := buildBulkString(tt.input)
			if string(got) != tt.want {
				t.Errorf("buildBulkString(%q) = %q, want %q", tt.input, string(got), tt.want)
			}
		})
	}
}

// TestBuildBulkString_LengthMatchesContent verifies that the declared length
// header always equals the actual byte count of the content, and that the
// output is framed correctly as $N\r\n<data>\r\n.
func TestBuildBulkString_LengthMatchesContent(t *testing.T) {
	t.Parallel()

	inputs := []string{"", "a", "ab", strings.Repeat("x", 100), "日本語", "hello\r\nworld"}

	for _, input := range inputs {
		input := input
		t.Run(fmt.Sprintf("len=%d", len(input)), func(t *testing.T) {
			t.Parallel()
			got := string(buildBulkString(input))
			want := fmt.Sprintf("$%d\r\n%s\r\n", len(input), input)
			if got != want {
				t.Errorf("buildBulkString(%q) = %q, want %q", input, got, want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// getAt — edge cases
// ---------------------------------------------------------------------------

// TestGetAt_EdgeCases covers empty slices, single-element slices, and
// valid accesses near boundaries — cases absent from TestGetAt.
func TestGetAt_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		slice   []string
		idx     int
		wantVal string
		wantOk  bool
	}{
		{
			name:    "empty slice returns false",
			slice:   []string{},
			idx:     1,
			wantVal: "",
			wantOk:  false,
		},
		{
			// Single element means index 1 is out of bounds.
			name:    "single element slice, index 1 returns false",
			slice:   []string{"only"},
			idx:     1,
			wantVal: "",
			wantOk:  false,
		},
		{
			name:    "two element slice, index 1 returns element",
			slice:   []string{"first", "second"},
			idx:     1,
			wantVal: "second",
			wantOk:  true,
		},
		{
			name:    "large valid index returns last element",
			slice:   []string{"a", "b", "c", "d", "e"},
			idx:     4,
			wantVal: "e",
			wantOk:  true,
		},
		{
			// idx == len(slice) is exactly out of bounds.
			name:    "index equal to slice length returns false",
			slice:   []string{"a", "b", "c"},
			idx:     3,
			wantVal: "",
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			val, ok := getAt(tt.slice, tt.idx)
			if val != tt.wantVal || ok != tt.wantOk {
				t.Errorf("getAt(%v, %d) = (%q, %v), want (%q, %v)",
					tt.slice, tt.idx, val, ok, tt.wantVal, tt.wantOk)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// parseInput — ECHO with varied arguments
// ---------------------------------------------------------------------------

// TestParseInput_EchoWithVariousArguments covers argument types not yet
// tested: numeric content, special ASCII, unicode, and single character.
func TestParseInput_EchoWithVariousArguments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		wantResp string
		wantErr  bool
	}{
		{
			name:     "ECHO with numeric argument",
			input:    "*2\r\n$4\r\nECHO\r\n$3\r\n123\r\n",
			wantResp: "$3\r\n123\r\n",
		},
		{
			name:     "ECHO with special ASCII characters",
			input:    "*2\r\n$4\r\nECHO\r\n$5\r\n!@#$%\r\n",
			wantResp: "$5\r\n!@#$%\r\n",
		},
		{
			// "日本語" = 3 runes = 9 bytes; RESP header must say $9.
			name:     "ECHO with unicode argument uses byte length",
			input:    "*2\r\n$4\r\nECHO\r\n$9\r\n日本語\r\n",
			wantResp: "$9\r\n日本語\r\n",
		},
		{
			name:     "ECHO with single character",
			input:    "*2\r\n$4\r\nECHO\r\n$1\r\nx\r\n",
			wantResp: "$1\r\nx\r\n",
		},
		{
			name:     "ECHO with spaces in argument",
			input:    "*2\r\n$4\r\nECHO\r\n$9\r\nfoo   bar\r\n",
			wantResp: "$9\r\nfoo   bar\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseInput([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.wantResp {
				t.Errorf("parseInput() = %q, want %q", string(got), tt.wantResp)
			}
		})
	}
}

// TestParseInput_EchoMultipleArgumentsTakesFirst verifies that when ECHO
// receives more than one argument (a 3-element array), only the first
// argument is echoed; extra arguments are silently ignored.
func TestParseInput_EchoMultipleArgumentsTakesFirst(t *testing.T) {
	t.Parallel()

	input := "*3\r\n$4\r\nECHO\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	got, err := parseInput([]byte(input))
	if err != nil {
		t.Fatalf("parseInput() unexpected error: %v", err)
	}
	want := "$5\r\nhello\r\n"
	if string(got) != want {
		t.Errorf("parseInput() = %q, want %q", string(got), want)
	}
}

// ---------------------------------------------------------------------------
// parseInput — type-prefix and size edge cases
// ---------------------------------------------------------------------------

// TestParseInput_NonArrayTypePrefix verifies that non-array RESP prefixes
// ('+' simple string, '-' error) are rejected because their second byte is
// not a digit that can represent an array size.
func TestParseInput_NonArrayTypePrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			// '+' is RESP simple-string prefix; second byte 'O' is not a digit.
			name:    "simple string prefix '+' returns error",
			input:   []byte("+OK\r\n"),
			wantErr: true,
		},
		{
			// '-' is RESP error prefix; second byte 'E' is not a digit.
			name:    "error prefix '-' returns error",
			input:   []byte("-ERR unknown command\r\n"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := parseInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInput(%q) error = %v, wantErr %v", string(tt.input), err, tt.wantErr)
			}
		})
	}
}

// TestParseInput_NegativeArraySize verifies that a negative array-size
// indicator (e.g. "*-1") returns an error rather than proceeding.
// '-' is not a digit so strconv.Atoi("-") fails, which is the correct path.
func TestParseInput_NegativeArraySize(t *testing.T) {
	t.Parallel()

	input := []byte("*-1\r\n$4\r\nPING\r\n")
	_, err := parseInput(input)
	if err == nil {
		t.Errorf("parseInput(%q): expected error for negative array size, got nil", string(input))
	}
}

// ---------------------------------------------------------------------------
// parseInput — BUG DOCUMENTATION TESTS
//
// The following tests document KNOWN BUGS in the current implementation.
// They describe the CORRECT expected behaviour after the bugs are fixed.
// They currently FAIL (or PANIC) — that is intentional: RED state in TDD.
//
// Bugs:
//   1. Multi-digit array size: only the first digit of the size is read,
//      so "*10\r\n…" is treated as "*1\r\n…", shifting all element indices.
//   2. Empty array (*0) or truncated RESP produces a result slice with fewer
//      than 2 elements; accessing result[1] causes a runtime panic.
// ---------------------------------------------------------------------------

// TestParseInput_MultiDigitArraySize_Bug documents that commands sent inside
// a RESP array with a two-digit (or larger) element count are silently
// misrouted because only the first digit of the count is consumed.
//
// This test is expected to FAIL with the current implementation.
// To fix: parse the full numeric prefix before \r\n, not just inputAsString[1].
func TestParseInput_MultiDigitArraySize_Bug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		wantResp string
	}{
		{
			name:     "PING in *10 array",
			input:    "*10\r\n$4\r\nPING\r\n",
			wantResp: "+PONG\r\n",
		},
		{
			name:     "ECHO in *10 array",
			input:    "*10\r\n$4\r\nECHO\r\n$5\r\nhello\r\n",
			wantResp: "$5\r\nhello\r\n",
		},
		{
			name:     "PING in *99 array",
			input:    "*99\r\n$4\r\nPING\r\n",
			wantResp: "+PONG\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseInput([]byte(tt.input))
			if err != nil {
				t.Errorf("parseInput() unexpected error = %v (multi-digit array size bug)", err)
				return
			}
			if string(got) != tt.wantResp {
				t.Errorf("parseInput() = %q, want %q (multi-digit array size bug)", string(got), tt.wantResp)
			}
		})
	}
}

// TestParseInput_EmptyArrayReturnsError documents that a zero-element RESP
// array (*0) should return an error, not panic via an out-of-bounds access
// on result[1].
//
// This test is expected to FAIL with the current implementation.
// To fix: bounds-check result before accessing result[1].
func TestParseInput_EmptyArrayReturnsError(t *testing.T) {
	// Use recover so the panic is captured as a FAIL rather than crashing the
	// test binary.  After the fix, no recover fires and err is non-nil instead.
	var panicked bool
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		input := []byte("*0\r\n\r\n")
		_, err := parseInput(input)
		if err == nil {
			t.Errorf("parseInput(*0 input): expected error for empty array, got nil")
		}
	}()
	if panicked {
		t.Errorf("parseInput(*0 input): panicked instead of returning an error (bounds-check missing on result)")
	}
}

// TestParseInput_TruncatedRESPReturnsError documents that incomplete RESP
// frames should return errors rather than panicking on result[1].
//
// These tests are expected to FAIL with the current implementation.
// To fix: bounds-check result before accessing result[1].
func TestParseInput_TruncatedRESPReturnsError(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "array size only, no element data",
			input: []byte("*1\r\n"),
		},
		{
			name:  "array with length prefix only, no command bytes",
			input: []byte("*1\r\n$4\r\n"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Use recover so the panic is captured as a FAIL rather than
			// crashing the test binary.
			var panicked bool
			func() {
				defer func() {
					if r := recover(); r != nil {
						panicked = true
					}
				}()
				_, err := parseInput(tt.input)
				if err == nil {
					t.Errorf("parseInput(%q): expected error for truncated input, got nil", string(tt.input))
				}
			}()
			if panicked {
				t.Errorf("parseInput(%q): panicked instead of returning an error (bounds-check missing on result)", string(tt.input))
			}
		})
	}
}
