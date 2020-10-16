package router

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/app/api"
)

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func init() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.WriteTpl("index.html")
	})
	s.SetRewriteMap(map[string]string{
		"/dashboard": "/",
		"/roles":     "/",
		"/monitor":   "/",
	})
	proxyInstanceApi := &api.ProxyInstanceApi{}
	proxyServerApi := &api.ProxyServerApi{}
	proxyRoleApi := &api.ProxyRoleApi{}

	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareCORS)
		// proxy instance api
		group.POST("/proxy/instance/create", proxyInstanceApi.Create)
		group.POST("/proxy/instance/update", proxyInstanceApi.Update)
		group.POST("/proxy/instance/remove", proxyInstanceApi.Remove)
		group.POST("/proxy/instance/pool/remove", proxyInstanceApi.RemoveFromPool)
		group.POST("/proxy/instance/pool/add", proxyInstanceApi.IntoPool)
		group.GET("/proxy/instances", proxyInstanceApi.Query)
		// proxy server api
		group.POST("/proxy/server/start", proxyServerApi.Start)
		group.POST("/proxy/server/stop", proxyServerApi.Stop)
		group.POST("/proxy/server/update", proxyServerApi.Update)
		group.GET("/proxy/server/info", proxyServerApi.Info)
		group.GET("/proxy/server/monitor", proxyServerApi.Monitor)
		// proxy url role
		group.POST("/proxy/role/add", proxyRoleApi.AddRole)
		group.POST("/proxy/role/remove", proxyRoleApi.Remove)
		group.GET("/proxy/roles", proxyRoleApi.All)
	})
}
