package config

import "fmt"

// RetrieverConf contains all the field to configure a retriever
type RetrieverConf struct {
	Kind           RetrieverKind       `mapstructure:"kind"`
	RepositorySlug string              `mapstructure:"repositorySlug"`
	Branch         string              `mapstructure:"branch"`
	Path           string              `mapstructure:"path"`
	GithubToken    string              `mapstructure:"githubToken"`
	URL            string              `mapstructure:"url"`
	Timeout        int64               `mapstructure:"timeout"`
	HTTPMethod     string              `mapstructure:"method"`
	HTTPBody       string              `mapstructure:"body"`
	HTTPHeaders    map[string][]string `mapstructure:"headers"`
	Bucket         string              `mapstructure:"bucket"`
	Object         string              `mapstructure:"object"`
	Item           string              `mapstructure:"item"`
}

// IsValid validate the configuration of the retriever
func (c *RetrieverConf) IsValid() error {
	if err := c.Kind.IsValid(); err != nil {
		return err
	}
	if c.Kind == GitHub && c.RepositorySlug == "" {
		return fmt.Errorf("invalid retriever: no \"repositorySlug\" property found for kind \"%s\"", c.Kind)
	}
	if c.Kind == S3 && c.Item == "" {
		return fmt.Errorf("invalid retriever: no \"item\" property found for kind \"%s\"", c.Kind)
	}
	if c.Kind == HTTP && c.URL == "" {
		return fmt.Errorf("invalid retriever: no \"url\" property found for kind \"%s\"", c.Kind)
	}
	if c.Kind == GoogleStorage && c.Object == "" {
		return fmt.Errorf("invalid retriever: no \"object\" property found for kind \"%s\"", c.Kind)
	}
	if (c.Kind == GitHub || c.Kind == File) && c.Path == "" {
		return fmt.Errorf("invalid retriever: no \"path\" property found for kind \"%s\"", c.Kind)
	}
	if (c.Kind == S3 || c.Kind == GoogleStorage) && c.Bucket == "" {
		return fmt.Errorf("invalid retriever: no \"bucket\" property found for kind \"%s\"", c.Kind)
	}
	return nil
}

// RetrieverKind is an enum containing all accepted Retriever kind
type RetrieverKind string

const (
	HTTP          RetrieverKind = "http"
	GitHub        RetrieverKind = "github"
	S3            RetrieverKind = "s3"
	File          RetrieverKind = "file"
	GoogleStorage RetrieverKind = "googleStorage"
)

// IsValid is checking if the value is part of the enum
func (r RetrieverKind) IsValid() error {
	switch r {
	case HTTP, GitHub, S3, File, GoogleStorage:
		return nil
	}
	return fmt.Errorf("invalid retriever: kind \"%s\" is not supported", r)
}
