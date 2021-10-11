package logip

import (
    "github.com/caddyserver/caddy"
    "github.com/caddyserver/caddy/caddyhttp/httpserver"
)

func init() {
    caddy.RegisterPlugin("logip", caddy.Plugin{
        ServerType: "http",
        Action:     Setup,
    })
}

func Setup(c *caddy.Controller) error {
    var path string

    for c.Next() {
        if !c.NextArg() {
            return c.ArgErr()
        }
        path = c.Val()

        if c.NextArg() {
            // only one argument allowed
            return c.ArgErr()
        }
    }

	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return LogIP{Next: next, Path: path}
	})

    return nil
}

