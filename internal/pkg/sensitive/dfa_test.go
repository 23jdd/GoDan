package sensitive

import "testing"

func TestDFAFilter_Check(t *testing.T) {
	words := []string{"fuck", "shit", "badword", "赌博", "色情"}
	f := NewDFAFilter(words)

	tests := []struct {
		input string
		want  bool
	}{
		{"hello world", false},
		{"what the fuck", true},
		{"FUCKING", false},
		{"this is shitty", true}, // "shit" is substring
		{"shit", true},
		{"badword here", true},
		{"bad words", false},
		{"这是赌博网站", true},
		{"赌博", true},
		{"正常评论", false},
		{"色情内容", true},
		{"包含色情关键字", true},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := f.Check(tt.input)
			if got != tt.want {
				t.Errorf("Check(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestDFAFilter_Replace(t *testing.T) {
	words := []string{"fuck", "shit", "赌博", "色情"}
	f := NewDFAFilter(words)

	tests := []struct {
		input string
		want  string
	}{
		{"hello world", "hello world"},
		{"what the fuck man", "what the **** man"},
		{"shit happens", "**** happens"},
		{"这是赌博网站啊", "这是**网站啊"},
		{"色情内容", "**内容"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := f.Replace(tt.input, '*')
			if got != tt.want {
				t.Errorf("Replace(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestDFAFilter_EmptyWords(t *testing.T) {
	f := NewDFAFilter(nil)
	if f.Check("anything") {
		t.Error("empty filter should not match")
	}
	if s := f.Replace("hello", '*'); s != "hello" {
		t.Errorf("empty filter should not replace, got %q", s)
	}
}

func TestDFAFilter_Overlap(t *testing.T) {
	words := []string{"ab", "bc"}
	f := NewDFAFilter(words)

	if !f.Check("abc") {
		t.Error("should detect overlapping patterns")
	}
	if s := f.Replace("abc", '*'); s != "***" {
		t.Errorf("overlap replace got %q, want %q", s, "***")
	}
}
