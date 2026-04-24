package input

import (
	"testing"

	"github.com/toon-format/toon-go"
)

func TestParseToToon(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "json object preserves field order",
			input: `{"name":"Ada","id":1,"active":true}`,
			want:  "name: Ada\nid: 1\nactive: true",
		},
		{
			name:  "json array",
			input: `[1,2,3]`,
			want:  "[3]: 1,2,3",
		},
		{
			name: "jsonc object",
			input: `{
				// preferred output order should match source
				"name": "Ada",
				"languages": ["go", "toon"],
			}`,
			want: "name: Ada\nlanguages[2]: go,toon",
		},
		{
			name:  "ndjson becomes top level array",
			input: "{\"id\":1,\"name\":\"Ada\"}\n{\"id\":2,\"name\":\"Bob\"}\n",
			want:  "[2]{id,name}:\n  1,Ada\n  2,Bob",
		},
		{
			name:  "large integers stay precise",
			input: `{"id":9007199254740993}`,
			want:  `id: "9007199254740993"`,
		},
		{
			name:  "huge integers stay precise",
			input: `{"id":1000000000000000000000000}`,
			want:  `id: "1000000000000000000000000"`,
		},
		{
			name:    "empty input fails",
			input:   " \n\t ",
			wantErr: true,
		},
		{
			name:    "invalid input fails",
			input:   `{"name":]`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Parse([]byte(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Fatal("Parse() error = nil, want error")
				}
				return
			}

			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			encoded, err := toon.MarshalString(got)
			if err != nil {
				t.Fatalf("toon.MarshalString() error = %v", err)
			}

			if encoded != tt.want {
				t.Fatalf("toon.MarshalString() = %q, want %q", encoded, tt.want)
			}
		})
	}
}
