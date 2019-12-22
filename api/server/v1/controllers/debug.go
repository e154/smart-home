package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http/pprof"
)

type ControllerDebug struct {
	*ControllerCommon
}

func NewControllerDebug(common *ControllerCommon) *ControllerDebug {
	return &ControllerDebug{ControllerCommon: common}
}

// swagger:operation GET /debug/pprof debugPprofIndex
// ---
// summary:
// parameters:
// description: Index responds with the pprof-formatted profile named by the request.  For example, "/debug/pprof/heap" serves the "heap" profile.  Index responds to a request for "/debug/pprof/" with an HTML page  listing the available profiles.
// tags:
// - debug
// responses:
//   "200":
//     description: Success
//     examples:
//       text/html:
//         "<html><body>some text</body></html>"
func (d *ControllerDebug) PprofIndex(ctx *gin.Context) {
	pprof.Index(ctx.Writer, ctx.Request)
}

// swagger:operation GET /debug/pprof/cmdline debugPprofCmdline
// ---
// summary:
// parameters:
// description: Cmdline responds with the running program's command line, with arguments separated by NUL bytes. The package initialization registers it as /debug/pprof/cmdline.
// tags:
// - debug
// responses:
//   "200":
//     description: Success
//     examples:
//       text/html:
//         "<html><body>some text</body></html>"
func (d *ControllerDebug) PprofCmdline(ctx *gin.Context) {
	pprof.Cmdline(ctx.Writer, ctx.Request)
}

// swagger:operation GET /debug/pprof/profile/{handler} debugPprofProfile
// ---
// summary:
// parameters:
// - description: handler
//   in: path
//   name: handler
//   required: false
//   type: string
// description: Profile responds with the pprof-formatted cpu profile. Profiling lasts for duration specified in seconds GET parameter, or for 30 seconds if not specified. The package initialization registers it as /debug/pprof/profile.
// tags:
// - debug
// responses:
//   "200":
//     description: Success
//     examples:
//       text/html:
//         "<html><body>some text</body></html>"
//   "404":
//	   $ref: '#/responses/Error'
//
func (d *ControllerDebug) PprofProfile(ctx *gin.Context) {
	handler := ctx.Param("handler")
	switch handler {
	case "":
	case "goroutine", "threadcreate", "heap", "block", "mutex":
	default:
		err404 := NewError(404)
		err404.AddField("common.field_not_valid", "Handler not found", "handler")
		err404.Send(ctx)
		return
	}
	if handler == "" {
		pprof.Profile(ctx.Writer, ctx.Request)
	} else {
		pprof.Handler(handler).ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// swagger:operation GET /debug/pprof/symbol debugPprofSymbol
// ---
// summary:
// parameters:
// description: Symbol looks up the program counters listed in the request, responding with a table mapping program counters to function names. The package initialization registers it as /debug/pprof/symbol.
// tags:
// - debug
// responses:
//   "200":
//     description: Success
//     examples:
//       text/html:
//         "<html><body>some text</body></html>"
func (d *ControllerDebug) PprofSymbol(ctx *gin.Context) {
	pprof.Symbol(ctx.Writer, ctx.Request)
}

// swagger:operation GET /debug/pprof/trace debugPprofTrace
// ---
// summary:
// parameters:
// description: Trace responds with the execution trace in binary form. Tracing lasts for duration specified in seconds GET parameter, or for 1 second if not specified. The package initialization registers it as /debug/pprof/trace.
// tags:
// - debug
// responses:
//   "200":
//     description: Success
//     examples:
//       text/html:
//         "<html><body>some text</body></html>"
func (d *ControllerDebug) PprofTrace(ctx *gin.Context) {
	pprof.Trace(ctx.Writer, ctx.Request)
}
