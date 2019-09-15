package jwt

import (
	"crypto/rsa"
	"net/http"
	"strings"
	"time"

	"github.com/micro/cli"
	"github.com/micro/micro/plugin"

	mlog "github.com/jinmukeji/go-pkg/log"
	cm "github.com/jinmukeji/plat-pkg/rpc/ctxmeta"
	j "github.com/jinmukeji/plat-pkg/rpc/jwt"
	"github.com/jinmukeji/plat-pkg/rpc/jwt/keystore"
	mc "github.com/jinmukeji/plat-pkg/rpc/jwt/keystore/micro-config"
)

var (
	log *mlog.Logger = mlog.StandardLogger()
)

type jwt struct {
	enabled         bool
	headerKey       string // HTTP Request Header 中的 jwt 使用的 key
	microConfigPath string
	maxExpInterval  time.Duration // 最大过期时间间隔
	store           keystore.Store
}

const (
	defaultMaxExpInterval  = 10 * time.Minute // 10分钟
	defaultMicroConfigPath = "platform/app-key"
)

func (p *jwt) Flags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:        "enable_jwt",
			Usage:       "Enable JWT validation",
			EnvVar:      "ENABLE_JWT",
			Destination: &(p.enabled),
		},
		cli.StringFlag{
			Name:        "jwt_key",
			Usage:       "JWT HTTP header key",
			EnvVar:      "JWT_KEY",
			Value:       cm.MetaJwtKey,
			Destination: &(p.headerKey),
		},
		cli.StringFlag{
			Name:        "jwt_config_path",
			Usage:       "Micro config path for JWT",
			EnvVar:      "JWT_CONFIG_PATH",
			Value:       defaultMicroConfigPath,
			Destination: &(p.microConfigPath),
		},
		cli.DurationFlag{
			Name:        "jwt_max_exp_interval",
			Usage:       "JWT max expiration interval",
			EnvVar:      "JWT_MAX_EXP_INTERVAL",
			Value:       defaultMaxExpInterval,
			Destination: &(p.maxExpInterval),
		},
	}
}

func (p *jwt) Commands() []cli.Command {
	return nil
}

const AppIdKey = "x-app-id"

func (p *jwt) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		if !p.enabled {
			return h
		}

		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

			token := r.Header.Get(p.headerKey)
			log.Debugf("Received JWT Token: %s", token)

			opt := j.VerifyOption{
				MaxExpInterval: p.maxExpInterval,
				GetPublicKeyFunc: func(iss string) *rsa.PublicKey {
					log.Debugf("Issuer from JWT: %s", iss)
					if key := p.store.Get(iss); key != nil {
						return key.PublicKey()
					}
					return nil
				},
			}

			valid, claims, err := j.RSAVerifyJWT(token, opt)
			if !valid {
				log.Warnf("failed to validate JWT: %v", err)
				http.Error(rw, "forbidden: JWT is invalid", 403)
				return
			}

			// 从 claims 中提取 iss，即 App ID
			r.Header.Set(AppIdKey, claims.Issuer)

			// serve the request
			h.ServeHTTP(rw, r)
		})
	}
}

func (p *jwt) Init(ctx *cli.Context) error {
	baseKeyPath := splitPath(p.microConfigPath)

	store := mc.NewMicroConfigStore(baseKeyPath...)
	p.store = store

	return nil
}

func splitPath(p string) []string {
	s := strings.Trim(p, "/")
	return strings.Split(s, "/")
}

func (p *jwt) String() string {
	return "jwt"
}

func NewPlugin() plugin.Plugin {
	return NewJWT()
}

func NewJWT() plugin.Plugin {
	// create plugin
	p := &jwt{
		enabled: false,
	}

	return p
}
