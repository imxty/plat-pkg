package main

import (
	"github.com/jinmukeji/plat-pkg/v2/jm-micro/plugins/cid"
	"github.com/jinmukeji/plat-pkg/v2/jm-micro/plugins/configloader"
	"github.com/jinmukeji/plat-pkg/v2/jm-micro/plugins/healthcheck"
	"github.com/jinmukeji/plat-pkg/v2/jm-micro/plugins/jwt"
	"github.com/jinmukeji/plat-pkg/v2/jm-micro/plugins/log"
	"github.com/jinmukeji/plat-pkg/v2/jm-micro/plugins/tcphealthcheck"
	"github.com/jinmukeji/plat-pkg/v2/jm-micro/plugins/tls-client"

	// "github.com/jinmukeji/plat-pkg/v2/jm-micro/plugins/whitelist"

	"github.com/micro/go-plugins/micro/cors/v2"
	// "github.com/micro/go-plugins/micro/gzip"
	"github.com/micro/go-plugins/micro/metadata/v2"

	"github.com/micro/micro/v2/api"
	"github.com/micro/micro/v2/plugin"
	"github.com/micro/micro/v2/proxy"
	"github.com/micro/micro/v2/web"
)

func init() {
	// 全局插件
	err := plugin.Register(log.NewPlugin(Name))
	die(err)

	err = plugin.Register(configloader.NewPlugin())
	die(err)

	// FIXME: metadata 插件不兼容 micro/cli/v2，等待 micro/go-plugins 修复
	err = plugin.Register(metadata.NewPlugin())
	die(err)

	// FIXME: cors 插件不兼容 micro/cli/v2，等待 micro/go-plugins 修复
	err = plugin.Register(cors.NewPlugin())
	die(err)

	// proxy 服务插件
	err = proxy.Register(tls.NewPlugin())
	die(err)

	// web 服务插件
	err = web.Register(tls.NewPlugin())
	die(err)

	// api 服务插件

	err = api.Register(tls.NewPlugin())
	die(err)

	// micro gzip 插件存在 bug，当 response 数据量过小的时候，压缩后的数据丢失
	// err = api.Register(gzip.NewPlugin())
	// die(err)
	err = api.Register(healthcheck.NewPlugin())
	die(err)

	err = api.Register(tcphealthcheck.NewPlugin())
	die(err)

	err = api.Register(cid.NewPlugin())
	die(err)

	err = api.Register(jwt.NewPlugin())
	die(err)

	// TODO: 白名单插件
	// err = api.Register(whitelist.NewRPCWhitelist("com.jinmuhealth.platform.srv.template-service1"))
	// die(err)
}
