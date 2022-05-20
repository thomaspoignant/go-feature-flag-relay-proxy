package service

import (
	"fmt"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/config"
	"github.com/thomaspoignant/go-feature-flag/ffexporter"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"time"
)

func NewGoFeatureFlagClient(proxyConf *config.Config, logger *zap.Logger) (*ffclient.GoFeatureFlag, error) {
	retriever, err := initRetriever(proxyConf.Retriever)
	if err != nil {
		return nil, err
	}

	var exp ffclient.DataExporter
	if proxyConf.Exporter != nil {
		exp, err = initExporter(proxyConf.Exporter)
		if err != nil {
			return nil, err
		}
	}

	f := ffclient.Config{
		PollingInterval: time.Duration(proxyConf.PollingInterval) * time.Millisecond,
		Logger:          zap.NewStdLog(logger),
		Context:         context.Background(),
		Retriever:       retriever,
		//Notifiers:               nil,
		FileFormat:              proxyConf.FileFormat,
		DataExporter:            exp,
		StartWithRetrieverError: proxyConf.StartWithRetrieverError,
	}

	return ffclient.New(f)
}

func initRetriever(c *config.RetrieverConf) (ffclient.Retriever, error) {
	// Conversions
	switch c.Kind {
	case config.GitHubRetriever:
		return &ffclient.GithubRetriever{
			RepositorySlug: c.RepositorySlug,
			Branch:         c.Branch,
			FilePath:       c.Path,
			GithubToken:    c.GithubToken,
			Timeout:        time.Duration(c.Timeout) * time.Millisecond,
		}, nil

	case config.FileRetriever:
		return &ffclient.FileRetriever{
			Path: c.Path,
		}, nil

	case config.S3Retriever:
		return &ffclient.S3Retriever{
			Bucket: c.Bucket,
			Item:   c.Item,
		}, nil

	case config.HTTPRetriever:
		return &ffclient.HTTPRetriever{
			URL:     c.URL,
			Method:  c.HTTPMethod,
			Body:    c.HTTPBody,
			Header:  c.HTTPHeaders,
			Timeout: time.Duration(c.Timeout) * time.Millisecond,
		}, nil

	case config.GoogleStorageRetriever:
		return &ffclient.GCStorageRetriever{
			Bucket: c.Bucket,
			Object: c.Object,
		}, nil

	default:
		return nil, fmt.Errorf("invalid retriever: kind \"%s\" "+
			"is not supported, accepted kind: [googleStorage, http, s3, file, github]", c.Kind)
	}
}

func initExporter(c *config.ExporterConf) (ffclient.DataExporter, error) {
	dataExp := ffclient.DataExporter{
		FlushInterval:    time.Duration(c.FlushInterval) * time.Millisecond,
		MaxEventInMemory: c.MaxEventInMemory,
	}

	switch c.Kind {
	case config.WebhookExporter:
		dataExp.Exporter = &ffexporter.Webhook{
			EndpointURL: c.EndpointURL,
			Secret:      c.Secret,
			Meta:        c.Meta,
		}
		return dataExp, nil

	case config.FileExporter:
		dataExp.Exporter = &ffexporter.File{
			Format:      c.Format,
			OutputDir:   c.OutputDir,
			Filename:    c.Filename,
			CsvTemplate: c.CsvTemplate,
		}
		return dataExp, nil

	case config.LogExporter:
		dataExp.Exporter = &ffexporter.Log{
			LogFormat: c.LogFormat,
		}
		return dataExp, nil

	case config.S3Exporter:
		dataExp.Exporter = &ffexporter.S3{
			Bucket:      c.Bucket,
			Format:      c.Format,
			S3Path:      c.Path,
			Filename:    c.Filename,
			CsvTemplate: c.CsvTemplate,
		}
		return dataExp, nil

	case config.GoogleStorageExporter:
		dataExp.Exporter = &ffexporter.GoogleCloudStorage{
			Bucket:      c.Bucket,
			Format:      c.Format,
			Path:        c.Path,
			Filename:    c.Filename,
			CsvTemplate: c.CsvTemplate,
		}
		return dataExp, nil

	default:
		return ffclient.DataExporter{}, fmt.Errorf("invalid exporter: kind \"%s\" is not supported", c.Kind)
	}
}
