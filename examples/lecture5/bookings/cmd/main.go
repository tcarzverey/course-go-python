package main

import (
	"context"
	"log"
	"os"

	"github.com/exaring/otelpgx"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/tcarzverey/bookings/internal/generated/api"
	"github.com/tcarzverey/bookings/internal/repository/rooms"
	"github.com/tcarzverey/bookings/internal/server"
	"github.com/tcarzverey/bookings/internal/usecases/list_rooms"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()

	cfg, err := pgxpool.ParseConfig(os.Getenv("DB_ADDR"))
	if err != nil {
		log.Fatal(err)
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := otelpgx.RecordStats(db); err != nil {
		log.Fatal(err)
	}

	roomsRepo := rooms.New(db)
	listRoomsUC := list_rooms.NewUsecase(roomsRepo)
	s := server.New(listRoomsUC)
	handler := api.NewStrictHandler(s, nil)

	router := gin.Default()
	router.Use(cors.Default())
	api.RegisterHandlers(router, handler)

	initSwagger(router)

	log.Println("listening on http://localhost:" + os.Getenv("SERVER_HTTP_PORT"))
	log.Fatal(router.Run(":" + os.Getenv("SERVER_HTTP_PORT")))
}
