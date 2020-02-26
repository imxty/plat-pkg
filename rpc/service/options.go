package service

import (
	"fmt"

	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/server"

	"github.com/micro/cli/v2"
)

type RegisterServerFunc func(srv server.Server) error

type Options struct {
	// 微服务的名称
	Name string
	// 微服务的命名空间
	Namespace string

	// ProductVersion is current product version.
	ProductVersion string
	// GitCommit is the git commit short hash
	GitCommit string
	// GoVersion is go compiler version `go version`
	GoVersion string
	// BuildTime is go build time
	BuildTime string

	// Flags are CLI flags
	Flags []cli.Flag

	// CliPreAction 在标准 Action 之前调用
	CliPreAction func(c *cli.Context)

	// CliPostAction 在标准 Action 之后调用
	CliPostAction func(c *cli.Context)

	// PreServerHandlerWrappers 自定义HandlerWrapper，在标准 HandlerWrapper 之前注册
	PreServerHandlerWrappers []server.HandlerWrapper

	// PostServerHandlerWrappers 自定义HandlerWrapper，在标准 HandlerWrapper 之后注册
	PostServerHandlerWrappers []server.HandlerWrapper

	// RegisterServer 注册 micro.Server
	RegisterServer RegisterServerFunc

	// PreClientWrappers 自定义 Client Wrapper，在标准 Wrapper 之前注册
	PreClientWrappers []client.Wrapper

	// PostClientWrappers 自定义 Client Wrapper，在标准 Wrapper 之前注册
	PostClientWrappers []client.Wrapper

	// ServiceOptions 其它 Service Option
	ServiceOptions []micro.Option
}

// ServiceFQDN 返回微服务的全名
func (opts *Options) FQDN() string {
	return fmt.Sprintf("%s.%s", opts.Namespace, opts.Name)
}

// ServiceMetadata 返回微服务的 metadata
func (opts *Options) ServiceMetadata() map[string]string {
	return map[string]string{
		"ProductVersion": opts.ProductVersion,
		"GitCommit":      opts.GitCommit,
		"GoVersion":      opts.GoVersion,
		"BuildTime":      opts.BuildTime,
	}
}
