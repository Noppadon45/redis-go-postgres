package main

import (
	"github.com/Noppadon/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func main() {
	app := fiber.New()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	
	app.Post("/users", handlers.CreateUser(rdb))
	app.Get("/users/:id", handlers.GetUser(rdb))
	app.Listen(":3000")
}

// type User struct {
// 	Fullname string `redis:"Fullname"`
// 	Age      string `redis:"Age"`
// 	Location string `redis:"Location"`
// 	Email    string `redis:"Email"`
// 	Zipcode  string `redis:"Zipcode"`
// }

// func main() {

// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 		Protocol: 2,
// 	})

// 	ctx := context.Background()

// 	// hashFields := []string{
// 	// 	"Fullname" , "Og Doe",
// 	// 	"age" , "50",
// 	// 	"location" , "Samutsakorn",
// 	// 	"email" , "Og.doe@example.com",
// 	// }

// 	// res1 , err := rdb.HSet(ctx , "user:2" , hashFields).Result()
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	err := SetUser(ctx, rdb, "user:2", User{
// 		Fullname: "Hello World",
// 		Age:      "30",
// 		Location: "Bankok",
// 		Email:    "New.doe@example.com",
// 		Zipcode: "10110",
// 	})

// 	res2, err := rdb.HGet(ctx, "user:1", "Fullname").Result()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("HGet result Key Fullname: ", res2)

// 	userData, err := GetAll(ctx, rdb, "user:1")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("HGetAll result: ", userData)

// }

// func SetUser(ctx context.Context, rdb *redis.Client, userID string, userData User) error {
// 	res, err := rdb.HSet(ctx, userID, userData).Result()
// 	fmt.Println("HSet result: ", res)
// 	return err
// }

// func GetAll(ctx context.Context, rdb *redis.Client, userID string) (User, error) {
// 	var userData User
// 	err := rdb.HGetAll(ctx, userID).Scan(&userData)
// 	return userData, err
// }
