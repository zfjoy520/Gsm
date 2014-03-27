package controllers

import (
	"fmt"
	"github.com/robfig/revel"
	"gsm/app/models"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	fmt.Printf("Index...\n")

	// gsm.GetDb().Incr("hello")
	// s, _ := gsm.GetDb().Get("hello")
	// log.Printf("> hello = %s\n", s)
	new(models.Test).Save()

	return c.Render()
}
