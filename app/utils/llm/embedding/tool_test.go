package embedding

import "testing"

func TestParseProviderModelID(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantProvID string
		wantModel  string
		wantErr    bool
	}{
		{
			name:       "normal",
			input:      "text-embedding-v4@Tongyi-Qianwen",
			wantProvID: "Tongyi-Qianwen",
			wantModel:  "text-embedding-v4",
		},
		{
			name:       "trim space",
			input:      " text-embedding-v4 @ Tongyi-Qianwen ",
			wantProvID: "Tongyi-Qianwen",
			wantModel:  "text-embedding-v4",
		},
		{
			name:    "missing separator",
			input:   "text-embedding-v4",
			wantErr: true,
		},
		{
			name:    "empty model",
			input:   "@Tongyi-Qianwen",
			wantErr: true,
		},
		{
			name:    "empty provider",
			input:   "text-embedding-v4@",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProvID, gotModel, err := ParseProviderModelID(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotProvID != tt.wantProvID {
				t.Fatalf("unexpected providerID: %s", gotProvID)
			}
			if gotModel != tt.wantModel {
				t.Fatalf("unexpected modelID: %s", gotModel)
			}
		})
	}
}
