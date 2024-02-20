package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
    "html/template"

	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
)

 func indexHandler(c *fiber.Ctx, db *sql.DB) error {
   var res string
   var items []string
   rows, err := db.Query("SELECT * FROM items")
   defer rows.Close()
   if err != nil {
       log.Fatalln(err)
       c.JSON("An error occured")
   }
   for rows.Next() {
       rows.Scan(&res)
       items = append(items, res)
   }
   return c.Render("index", fiber.Map{
       "Items": items,
   })
}

func postHandler(c *fiber.Ctx, db *sql.DB) error {
   return c.SendString("Hello")
}

func putHandler(c *fiber.Ctx, db *sql.DB) error {
   return c.SendString("Hello")
}

func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
   return c.SendString("Hello")
}

func main() {
    
   connStr := //
   // Connect to database
   db, err := sql.Open("postgres", connStr)
   if err != nil {
       log.Fatal(err)
   }
   
   engine := html.New("./views", ".html")
   app := fiber.New(fiber.Config{
       Views: engine,
   })

   app.Get("/", func(c *fiber.Ctx) error {
       return indexHandler(c, db)
   })

   app.Post("/", func(c *fiber.Ctx) error {
       return postHandler(c, db)
   })

   app.Put("/update", func(c *fiber.Ctx) error {
       return putHandler(c, db)
   })

   app.Delete("/delete", func(c *fiber.Ctx) error {
       return deleteHandler(c, db)
   })

   port := os.Getenv("PORT")
   if port == "" {
       port = "3000"
   }

   
   app.Static("/", "/views/public/index.html") 
   log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}