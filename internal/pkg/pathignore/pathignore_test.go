package pathignore

import "testing"

func TestMatcherSupportsGitignoreStylePatterns(t *testing.T) {
	matcher, err := CompileCSV(`node_modules,**/.cache,tmp-?,generated/**,!generated/keep.go,logs/`)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		isDir bool
		want  bool
	}{
		{"node_modules/pkg/index.js", false, true},
		{"app/admin/web/node_modules/pkg/index.js", false, true},
		{"app/.cache/data", false, true},
		{"tmp-a/file.go", false, true},
		{"tmp-long/file.go", false, false},
		{"generated", true, false},
		{"generated/output.go", false, true},
		{"generated/keep.go", false, false},
		{"logs", true, true},
		{"logs", false, false},
	}
	for _, test := range tests {
		if got := matcher.Match(test.name, test.isDir); got != test.want {
			t.Errorf("Match(%q, %v) = %v, want %v", test.name, test.isDir, got, test.want)
		}
	}
}

func TestMatcherAcceptsWindowsSeparators(t *testing.T) {
	matcher, err := CompileCSV(`app\admin\web\node_modules`)
	if err != nil {
		t.Fatal(err)
	}
	if !matcher.Match(`app\admin\web\node_modules\pkg\index.js`, false) {
		t.Fatal("Windows-style pattern did not match Windows-style path")
	}
}

func TestMatcherRejectsMalformedPattern(t *testing.T) {
	if _, err := CompileCSV("broken["); err == nil {
		t.Fatal("CompileCSV() accepted a malformed pattern")
	}
}
