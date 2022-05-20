package service

import (
	"fmt"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/config"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"time"
)

func NewGoFeatureFlagClient(proxyConf *config.Config, logger *zap.Logger) (*ffclient.GoFeatureFlag, error) {
	retriever, err := initRetriever(proxyConf.Retriever)
	if err != nil {
		return nil, err
	}

	f := ffclient.Config{
		PollingInterval: time.Duration(proxyConf.PollingInterval) * time.Millisecond,
		Logger:          zap.NewStdLog(logger),
		Context:         context.Background(),
		Retriever:       retriever,
		//Notifiers:               nil,
		FileFormat: proxyConf.FileFormat,
		//DataExporter:            ffclient.DataExporter{},
		StartWithRetrieverError: proxyConf.StartWithRetrieverError,
	}

	return ffclient.New(f)
}

func initRetriever(c config.RetrieverConf) (ffclient.Retriever, error) {
	// Conversions
	switch c.Kind {
	case config.GitHub:
		return &ffclient.GithubRetriever{
			RepositorySlug: c.RepositorySlug,
			Branch:         c.Branch,
			FilePath:       c.Path,
			GithubToken:    c.GithubToken,
			Timeout:        time.Duration(c.Timeout) * time.Millisecond,
		}, nil

	case config.File:
		return &ffclient.FileRetriever{
			Path: c.Path,
		}, nil

	case config.S3:
		return &ffclient.S3Retriever{
			Bucket: c.Bucket,
			Item:   c.Item,
		}, nil

	case config.HTTP:
		return &ffclient.HTTPRetriever{
			URL:     c.URL,
			Method:  c.HTTPMethod,
			Body:    c.HTTPBody,
			Header:  c.HTTPHeaders,
			Timeout: time.Duration(c.Timeout) * time.Millisecond,
		}, nil

	case config.GoogleStorage:
		return &ffclient.GCStorageRetriever{
			Bucket: c.Bucket,
			Object: c.Object,
		}, nil

	default:
		return nil, fmt.Errorf("invalid retriever: kind \"%s\" "+
			"is not supported, accepted kind: [googleStorage, http, s3, file, github]", c.Kind)
	}
}
