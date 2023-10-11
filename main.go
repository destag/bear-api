package main

import (
	"fmt"
	"net/http"

	"github.com/destag/bear-api/controllers"
	"github.com/destag/bear-api/database"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"context"
	"log"
	"os"

	"github.com/destag/bear-api/metrics"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Bear struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func initDatabase() *sqlx.DB {
	db := database.Connect(":memory:")
	database.Migrate(db)
	return db
}

func mainC() {
	db := initDatabase()

	r := gin.Default()

	r.Use(otelgin.Middleware("bears-api"))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/bears", func(c *gin.Context) {
		bears := []Bear{}
		fmt.Printf("%#v\n", bears)
		fmt.Println(bears)
		err := db.Select(&bears, "SELECT * FROM bears;")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		fmt.Println(bears)
		fmt.Println([]Bear{})
		fmt.Printf("%T\n", bears)
		fmt.Printf("%+v\n", bears)
		fmt.Printf("%#v\n", bears)
		fmt.Printf("%#v\n", []Bear{})

		c.JSON(http.StatusOK, bears)
	})

	r.POST("/bears", func(c *gin.Context) {
		input := Bear{}
		if err := c.BindJSON(&input); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		fmt.Printf("input: %+v\n", input)

		res, err := db.NamedExec("INSERT INTO bears (name) VALUES (:name)", input)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		input.ID = uint(id)

		c.JSON(http.StatusOK, input)
	})
	r.Run()
}

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
)

func initTracer() func(context.Context) error {
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

func main() {
	cleanup := initTracer()
	defer cleanup(context.Background())

	provider := metrics.InitMeter()
	defer provider.Shutdown(context.Background())

	meter := provider.Meter("sample-golang-app")
	metrics.GenerateMetrics(meter)

	r := gin.Default()
	r.Use(otelgin.Middleware(serviceName))
	// Connect to database
	db := initDatabase()
	// models.ConnectDatabase()

	bearController := controllers.BearController{
		DB: db,
	}

	// Routes
	r.GET("/bears", bearController.ListBears)

	// Run the server
	r.Run(":8090")
}
