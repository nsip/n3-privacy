package webapi

import (
	"net/http"

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
			fSf("GET    %-50s->  %s\n", fullIP+route.PeekPolicy, "Get policy's HashCode. If no policy, return empty")+
				fSf("GET    %-50s->  %s\n", fullIP+route.GetPolicy, "Get policy's JSON file. If no policy, return empty")+
				fSf("POST   %-50s->  %s\n", fullIP+route.UpdatePolicy, "Update policy. If no policy exists, add it")+
				fSf("GET    %-50s->  %s\n", fullIP+route.GetJM, "Get JSON enforcement tool (jm). This tool is dependent on (jq)")+
				fSf("GET    %-50s->  %s\n", fullIP+route.GetJQ, "Get JQ1.6. Put (jq) into (jm) directory"))
	})

	e.GET(route.PeekPolicy, func(c echo.Context) error {
		defer func() { mMtx[route.PeekPolicy].Unlock() }()
		mMtx[route.PeekPolicy].Lock()
		g.WDCheck()

		// policy := pp.FmtJSONFile("../../Server/config/mask.json", "../preprocess/utils")
		// db.UpdatePolicy("qm", "ctx1", "r", policy)
		// policy, ok := db.GetPolicy("qm", "ctx1", "inquiry_skills", "r")
		// fmt.Println(policy, ok)

		params := c.QueryParams()
		if uid, ok := params["uid"]; ok {
			if ctx, ok := params["ctx"]; ok {
				if object, ok := params["object"]; ok {
					if rw, ok := params["rw"]; ok {
						if mCode, ok := db.GetPolicyCode(uid[0], ctx[0], object[0], rw[0]); ok {
							return c.JSON(http.StatusOK, mCode)
						}
						return c.JSON(http.StatusNotFound, "No Policy as your request")
					}
				}
			}
		}
		return c.JSON(http.StatusBadRequest, "<uid>, <ctx>, <object>, and <rw> parameters must be provided")
	})

	e.Start(fSf(":%d", port))
}
