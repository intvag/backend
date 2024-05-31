package main

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	decisions "github.com/intvag/decision-engine/service"
	"gorm.io/gorm"
)

type userIDContextKey struct{}

type Server struct {
	r      *gin.Engine
	db     *gorm.DB
	quotes decisions.QuotesClient
}

func New(d *gorm.DB, qc decisions.QuotesClient) (s Server, err error) {
	s.db = d
	s.quotes = qc

	s.r = gin.New()

	s.r.Use(gin.Logger())
	s.r.Use(gin.Recovery())
	s.r.Use(cors.Default())
	s.r.Use(gzip.Gzip(gzip.DefaultCompression))
	s.r.Use(requestid.New())

	quotes := s.r.Group("/quote")
	quotes.GET("", s.newQuote)
	quotes.GET("/:quote", s.quote)
	quotes.POST("/:quote", s.newQuoteItem)
	quotes.DELETE("/:quote/item/:item", s.deleteQuoteItem)

	v1 := s.r.Group("/v1", s.validateToken)

	v1.GET("/policy", s.getUserPolicies)
	v1.GET("/policy/:policy", s.getUserPolicy)

	v1.POST("/policy", s.createPolicy)

	return
}

func (s Server) validateToken(g *gin.Context) {
	// Do something here to validate jwt, but for now
	// just assume it's correct.
	//
	// Even better, use an API Gateway to do this
	g.Request = g.Request.WithContext(context.WithValue(g.Request.Context(), userIDContextKey{}, "204fa646-8f27-487f-a406-966d1f3992de"))

	g.Next()
}
