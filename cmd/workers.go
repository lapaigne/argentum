package main

import (
	"argentum/db"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func (e Endpoints) workers_GET(c echo.Context) error {
	if getClaims(c).Level < ACC_DISPATCHER {
		return c.Redirect(303, "/menu/")
	}

	if uws, ok := w_filter(c); ok {
		return c.Render(200, "workers", uws)
	}

	return c.Render(200, "workers", data.UWs)
}

func (e Endpoints) workers_POST(c echo.Context) error {
	if getClaims(c).Level < ACC_DISPATCHER {
		return c.Redirect(303, "/menu/")
	}

	url := c.Request().URL.Path

	if id, ok := strings.CutPrefix(url, "/workers/edit-"); ok {
		c.Set("id", id)
		return e.workers_edit(c)
	}

	if id, ok := strings.CutPrefix(url, "/workers/upd-"); ok {
		c.Set("id", id)
		return e.workers_upd(c)
	}

	if id, ok := strings.CutPrefix(url, "/workers/del-"); ok {
		c.Set("id", id)
		return e.workers_del(c)
	}

	switch url {
	case "/workers/add-panel":
		return e.workers_addpanel(c)
	case "/workers/add":
		return e.workers_add(c)
	case "/workers/tel":
		return e.workers_tel(c)
	case "/workers/gen":
		return e.workers_gen(c)
	case "/workers/filter":
		return e.workers_filter(c)
	default:
		return echo.ErrNotFound
	}
}

func (e Endpoints) workers_add(c echo.Context) error {
	f := c.FormValue("f_name")
	i := c.FormValue("i_name")
	o := c.FormValue("o_name")

	login := c.FormValue("login")
	pass := c.FormValue("pass")

	if !validLogin(login) {
		res := fmt.Sprintf(`<input type="text" name="login" id="workers-login" hx-post="/workers/tel" hx-target="#workers-login" hx-swap="outerHTML" placeholder="Логин" value="%s" class="reoc form-error">`, login)
		c.Response().Header().Set("HX-Retarget", "#workers-login")
		c.Response().Header().Set("HX-Reswap", "outerHTML")
		return c.HTML(200, res)
	}

	lvl, err := strconv.Atoi(c.FormValue("workers-lvl-sel"))
	if err != nil {
		return err
	}

	w := db.Worker{
		F_name: f,
		I_name: i,
		O_name: o,
	}

	hash, err := hash(pass)
	if err != nil {
		return err
	}

	u := db.User{
		Login: login,
		Hash:  string(hash),
		Level: lvl,
	}

	wid, err := db.AddWorker(w)
	if err != nil {
		return err
	}
	w.Id = wid
	u.Worker = wid

	uid, err := db.AddUser(u)
	if err != nil {
		return err
	}
	u.Id = uid

	uw := UserWorker{
		U:  u,
		W:  w,
		Id: w.Id,
	}

	data.UWs[w.Id] = &uw

	return c.Render(200, "workers", data.UWs)
}

func (e Endpoints) workers_gen(c echo.Context) error {
	pass, err := genPass(10)
	if err != nil {
		return err
	}
	res := fmt.Sprintf(`<input type="text" name="pass" id="workers-pass" value="%s">`, pass)
	return c.HTML(200, res)
}

func (e Endpoints) workers_del(c echo.Context) error {
	wid, err := strconv.Atoi(c.Get("id").(string))
	if err != nil {
		return err
	}

	uw, ok := data.UWs[wid]
	if !ok {
		return echo.ErrNotFound
	}

	uid := uw.U.Id

	if err := db.DeleteUser(uid); err != nil {
		return err
	}

	uw.SoftDelete()

	return c.Render(200, "workers-row", uw)
}

func (e Endpoints) workers_upd(c echo.Context) error {
	f := c.FormValue("f_name")
	i := c.FormValue("i_name")
	o := c.FormValue("o_name")
	l := c.FormValue("workers-lvl-sel")

	wid, err := strconv.Atoi(c.Get("id").(string))
	if err != nil {
		return err
	}

	uw := data.UWs[wid]
	uid := uw.U.Id

	level, err := strconv.Atoi(l)
	if err != nil {
		return err
	}

	d := UserWorker{
		W: db.Worker{
			F_name: f,
			I_name: i,
			O_name: o},
		U: db.User{
			Level: level,
		},
		Id: wid,
	}

	dw, du := Diffs(*uw, d)
	if dw {
		if err = db.UpdateWorker(wid, d.W); err != nil {
			fmt.Println(err)
		}
	}
	if du {
		if err = db.UpdateUserLevel(uid, d.U.Level); err != nil {
			fmt.Println(err)
		}
	}

	// TODO: warn if editing own account details

	data.UWs[wid].SoftReplace(d)

	return c.Render(200, "workers-row", d)
}

func (e Endpoints) workers_edit(c echo.Context) error {
	id, err := strconv.Atoi(c.Get("id").(string))
	if err != nil {
		return err
	}

	type Opt struct {
		Level int
		Name  string
		Sel   bool
	}

	uw := data.UWs[id]
	var opts []Opt
	switch uw.U.Level {
	case 0:
		opts = []Opt{
			{Level: 0, Name: "Отключенный аккаунт", Sel: true},
		}
	case 10:
		opts = []Opt{
			{Level: 10, Name: "Работник", Sel: true},
			{Level: 50, Name: "Диспетчер", Sel: false},
			{Level: 100, Name: "Админ", Sel: false},
		}
	case 50:
		opts = []Opt{
			{Level: 10, Name: "Работник", Sel: false},
			{Level: 50, Name: "Диспетчер", Sel: true},
			{Level: 100, Name: "Админ", Sel: false},
		}
	case 100:
		opts = []Opt{
			{Level: 10, Name: "Работник", Sel: false},
			{Level: 50, Name: "Диспетчер", Sel: false},
			{Level: 100, Name: "Админ", Sel: true},
		}
	}

	d := struct {
		W    db.Worker
		U    db.User
		Id   int
		Opts []Opt
	}{
		W:    uw.W,
		U:    uw.U,
		Id:   id,
		Opts: opts,
	}

	return c.Render(200, "workers-upd", d)
}

func (e Endpoints) workers_addpanel(c echo.Context) error {
	return c.Render(200, "workers-add", nil)
}

func (e Endpoints) workers_tel(c echo.Context) error {
	login := c.FormValue("login")

	str := fmt.Sprintf(`<input type="text" name="login" id="workers-login" hx-post="/workers/tel" hx-target="#workers-login" hx-swap="outerHTML" placeholder="Логин" value="%s"`, login)
	end := ">"

	if validLogin(login) {
		return c.HTML(200, str+end)
	}

	return c.HTML(200, str+`class=form-error`+end)
}

func validLogin(login string) bool {
	if len(login) < 4 {
		return false
	}

	if _, err := strconv.Atoi(login); err != nil {
		return false
	}

	if slices.Contains(data.Logins, login) {
		return false
	}

	return true
}

func w_filter(c echo.Context) ([]UserWorker, bool) {
	if c.FormValue("w-hide") == "1" {
		uws := []UserWorker{}
		for _, v := range data.UWs {
			if v.U.Level != ACC_NONE {
				uws = append(uws, *v)
			}
		}
		return uws, true
	}
	return nil, false
}

func (e Endpoints) workers_filter(c echo.Context) error {
	if uws, ok := w_filter(c); ok {
		return c.Render(200, "workers-tbody", uws)
	}

	return c.Render(200, "workers-tbody", data.UWs)
}
