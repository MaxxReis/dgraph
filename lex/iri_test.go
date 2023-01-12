package lex

import (
	"testing"
)

type test struct {
	data int32
	want bool
}

type testLexer struct {
	// NOTE: Using a text scanner wouldn't work because it's designed for parsing
	// Golang. It won't keep track of Start Position, or allow us to retrieve
	// slice from [Start:Pos]. Better to just use normal string.
	Input      string // string being scanned.
	Start      int    // Start Position of this item.
	Pos        int    // current Position of this item.
	Width      int    // Width of last rune read from input.
	widthStack []*RuneWidth
	items      []Item  // channel of scanned items.
	Depth      int     // nesting of {}
	BlockDepth int     // nesting of blocks (e.g. mutation block inside upsert block)
	ArgDepth   int     // nesting of ()
	Mode       StateFn // Default state to go back to after reading a token.
	Line       int     // the current line number corresponding to Start
	Column     int     // the current column number corresponding to Start
}

type testChars struct {
	r rune
	l *Lexer
}

// func TestIsHex(t *testing.T) {
// 	got := '0'
// 	r := isHex(got)
// 	if r == false {
// 		t.Error("Expected: a character beteween '0' and '9", "got: ", r)
// 	}

// }

func Test_isHex(t *testing.T) {
	tests := []test{
		{data: '0', want: true},
		{data: '1', want: true},
		{data: '2', want: true},
		{data: '3', want: true},
		{data: '4', want: true},
		{data: '5', want: true},
		{data: '6', want: true},
		{data: '7', want: true},
		{data: '8', want: true},
		{data: '9', want: true},
		{data: 'a', want: true},
		{data: 'b', want: true},
		{data: 'c', want: true},
		{data: 'd', want: true},
		{data: 'e', want: true},
		{data: 'f', want: true},
		{data: 'A', want: true},
		{data: 'B', want: true},
		{data: 'C', want: true},
		{data: 'D', want: true},
		{data: 'E', want: true},
		{data: 'F', want: true},
	}
	for _, value := range tests {
		got := isHex(value.data)
		if got != value.want {
			t.Error("Expected: ", value.want, "got: ", got)
		}
	}
}

func Test_isHexTableFalse(t *testing.T) {
	tests := []test{
		{data: 'G', want: false},
		{data: 'g', want: false},
		{data: 'H', want: false},
		{data: 'h', want: false},
	}
	for _, value := range tests {
		got := isHex(value.data)
		if got != value.want {
			t.Error("Expected: ", value.want, "got: ", got)
		}
	}
}

// Test IsIRIRefChar ok (but need review)
func Test_isIRIRefChar(t *testing.T) {
	type args struct {
		r rune
		l *Lexer
	}
	tests := []args{
		{r: 'a', l: &Lexer{}},
		{r: 'u', l: &Lexer{}},
		{r: 'U', l: &Lexer{}},
	}
	for _, tt := range tests {
		got := isIRIRefChar(tt.r, tt.l)
		if got == false {
			t.Errorf("isIRIRefChar() = %v", got)
		}
	}
}

func TestIRIRef(t *testing.T) {
	type args struct {
		l    *Lexer
		styp ItemType
	}
	tests := []args{
		{l: &Lexer{
			Input:      ">",
			Start:      0,
			Pos:        0,
			Width:      0,
			widthStack: []*RuneWidth{},
			items:      []Item{},
			Depth:      0,
			BlockDepth: 0,
			ArgDepth:   0,
			Line:       0,
			Column:     0,
		}, styp: 5},
	}
	for _, tt := range tests {
		got := IRIRef(tt.l, tt.styp)
		if got != nil {
			t.Error("Expected: ", nil, " got: ", got)
		}
		// print("Got: ", got)
	}
}

// func TestHasXChars(t *testing.T) {
// 	type chars struct {
// 		r rune
// 		l *Lexer
// 	}

// 	tests := []chars{
// 		{r: 'x',
// 			l: &Lexer{
// 				Input:      "ex machine",
// 				Start:      0,
// 				Pos:        0,
// 				Width:      0,
// 				widthStack: []*RuneWidth{},
// 				items:      []Item{},
// 				Depth:      0,
// 				BlockDepth: 0,
// 				ArgDepth:   0,
// 				Line:       2,
// 				Column:     0,
// 			}},
// 	}
// 	for _, test := range tests {
// 		got := HasXChars(test.r, test.l)
// 		if got == false {
// 			t.Error("Expected: ", true, " got: ", got)
// 		}
// 	}
// }

// func TestHasUChars(t *testing.T) {
// 	type args struct {
// 		r rune
// 		l *Lexer
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want bool
// 	}{
// 		{
// 			name: "TestHasUChars 1",
// 			args: args{r: 'u' + ' ' + 'U',
// 				l: &Lexer{
// 					Input:      "u" + " " + "U",
// 					Start:      0,
// 					Pos:        0,
// 					Width:      0,
// 					widthStack: []*RuneWidth{},
// 					items:      []Item{},
// 					Depth:      0,
// 					BlockDepth: 0,
// 					ArgDepth:   0,
// 					Line:       0,
// 					Column:     0,
// 				}},
// 			want: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := HasUChars(tt.args.r, tt.args.l); got != tt.want {
// 				t.Errorf("HasUChars() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
