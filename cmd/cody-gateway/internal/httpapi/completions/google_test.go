package completions

import (
	"strings"
	"testing"

	"bytes"

	"github.com/sourcegraph/log/logtest"
	"github.com/stretchr/testify/assert"
)

func TestGoogleRequestGetTokenCount(t *testing.T) {
	logger := logtest.Scoped(t)

	t.Run("streaming", func(t *testing.T) {
		req := googleRequest{}
		r := strings.NewReader(googleStreamingResponse)
		handler := &GoogleHandlerMethods{}
		promptUsage, completionUsage := handler.parseResponseAndUsage(logger, req, r)

		assert.Equal(t, 21, promptUsage.tokens)
		assert.Equal(t, 87, completionUsage.tokens)
	})
}

var googleStreamingResponse = `data: {"candidates": [{"content": {"parts": [{"text": "def"}],"role": "model"},"finishReason": "STOP","index": 0}],"usageMetadata": {"promptTokenCount": 21,"candidatesTokenCount": 1,"totalTokenCount": 22}}

data: {"candidates": [{"content": {"parts": [{"text": " bubble_sort(list1):\n  n = len(list1)"}],"role": "model"},"finishReason": "STOP","index": 0,"safetyRatings": [{"category": "HARM_CATEGORY_SEXUALLY_EXPLICIT","probability": "NEGLIGIBLE"},{"category": "HARM_CATEGORY_HATE_SPEECH","probability": "NEGLIGIBLE"},{"category": "HARM_CATEGORY_HARASSMENT","probability": "NEGLIGIBLE"},{"category": "HARM_CATEGORY_DANGEROUS_CONTENT","probability": "NEGLIGIBLE"}]}],"usageMetadata": {"promptTokenCount": 21,"candidatesTokenCount": 17,"totalTokenCount": 38}}

data: {"candidates": [{"content": {"parts": [{"text": "\n  for i in range(n-1):\n    for j in"}],"role": "model"},"finishReason": "STOP","index": 0,"safetyRatings": [{"category": "HARM_CATEGORY_SEXUALLY_EXPLICIT","probability": "NEGLIGIBLE"},{"category": "HARM_CATEGORY_HATE_SPEECH","probability": "NEGLIGIBLE"},{"category": "HARM_CATEGORY_HARASSMENT","probability": "NEGLIGIBLE"},{"category": "HARM_CATEGORY_DANGEROUS_CONTENT","probability": "NEGLIGIBLE"}]}],"usageMetadata": {"promptTokenCount": 21,"candidatesTokenCount": 31,"totalTokenCount": 52}}

data: {"candidates": [{"content": {"parts": [{"text": " range(n-i-1):\n      if list1[j] \u003e list1[j+1]:\n        list1[j], list"}],"role": "model"},"finishReason": "STOP","index": 0,"safetyRatings": [{"category": "HARM_CATEGORY_SEXUALLY_EXPLICIT","probability": "NEGLIGIBLE"},{"category": "HARM_CATEGORY_HATE_SPEECH","probability": "NEGLIGIBLE"},{"category": "HARM_CATEGORY_HARASSMENT","probability": "NEGLIGIBLE"},{"category": "HARM_CATEGORY_DANGEROUS_CONTENT","probability": "NEGLIGIBLE"}]}],"usageMetadata": {"promptTokenCount": 21,"candidatesTokenCount": 63,"totalTokenCount": 84}}

data: {"candidates": [{"content": {"parts": [{"text": "1[j+1] = list1[j+1], list1[j]\n  return list1\n"}],"role": "model"},"finishReason": "STOP","index": 0,"safetyRatings": [{"category": "HARM_CATEGORY_SEXUALLY_EXPLICIT","probability": "NEGLIGIBLE"},{"category": "HARM_CATEGORY_HATE_SPEECH","probability": "NEGLIGIBLE"},{"category": "HARM_CATEGORY_HARASSMENT","probability": "NEGLIGIBLE"},{"category": "HARM_CATEGORY_DANGEROUS_CONTENT","probability": "NEGLIGIBLE"}],"citationMetadata": {"citationSources": [{"startIndex": 1,"endIndex": 185,"uri": "https://github.com/Feng080412/Searches-and-sorts","license": ""}]}}],"usageMetadata": {"promptTokenCount": 21,"candidatesTokenCount": 87,"totalTokenCount": 108}}`

func TestParseGoogleTokenUsage(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *googleResponse
		wantErr bool
	}{
		{
			name:  "valid response",
			input: `data: {"candidates": [{"content": {"parts": [{"text": "def"}],"role": "model"},"finishReason": "STOP","index": 0}],"usageMetadata": {"promptTokenCount": 21,"candidatesTokenCount": 1,"totalTokenCount": 22}}`,
			want: &googleResponse{
				UsageMetadata: googleUsage{
					PromptTokenCount:     21,
					CompletionTokenCount: 1,
					TotalTokenCount:      0,
				},
			},
			wantErr: false,
		},
		{
			name:  "valid response - with candidates",
			input: `data: {"usageMetadata": {"promptTokenCount": 10, "candidatesTokenCount": 20}}`,
			want: &googleResponse{
				UsageMetadata: googleUsage{
					PromptTokenCount:     10,
					CompletionTokenCount: 20,
					TotalTokenCount:      0,
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `data: {"usageMetadata": {"promptTokenCount": 10, "candidatesTokenCount": 20}`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no prefix",
			input:   `{"usageMetadata": {"promptTokenCount": 10, "candidatesTokenCount": 20}}`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   ``,
			want:    nil,
			wantErr: true,
		},
		{
			name: "multiple lines with one valid",
			input: `data: {"usageMetadata": {"promptTokenCount": 5, "candidatesTokenCount": 15}}

data: {"usageMetadata": {"promptTokenCount": 10, "candidatesTokenCount": 20}}`,
			want: &googleResponse{
				UsageMetadata: googleUsage{
					PromptTokenCount:     10,
					CompletionTokenCount: 20,
					TotalTokenCount:      0,
				},
			},
			wantErr: false,
		},
		{
			name:    "non-JSON data",
			input:   `data: not-a-json`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "partial data",
			input:   `data: {"usageMetadata": {"promptTokenCount": 10`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader([]byte(tt.input))
			logger := logtest.Scoped(t)
			promptTokens, completionTokens, err := parseGoogleTokenUsage(r, logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGoogleTokenUsage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				got := &googleResponse{
					UsageMetadata: googleUsage{
						PromptTokenCount:     promptTokens,
						CompletionTokenCount: completionTokens,
					},
				}
				if !assert.ObjectsAreEqual(got, tt.want) {
					t.Errorf("parseGoogleTokenUsage() mismatch (-want +got):\n%v", assert.ObjectsAreEqual(got, tt.want))
				}
			}
		})
	}
}
