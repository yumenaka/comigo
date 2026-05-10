package tools

import "testing"

func TestShortName(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  string
	}{
		{
			name:  "domain bracket and extension",
			title: "example.com/[Group] My Comic.cbz",
			want:  "My Comic",
		},
		{
			name:  "leading punctuation after cleanup",
			title: "---  漫画标题.rar",
			want:  "漫画标题",
		},
		{
			name:  "fallback to original when cleanup is empty",
			title: "(x).cbz",
			want:  "(x).cbz",
		},
		{
			name:  "fallback keeps exact length title",
			title: "(1234567890123)",
			want:  "(1234567890123)",
		},
		{
			name:  "truncate cleaned title",
			title: "1234567890123456.cbz",
			want:  "123456789012345…",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShortName(tt.title); got != tt.want {
				t.Fatalf("ShortName() = %q, want %q", got, tt.want)
			}
		})
	}
}
