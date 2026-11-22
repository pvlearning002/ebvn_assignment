package main

import (
	"github.com/pvlearning002/ebvn_assignment/pkg/utils/utils_log"
)

func main() {
	// test_config()
	// test_redis()
	test_log()
}

// func test_config() {
// 	cfg := config.CreateNewConfig()
// 	fmt.Println("App Port:", cfg.AppPort)
// 	genPassService := service.NewGenPass(*cfg)
// 	generatedPass, err := genPassService.GeneratePassword()
// 	if err != nil {
// 		fmt.Println("Error generating password:", err)
// 		return
// 	}
// 	fmt.Println("Generated Password:", generatedPass)
// }

// func test_redis() {
// 	cfg := config.CreateNewConfig()
// 	redisClient, err := redis.NewClient("REDIS")
// 	if err != nil {
// 		fmt.Println("Error creating Redis client:", err)
// 		return
// 	}
// 	urlStorage := repository.NewURLStorage(redisClient, cfg)

// 	// Store a URL in Redis
// 	key := "example_key"
// 	url := "https://example.com"
// 	err = urlStorage.StoreURL(context.Background(), key, url)
// 	if err != nil {
// 		fmt.Println("Error storing URL:", err)
// 		return
// 	}
// 	fmt.Println("URL stored successfully.")

// 	// Retrieve the URL from Redis
// 	retrievedURL, err := urlStorage.GetURL(context.Background(), key)
// 	if err != nil {
// 		fmt.Println("Error retrieving URL:", err)
// 		return
// 	}
// 	fmt.Println("Retrieved URL:", retrievedURL)
// }

func test_log() {
	my_logger := utils_log.NewCustomLogger()
	my_logger.WriteLogInfo("This is info text")

	// utils_log.WriteLogInfo("This is info text")
	// utils_log.WriteLogError("This is error text")
	// utils_log.WriteLogDebug("This is debug text")
	// utils_log.WriteLogWarning("This is warn text")
}
