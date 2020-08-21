package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/cdutwhu/n3-util/n3cfg"
	"github.com/cdutwhu/n3-util/n3err"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	enf "github.com/nsip/n3-privacy/Enforcer/process"
	"github.com/nsip/n3-privacy/Server/storage"
)

func shutdownAsync(e *echo.Echo, sig <-chan os.Signal, done chan<- string) {
	<-sig
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	failOnErr("%v", e.Shutdown(ctx))
	time.Sleep(20 * time.Millisecond)
	done <- "Shutdown Successfully"
}

// HostHTTPAsync : Host a HTTP Server for providing policy json
func HostHTTPAsync(sig <-chan os.Signal, done chan<- string) {
	defer func() { logGrp.Do("HostHTTPAsync Exit") }()

	e := echo.New()
	defer e.Close()

	// waiting for shutdown
	go shutdownAsync(e, sig, done)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("2G"))

	// Add Jaeger Tracer into Middleware
	c := jaegertracing.New(e, nil)
	defer c.Close()

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.DELETE},
		AllowCredentials: true,
	}))

	var (
		Cfg      = n3cfg.FromEnvN3privacyServer(envKey)
		port     = Cfg.WebService.Port
		fullIP   = localIP() + fSf(":%d", port)
		route    = Cfg.Route
		database = Cfg.Storage.DB
		mMtx     = initMutex(&route)
		db       = storage.NewDB(database)
	)

	// Tracing: Middleware for DB-tracing
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			db.SetContext(c.Request().Context())
			return next(c)
		}
	})

	logGrp.Do("Echo Service is Starting")
	defer e.Start(fSf(":%d", port))

	// *************************************** List all API, FILE *************************************** //

	path := route.HELP
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		return c.String(
			http.StatusOK,
			fSf("wget   %-55s-> %s\n", fullIP+"/enforcer-(linux64|mac|win64)", "Get Enforcer(linux64|mac|win64)")+
				fSf("GET    %-55s-> %s\n", fullIP+route.GetID, "Get policy's ID. If no policy, return empty")+
				fSf("GET    %-55s-> %s\n", fullIP+route.GetHash, "Get policy's Hash string. If no policy, return empty")+
				fSf("GET    %-55s-> %s\n", fullIP+route.Get, "Get policy's JSON file. If no policy, return empty")+
				fSf("POST   %-55s-> %s\n", fullIP+route.Update, "Update policy. If no policy, add it")+
				fSf("DELETE %-55s-> %s\n", fullIP+route.Delete, "Delete policy")+
				fSf("GET    %-55s-> %s\n", fullIP+route.LsID, "Get list of policy id. If no user or ctx, return all policy id")+
				fSf("GET    %-55s-> %s\n", fullIP+route.LsUser, "Get list of user. If no ctx, return all user")+
				fSf("GET    %-55s-> %s\n", fullIP+route.LsContext, "Get list of context. If no user, return all context")+
				fSf("GET    %-55s-> %s\n", fullIP+route.LsObject, "Get list of object. If neither user nor ctx, return all object")+
				fSf("POST   %-55s-> %s\n", fullIP+route.Enforce, "Send json, return enforced result. If no policy, return empty"))
	})

	// -------------------------------------------------------------------------- //

	mRouteRes := map[string]string{
		"/enforcer-linux64": Cfg.File.EnforcerLinux64,
		"/enforcer-mac":     Cfg.File.EnforcerMac,
		"/enforcer-win64":   Cfg.File.EnforcerWin64,
	}

	routeFun := func(rt, res string) func(c echo.Context) error {
		return func(c echo.Context) (err error) {
			if _, err = os.Stat(res); err == nil {
				fPln(rt, res)
				return c.File(res)
			}
			return warnOnErr("%v: [%s]  get [%s]", n3err.FILE_NOT_FOUND, rt, res)
		}
	}

	for rt, res := range mRouteRes {
		e.GET(rt, routeFun(rt, res))
	}

	// ---------------------------------------------------- Basic ---------------------------------------------------- //

	path = route.GetID
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		var (
			status = http.StatusOK
			ret    string
		)

		logGrp.Do("Parsing Params")
		ok, user, n3ctx, object, rw := url4Values(c.QueryParams(), 0, "user", "ctx", "object", "rw")
		if !ok {
			status = http.StatusBadRequest
			ret = "[user], [ctx], [object], [rw] all are required"
			goto RET
		}
		if ret = mustInvokeWithMW(db, "PolicyID", user, n3ctx, rw, object)[0].(string); ret == "" {
			status = http.StatusNotFound
			ret = fSf("No Policy ID @ user:[%s] ctx:[%s] rw:[%s] object:[%s]", user, n3ctx, rw, object)
			goto RET
		}
	RET:
		if status != http.StatusOK {
			warnGrp.Do(ret)
		} else {
			logGrp.Do("--> Finish GetID")
		}
		return c.String(status, ret)
	})

	// ------------------------------------- //

	path = route.GetHash
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		var (
			status = http.StatusOK
			ret    string
			rets   []interface{}
		)

		logGrp.Do("Parsing Params")
		ok, id := url1Value(c.QueryParams(), 0, "id")
		if !ok {
			status = http.StatusBadRequest
			ret = "policy [id] is required"
			goto RET
		}
		rets = mustInvokeWithMW(db, "PolicyHash", id)
		ret, ok = rets[0].(string), rets[1].(bool)
		if !ok || ret == "" {
			status = http.StatusNotFound
			ret = fSf("No Policy Hash @ id:[%s]", id)
			goto RET
		}
	RET:
		if status != http.StatusOK {
			warnGrp.Do(ret)
		} else {
			logGrp.Do("--> Finish GetHash")
		}
		return c.String(status, ret)
	})

	// ------------------------------------- //

	path = route.Get
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		var (
			status = http.StatusOK
			ret    string
			rets   []interface{}
		)

		logGrp.Do("Parsing Params")
		ok, id := url1Value(c.QueryParams(), 0, "id")
		if !ok {
			status = http.StatusBadRequest
			ret = "policy [id] is required"
			goto RET
		}
		rets = mustInvokeWithMW(db, "Policy", id)
		ret, ok = rets[0].(string), rets[1].(bool)
		if !ok || ret == "" {
			status = http.StatusNotFound
			ret = fSf("No Policy Found @ id:[%s]", id)
			goto RET
		}
	RET:
		if status != http.StatusOK {
			warnGrp.Do(ret)
		} else {
			logGrp.Do("--> Finish GetPolicy")
		}
		return c.String(status, ret)
	})

	// ------------------------------------- //

	path = route.Delete
	e.DELETE(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		var (
			status = http.StatusOK
			ret    string
		)

		logGrp.Do("Parsing Params")
		ok, id := url1Value(c.QueryParams(), 0, "id")
		if !ok {
			status = http.StatusBadRequest
			ret = "policy [id] is required"
			goto RET
		}
		if mustInvokeWithMW(db, "DeletePolicy", id)[0] != nil {
			status = http.StatusInternalServerError
			ret = "Policy Delete Error"
			goto RET
		}
	RET:
		if status != http.StatusOK {
			warnGrp.Do(ret)
		} else {
			logGrp.Do("--> Finish Delete Policy")
		}
		return c.String(status, ret)
	})

	// ------------------------------------- //

	path = route.Update
	e.POST(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		var (
			status  = http.StatusOK
			ret     string
			rets    []interface{}
			body    []byte
			err     error
			jsonstr string
		)

		logGrp.Do("Parsing Params")
		qryParams, req := c.QueryParams(), c.Request()
		name, user, n3ctx, rw := "", "", "", ""
		if ok, Name, User, n3Ctx, Rw := url4Values(qryParams, 0, "name", "user", "ctx", "rw"); ok {
			name, user, n3ctx, rw = Name, User, n3Ctx, Rw
		} else if ok, User, n3Ctx, Rw := url3Values(qryParams, 0, "user", "ctx", "rw"); ok {
			user, n3ctx, rw = User, n3Ctx, Rw
		} else {
			status = http.StatusBadRequest
			ret = "At least [user] [ctx] [rw] are required"
			goto RET
		}

		logGrp.Do("Reading Request Body")
		body, err = ioutil.ReadAll(req.Body)
		jsonstr = string(body)
		if err != nil || !isJSON(jsonstr) {
			status = http.StatusBadRequest
			ret = "Policy is NOT in Request BODY, or invalid JSON"
			goto RET
		}

		logGrp.Do("UpdatePolicy")
		rets = mustInvokeWithMW(db, "UpdatePolicy", jsonstr, name, user, n3ctx, rw)
		if rets[2] != nil {
			status = http.StatusInternalServerError
			ret = "Update DB error"
			goto RET
		}

		ret = rets[0].(string) + " @ " + rets[1].(string)

	RET:
		if status != http.StatusOK {
			warnGrp.Do(ret)
		} else {
			logGrp.Do("--> Finish Update Policy")
		}
		return c.String(status, ret)
	})

	// ------------------------------------------ Optional ------------------------------------------ //

	path = route.LsID
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		fn := "MapRW2lsPID"
		qryParams := c.QueryParams()
		if ok, user, n3ctx := url2Values(qryParams, 0, "user", "ctx"); ok {
			rets := mustInvokeWithMW(db, fn, user, n3ctx)
			return c.JSON(http.StatusOK, rets[0])
		}
		if ok, user := url1Value(qryParams, 0, "user"); ok {
			rets := mustInvokeWithMW(db, fn, user, "")
			return c.JSON(http.StatusOK, rets[0])
		}
		if ok, n3ctx := url1Value(qryParams, 0, "ctx"); ok {
			rets := mustInvokeWithMW(db, fn, "", n3ctx)
			return c.JSON(http.StatusOK, rets[0])
		}
		rets := mustInvokeWithMW(db, fn, "", "")
		return c.JSON(http.StatusOK, rets[0])
	})

	path = route.LsUser
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		fn := "MapCtx2lsUser"
		qryParams := c.QueryParams()
		if ok, lsValues := urlValues(qryParams, "ctx"); ok {
			rets := mustInvokeWithMW(db, fn, toGeneralSlc(lsValues[0])...)
			return c.JSON(http.StatusOK, rets[0])
		}
		rets := mustInvokeWithMW(db, fn)
		return c.JSON(http.StatusOK, rets[0])
	})

	path = route.LsContext
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		fn := "MapUser2lsCtx"
		qryParams := c.QueryParams()
		if ok, lsValues := urlValues(qryParams, "user"); ok {
			rets := mustInvokeWithMW(db, fn, toGeneralSlc(lsValues[0])...)
			return c.JSON(http.StatusOK, rets[0])
		}
		rets := mustInvokeWithMW(db, fn)
		return c.JSON(http.StatusOK, rets[0])
	})

	path = route.LsObject
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		fn := "MapUC2lsObject"
		qryParams := c.QueryParams()
		if ok, user, n3ctx := url2Values(qryParams, 0, "user", "ctx"); ok {
			rets := mustInvokeWithMW(db, fn, user, n3ctx)
			return c.JSON(http.StatusOK, rets[0])
		}
		if ok, user := url1Value(qryParams, 0, "user"); ok {
			rets := mustInvokeWithMW(db, fn, user, "")
			return c.JSON(http.StatusOK, rets[0])
		}
		if ok, n3ctx := url1Value(qryParams, 0, "ctx"); ok {
			rets := mustInvokeWithMW(db, fn, "", n3ctx)
			return c.JSON(http.StatusOK, rets[0])
		}
		rets := mustInvokeWithMW(db, fn, "", "")
		return c.JSON(http.StatusOK, rets[0])
	})

	// -------------------------------------------------------------------------- //

	path = route.Enforce
	e.POST(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		var (
			status = http.StatusOK
			ret    string
			rets   []interface{}

			body   []byte
			err    error
			json   string
			pid    string
			policy string
		)

		logGrp.Do("Parsing Params")
		qryParams, req := c.QueryParams(), c.Request()
		ok, object, user, n3ctx, rw := url4Values(qryParams, 0, "name", "user", "ctx", "rw")
		if !ok {
			ok, user, n3ctx, rw = url3Values(qryParams, 0, "user", "ctx", "rw")
			if !ok {
				status = http.StatusBadRequest
				ret = "At least [user] [ctx] [rw] are required"
				goto RET
			}
		}

		body, err = ioutil.ReadAll(req.Body)
		if err != nil {
			status = http.StatusBadRequest
			ret = "Error in reading Request Body"
			goto RET
		}
		if json = string(body); !isJSON(json) {
			status = http.StatusBadRequest
			ret = "Request Body is invalid JSON"
			goto RET
		}
		if object == "" {
			object = jsonRoot(json)
		}

		rets = mustInvokeWithMW(db, "PolicyID", user, n3ctx, rw, object)
		if pid = rets[0].(string); pid == "" {
			status = http.StatusNotFound
			ret = fSf("No policies @ user:[%s] context:[%s] read/write:[%s] object:[%s]", user, n3ctx, rw, object)
			goto RET
		}

		rets = mustInvokeWithMW(db, "Policy", pid)
		policy, ok = rets[0].(string), rets[1].(bool)
		if !ok {
			status = http.StatusInternalServerError
			ret = fSf("Could NOT get policy @ pid %s", pid)
			goto RET
		}

		ret = jaegertracing.TraceFunction(c, enf.Execute, json, policy)[0].Interface().(string)

	RET:
		if status != http.StatusOK {
			warnGrp.Do(ret)
		} else {
			logGrp.Do("--> Finish Enforce")
		}
		return c.String(status, ret)
	})
}
