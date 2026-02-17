package main

import (
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
