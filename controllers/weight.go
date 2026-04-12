package controllers

import (
	"myGymPal/models"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type weightAverage struct {
	Days   int
	Weight float64
	Unit   string
}

func computeWeightAverage(entries []*models.BodyWeight, targetUnit string) *weightAverage {
	n := len(entries)
	if n == 0 {
		return nil
	}
	if n > 3 {
		n = 3
	}
	var sum float64
	for i := 0; i < n; i++ {
		sum += models.ConvertWeight(entries[i].Weight, entries[i].WeightUnit, targetUnit)
	}
	return &weightAverage{
		Days:   n,
		Weight: sum / float64(n),
		Unit:   targetUnit,
	}
}

type WeightController struct {
	beego.Controller
}

func (c *WeightController) Index() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	entries, err := BodyWeights.GetAllByUser(userID.(int64))
	if err != nil {
		logs.Error("WeightController.Index: GetAllByUser: %v", err)
	}

	user, _ := Users.GetByID(userID.(int64))
	weightUnit := "lb"
	if user != nil {
		weightUnit = user.WeightUnit
	}

	flash := beego.ReadFromRequest(&c.Controller)
	if msg, ok := flash.Data["success"]; ok {
		c.Data["Success"] = msg
	}

	// Compute average first (before in-place conversion, using stored units).
	avg := computeWeightAverage(entries, weightUnit)

	// Convert all stored entries to the user's preferred unit for display.
	for _, e := range entries {
		e.Weight = models.ConvertWeight(e.Weight, e.WeightUnit, weightUnit)
		e.WeightUnit = weightUnit
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "weight"
	c.Data["Entries"] = entries
	c.Data["WeightUnit"] = weightUnit
	c.Data["DefaultDate"] = time.Now().Format("2006-01-02")
	c.Data["WeightAvg"] = avg
	c.TplName = "weight/index.tpl"
}

func (c *WeightController) Create() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	dateStr := c.GetString("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.Redirect("/weight", 302)
		return
	}

	weightStr := c.GetString("weight")
	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil || weight <= 0 {
		c.Redirect("/weight", 302)
		return
	}

	weightUnit := c.GetString("weight_unit")
	if weightUnit != "kg" {
		weightUnit = "lb"
	}

	if _, err := BodyWeights.Create(userID.(int64), date, weight, weightUnit); err != nil {
		logs.Error("WeightController.Create: %v", err)
	}

	c.Redirect("/weight", 302)
}

func (c *WeightController) Update() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Redirect("/weight", 302)
		return
	}

	weightStr := c.GetString("weight")
	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil || weight <= 0 {
		c.Redirect("/weight", 302)
		return
	}

	weightUnit := c.GetString("weight_unit")
	if weightUnit != "kg" {
		weightUnit = "lb"
	}

	if _, err := BodyWeights.Update(id, userID.(int64), weight, weightUnit); err != nil {
		logs.Error("WeightController.Update: %v", err)
	}

	c.Redirect("/weight", 302)
}

func (c *WeightController) Delete() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Redirect("/weight", 302)
		return
	}

	if err := BodyWeights.Delete(id, userID.(int64)); err != nil {
		logs.Error("WeightController.Delete: %v", err)
	}

	c.Redirect("/weight", 302)
}
