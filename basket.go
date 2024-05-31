package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	decisions "github.com/intvag/decision-engine/service"
)

type Quote struct {
	ID        uuid.UUID  `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" sql:"index"`

	Items []QuoteItem `json:"items" gorm:"foreignKey:QuoteID"`
}

type QuoteItem struct {
	PolicyItem

	QuoteID uuid.UUID `json:"quote"`
}

type QuoteInput struct {
	Manufacturer  string  `json:"manufacturer"`
	ModelName     string  `json:"model"`
	Age           float64 `json:"age"`
	PurchasePrice float64 `json:"purchase_price"`
}

func (s Server) newQuote(g *gin.Context) {
	now := time.Now()
	b := &Quote{ID: uuid.Must(uuid.NewV4()), CreatedAt: now, UpdatedAt: now}

	s.db.Create(b)

	g.JSON(http.StatusCreated, b)
}

func (s Server) quote(g *gin.Context) {
	b := new(Quote)

	err := s.db.Model(b).Preload("Items").First(b, "id = ?", g.Param("quote")).Error
	if err != nil {
		g.AbortWithError(http.StatusNotFound, err)

		return
	}

	g.JSON(http.StatusOK, b)
}

func (s Server) newQuoteItem(g *gin.Context) {
	qi := new(QuoteInput)

	err := g.BindJSON(qi)
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)

		return
	}

	item := new(Item)
	s.db.First(item, "manufacturer = ? AND model_name = ?", qi.Manufacturer, qi.ModelName)

	quote, err := s.quotes.GetQuote(g.Request.Context(), &decisions.Input{
		ExpectedLifetime: float64(item.AverageLifetime),
		Age:              qi.Age,
		PurchasePrice:    qi.PurchasePrice,
		Lastability:      item.Lastability,
		Repairability:    item.Repairability,
	})
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	b := new(Quote)

	s.db.First(b, "id = ?", g.Param("quote"))
	b.Items = append(b.Items, QuoteItem{
		QuoteID: b.ID,
		PolicyItem: PolicyItem{
			ItemID:        item.ID,
			Cost:          quote.Monthly,
			Age:           qi.Age,
			OriginalPrice: qi.PurchasePrice,
		},
	})

	s.db.Save(b)

	g.JSON(http.StatusCreated, b)
}

func (s Server) deleteQuoteItem(g *gin.Context) {
	q := new(QuoteItem)
	s.db.Delete(q, g.Param("quote"))

	g.Status(http.StatusNoContent)
}
