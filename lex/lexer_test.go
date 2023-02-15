package lex

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type fields struct {
	Typ    ItemType
	Val    string
	line   int
	column int
}

type args struct {
	format string
	args   []interface{}
}

func TestItem_Errorf(t *testing.T) {
	testCase := []*Item{{0, "test", 0, 0}}
	argsCase := args{"", []interface{}{}}
	for _, test := range testCase {
		if err := test.Errorf(argsCase.format, argsCase.args...); err == nil {
			require.Equal(t, test.Errorf("line %v column %v", test.line, test.column), err)
		}
	}
}

func TestItem_String(t *testing.T) {
	var result string
	result = "EOF"
	testCase := []*Item{{0, "test", 0, 0}, {5, "test", 0, 0}}
	for _, test := range testCase {
		if err := test.String(); err != result {
			fmt.Sprintf("lex.Item [%v] %q at %d:%d", test.Typ, test.Val, test.line, test.column)
			// test.Errorf("lex.Item [%v] %q at %d:%d", test.Typ, test.Val, test.line, test.column)
		}
		// require.Equal(t, "EOF", result)
	}
}

func TestLexer_NewIterator(t *testing.T) {
	testCase := []*ItemIterator{{&Lexer{Input: "test", Start: 1, Pos: 1}, 1}}
	for _, test := range testCase {
		result := test.l.NewIterator()
		if result == nil {
			require.Equal(t, test.Errorf(""), result)
		}
		require.Equal(t, test.l.NewIterator(), result)
	}
}

func TestItemIterator_Errorf(t *testing.T) {
	testCase := []*ItemIterator{{&Lexer{items: []Item{{0, "test", 0, 0}, {0, "test", 0, 0}}}, 2}}
	argsCase := args{"", []interface{}{}}
	for _, test := range testCase {
		err := test.Errorf(argsCase.format, argsCase.args...)
		if err == nil {
			test.Errorf("Unknown error ItemIterator")
		}
		test.Errorf("Error %v", err)
	}
}

func TestItemIterator_Next(t *testing.T) {
	testCase := []*ItemIterator{{&Lexer{Input: "", items: []Item{{5, "test", 0, 0}, {5, "test", 0, 0}}}, -1}}
	for _, test := range testCase {
		if result := test.Next(); result == false {
			test.Errorf("Error index")
		}
	}
}

func TestItemIterator_Item(t *testing.T) {
	testCase := []*ItemIterator{
		{&Lexer{items: []Item{{0, "test", 0, 0}, {0, "test", 0, 0}}}, -1},
		{&Lexer{items: []Item{{0, "test", 0, 0}, {0, "test", 0, 0}}}, 0},
	}
	for _, test := range testCase {
		if result := test.Item(); result.line == -1 {
			errors.New("out-of-range item")
		}
	}
}

func TestItemIterator_Prev(t *testing.T) {
	testCase := []*ItemIterator{{&Lexer{}, -1}, {&Lexer{}, 1}}
	for _, test := range testCase {
		if result := test.Prev(); result == false {
			errors.New("Error: index negative")
		}
	}
}

func TestItemIterator_Restore(t *testing.T) {
	testCase := []*ItemIterator{
		{&Lexer{items: []Item{{0, "test", 0, 0}, {0, "test", 0, 0}}}, 0},
		{&Lexer{items: []Item{{0, "test", 0, 0}, {0, "test", 0, 0}}}, 1},
	}
	for _, test := range testCase {
		test.Restore(1)
	}
}

func TestItemIterator_Save(t *testing.T) {
	testCase := []*ItemIterator{
		{&Lexer{items: []Item{{0, "test", 0, 0}, {0, "test", 0, 0}}}, 0},
		{&Lexer{items: []Item{{0, "test", 0, 0}, {0, "test", 0, 0}}}, -2},
	}
	for _, test := range testCase {
		if result := test.Save(); result < 0 {
			errors.New("Error: Index invalid")
		}
	}
}

func TestItemIterator_Peek(t *testing.T) {
	testCase := []*ItemIterator{
		{&Lexer{items: []Item{{0, "test", 0, 0}, {0, "test", 0, 0}, {0, "test", 0, 0}}}, 0},
		{&Lexer{items: []Item{{0, "test", 0, 0}, {0, "test", 0, 0}}}, 1},
	}
	for _, test := range testCase {
		result, err := test.Peek(1)
		if (result == nil) && (err != nil) {
			require.Error(t, err, "Out of range for peek")
		}
	}
}

func TestItemIterator_PeekOne(t *testing.T) {
	testCase := []*ItemIterator{
		{&Lexer{items: []Item{{0, "test", 0, 0}, {0, "test", 0, 0}, {0, "test", 0, 0}}}, 1},
		{&Lexer{items: []Item{{0, "test", 0, 0}, {0, "test", 0, 0}}}, 1},
	}
	for _, test := range testCase {
		result, err := test.PeekOne()
		if err == false {
			// fmt.Sprintf("line %v column %v", result.line, result.column)
			require.NotNil(t, result)

		}
	}
}

func TestLexer_Reset(t *testing.T) {
	testCase := []*Lexer{{Input: "test"}}
	for _, test := range testCase {
		test.Reset("reset")
	}
}

func TestLexer_ValidateResult(t *testing.T) {
	testCase := []*Lexer{
		{items: []Item{{0, "test", 0, 0}}},
		{items: []Item{{1, "test", 0, 0}}},
	}
	for _, test := range testCase {
		err := test.ValidateResult()
		if err != nil {
			// errors.New("Error to validate result")
			require.Error(t, err)
		}
	}
}

// review
func TestLexer_Run(t *testing.T) {
	testCase := []*Lexer{
		{Mode: func(l *Lexer) StateFn { return nil }},
		{Mode: nil},
	}
	for _, test := range testCase {
		result := test.Run(test.Mode)
		if result == nil {
			require.Empty(t, result)
		}
	}
}

// review
func TestLexer_Errorf(t *testing.T) {
	testCase := []*Lexer{
		{Input: "test1", items: nil},
		{Input: "test2", items: []Item{{1, "test", 0, 0}}},
	}
	argsCase := args{"", []interface{}{}}
	for _, test := range testCase {
		result := test.Errorf(argsCase.format, argsCase.args...)
		if result != nil {
			errors.New("")
		}
	}
}

// check how to pass Item.Typ dynamically
func TestLexer_Emit(t *testing.T) {
	testCase := []*Lexer{
		{Start: 1, Pos: 0, items: []Item{{5, "test", 0, 0}}},
		{Input: "test", Start: 0, Pos: 1, items: []Item{{5, "test", 0, 0}}},
	}
	for _, test := range testCase {
		test.Emit(5)
	}
}

func TestLexer_pushWidth(t *testing.T) {
	testCase := []*Lexer{
		{widthStack: []*RuneWidth{{width: 0, count: 0}}},
		{widthStack: nil},
	}
	for _, test := range testCase {
		test.pushWidth(test.Width)
	}
}

func TestLexer_Next(t *testing.T) {
	testCase := []*Lexer{{Input: "test", Pos: 4}, {Input: "test", Pos: 1}}
	for _, test := range testCase {
		result := test.Next()
		if result == EOF {
			errors.New("EOF")
		}
	}
}

// func TestLexer_Backup(t *testing.T) {
// 	testCase := []*Lexer{
// 		{widthStack: []*RuneWidth{{width: 0, count: 1}}},
// 		{widthStack: []*RuneWidth{{width: 0, count: 0}}},
// 	}
// 	for err, test := range testCase {
// 		test.Backup()
// 	}
// }

// func TestLexer_Backup(t *testing.T) {
// 	type fields struct {
// 		Input      string
// 		Start      int
// 		Pos        int
// 		Width      int
// 		widthStack []*RuneWidth
// 		items      []Item
// 		Depth      int
// 		BlockDepth int
// 		ArgDepth   int
// 		Mode       StateFn
// 		Line       int
// 		Column     int
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 	}{
// 		{
// 			name: "TestLexer_Backup 1",
// 			fields: fields{
// 				Input: "Test",
// 				Start: 0,
// 				Pos:   3,
// 				Width: 0,
// 				widthStack: []*RuneWidth{
// 					{width: 1, count: 1},
// 				},
// 				items:      []Item{},
// 				Depth:      0,
// 				BlockDepth: 0,
// 				ArgDepth:   0,
// 				Line:       0,
// 				Column:     0,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			l := &Lexer{
// 				Input:      tt.fields.Input,
// 				Start:      tt.fields.Start,
// 				Pos:        tt.fields.Pos,
// 				Width:      tt.fields.Width,
// 				widthStack: tt.fields.widthStack,
// 				items:      tt.fields.items,
// 				Depth:      tt.fields.Depth,
// 				BlockDepth: tt.fields.BlockDepth,
// 				ArgDepth:   tt.fields.ArgDepth,
// 				Mode:       tt.fields.Mode,
// 				Line:       tt.fields.Line,
// 				Column:     tt.fields.Column,
// 			}
// 			l.Backup()
// 		})
// 	}
// }

// func TestLexer_Peek(t *testing.T) {
// 	type fields struct {
// 		Input      string
// 		Start      int
// 		Pos        int
// 		Width      int
// 		widthStack []*RuneWidth
// 		items      []Item
// 		Depth      int
// 		BlockDepth int
// 		ArgDepth   int
// 		Mode       StateFn
// 		Line       int
// 		Column     int
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   rune
// 	}{
// 		{
// 			name: "TestLexer_Peek 1",
// 			fields: fields{
// 				Input:      "test",
// 				Start:      0,
// 				Pos:        0,
// 				Width:      0,
// 				widthStack: []*RuneWidth{},
// 				items:      []Item{},
// 				Depth:      0,
// 				BlockDepth: 0,
// 				ArgDepth:   0,
// 				Line:       0,
// 				Column:     0,
// 			},
// 			want: 't',
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			l := &Lexer{
// 				Input:      tt.fields.Input,
// 				Start:      tt.fields.Start,
// 				Pos:        tt.fields.Pos,
// 				Width:      tt.fields.Width,
// 				widthStack: tt.fields.widthStack,
// 				items:      tt.fields.items,
// 				Depth:      tt.fields.Depth,
// 				BlockDepth: tt.fields.BlockDepth,
// 				ArgDepth:   tt.fields.ArgDepth,
// 				Mode:       tt.fields.Mode,
// 				Line:       tt.fields.Line,
// 				Column:     tt.fields.Column,
// 			}
// 			if got := l.Peek(); got != tt.want {
// 				t.Errorf("Lexer.Peek() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestLexer_PeekTwo(t *testing.T) {
	type fields struct {
		Input      string
		Start      int
		Pos        int
		Width      int
		widthStack []*RuneWidth
		items      []Item
		Depth      int
		BlockDepth int
		ArgDepth   int
		Mode       StateFn
		Line       int
		Column     int
	}
	tests := []struct {
		name   string
		fields fields
		want   []rune
	}{
		{
			name: "TestLexer_PeekTwo 1",
			fields: fields{
				Input:      "test",
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
			},
			want: []rune{'t', 'e'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				Input:      tt.fields.Input,
				Start:      tt.fields.Start,
				Pos:        tt.fields.Pos,
				Width:      tt.fields.Width,
				widthStack: tt.fields.widthStack,
				items:      tt.fields.items,
				Depth:      tt.fields.Depth,
				BlockDepth: tt.fields.BlockDepth,
				ArgDepth:   tt.fields.ArgDepth,
				Mode:       tt.fields.Mode,
				Line:       tt.fields.Line,
				Column:     tt.fields.Column,
			}
			if got := l.PeekTwo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexer.PeekTwo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_PeekTwoEOF(t *testing.T) {
	type fields struct {
		Input      string
		Start      int
		Pos        int
		Width      int
		widthStack []*RuneWidth
		items      []Item
		Depth      int
		BlockDepth int
		ArgDepth   int
		Mode       StateFn
		Line       int
		Column     int
	}
	tests := []struct {
		name   string
		fields fields
		want   []rune
	}{
		{
			name: "TestLexer_PeekTwo 1",
			fields: fields{
				Input:      "",
				Start:      0,
				Pos:        0,
				Width:      0,
				widthStack: []*RuneWidth{},
				items: []Item{
					{
						Typ:    -1,
						Val:    "",
						line:   0,
						column: 0,
					},
				},
				Depth:      0,
				BlockDepth: 0,
				ArgDepth:   0,
				Line:       0,
				Column:     0,
			},
			want: []rune{'e', 'e'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				Input:      tt.fields.Input,
				Start:      tt.fields.Start,
				Pos:        tt.fields.Pos,
				Width:      tt.fields.Width,
				widthStack: tt.fields.widthStack,
				items:      tt.fields.items,
				Depth:      tt.fields.Depth,
				BlockDepth: tt.fields.BlockDepth,
				ArgDepth:   tt.fields.ArgDepth,
				Mode:       tt.fields.Mode,
				Line:       tt.fields.Line,
				Column:     tt.fields.Column,
			}
			if got := l.PeekTwo(); reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexer.PeekTwo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_moveStartToPos(t *testing.T) {
	type fields struct {
		Input      string
		Start      int
		Pos        int
		Width      int
		widthStack []*RuneWidth
		items      []Item
		Depth      int
		BlockDepth int
		ArgDepth   int
		Mode       StateFn
		Line       int
		Column     int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "TestLexer_moveStartToPos 1",
			fields: fields{
				Input:      "This a simple test",
				Start:      0,
				Pos:        2,
				Width:      0,
				widthStack: []*RuneWidth{},
				items:      []Item{},
				Depth:      0,
				BlockDepth: 0,
				ArgDepth:   0,
				Line:       0,
				Column:     0,
			},
		},
		{
			name: "TestLexer_moveStartToPos 2",
			fields: fields{
				Input:      " ",
				Start:      0,
				Pos:        1,
				Width:      0,
				widthStack: []*RuneWidth{},
				items:      []Item{},
				Depth:      0,
				BlockDepth: 0,
				ArgDepth:   0,
				Line:       0,
				Column:     0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				Input:      tt.fields.Input,
				Start:      tt.fields.Start,
				Pos:        tt.fields.Pos,
				Width:      tt.fields.Width,
				widthStack: tt.fields.widthStack,
				items:      tt.fields.items,
				Depth:      tt.fields.Depth,
				BlockDepth: tt.fields.BlockDepth,
				ArgDepth:   tt.fields.ArgDepth,
				Mode:       tt.fields.Mode,
				Line:       tt.fields.Line,
				Column:     tt.fields.Column,
			}
			l.moveStartToPos()
		})
	}
}

func TestLexer_moveStartToPosEOL(t *testing.T) {
	type fields struct {
		Input      string
		Start      int
		Pos        int
		Width      int
		widthStack []*RuneWidth
		items      []Item
		Depth      int
		BlockDepth int
		ArgDepth   int
		Mode       StateFn
		Line       int
		Column     int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "TestLexer_moveStartToPos 1",
			fields: fields{
				Input:      "t\u000A",
				Start:      0,
				Pos:        2,
				Width:      0,
				widthStack: []*RuneWidth{},
				items:      []Item{},
				Depth:      0,
				BlockDepth: 0,
				ArgDepth:   0,
				Line:       0,
				Column:     0,
			},
		},
		{
			name: "TestLexer_moveStartToPos 1",
			fields: fields{
				Input:      "t\u000D",
				Start:      0,
				Pos:        2,
				Width:      0,
				widthStack: []*RuneWidth{},
				items:      []Item{},
				Depth:      0,
				BlockDepth: 0,
				ArgDepth:   0,
				Line:       0,
				Column:     0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				Input:      tt.fields.Input,
				Start:      tt.fields.Start,
				Pos:        tt.fields.Pos,
				Width:      tt.fields.Width,
				widthStack: tt.fields.widthStack,
				items:      tt.fields.items,
				Depth:      tt.fields.Depth,
				BlockDepth: tt.fields.BlockDepth,
				ArgDepth:   tt.fields.ArgDepth,
				Mode:       tt.fields.Mode,
				Line:       tt.fields.Line,
				Column:     tt.fields.Column,
			}
			l.moveStartToPos()
		})
	}
}

func TestLexer_Ignore(t *testing.T) {
	type fields struct {
		Input      string
		Start      int
		Pos        int
		Width      int
		widthStack []*RuneWidth
		items      []Item
		Depth      int
		BlockDepth int
		ArgDepth   int
		Mode       StateFn
		Line       int
		Column     int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "TestLexer_Ignore 1",
			fields: fields{
				Input:      "Test",
				Start:      0,
				Pos:        1,
				Width:      0,
				widthStack: []*RuneWidth{},
				items:      []Item{},
				Depth:      0,
				BlockDepth: 0,
				ArgDepth:   0,
				Line:       0,
				Column:     0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				Input:      tt.fields.Input,
				Start:      tt.fields.Start,
				Pos:        tt.fields.Pos,
				Width:      tt.fields.Width,
				widthStack: tt.fields.widthStack,
				items:      tt.fields.items,
				Depth:      tt.fields.Depth,
				BlockDepth: tt.fields.BlockDepth,
				ArgDepth:   tt.fields.ArgDepth,
				Mode:       tt.fields.Mode,
				Line:       tt.fields.Line,
				Column:     tt.fields.Column,
			}
			l.Ignore()
		})
	}
}

func TestLexer_AcceptRun(t *testing.T) {
	type fields struct {
		Input      string
		Start      int
		Pos        int
		Width      int
		widthStack []*RuneWidth
		items      []Item
		Depth      int
		BlockDepth int
		ArgDepth   int
		Mode       StateFn
		Line       int
		Column     int
	}
	type args struct {
		c CheckRune
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantLastr  rune
		wantValidr bool
	}{
		{
			name: "TestLexer_AcceptRun 1",
			fields: fields{
				Input:      "test",
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
			},
			args: args{
				c: func(r rune) bool {
					if r == 't' {
						return true
					}
					return false
				},
			},
			wantLastr:  't',
			wantValidr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				Input:      tt.fields.Input,
				Start:      tt.fields.Start,
				Pos:        tt.fields.Pos,
				Width:      tt.fields.Width,
				widthStack: tt.fields.widthStack,
				items:      tt.fields.items,
				Depth:      tt.fields.Depth,
				BlockDepth: tt.fields.BlockDepth,
				ArgDepth:   tt.fields.ArgDepth,
				Mode:       tt.fields.Mode,
				Line:       tt.fields.Line,
				Column:     tt.fields.Column,
			}
			gotLastr, gotValidr := l.AcceptRun(tt.args.c)
			if gotLastr != tt.wantLastr {
				t.Errorf("Lexer.AcceptRun() gotLastr = %v, want %v", gotLastr, tt.wantLastr)
			}
			if gotValidr != tt.wantValidr {
				t.Errorf("Lexer.AcceptRun() gotValidr = %v, want %v", gotValidr, tt.wantValidr)
			}
		})
	}
}

func TestLexer_AcceptRunRec(t *testing.T) {
	type fields struct {
		Input      string
		Start      int
		Pos        int
		Width      int
		widthStack []*RuneWidth
		items      []Item
		Depth      int
		BlockDepth int
		ArgDepth   int
		Mode       StateFn
		Line       int
		Column     int
	}
	type args struct {
		c CheckRuneRec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestLexer_AcceptRunRec 1",
			fields: fields{
				Input:      "test",
				Start:      0,
				Pos:        1,
				Width:      4,
				widthStack: []*RuneWidth{},
				items: []Item{
					{
						Typ:    5,
						Val:    "t",
						line:   0,
						column: 0,
					},
					{
						Typ:    5,
						Val:    "e",
						line:   0,
						column: 0,
					},
				},
				Depth:      0,
				BlockDepth: 0,
				ArgDepth:   0,
				Line:       0,
				Column:     0,
			},
			args: args{
				c: func(r rune, l *Lexer) bool {
					if r == 'e' {
						return true
					}
					return false
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				Input:      tt.fields.Input,
				Start:      tt.fields.Start,
				Pos:        tt.fields.Pos,
				Width:      tt.fields.Width,
				widthStack: tt.fields.widthStack,
				items:      tt.fields.items,
				Depth:      tt.fields.Depth,
				BlockDepth: tt.fields.BlockDepth,
				ArgDepth:   tt.fields.ArgDepth,
				Mode:       tt.fields.Mode,
				Line:       tt.fields.Line,
				Column:     tt.fields.Column,
			}
			l.AcceptRunRec(tt.args.c)
		})
	}
}

func TestLexer_AcceptUntil(t *testing.T) {
	type fields struct {
		Input      string
		Start      int
		Pos        int
		Width      int
		widthStack []*RuneWidth
		items      []Item
		Depth      int
		BlockDepth int
		ArgDepth   int
		Mode       StateFn
		Line       int
		Column     int
	}
	type args struct {
		c CheckRune
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestLexer_AcceptUntil 1",
			fields: fields{
				Input:      "test",
				Start:      0,
				Pos:        1,
				Width:      4,
				widthStack: []*RuneWidth{},
				items: []Item{
					{
						Typ:    5,
						Val:    "t",
						line:   0,
						column: 0,
					},
					{
						Typ:    5,
						Val:    "e",
						line:   0,
						column: 0,
					},
				},
				Depth:      0,
				BlockDepth: 0,
				ArgDepth:   0,
				Line:       0,
				Column:     0,
			},
			args: args{
				c: func(r rune) bool {
					if r == 'e' {
						return true
					}
					return false
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				Input:      tt.fields.Input,
				Start:      tt.fields.Start,
				Pos:        tt.fields.Pos,
				Width:      tt.fields.Width,
				widthStack: tt.fields.widthStack,
				items:      tt.fields.items,
				Depth:      tt.fields.Depth,
				BlockDepth: tt.fields.BlockDepth,
				ArgDepth:   tt.fields.ArgDepth,
				Mode:       tt.fields.Mode,
				Line:       tt.fields.Line,
				Column:     tt.fields.Column,
			}
			l.AcceptUntil(tt.args.c)
		})
	}
}

func TestLexer_AcceptRunTimes(t *testing.T) {
	type fields struct {
		Input      string
		Start      int
		Pos        int
		Width      int
		widthStack []*RuneWidth
		items      []Item
		Depth      int
		BlockDepth int
		ArgDepth   int
		Mode       StateFn
		Line       int
		Column     int
	}
	type args struct {
		c     CheckRune
		times int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "TestLexer_AcceptRunTimes 1",
			fields: fields{
				Input:      "test",
				Start:      0,
				Pos:        0,
				Width:      0,
				widthStack: []*RuneWidth{},
				items: []Item{
					{
						Typ:    5,
						Val:    "t",
						line:   0,
						column: 0,
					},
					{
						Typ:    5,
						Val:    "e",
						line:   0,
						column: 0,
					},
				},
				Depth:      0,
				BlockDepth: 0,
				ArgDepth:   0,
				Mode: func(*Lexer) StateFn {
					return nil
				},
				Line:   0,
				Column: 0,
			},
			args: args{
				c: func(r rune) bool {
					if r == 'e' {
						return true
					}
					return false
				},
				times: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				Input:      tt.fields.Input,
				Start:      tt.fields.Start,
				Pos:        tt.fields.Pos,
				Width:      tt.fields.Width,
				widthStack: tt.fields.widthStack,
				items:      tt.fields.items,
				Depth:      tt.fields.Depth,
				BlockDepth: tt.fields.BlockDepth,
				ArgDepth:   tt.fields.ArgDepth,
				Mode:       tt.fields.Mode,
				Line:       tt.fields.Line,
				Column:     tt.fields.Column,
			}
			if got := l.AcceptRunTimes(tt.args.c, tt.args.times); got != tt.want {
				t.Errorf("Lexer.AcceptRunTimes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_IgnoreRun(t *testing.T) {
	type fields struct {
		Input      string
		Start      int
		Pos        int
		Width      int
		widthStack []*RuneWidth
		items      []Item
		Depth      int
		BlockDepth int
		ArgDepth   int
		Mode       StateFn
		Line       int
		Column     int
	}
	type args struct {
		c CheckRune
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestLexer_IgnoreRun 1",
			fields: fields{
				Input:      "Test string",
				Start:      0,
				Pos:        0,
				Width:      0,
				widthStack: []*RuneWidth{},
				items: []Item{
					{
						Typ:    5,
						Val:    "t",
						line:   0,
						column: 0,
					},
					{
						Typ:    5,
						Val:    "e",
						line:   0,
						column: 0,
					},
				},
				Depth:      0,
				BlockDepth: 0,
				ArgDepth:   0,
				Line:       0,
				Column:     0,
			},
			args: args{
				c: func(r rune) bool {
					if r == 'e' {
						return true
					}
					return false
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				Input:      tt.fields.Input,
				Start:      tt.fields.Start,
				Pos:        tt.fields.Pos,
				Width:      tt.fields.Width,
				widthStack: tt.fields.widthStack,
				items:      tt.fields.items,
				Depth:      tt.fields.Depth,
				BlockDepth: tt.fields.BlockDepth,
				ArgDepth:   tt.fields.ArgDepth,
				Mode:       tt.fields.Mode,
				Line:       tt.fields.Line,
				Column:     tt.fields.Column,
			}
			l.IgnoreRun(tt.args.c)
		})
	}
}

func TestLexer_IsEscChar(t *testing.T) {
	type fields struct {
		Input      string
		Start      int
		Pos        int
		Width      int
		widthStack []*RuneWidth
		items      []Item
		Depth      int
		BlockDepth int
		ArgDepth   int
		Mode       StateFn
		Line       int
		Column     int
	}
	type args struct {
		r rune
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "TestLexer_IsEscChar 1",
			fields: fields{
				Input:      "u",
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
			},
			args: args{r: 'u'},
			want: true,
		},
		{
			name: "TestLexer_IsEscChar 2",
			fields: fields{
				Input:      "a",
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
			},
			args: args{r: 'a'},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				Input:      tt.fields.Input,
				Start:      tt.fields.Start,
				Pos:        tt.fields.Pos,
				Width:      tt.fields.Width,
				widthStack: tt.fields.widthStack,
				items:      tt.fields.items,
				Depth:      tt.fields.Depth,
				BlockDepth: tt.fields.BlockDepth,
				ArgDepth:   tt.fields.ArgDepth,
				Mode:       tt.fields.Mode,
				Line:       tt.fields.Line,
				Column:     tt.fields.Column,
			}
			if got := l.IsEscChar(tt.args.r); got != tt.want {
				t.Errorf("Lexer.IsEscChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEndOfLine(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "TestIsEndOfLine 1",
			args: args{r: '\u000A'},
			want: true,
		},
		{
			name: "TestIsEndOfLine 2",
			args: args{r: '\u000D'},
			want: true,
		},
		{
			name: "TestIsEndOfLine 3",
			args: args{r: 't'},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEndOfLine(tt.args.r); got != tt.want {
				t.Errorf("IsEndOfLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_LexQuotedString(t *testing.T) {
	type fields struct {
		Input      string
		Start      int
		Pos        int
		Width      int
		widthStack []*RuneWidth
		items      []Item
		Depth      int
		BlockDepth int
		ArgDepth   int
		Mode       StateFn
		Line       int
		Column     int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestLexer_LexQuotedString 1",
			fields: fields{
				Input: "string",
				Start: 0,
				Pos:   4,
				Width: 4,
				widthStack: []*RuneWidth{
					{
						width: 4,
						count: 0,
					},
				},
				items: []Item{
					{
						Typ:    5,
						Val:    "t",
						line:   0,
						column: 0,
					},
					{
						Typ:    5,
						Val:    "e",
						line:   0,
						column: 0,
					},
					{
						Typ:    5,
						Val:    "s",
						line:   0,
						column: 0,
					},
					{
						Typ:    5,
						Val:    "t",
						line:   0,
						column: 0,
					},
				},
				Depth:      0,
				BlockDepth: 0,
				ArgDepth:   0,
				Line:       0,
				Column:     0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				Input:      tt.fields.Input,
				Start:      tt.fields.Start,
				Pos:        tt.fields.Pos,
				Width:      tt.fields.Width,
				widthStack: tt.fields.widthStack,
				items:      tt.fields.items,
				Depth:      tt.fields.Depth,
				BlockDepth: tt.fields.BlockDepth,
				ArgDepth:   tt.fields.ArgDepth,
				Mode:       tt.fields.Mode,
				Line:       tt.fields.Line,
				Column:     tt.fields.Column,
			}
			if err := l.LexQuotedString(); (err != nil) != tt.wantErr {
				t.Errorf("Lexer.LexQuotedString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLexer_LexQuotedStringForEOF(t *testing.T) {
	type fields struct {
		Input      string
		Start      int
		Pos        int
		Width      int
		widthStack []*RuneWidth
		items      []Item
		Depth      int
		BlockDepth int
		ArgDepth   int
		Mode       StateFn
		Line       int
		Column     int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestLexer_LexQuotedString 1",
			fields: fields{
				Input: `"aaa`,
				Start: 0,
				Pos:   4,
				Width: 4,
				widthStack: []*RuneWidth{
					{
						width: 4,
						count: 3,
					},
				},
				items:      []Item{},
				Depth:      0,
				BlockDepth: 0,
				ArgDepth:   0,
				Line:       0,
				Column:     0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				Input:      tt.fields.Input,
				Start:      tt.fields.Start,
				Pos:        tt.fields.Pos,
				Width:      tt.fields.Width,
				widthStack: tt.fields.widthStack,
				items:      tt.fields.items,
				Depth:      tt.fields.Depth,
				BlockDepth: tt.fields.BlockDepth,
				ArgDepth:   tt.fields.ArgDepth,
				Mode:       tt.fields.Mode,
				Line:       tt.fields.Line,
				Column:     tt.fields.Column,
			}
			if err := l.LexQuotedString(); (err != nil) != tt.wantErr {
				t.Errorf("Lexer.LexQuotedString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLexer_LexQuotedStringFor1(t *testing.T) {
	type fields struct {
		Input      string
		Start      int
		Pos        int
		Width      int
		widthStack []*RuneWidth
		items      []Item
		Depth      int
		BlockDepth int
		ArgDepth   int
		Mode       StateFn
		Line       int
		Column     int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestLexer_LexQuotedString 1",
			fields: fields{
				Input: `"aa"`,
				Start: 0,
				Pos:   4,
				Width: 4,
				widthStack: []*RuneWidth{
					{
						width: 4,
						count: 3,
					},
				},
				items:      []Item{},
				Depth:      0,
				BlockDepth: 0,
				ArgDepth:   0,
				Line:       0,
				Column:     0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				Input:      tt.fields.Input,
				Start:      tt.fields.Start,
				Pos:        tt.fields.Pos,
				Width:      tt.fields.Width,
				widthStack: tt.fields.widthStack,
				items:      tt.fields.items,
				Depth:      tt.fields.Depth,
				BlockDepth: tt.fields.BlockDepth,
				ArgDepth:   tt.fields.ArgDepth,
				Mode:       tt.fields.Mode,
				Line:       tt.fields.Line,
				Column:     tt.fields.Column,
			}
			if err := l.LexQuotedString(); (err != nil) == tt.wantErr {
				t.Errorf("Lexer.LexQuotedString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLexer_LexQuotedStringFor2(t *testing.T) {
	type fields struct {
		Input      string
		Start      int
		Pos        int
		Width      int
		widthStack []*RuneWidth
		items      []Item
		Depth      int
		BlockDepth int
		ArgDepth   int
		Mode       StateFn
		Line       int
		Column     int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestLexer_LexQuotedString 1",
			fields: fields{
				Input: `"\\`,
				Start: 0,
				Pos:   4,
				Width: 4,
				widthStack: []*RuneWidth{
					{
						width: 4,
						count: 3,
					},
				},
				items:      []Item{},
				Depth:      0,
				BlockDepth: 0,
				ArgDepth:   0,
				Line:       0,
				Column:     0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				Input:      tt.fields.Input,
				Start:      tt.fields.Start,
				Pos:        tt.fields.Pos,
				Width:      tt.fields.Width,
				widthStack: tt.fields.widthStack,
				items:      tt.fields.items,
				Depth:      tt.fields.Depth,
				BlockDepth: tt.fields.BlockDepth,
				ArgDepth:   tt.fields.ArgDepth,
				Mode:       tt.fields.Mode,
				Line:       tt.fields.Line,
				Column:     tt.fields.Column,
			}
			if err := l.LexQuotedString(); (err != nil) != tt.wantErr {
				t.Errorf("Lexer.LexQuotedString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLexer_LexQuotedStringFor3(t *testing.T) {
	type fields struct {
		Input      string
		Start      int
		Pos        int
		Width      int
		widthStack []*RuneWidth
		items      []Item
		Depth      int
		BlockDepth int
		ArgDepth   int
		Mode       StateFn
		Line       int
		Column     int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestLexer_LexQuotedString 1",
			fields: fields{
				Input: `"'\\'a\`,
				Start: 0,
				Pos:   3,
				Width: 3,
				widthStack: []*RuneWidth{
					{
						width: 3,
						count: 2,
					},
				},
				items:      []Item{},
				Depth:      0,
				BlockDepth: 0,
				ArgDepth:   0,
				Line:       0,
				Column:     0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				Input:      tt.fields.Input,
				Start:      tt.fields.Start,
				Pos:        tt.fields.Pos,
				Width:      tt.fields.Width,
				widthStack: tt.fields.widthStack,
				items:      tt.fields.items,
				Depth:      tt.fields.Depth,
				BlockDepth: tt.fields.BlockDepth,
				ArgDepth:   tt.fields.ArgDepth,
				Mode:       tt.fields.Mode,
				Line:       tt.fields.Line,
				Column:     tt.fields.Column,
			}
			if err := l.LexQuotedString(); (err != nil) != tt.wantErr {
				t.Errorf("Lexer.LexQuotedString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
