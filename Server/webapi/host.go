package webapi

import (
	"io/ioutil"
	"net/http"
	"os"

	eg "github.com/cdutwhu/n3-util/n3errs"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/middleware"
	enf "github.com/nsip/n3-privacy/Enforcer/process"
	cfg "github.com/nsip/n3-privacy/Server/config"
	"github.com/nsip/n3-privacy/Server/storage"
)

// HostHTTPAsync : Host a HTTP Server for providing policy json
func HostHTTPAsync() {
	e := echo.New()
	defer e.Close()

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
		Cfg      = env2Struct("Cfg", &cfg.Config{}).(*cfg.Config)
		port     = Cfg.WebService.Port
		fullIP   = localIP() + fSf(":%d", port)
		route    = Cfg.Route
		file     = Cfg.File
		database = Cfg.Storage.DataBase
		tracing  = Cfg.Storage.Tracing
		mMtx     = initMutex(route)
		db       = storage.NewDB(database, tracing).(storage.DBTr) // DBTr covers DB
	)

	// Tracing: Middleware for DB-tracing
	if tracing {
		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				db.SetContext(c.Request().Context())
				return next(c)
			}
		})
	}

	defer e.Start(fSf(":%d", port))

	// *************************************** List all API, FILE *************************************** //

	path := route.HELP
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		return c.String(
			http.StatusOK,
			fSf("wget %-55s-> %s\n", fullIP+"/enforcer-linux64", "Get Enforcer(linux64)")+
				fSf("wget %-55s-> %s\n", fullIP+"/enforcer-mac", "Get Enforcer(mac)")+
				fSf("wget %-55s-> %s\n", fullIP+"/enforcer-win64", "Get Enforcer(windows64)")+
				fSf("\n")+
				fSf("wget %-55s-> %s\n", fullIP+"/client-linux64", "Get Client(linux64)")+
				fSf("wget %-55s-> %s\n", fullIP+"/client-mac", "Get Client(mac)")+
				fSf("wget %-55s-> %s\n", fullIP+"/client-win64", "Get Client(windows64)")+
				fSf("wget -O config.toml %-40s-> %s\n", fullIP+"/client-config", "Get Client Config")+
				fSf("\n")+
				fSf("GET    %-55s-> %s\n", fullIP+route.GetID, "Get policy's ID. If no policy, return empty")+
				fSf("GET    %-55s-> %s\n", fullIP+route.GetHash, "Get policy's Hash String. If no policy, return empty")+
				fSf("GET    %-55s-> %s\n", fullIP+route.Get, "Get policy's JSON file. If no policy, return empty")+
				fSf("POST   %-55s-> %s\n", fullIP+route.Update, "Update policy. If no policy exists, add it")+
				fSf("DELETE %-55s-> %s\n", fullIP+route.Delete, "Delete policy")+
				fSf("\n")+
				fSf("GET    %-55s-> %s\n", fullIP+route.LsID, "Get a list of policy id. If no user or ctx restriction, return all policy id")+
				fSf("GET    %-55s-> %s\n", fullIP+route.LsUser, "Get a list of user. If no ctx restriction, return all user")+
				fSf("GET    %-55s-> %s\n", fullIP+route.LsContext, "Get a list of context. If no user restriction, return all context")+
				fSf("GET    %-55s-> %s\n", fullIP+route.LsObject, "Get a list of object. If no user or ctx restriction, return all object")+
				fSf("\n")+
				fSf("POST   %-55s-> %s\n", fullIP+route.Enforce, "Send json, return its enforced result. If its policy does not exist, return empty"))
	})

	// -------------------------------------------------------------------------- //

	mRouteRes := map[string]string{
		"/enforcer-linux64": file.EnforcerLinux64,
		"/enforcer-mac":     file.EnforcerMac,
		"/enforcer-win64":   file.EnforcerWin64,
		"/enforcer-config":  file.EnforcerConfig,
		"/client-linux64":   file.ClientLinux64,
		"/client-mac":       file.ClientMac,
		"/client-win64":     file.ClientWin64,
		"/client-config":    file.ClientConfig,
	}

	routeFun := func(rt, res string) func(c echo.Context) error {
		return func(c echo.Context) (err error) {
			if _, err = os.Stat(res); err == nil {
				fPln(rt, res)
				return c.File(res)
			}
			fPf("%v\n", warnOnErr("%v: [%s]  get [%s]", eg.FILE_NOT_FOUND, rt, res))
			return err
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

		qryParams := c.QueryParams()
		if ok, user, n3ctx, object, rw := url4Values(qryParams, 0, "user", "ctx", "object", "rw"); ok {
			if pid := db.PolicyID(user, n3ctx, rw, object); pid != "" {
				return c.JSON(http.StatusOK, result{
					Data:  pid,
					Empty: false,
					Error: "",
				})
			}
			return c.JSON(http.StatusNotFound, result{
				Data:  "",
				Empty: true,
				Error: "",
			})
		}
		return c.JSON(http.StatusBadRequest, result{
			Data:  "",
			Empty: true,
			Error: fSf("[user] [ctx] [object] [rw] all are required"),
		})
	})

	path = route.GetHash
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		qryParams := c.QueryParams()
		if ok, id := url1Value(qryParams, 0, "id"); ok {
			if hashstr, ok := db.PolicyHash(id); ok {
				return c.JSON(http.StatusOK, result{
					Data:  hashstr,
					Empty: false,
					Error: "",
				})
			}
			return c.JSON(http.StatusNotFound, result{
				Data:  "",
				Empty: true,
				Error: "",
			})
		}
		return c.JSON(http.StatusBadRequest, result{
			Data:  "",
			Empty: true,
			Error: fSf("policy [id] is required"),
		})
	})

	path = route.Get
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		qryParams := c.QueryParams()
		if ok, id := url1Value(qryParams, 0, "id"); ok {
			if policy, ok := db.Policy(id); ok {
				return c.JSON(http.StatusOK, result{
					Data:  policy,
					Empty: false,
					Error: "",
				})
			}
			return c.JSON(http.StatusNotFound, result{
				Data:  "",
				Empty: true,
				Error: "",
			})
		}
		return c.JSON(http.StatusBadRequest, result{
			Data:  "",
			Empty: true,
			Error: fSf("policy [id] is required"),
		})
	})

	path = route.Delete
	e.DELETE(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		qryParams := c.QueryParams()
		if ok, id := url1Value(qryParams, 0, "id"); ok {
			if db.DeletePolicy(id) == nil {
				// fPln(db.PolicyCount(), ": exist in db")
				return c.JSON(http.StatusOK, result{
					Data:  id,
					Empty: false,
					Error: "",
				})
			}
			return c.JSON(http.StatusInternalServerError, result{
				Data:  "",
				Empty: true,
				Error: fSf("Policy Delete Error"),
			})
		}
		return c.JSON(http.StatusBadRequest, result{
			Data:  "",
			Empty: true,
			Error: fSf("policy [id] is required"),
		})
	})

	path = route.Update
	e.POST(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		qryParams := c.QueryParams()
		req := c.Request()
		name, user, n3ctx, rw := "", "", "", ""
		if Ok, Name, User, n3Ctx, Rw := url4Values(qryParams, 0, "name", "user", "ctx", "rw"); Ok {
			name, user, n3ctx, rw = Name, User, n3Ctx, Rw
		} else if Ok, User, n3Ctx, Rw := url3Values(qryParams, 0, "user", "ctx", "rw"); Ok {
			user, n3ctx, rw = User, n3Ctx, Rw
		} else {
			return c.JSON(http.StatusBadRequest, result{
				Data:  "",
				Empty: true,
				Error: fSf("[user] [ctx] [rw] are required"),
			})
		}

		if bytes, err := ioutil.ReadAll(req.Body); err == nil && isJSON(string(bytes)) {
			if _, obj, err := db.UpdatePolicy(string(bytes), name, user, n3ctx, rw); err == nil {
				// fPln(db.PolicyCount(), ": exist in db")
				// return c.String(http.StatusOK, id+" - "+obj)
				return c.JSON(http.StatusOK, result{
					Data:  obj,
					Empty: false,
					Error: "",
				})
			}
			return c.JSON(http.StatusInternalServerError, result{
				Data:  "",
				Empty: true,
				Error: fSf("Update DB error"),
			})
		}
		return c.JSON(http.StatusBadRequest, result{
			Data:  "",
			Empty: true,
			Error: fSf("Policy is NOT in Request BODY, or invalid JSON"),
		})
	})

	// ---------------------------------------------------- Optional ---------------------------------------------------- //

	path = route.LsID
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		qryParams := c.QueryParams()
		if ok, user, n3ctx := url2Values(qryParams, 0, "user", "ctx"); ok {
			return c.JSON(http.StatusOK, db.MapRW2lsPID(user, n3ctx))
		}
		if ok, user := url1Value(qryParams, 0, "user"); ok {
			return c.JSON(http.StatusOK, db.MapRW2lsPID(user, ""))
		}
		if ok, n3ctx := url1Value(qryParams, 0, "ctx"); ok {
			return c.JSON(http.StatusOK, db.MapRW2lsPID("", n3ctx))
		}
		return c.JSON(http.StatusOK, db.MapRW2lsPID("", ""))
	})

	path = route.LsUser
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		qryParams := c.QueryParams()
		if ok, lsValues := urlValues(qryParams, "ctx"); ok {
			return c.JSON(http.StatusOK, db.MapCtx2lsUser(lsValues[0]...))
		}
		return c.JSON(http.StatusOK, db.MapCtx2lsUser())
	})

	path = route.LsContext
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		qryParams := c.QueryParams()
		if ok, lsValues := urlValues(qryParams, "user"); ok {
			return c.JSON(http.StatusOK, db.MapUser2lsCtx(lsValues[0]...))
		}
		return c.JSON(http.StatusOK, db.MapUser2lsCtx())
	})

	path = route.LsObject
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		qryParams := c.QueryParams()
		if ok, user, n3ctx := url2Values(qryParams, 0, "user", "ctx"); ok {
			return c.JSON(http.StatusOK, db.MapUC2lsObject(user, n3ctx))
		}
		if ok, user := url1Value(qryParams, 0, "user"); ok {
			return c.JSON(http.StatusOK, db.MapUC2lsObject(user, ""))
		}
		if ok, n3ctx := url1Value(qryParams, 0, "ctx"); ok {
			return c.JSON(http.StatusOK, db.MapUC2lsObject("", n3ctx))
		}
		return c.JSON(http.StatusOK, db.MapUC2lsObject("", ""))
	})

	// -------------------------------------------------------------------------- //

	path = route.Enforce
	e.POST(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		qryParams := c.QueryParams()
		req := c.Request()
		name, user, n3ctx, rw := "", "", "", ""
		if Ok, Name, User, n3Ctx, Rw := url4Values(qryParams, 0, "name", "user", "ctx", "rw"); Ok {
			name, user, n3ctx, rw = Name, User, n3Ctx, Rw
		} else if Ok, User, n3Ctx, Rw := url3Values(qryParams, 0, "user", "ctx", "rw"); Ok {
			user, n3ctx, rw = User, n3Ctx, Rw
		} else {
			return c.JSON(http.StatusBadRequest, result{
				Data:  "",
				Empty: true,
				Error: fSf("[user] [ctx] [rw] are required"),
			})
		}

		// get uploaded json and object
		json, object := "", name
		if bytes, err := ioutil.ReadAll(req.Body); err == nil {
			json = string(bytes)
			if isJSON(json) {
				if object == "" {
					object = jsonRoot(json)
				}
			} else {
				return c.JSON(http.StatusBadRequest, result{
					Data:  "",
					Empty: true,
					Error: fSf("POST Body Content is invalid JSON"),
				})
			}
		} else {
			return c.JSON(http.StatusBadRequest, result{
				Data:  "",
				Empty: true,
				Error: fSf("Error occurred when reading POST Body Content"),
			})
		}

		if pid := db.PolicyID(user, n3ctx, rw, object); pid != "" {
			if policy, ok := db.Policy(pid); ok {

				// ret := enf.Execute(json, policy)

				// Trace [enf.Execute]
				results := jaegertracing.TraceFunction(c, enf.Execute, json, policy)
				ret := results[0].Interface().(string)

				return c.JSON(http.StatusOK, result{
					Data:  ret,
					Empty: false,
					Error: "",
				})
			}
		}
		return c.JSON(http.StatusNotFound, result{
			Data:  "",
			Empty: true,
			Error: fSf("No policies@ user-[%s] context-[%s] read/write-[%s] object-[%s]", user, n3ctx, rw, object),
		})
	})
}
