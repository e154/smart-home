package controllers

//
// Beego supports 404, 401, 403, 500, 503 error handling by default.
// You can also define custom error handling page.
//
// http://beego.me/docs/mvc/controller/errors.md
//
type ErrorController struct {
	CommonController
}

func (c *ErrorController) Prepare() {

}

func (c *ErrorController) Error401() {
	c.ErrHan(401, "Forbidden")
}

func (c *ErrorController) Error403() {
	c.ErrHan(403, "Forbidden")
}

func (c *ErrorController) Error404() {
	c.ErrHan(404, "Page Not Found")
}

func (c *ErrorController) Error500() {
	c.ErrHan(500, "internal server error")
}

func (c *ErrorController) ErrorDb() {
	c.ErrHan(500, "database is now down")
}

