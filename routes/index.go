package routes

import (
	"bookbalance/app/controllers"
	"bookbalance/app/middleware"
	"bookbalance/app/models"
	"bookbalance/app/repository"
	"bookbalance/app/services"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

func InitRoutes(db *sql.DB, redis *redis.Client) *echo.Echo {
	app := echo.New();
	app.Use(middleware.ContentTypeResponse);
	app.Validator = &models.PayloadValidator{Validator: validator.New()};

	router := app.Group("/api/v1");
	router.Use(middleware.AuthBasic(db));

	// Fare
	fareRepo := repository.NewFareRepository(db);
	fareService := services.NewFareService(fareRepo);
	fareController := controllers.NewFareController(app, fareService, db);
	router.GET("/fare/origin/:origin/dest/:destination", fareController.Fare);

	// Ticket
	ticket := router.Group("/ticket");
	ticketRepo := repository.NewTicketRepo(db, redis);
	ticketService := services.NewTicketService(ticketRepo);
	ticketController := controllers.NewTicketController(app, ticketService, db);
	ticket.POST("/add", ticketController.AddTicket);

	return app;
}