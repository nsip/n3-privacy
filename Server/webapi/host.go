package webapi

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	glb "github.com/nsip/n3-privacy/Server/global"
	cmn "github.com/nsip/n3-privacy/common"
	"github.com/nsip/n3-privacy/jkv"
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

	port := glb.Cfg.WebService.Port
	fullIP := cmn.LocalIP() + fSf(":%d", port)
	route := glb.Cfg.Route
	initMutex()
	initDB()

	// *************************************** List all APP, API *************************************** //
	e.GET("/", func(c echo.Context) error {
		return c.String(
			http.StatusOK,
			fSf("GET    %-50s->  %s\n", fullIP+route.GetID, "Get policy's Fetch ID. If no policy, return empty")+
				fSf("GET    %-50s->  %s\n", fullIP+route.GetHash, "Get policy's Hash String. If no policy, return empty")+
				fSf("GET    %-50s->  %s\n", fullIP+route.Get, "Get policy's JSON file. If no policy, return empty")+
				fSf("POST   %-50s->  %s\n", fullIP+route.Update, "Update policy. If no policy exists, add it")+
				fSf("GET    %-50s->  %s\n", fullIP+route.GetJM, "Get JSON enforcement tool (jm). This tool is dependent on (jq)")+
				fSf("GET    %-50s->  %s\n", fullIP+route.GetJQ, "Get JQ1.6. Put (jq) into (jm) directory"))
	})

	e.GET(route.GetID, func(c echo.Context) error {
		defer func() { mMtx[route.GetID].Unlock() }()
		mMtx[route.GetID].Lock()
		glb.WDCheck()
		if ok, uid, ctx, object, rw := url1stValuesOf4(c.QueryParams(), "uid", "ctx", "object", "rw"); ok {
			if mCodes := db.PolicyID(uid, ctx, rw, object); len(mCodes) > 0 {
				return c.JSON(http.StatusOK, mCodes)
			}
			return c.String(http.StatusNotFound, "No Policy as your request")
		}
		return c.String(http.StatusBadRequest, "<uid>, <ctx>, <object>, and <rw> parameters must be provided")
	})

	e.GET(route.GetHash, func(c echo.Context) error {
		defer func() { mMtx[route.GetHash].Unlock() }()
		mMtx[route.GetHash].Lock()
		glb.WDCheck()
		if ok, id := url1stValuesOf1(c.QueryParams(), "id"); ok {
			if hashstr, ok := db.PolicyHash(id); ok {
				return c.String(http.StatusOK, hashstr)
			}
			return c.String(http.StatusNotFound, "No Policy as your request")
		}
		return c.String(http.StatusBadRequest, "policy <id> parameters must be provided")
	})

	e.GET(route.Get, func(c echo.Context) error {
		defer func() { mMtx[route.Get].Unlock() }()
		mMtx[route.Get].Lock()
		glb.WDCheck()
		if ok, id := url1stValuesOf1(c.QueryParams(), "id"); ok {
			if policy, ok := db.Policy(id); ok {
				return c.String(http.StatusOK, policy)
			}
			return c.String(http.StatusNotFound, "No Policy as your request")
		}
		return c.String(http.StatusBadRequest, "policy <id> parameters must be provided")
	})

	e.POST(route.Update, func(c echo.Context) error {
		defer func() { mMtx[route.Update].Unlock() }()
		mMtx[route.Update].Lock()
		glb.WDCheck()
		if ok, uid, ctx, rw := url1stValuesOf3(c.QueryParams(), "uid", "ctx", "rw"); ok {
			if bPolicy, err := ioutil.ReadAll(c.Request().Body); err == nil && jkv.IsJSON(string(bPolicy)) {
				if id, err := db.UpdatePolicy(string(bPolicy), uid, ctx, rw); err == nil {
					return c.String(http.StatusOK, id)
				}
				return c.String(http.StatusInternalServerError, "Update DB error")
			}
			return c.String(http.StatusBadRequest, "Policy is not in BODY, or is not valid JSON")
		}
		return c.String(http.StatusBadRequest, "<uid>, <ctx> and <rw> parameters must be provided")
	})

	e.Start(fSf(":%d", port))
}
