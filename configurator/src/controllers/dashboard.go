package controllers

type DashboardController struct {
	BaseController
}

// URLMapping ...
func (c *DashboardController) URLMapping() {
	c.Mapping("Index", c.Index)
	c.Mapping("Login", c.Login)
	c.Mapping("Logout", c.Logout)
	c.Mapping("Recovery", c.Recovery)
	c.Mapping("Reset", c.Reset)
	c.Mapping("ResetPost", c.ResetPost)
}

func (h *DashboardController) Index() {

	h.Layout = h.GetTemplate() + "/private/base.tpl.html"
	h.TplName = h.GetTemplate() + "/private/index.tpl.html"
	h.UpdateTemplate()
}

func (h *DashboardController) Login() {

	h.Layout = h.GetTemplate() + "/public/base.tpl.html"
	h.TplName = h.GetTemplate() + "/public/login.tpl.html"
	h.UpdateTemplate()
}

func (h *DashboardController) Logout() {}

func (h *DashboardController) Recovery() {}

func (h *DashboardController) Reset() {}

func (h *DashboardController) ResetPost() {}