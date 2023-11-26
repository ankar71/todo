package main

import (
	// "fmt"

	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/ankar71/todo/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const tmplDir = "templates/"
const tmplFile = "main.tmpl.html"
const mainPage = "main-page"

const taskList = "task-list"

const combo = "combo"

type TaskInput struct {
	Task string `form:"task"`
}

type ViewModel struct {
	TaskList  []*model.Todo
	AllDone   bool
	ItemsLeft int
}

func main() {
	e := echo.New()

	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// e.Logger.SetLevel(log.DEBUG)

	db := model.DefaultDict()
	e.Renderer = parseTmpl()

	currentViewModel := func() *ViewModel {
		return &ViewModel{
			TaskList:  db.List(),
			AllDone:   db.AllDone(),
			ItemsLeft: db.ItemsLeft(),
		}
	}
	e.GET("/", func(ctx echo.Context) error {
		return ctx.Render(http.StatusOK, mainPage, currentViewModel())
	})

	e.GET("/show/:only", func(ctx echo.Context) error {
		only := ctx.Param("only")
		switch only {
		case "active":
			return ctx.Render(http.StatusOK, taskList, db.FilteredList(false))
		case "completed":
			return ctx.Render(http.StatusOK, taskList, db.FilteredList(true))
		}
		return ctx.Render(http.StatusOK, taskList, db.List())
	})

	e.POST("/task", func(ctx echo.Context) error {
		task := new(TaskInput)

		if err := ctx.Bind(task); err != nil {
			return err
		}
		if task.Task != "" {
			db.AddTask(task.Task)
		}
		return ctx.Render(http.StatusOK, combo, currentViewModel())
	})

	e.PUT("/task/:id/status", func(ctx echo.Context) error {
		// read response
		id := ctx.Param("id")
		chkbox := "checkbox-" + id
		isDone := ctx.FormValue(chkbox) == "on"

		// update db
		i64, _ := strconv.ParseInt(id, 10, strconv.IntSize)
		db.SetStatus(int(i64), isDone)
		return ctx.Render(http.StatusOK, combo, currentViewModel())
	})

	e.PUT("/tasks/status", func(ctx echo.Context) error {
		if db.AllDone() {
			db.SetStatusAll(false)
		} else {
			db.SetStatusAll(true)
		}
		return ctx.Render(http.StatusOK, combo, currentViewModel())
	})

	e.DELETE("/task/:id", func(ctx echo.Context) error {
		id := ctx.Param("id")
		i64, _ := strconv.ParseInt(id, 10, strconv.IntSize)
		db.RemoveTask(int(i64))
		return ctx.Render(http.StatusOK, combo, currentViewModel())
	})

	e.DELETE("/tasks/completed", func(ctx echo.Context) error {
		db.RemoveCompleted()
		return ctx.Render(http.StatusOK, combo, currentViewModel())
	})

	e.Logger.Fatal(e.Start(":3000"))
}

type Template struct {
	templates *template.Template
}

func (tmpl *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return tmpl.templates.ExecuteTemplate(w, name, data)
}

func parseTmpl() *Template {
	tmpl := template.Must(template.ParseFiles(tmplDir + tmplFile))
	return &Template{templates: tmpl}
}
