package webapi

import (
	"io/ioutil"
	"net/http"

	"github.com/nsip/n3-privacy/jkv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	g "github.com/nsip/n3-privacy/Server/global"
	db "github.com/nsip/n3-privacy/Server/storage"
	cmn "github.com/nsip/n3-privacy/common"
)

// HostHTTPAsync : Host a HTTP Server for providing policy json
func HostHTTPAsync() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("2G"))

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.DELETE},
		AllowCredentials: true,
	}))

	port := g.Cfg.WebService.Port
	fullIP := cmn.LocalIP() + fSf(":%d", port)

	route := g.Cfg.Route
	initMutex()

	// *************************************** List all APP, API *************************************** //
	e.GET("/", func(c echo.Context) error {
		return c.String(
			http.StatusOK,
			fSf("GET    %-50s->  %s\n", fullIP+route.Peek, "Get policy's HashCode. If no policy, return empty")+
				fSf("GET    %-50s->  %s\n", fullIP+route.Get, "Get policy's JSON file. If no policy, return empty")+
				fSf("POST   %-50s->  %s\n", fullIP+route.Update, "Update policy. If no policy exists, add it")+
				fSf("GET    %-50s->  %s\n", fullIP+route.GetJM, "Get JSON enforcement tool (jm). This tool is dependent on (jq)")+
				fSf("GET    %-50s->  %s\n", fullIP+route.GetJQ, "Get JQ1.6. Put (jq) into (jm) directory"))
	})

	e.GET(route.Peek, func(c echo.Context) error {
		defer func() { mMtx[route.Peek].Unlock() }()
		mMtx[route.Peek].Lock()
		g.WDCheck()
		params := c.QueryParams()
		if uid, ok := params["uid"]; ok {
			if ctx, ok := params["ctx"]; ok {
				if object, ok := params["object"]; ok {
					if rw, ok := params["rw"]; ok {
						if mCodes := db.GetPolicyCode(uid[0], ctx[0], object[0], rw[0]); len(mCodes) > 0 {
							return c.JSON(http.StatusOK, mCodes)
						}
						return c.JSON(http.StatusNotFound, "No Policy as your request")
					}
				}
			}
		}
		return c.JSON(http.StatusBadRequest, "<uid>, <ctx>, <object>, and <rw> parameters must be provided")
	})

	e.GET(route.Get, func(c echo.Context) error {
		defer func() { mMtx[route.Get].Unlock() }()
		mMtx[route.Get].Lock()
		g.WDCheck()
		params := c.QueryParams()
		if code, ok := params["code"]; ok {
			if policy, ok := db.GetPolicy(code[0]); ok {
				return c.JSON(http.StatusOK, policy)
			}
			return c.JSON(http.StatusNotFound, "No Policy as your request")
		}
		return c.JSON(http.StatusBadRequest, "policy <code> parameters must be provided")
	})

	e.POST(route.Update, func(c echo.Context) error {
		defer func() { mMtx[route.Update].Unlock() }()
		mMtx[route.Update].Lock()
		g.WDCheck()
		params := c.QueryParams()
		if uid, ok := params["uid"]; ok {
			if ctx, ok := params["ctx"]; ok {
				if rw, ok := params["rw"]; ok {
					if bPolicy, err := ioutil.ReadAll(c.Request().Body); err == nil && jkv.IsJSON(string(bPolicy)) {
						db.UpdatePolicy(uid[0], ctx[0], rw[0], string(bPolicy))
						return c.JSON(http.StatusOK, "OK")
					}
					return c.String(http.StatusBadRequest, "Policy is not in BODY, or is not valid JSON")
				}
			}
		}
		return c.JSON(http.StatusBadRequest, "<uid>, <ctx> and <rw> parameters must be provided")
	})

	e.Start(fSf(":%d", port))
}
