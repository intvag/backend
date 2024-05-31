package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type Policy struct {
	ID        uuid.UUID  `json:"id" gorm:"primary_key"`
	PersonID  uuid.UUID  `json:"person"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" sql:"index"`
	IsLive    bool       `json:"is_live"`
	DocLink   string     `json:"-"`

	Items []PolicyItem `json:"items" gorm:"foreignKey:PolicyID"`
}

type PolicyItem struct {
	ItemID        uuid.UUID `json:"-"`
	PolicyID      uuid.UUID `json:"-"`
	Cost          float64   `json:"cost"`
	Age           float64   `json:"age"`
	OriginalPrice float64   `json:"original_price"`
}

type PolicyInput struct {
	QuotetID uuid.UUID `json:"quote_id"`
}

func (s Server) getUserPolicies(g *gin.Context) {
	id := uuid.Must(uuid.FromString(g.Request.Context().Value(userIDContextKey{}).(string)))

	policies := new([]Policy)

	err := s.db.Model(new(Policy)).Preload("Items").Find(policies, "person_id = ?", id).Error
	if err != nil {
		g.AbortWithError(http.StatusNotFound, err)

		return
	}

	g.JSON(http.StatusOK, policies)
}

func (s Server) getUserPolicy(g *gin.Context) {
	policy := g.Param("policy")

	p := new(Policy)

	err := s.db.Model(p).Preload("Items").Find(p, "id = ?", policy).Error
	if err != nil {
		g.AbortWithError(http.StatusNotFound, err)

		return
	}

	g.JSON(http.StatusOK, p)
}

func (s Server) createPolicy(g *gin.Context) {
	pi := new(PolicyInput)

	err := g.BindJSON(pi)
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)

		return
	}

	p := new(Policy)
	p.ID = uuid.Must(uuid.NewV4())
	p.PersonID = uuid.Must(uuid.FromString(g.Request.Context().Value(userIDContextKey{}).(string)))
	p.CreatedAt = time.Now()

	quote := new(Quote)
	err = s.db.Model(quote).Preload("Items").Find(quote, "id = ?", pi.QuotetID).Error
	if err != nil {
		g.AbortWithError(http.StatusNotFound, err)

		return
	}

	p.Items = make([]PolicyItem, len(quote.Items))

	for idx, q := range quote.Items {
		p.Items[idx] = PolicyItem{
			ItemID:        q.ItemID,
			PolicyID:      p.ID,
			Cost:          q.Cost,
			Age:           q.Age,
			OriginalPrice: q.OriginalPrice,
		}
	}

	s.db.Create(p)
	s.db.Delete(quote, quote.ID)

	g.JSON(http.StatusCreated, p)
}

func (s Server) policyPayment(g *gin.Context) {
	// In reality this endpoint will be passed to a payment provider as
	// a callback address
	//
	// The payment provider calls this endpoint, and then this endpoint
	// does something to validate that ahead of updating the database
	// to mark the policy as live
	//
	// This stage is also where we trigger jobs to generate policy documents
	// and email them to users
	id, err := uuid.FromString(g.Param("policy"))
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)

		return
	}

	err = s.db.Model(new(Policy)).Where("id = ?", id).Update("IsLive", true).Error
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)

		return
	}

	go s.generatePolicyDocs(id)

	g.Status(http.StatusAccepted)
}

func (s Server) generatePolicyDocs(id uuid.UUID) {
	policy := new(Policy)

	err := s.db.Model(policy).Preload("Items").Find(policy, "id = ?", id).Error
	if err != nil {
		log.Print(err)

		return
	}

	person := new(Person)
	err = s.db.Model(person).Find(person, "id = ?", policy.PersonID).Error
	if err != nil {
		log.Print(err)

		return
	}

	doc, err := policy.GeneratePDF(*person)
	if err != nil {
		log.Print(err)

		return
	}

	err = s.db.Model(policy).Where("id = ?", policy.ID).Update("DocLink", doc).Error
	if err != nil {
		log.Print(err)

	}
}
