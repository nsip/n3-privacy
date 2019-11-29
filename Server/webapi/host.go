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
	path := "/"
	e.GET(path, func(c echo.Context) error {
		return c.String(
			http.StatusOK,
			fSf("GET    %-55s-> %s\n", fullIP+route.GetID, "Get policy's ID. If no policy, return empty")+
				fSf("GET    %-55s-> %s\n", fullIP+route.GetHash, "Get policy's Hash String. If no policy, return empty")+
				fSf("GET    %-55s-> %s\n", fullIP+route.Get, "Get policy's JSON file. If no policy, return empty")+
				fSf("POST   %-55s-> %s\n", fullIP+route.Update, "Update policy. If no policy exists, add it")+
				fSf("DELETE %-55s-> %s\n", fullIP+route.Delete, "Delete policy")+
				fSf("GET    %-55s-> %s\n", fullIP+route.ListOfPID, "Get a list of policy id. If no user or ctx restriction, return all policy id")+
				fSf("GET    %-55s-> %s\n", fullIP+route.ListOfUser, "Get a list of user. If no ctx restriction, return all user")+
				fSf("GET    %-55s-> %s\n", fullIP+route.ListOfCtx, "Get a list of context. If no user restriction, return all context")+
				fSf("GET    %-55s-> %s\n", fullIP+route.ListOfObject, "Get a list of object. If no user or ctx restriction, return all object"))
		// fSf("GET    %-55s-> %s\n", fullIP+route.GetJM, "Get JSON enforcement tool (jm). This tool is dependent on (jq)")+
		// fSf("GET    %-55s-> %s\n", fullIP+route.GetJQ, "Get JQ1.6. Put (jq) into (jm) directory")+
	})

	// -------------------------- Basic -------------------------- //

	path = route.GetID
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()
		glb.WDCheck()
		if ok, user, ctx, object, rw := url4Values(c.QueryParams(), 0, "user", "ctx", "object", "rw"); ok {
			if mCodes := db.PolicyID(user, ctx, rw, object); len(mCodes) > 0 {
				return c.JSON(http.StatusOK, mCodes)
			}
			return c.String(http.StatusNotFound, "No Policy as your request")
		}
		return c.String(http.StatusBadRequest, "<user>, <ctx>, <object>, and <rw> parameters must be provided")
	})

	path = route.GetHash
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()
		glb.WDCheck()
		if ok, id := url1Value(c.QueryParams(), 0, "id"); ok {
			if hashstr, ok := db.PolicyHash(id); ok {
				return c.String(http.StatusOK, hashstr)
			}
			return c.String(http.StatusNotFound, "No Policy as your request")
		}
		return c.String(http.StatusBadRequest, "policy <id> parameters must be provided")
	})

	path = route.Get
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()
		glb.WDCheck()
		if ok, id := url1Value(c.QueryParams(), 0, "id"); ok {
			if policy, ok := db.Policy(id); ok {
				return c.String(http.StatusOK, policy)
			}
			return c.String(http.StatusNotFound, "No Policy as your request")
		}
		return c.String(http.StatusBadRequest, "policy <id> parameters must be provided")
	})

	path = route.Delete
	e.DELETE(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()
		glb.WDCheck()
		if ok, id := url1Value(c.QueryParams(), 0, "id"); ok {
			if db.DeletePolicy(id) == nil {
				fPln(db.PolicyCount(), ": exist in db")
				return c.String(http.StatusOK, "Policy: "+id+" is deleted")
			}
			return c.String(http.StatusNotFound, "Policy Delete Error")
		}
		return c.String(http.StatusBadRequest, "policy <id> parameters must be provided")
	})

	path = route.Update
	e.POST(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()
		glb.WDCheck()
		if ok, user, ctx, rw := url3Values(c.QueryParams(), 0, "user", "ctx", "rw"); ok {
			if bPolicy, err := ioutil.ReadAll(c.Request().Body); err == nil && jkv.IsJSON(string(bPolicy)) {
				if id, obj, err := db.UpdatePolicy(string(bPolicy), user, ctx, rw); err == nil {
					fPln(db.PolicyCount(), ": exist in db")
					return c.String(http.StatusOK, id+" - "+obj)
				}
				return c.String(http.StatusInternalServerError, "Update DB error")
			}
			return c.String(http.StatusBadRequest, "Policy is not in BODY, or is not valid JSON")
		}
		return c.String(http.StatusBadRequest, "<user>, <ctx> and <rw> parameters must be provided")
	})

	// -------------------------- Optional -------------------------- //

	path = route.ListOfPID
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()
		glb.WDCheck()
		if ok, user, ctx := url2Values(c.QueryParams(), 0, "user", "ctx"); ok {
			return c.JSON(http.StatusOK, db.MapRWListOfPID(user, ctx))
		}
		if ok, user := url1Value(c.QueryParams(), 0, "user"); ok {
			return c.JSON(http.StatusOK, db.MapRWListOfPID(user, ""))
		}
		if ok, ctx := url1Value(c.QueryParams(), 0, "ctx"); ok {
			return c.JSON(http.StatusOK, db.MapRWListOfPID("", ctx))
		}
		return c.JSON(http.StatusOK, db.MapRWListOfPID("", ""))
	})

	path = route.ListOfUser
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()
		glb.WDCheck()
		if ok, lsValues := urlValues(c.QueryParams(), "ctx"); ok {
			return c.JSON(http.StatusOK, db.MapCtxListOfUser(lsValues[0]...))
		}
		return c.JSON(http.StatusOK, db.MapCtxListOfUser())
	})

	path = route.ListOfCtx
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()
		glb.WDCheck()
		if ok, lsValues := urlValues(c.QueryParams(), "user"); ok {
			return c.JSON(http.StatusOK, db.MapUserListOfCtx(lsValues[0]...))
		}
		return c.JSON(http.StatusOK, db.MapUserListOfCtx())
	})

	path = route.ListOfObject
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()
		glb.WDCheck()
		if ok, user, ctx := url2Values(c.QueryParams(), 0, "user", "ctx"); ok {
			return c.JSON(http.StatusOK, db.MapUCListOfObject(user, ctx))
		}
		if ok, user := url1Value(c.QueryParams(), 0, "user"); ok {
			return c.JSON(http.StatusOK, db.MapUCListOfObject(user, ""))
		}
		if ok, ctx := url1Value(c.QueryParams(), 0, "ctx"); ok {
			return c.JSON(http.StatusOK, db.MapUCListOfObject("", ctx))
		}
		return c.JSON(http.StatusOK, db.MapUCListOfObject("", ""))
	})

	e.Start(fSf(":%d", port))
}
