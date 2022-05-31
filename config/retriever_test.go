package config_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/config"
	"net/http"
	"testing"
)

func TestRetrieverConf_IsValid(t *testing.T) {
	type fields struct {
		Kind           string
		RepositorySlug string
		Branch         string
		Path           string
		GithubToken    string
		URL            string
		Timeout        int64
		HTTPMethod     string
		HTTPBody       string
		HTTPHeaders    map[string][]string
		Bucket         string
		Object         string
		Item           string
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		errValue string
	}{
		{
			name:     "no fields",
			fields:   fields{},
			wantErr:  true,
			errValue: "invalid retriever: kind \"\" is not supported",
		},
		{
			name: "invalid kind",
			fields: fields{
				Kind: "invalid",
			},
			wantErr:  true,
			errValue: "invalid retriever: kind \"invalid\" is not supported",
		},
		{
			name: "kind GitHubRetriever without repo slug",
			fields: fields{
				Kind: "github",
			},
			wantErr:  true,
			errValue: "invalid retriever: no \"repositorySlug\" property found for kind \"github\"",
		},
		{
			name: "kind S3Retriever without item",
			fields: fields{
				Kind: "s3",
			},
			wantErr:  true,
			errValue: "invalid retriever: no \"item\" property found for kind \"s3\"",
		},
		{
			name: "kind HTTPRetriever without URL",
			fields: fields{
				Kind: "http",
			},
			wantErr:  true,
			errValue: "invalid retriever: no \"url\" property found for kind \"http\"",
		},
		{
			name: "kind GCP without Object",
			fields: fields{
				Kind: "googleStorage",
			},
			wantErr:  true,
			errValue: "invalid retriever: no \"object\" property found for kind \"googleStorage\"",
		},
		{
			name: "kind GitHubRetriever without path",
			fields: fields{
				Kind:           "github",
				RepositorySlug: "thomaspoignant/go-feature-flag",
			},
			wantErr:  true,
			errValue: "invalid retriever: no \"path\" property found for kind \"github\"",
		},
		{
			name: "kind file without path",
			fields: fields{
				Kind: "file",
			},
			wantErr:  true,
			errValue: "invalid retriever: no \"path\" property found for kind \"file\"",
		},
		{
			name: "kind s3 without bucket",
			fields: fields{
				Kind: "s3",
				Item: "test.yaml",
			},
			wantErr:  true,
			errValue: "invalid retriever: no \"bucket\" property found for kind \"s3\"",
		},
		{
			name: "kind google storage without bucket",
			fields: fields{
				Kind:   "googleStorage",
				Object: "test.yaml",
			},
			wantErr:  true,
			errValue: "invalid retriever: no \"bucket\" property found for kind \"googleStorage\"",
		},
		{
			name: "valid s3",
			fields: fields{
				Kind:   "s3",
				Item:   "test.yaml",
				Bucket: "testBucket",
			},
		},
		{
			name: "valid googleStorage",
			fields: fields{
				Kind:   "googleStorage",
				Object: "test.yaml",
				Bucket: "testBucket",
			},
		},
		{
			name: "valid github",
			fields: fields{
				Kind:           "github",
				RepositorySlug: "thomaspoignant/go-feature-flag",
				Branch:         "main",
				Path:           "testdata/config.yaml",
				GithubToken:    "XXX",
				Timeout:        5000,
			},
		},
		{
			name: "valid file",
			fields: fields{
				Kind: "file",
				Path: "testdata/config.yaml",
			},
		},
		{
			name: "valid http",
			fields: fields{
				Kind:       "http",
				URL:        "http://perdu.com/flags",
				HTTPMethod: http.MethodGet,
				HTTPBody:   `{"yo"": "yo"}`,
				HTTPHeaders: map[string][]string{
					"Test": {"Val1"},
				},
				Timeout: 5000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &config.RetrieverConf{
				Kind:           config.RetrieverKind(tt.fields.Kind),
				RepositorySlug: tt.fields.RepositorySlug,
				Branch:         tt.fields.Branch,
				Path:           tt.fields.Path,
				GithubToken:    tt.fields.GithubToken,
				URL:            tt.fields.URL,
				Timeout:        tt.fields.Timeout,
				HTTPMethod:     tt.fields.HTTPMethod,
				HTTPBody:       tt.fields.HTTPBody,
				HTTPHeaders:    tt.fields.HTTPHeaders,
				Bucket:         tt.fields.Bucket,
				Object:         tt.fields.Object,
				Item:           tt.fields.Item,
			}
			err := c.IsValid()
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				assert.Equal(t, tt.errValue, err.Error())
			}
		})
	}
}
