package rules

import (
	"go/token"
	"regexp"
	"testing"

	"log-linter/rule"
)

func pos(n int) token.Pos {
	return token.Pos(n)
}

func TestLowercaseRule_Check(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		enabled  bool
		wantDiag bool
		wantPos  token.Pos
	}{
		{
			name:     "строчная буква",
			msg:      "starting server",
			enabled:  true,
			wantDiag: false,
		},
		{
			name:     "начинается с цифры",
			msg:      "404 Not found",
			enabled:  true,
			wantDiag: false,
		},
		{
			name:     "начинается с пробелов и строчной",
			msg:      "   starting server",
			enabled:  true,
			wantDiag: false,
		},
		{
			name:     "заглавная буква",
			msg:      "Starting server",
			enabled:  true,
			wantDiag: true,
			wantPos:  pos(0),
		},
		{
			name:     "пробелы и заглавная",
			msg:      "   Starting server",
			enabled:  true,
			wantDiag: true,
			wantPos:  pos(3),
		},
		{
			name:     "правило отключено",
			msg:      "Starting server",
			enabled:  false,
			wantDiag: false,
		},
		{
			name:     "пустое сообщение",
			msg:      "",
			enabled:  true,
			wantDiag: false,
		},
		{
			name:     "только пробелы",
			msg:      "   \t  ",
			enabled:  true,
			wantDiag: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &LowercaseRule{
				BaseRule: rule.NewBaseRule(LowercaseName, tt.enabled),
			}

			diags := rule.Check(tt.msg, pos(0))

			if (len(diags) > 0) != tt.wantDiag {
				t.Errorf("LowercaseRule.Check() got %v diagnostics, wantDiag %v",
					len(diags), tt.wantDiag)
			}

			if tt.wantDiag && len(diags) > 0 {
				if diags[0].Pos != tt.wantPos {
					t.Errorf("wrong position: got %v, want %v", diags[0].Pos, tt.wantPos)
				}
				if diags[0].Message == "" {
					t.Error("diagnostic message is empty")
				}
			}
		})
	}
}

func TestEnglishRule_Check(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		enabled  bool
		wantDiag bool
		wantPos  token.Pos
	}{
		{
			name:     "только английские буквы",
			msg:      "starting server",
			enabled:  true,
			wantDiag: false,
		},
		{
			name:     "с цифрами и пунктуацией",
			msg:      "port 8080 is open!",
			enabled:  true,
			wantDiag: false,
		},
		{
			name:     "с пробелами в начале",
			msg:      "   starting server",
			enabled:  true,
			wantDiag: false,
		},
		{
			name:     "русская буква в начале",
			msg:      "привет world",
			enabled:  true,
			wantDiag: true,
			wantPos:  pos(0),
		},
		{
			name:     "русская буква в середине",
			msg:      "hello мир",
			enabled:  true,
			wantDiag: true,
			wantPos:  pos(6),
		},
		{
			name:     "русская буква после пробелов",
			msg:      "   привет",
			enabled:  true,
			wantDiag: true,
			wantPos:  pos(3),
		},
		{
			name:     "латинская буква с диакритикой",
			msg:      "café",
			enabled:  true,
			wantDiag: true,
			wantPos:  pos(3),
		},
		{
			name:     "китайские иероглифы",
			msg:      "hello 世界",
			enabled:  true,
			wantDiag: true,
			wantPos:  pos(6),
		},
		{
			name:     "правило отключено",
			msg:      "привет",
			enabled:  false,
			wantDiag: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &EnglishRule{
				BaseRule: rule.NewBaseRule(EnglishName, tt.enabled),
			}

			diags := rule.Check(tt.msg, pos(0))

			if (len(diags) > 0) != tt.wantDiag {
				t.Errorf("EnglishRule.Check() got %v diagnostics, wantDiag %v",
					len(diags), tt.wantDiag)
			}

			if tt.wantDiag && len(diags) > 0 {
				if diags[0].Pos != tt.wantPos {
					t.Errorf("wrong position: got %v, want %v", diags[0].Pos, tt.wantPos)
				}
			}
		})
	}
}

func TestSymbolsRule_Check(t *testing.T) {
	t.Run("default settings", func(t *testing.T) {
		rule := &SymbolsRule{
			BaseRule: rule.NewBaseRule(SymbolsName, true),
			allowed: map[rune]bool{
				'.': true,
				'-': true,
				'_': true,
				'/': true,
				':': true,
				'=': true,
				',': true,
			},
		}

		tests := []struct {
			name     string
			msg      string
			wantDiag bool
			wantPos  token.Pos
		}{
			{
				name:     "только буквы и пробелы",
				msg:      "starting server",
				wantDiag: false,
			},
			{
				name:     "разрешенная точка",
				msg:      "config.json",
				wantDiag: false,
			},
			{
				name:     "разрешенный дефис",
				msg:      "user-auth",
				wantDiag: false,
			},
			{
				name:     "разрешенное подчеркивание",
				msg:      "api_key",
				wantDiag: false,
			},
			{
				name:     "разрешенный слэш",
				msg:      "path/to/file",
				wantDiag: false,
			},
			{
				name:     "разрешенное двоеточие",
				msg:      "status: ok",
				wantDiag: false,
			},
			{
				name:     "запрещенный восклицательный знак",
				msg:      "hey!",
				wantDiag: true,
				wantPos:  pos(3),
			},
			{
				name:     "запрещенный вопросительный знак",
				msg:      "sure?",
				wantDiag: true,
				wantPos:  pos(4),
			},
			{
				name:     "эмодзи",
				msg:      "server started 🚀",
				wantDiag: true,
				wantPos:  pos(15),
			},
			{
				name:     "несколько запрещенных символов",
				msg:      "hello!!!",
				wantDiag: true,
				wantPos:  pos(5),
			},
			{
				name:     "смесь разрешенных и запрещенных",
				msg:      "config.json!",
				wantDiag: true,
				wantPos:  pos(11),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				diags := rule.Check(tt.msg, pos(0))

				if (len(diags) > 0) != tt.wantDiag {
					t.Errorf("got %v diagnostics, wantDiag %v", len(diags), tt.wantDiag)
				}
				if tt.wantDiag && len(diags) > 0 {
					if diags[0].Pos != tt.wantPos {
						t.Errorf("wrong position: got %v, want %v", diags[0].Pos, tt.wantPos)
					}
				}
			})
		}
	})

	t.Run("empty allowed list", func(t *testing.T) {
		rule := &SymbolsRule{
			BaseRule: rule.NewBaseRule(SymbolsName, true),
			allowed:  map[rune]bool{},
		}

		diags := rule.Check("config.json", pos(0))
		if len(diags) == 0 {
			t.Error("expected error for '.' with empty allowed list")
		}
	})

	t.Run("rule disabled", func(t *testing.T) {
		rule := &SymbolsRule{
			BaseRule: rule.NewBaseRule(SymbolsName, false),
			allowed: map[rune]bool{
				'.': true,
			},
		}
		diags := rule.Check("hello!!!", pos(0))
		if len(diags) > 0 {
			t.Error("expected no diagnostics for disabled rule")
		}
	})
}

func TestSensitiveRule_Check(t *testing.T) {
	patterns := []*regexp.Regexp{}
	patternStrings := []string{
		"(?i)\\b(password|passwd|pwd)\\b",
		"(?i)\\b(api[_-]?key|apikey)\\b",
		"(?i)\\btoken\\b",
	}

	for _, p := range patternStrings {
		re, _ := regexp.Compile(p)
		patterns = append(patterns, re)
	}

	t.Run("default patterns", func(t *testing.T) {
		rule := &SensitiveRule{
			BaseRule: rule.NewBaseRule(SensitiveName, true),
			patterns: patterns,
		}

		tests := []struct {
			name      string
			msg       string
			wantDiag  bool
			wantCount int
			wantPos   token.Pos
		}{
			{
				name:      "без чувствительных данных",
				msg:       "user authenticated successfully",
				wantDiag:  false,
				wantCount: 0,
			},
			{
				name:      "одно - password",
				msg:       "user password: secret123",
				wantDiag:  true,
				wantCount: 1,
				wantPos:   pos(5),
			},
			{
				name:      "одно - pwd",
				msg:       "pwd=12345",
				wantDiag:  true,
				wantCount: 1,
				wantPos:   pos(0), // позиция 'p' в "pwd"
			},
			{
				name:      "одно - api_key",
				msg:       "api_key=abc123",
				wantDiag:  true,
				wantCount: 1,
				wantPos:   pos(0), // позиция 'a' в "api_key"
			},
			{
				name:      "одно - apikey",
				msg:       "apikey=abc123",
				wantDiag:  true,
				wantCount: 1,
				wantPos:   pos(0),
			},
			{
				name:      "одно - api-key",
				msg:       "api-key=abc123",
				wantDiag:  true,
				wantCount: 1,
				wantPos:   pos(0),
			},
			{
				name:      "одно - token",
				msg:       "token: xyz",
				wantDiag:  true,
				wantCount: 1,
				wantPos:   pos(0),
			},
			{
				name:      "два в одной строке",
				msg:       "api_key=abc123, token=xyz",
				wantDiag:  true,
				wantCount: 2,
			},
			{
				name:      "три в одной строке",
				msg:       "password=123, api_key=abc, token=xyz",
				wantDiag:  true,
				wantCount: 3,
			},
			{
				name:      "разные регистры",
				msg:       "PASSWORD=123, ApiKey=456, TOKEN=789",
				wantDiag:  true,
				wantCount: 3,
			},
			{
				name:      "с пробелами в начале",
				msg:       "   password: 123",
				wantDiag:  true,
				wantCount: 1,
				wantPos:   pos(3),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				diags := rule.Check(tt.msg, pos(0))

				if (len(diags) > 0) != tt.wantDiag {
					t.Errorf("got %v diagnostics, wantDiag %v", len(diags), tt.wantDiag)
				}
				if tt.wantDiag {
					if len(diags) != tt.wantCount {
						t.Errorf("wrong number of diagnostics: got %v, want %v",
							len(diags), tt.wantCount)
					}
					if tt.wantCount == 1 && diags[0].Pos != tt.wantPos {
						t.Errorf("wrong position: got %v, want %v", diags[0].Pos, tt.wantPos)
					}
					for i, d := range diags {
						if d.Message == "" {
							t.Errorf("diagnostic %d has empty message", i)
						}
					}
				}
			})
		}
	})

	t.Run("empty patterns", func(t *testing.T) {
		rule := &SensitiveRule{
			BaseRule: rule.NewBaseRule(SensitiveName, true),
			patterns: []*regexp.Regexp{},
		}
		diags := rule.Check("password: 123", pos(0))
		if len(diags) > 0 {
			t.Error("expected no diagnostics with empty patterns")
		}
	})

	t.Run("rule disabled", func(t *testing.T) {
		rule := &SensitiveRule{
			BaseRule: rule.NewBaseRule(SensitiveName, false),
			patterns: patterns,
		}
		diags := rule.Check("password: 123", pos(0))
		if len(diags) > 0 {
			t.Error("expected no diagnostics for disabled rule")
		}
	})
}

func TestRuleRegistration(t *testing.T) {
	expectedRules := []string{
		LowercaseName,
		EnglishName,
		SymbolsName,
		SensitiveName,
	}

	registered := rule.Global.List()

	for _, expected := range expectedRules {
		found := false
		for _, r := range registered {
			if r == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("rule %q not registered", expected)
		}
	}
}
