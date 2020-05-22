package webapi

import (
	"io/ioutil"
	"net/http"
	"os"

	eg "github.com/cdutwhu/json-util/n3errs"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/middleware"
	enf "github.com/nsip/n3-privacy/Enforcer/process"
	glb "github.com/nsip/n3-privacy/Server/global"
)

// HostHTTPAsync : Host a HTTP Server for providing policy json
func HostHTTPAsync() {
	e := echo.New()

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

	cfg := glb.Cfg
	port := cfg.WebService.Port
	fullIP := localIP() + fSf(":%d", port)
	route := cfg.Route
	file := cfg.File

	initMutex()
	initDB()

	defer e.Start(fSf(":%d", port))

	// *************************************** List all API, FILE *************************************** //

	path := route.HELP
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		return c.String(
			http.StatusOK,
			fSf("wget %-55s-> %s\n", fullIP+"/enforcer-linux64", "Get Enforcer(Linux64)")+
				fSf("wget %-55s-> %s\n", fullIP+"/enforcer-mac", "Get Enforcer(Mac)")+
				fSf("wget %-55s-> %s\n", fullIP+"/enforcer-win64", "Get Enforcer(Windows64)")+
				fSf("wget -O config.toml %-40s-> %s\n", fullIP+"/enforcer-config", "Get Enforcer config")+
				fSf("\n")+
				fSf("wget %-55s-> %s\n", fullIP+"/client-linux64", "Get Client(Linux64)")+
				fSf("wget %-55s-> %s\n", fullIP+"/client-mac", "Get Client(Mac)")+
				fSf("wget %-55s-> %s\n", fullIP+"/client-win64", "Get Client(Windows64)")+
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
				fSf("POST   %-55s-> %s\n", fullIP+route.GetEnforced, "Send json, return its enforced result. If its policy does not exist, return empty"))
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
		if ok, user, ctx, object, rw := url4Values(c.QueryParams(), 0, "user", "ctx", "object", "rw"); ok {
			if pid := db.PolicyID(user, ctx, rw, object); pid != "" {
				return c.JSON(http.StatusOK, result{
					Data:  &pid,
					Empty: False,
					Error: "",
				})
			}
			return c.JSON(http.StatusNotFound, result{
				Data:  EmptyStr,
				Empty: True,
				Error: "",
			})
		}
		return c.JSON(http.StatusBadRequest, result{
			Data:  nil,
			Empty: nil,
			Error: "[user], [ctx], [object], and [rw] must be provided",
		})
	})

	path = route.GetHash
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()
		if ok, id := url1Value(c.QueryParams(), 0, "id"); ok {
			if hashstr, ok := db.PolicyHash(id); ok {
				return c.JSON(http.StatusOK, result{
					Data:  &hashstr,
					Empty: False,
					Error: "",
				})
			}
			return c.JSON(http.StatusNotFound, result{
				Data:  EmptyStr,
				Empty: True,
				Error: "",
			})
		}
		return c.JSON(http.StatusBadRequest, result{
			Data:  nil,
			Empty: nil,
			Error: "policy [id] must be provided",
		})
	})

	path = route.Get
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()
		if ok, id := url1Value(c.QueryParams(), 0, "id"); ok {
			if policy, ok := db.Policy(id); ok {
				return c.JSON(http.StatusOK, result{
					Data:  &policy,
					Empty: False,
					Error: "",
				})
			}
			return c.JSON(http.StatusNotFound, result{
				Data:  EmptyStr,
				Empty: True,
				Error: "",
			})
		}
		return c.JSON(http.StatusBadRequest, result{
			Data:  nil,
			Empty: nil,
			Error: "policy [id] must be provided",
		})
	})

	path = route.Delete
	e.DELETE(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()
		if ok, id := url1Value(c.QueryParams(), 0, "id"); ok {
			if db.DeletePolicy(id) == nil {
				fPln(db.PolicyCount(), ": exist in db")
				return c.JSON(http.StatusOK, result{
					Data:  &id,
					Empty: nil,
					Error: "",
				})
			}
			return c.JSON(http.StatusInternalServerError, result{
				Data:  nil,
				Empty: nil,
				Error: "Policy Delete Error",
			})
		}
		return c.JSON(http.StatusBadRequest, result{
			Data:  nil,
			Empty: nil,
			Error: "policy [id] must be provided",
		})
	})

	path = route.Update
	e.POST(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		name, user, ctx, rw := "", "", "", ""
		if Ok, Name, User, Ctx, Rw := url4Values(c.QueryParams(), 0, "name", "user", "ctx", "rw"); Ok {
			name, user, ctx, rw = Name, User, Ctx, Rw
		} else if Ok, User, Ctx, Rw := url3Values(c.QueryParams(), 0, "user", "ctx", "rw"); Ok {
			user, ctx, rw = User, Ctx, Rw
		} else {
			return c.JSON(http.StatusBadRequest, result{
				Data:  nil,
				Empty: nil,
				Error: "at least, [user], [ctx] and [rw] must be provided",
			})
		}

		if bytes, err := ioutil.ReadAll(c.Request().Body); err == nil && isJSON(string(bytes)) {
			if id, _, err := db.UpdatePolicy(string(bytes), name, user, ctx, rw); err == nil {
				fPln(db.PolicyCount(), ": exist in db")
				// return c.String(http.StatusOK, id+" - "+obj)
				return c.JSON(http.StatusOK, result{
					Data:  &id,
					Empty: nil,
					Error: "",
				})
			}
			return c.JSON(http.StatusInternalServerError, result{
				Data:  nil,
				Empty: nil,
				Error: "Update DB error",
			})
		}
		return c.JSON(http.StatusBadRequest, result{
			Data:  nil,
			Empty: nil,
			Error: "Policy is NOT in Request BODY, or NOT valid JSON",
		})

		// if ok, name, user, ctx, rw := url4Values(c.QueryParams(), 0, "name", "user", "ctx", "rw"); ok {
		// 	if bytes, err := ioutil.ReadAll(c.Request().Body); err == nil && isJSON(string(bytes)) {
		// 		if id, _, err := db.UpdatePolicy(string(bytes), name, user, ctx, rw); err == nil {
		// 			fPln(db.PolicyCount(), ": exist in db")
		// 			// return c.String(http.StatusOK, id+" - "+obj)
		// 			return c.JSON(http.StatusOK, result{
		// 				Data:  &id,
		// 				Empty: nil,
		// 				Error: "",
		// 			})
		// 		}
		// 		return c.JSON(http.StatusInternalServerError, result{
		// 			Data:  nil,
		// 			Empty: nil,
		// 			Error: "Update DB error",
		// 		})
		// 	}
		// 	return c.JSON(http.StatusBadRequest, result{
		// 		Data:  nil,
		// 		Empty: nil,
		// 		Error: "Policy is NOT in Request BODY, or NOT valid JSON",
		// 	})
		// }

		// return c.JSON(http.StatusBadRequest, result{
		// 	Data:  nil,
		// 	Empty: nil,
		// 	Error: "[user], [ctx] and [rw] must be provided",
		// })
	})

	// ---------------------------------------------------- Optional ---------------------------------------------------- //

	path = route.LsID
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()
		if ok, user, ctx := url2Values(c.QueryParams(), 0, "user", "ctx"); ok {
			return c.JSON(http.StatusOK, db.MapRW2lsPID(user, ctx))
		}
		if ok, user := url1Value(c.QueryParams(), 0, "user"); ok {
			return c.JSON(http.StatusOK, db.MapRW2lsPID(user, ""))
		}
		if ok, ctx := url1Value(c.QueryParams(), 0, "ctx"); ok {
			return c.JSON(http.StatusOK, db.MapRW2lsPID("", ctx))
		}
		return c.JSON(http.StatusOK, db.MapRW2lsPID("", ""))
	})

	path = route.LsUser
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()
		if ok, lsValues := urlValues(c.QueryParams(), "ctx"); ok {
			return c.JSON(http.StatusOK, db.MapCtx2lsUser(lsValues[0]...))
		}
		return c.JSON(http.StatusOK, db.MapCtx2lsUser())
	})

	path = route.LsContext
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()
		if ok, lsValues := urlValues(c.QueryParams(), "user"); ok {
			return c.JSON(http.StatusOK, db.MapUser2lsCtx(lsValues[0]...))
		}
		return c.JSON(http.StatusOK, db.MapUser2lsCtx())
	})

	path = route.LsObject
	e.GET(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()
		if ok, user, ctx := url2Values(c.QueryParams(), 0, "user", "ctx"); ok {
			return c.JSON(http.StatusOK, db.MapUC2lsObject(user, ctx))
		}
		if ok, user := url1Value(c.QueryParams(), 0, "user"); ok {
			return c.JSON(http.StatusOK, db.MapUC2lsObject(user, ""))
		}
		if ok, ctx := url1Value(c.QueryParams(), 0, "ctx"); ok {
			return c.JSON(http.StatusOK, db.MapUC2lsObject("", ctx))
		}
		return c.JSON(http.StatusOK, db.MapUC2lsObject("", ""))
	})

	// -------------------------------------------------------------------------- //

	path = route.GetEnforced
	e.POST(path, func(c echo.Context) error {
		defer func() { mMtx[path].Unlock() }()
		mMtx[path].Lock()

		name, user, ctx, rw := "", "", "", ""
		if Ok, Name, User, Ctx, Rw := url4Values(c.QueryParams(), 0, "name", "user", "ctx", "rw"); Ok {
			name, user, ctx, rw = Name, User, Ctx, Rw
		} else if Ok, User, Ctx, Rw := url3Values(c.QueryParams(), 0, "user", "ctx", "rw"); Ok {
			user, ctx, rw = User, Ctx, Rw
		} else {
			return c.JSON(http.StatusBadRequest, result{
				Data:  nil,
				Empty: nil,
				Error: "at least, [user], [ctx] and [rw] must be provided",
			})
		}

		// get uploaded json and object
		json, object := "", name
		if bytes, err := ioutil.ReadAll(c.Request().Body); err == nil {
			json = string(bytes)
			if isJSON(json) {
				if object == "" {
					object = jsonRoot(json)
				}
			} else {
				return c.JSON(http.StatusBadRequest, result{
					Data:  nil,
					Empty: nil,
					Error: "POST Body Content is invalid JSON",
				})
			}
		} else {
			return c.JSON(http.StatusBadRequest, result{
				Data:  nil,
				Empty: nil,
				Error: "Error occurred when reading POST Body Content",
			})
		}

		if pid := db.PolicyID(user, ctx, rw, object); pid != "" {
			if policy, ok := db.Policy(pid); ok {

				// ret := enf.Execute(json, policy)

				// Trace [enf.Execute]
				results := jaegertracing.TraceFunction(c, enf.Execute, json, policy)
				ret := results[0].Interface().(string)

				return c.JSON(http.StatusOK, result{
					Data:  &ret,
					Empty: False,
					Error: "",
				})
			}
		}
		return c.JSON(http.StatusNotFound, result{
			Data:  EmptyStr,
			Empty: True,
			Error: fSf("No policies for uploaded JSON @ user-[%s] context-[%s] read/write-[%s] object-[%s]", user, ctx, rw, object),
		})
	})

}
