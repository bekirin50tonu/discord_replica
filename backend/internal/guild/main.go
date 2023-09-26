package guild

import (
	"backend/pkg/eventbus/rabbitmq"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Initialize() {

}

type UserAddIntegrationEvent struct {
	Name string `json:"name"`
}

func UserAddIntegrationEventHandler(d amqp.Delivery) {
	var user UserAddIntegrationEvent
	err := json.Unmarshal(d.Body, &user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("test geldi gibi ama : %v", user.Name)
}

func Run() {

	/* _, db := database.Initialize(database.DatabaseConfig{Url: "mongodb://localhost:27017", DatabaseName: "tests"})
	baseModel := models.Initialize(db)
	// model := models.NewAccount("bekirin50tonu", "password")
	// baseModel.Create(*model)

	mdl := baseModel.Find(models.Account{Username: "bekirin50tonu"})

	fmt.Printf("Gelen:%v", mdl) */

	broker, err := rabbitmq.NewRabbitMQ(rabbitmq.RabbitMQConfig{URL: "amqp://localhost:5672", ExchangeName: "GuildService", AcK: true, ExchangeSuffix: "IntegrationEvent"})

	if err != nil {
		panic(err)

	}

	_, err = broker.AddSubscription(UserAddIntegrationEvent{}, UserAddIntegrationEventHandler)

	if err != nil {
		panic(err)
	}

	defer broker.Close()

	// Ctrl+C sinyalini yakala
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Programın ana işlemi
	go func() {
		for {
			select {
			case <-sigCh:
				fmt.Println("\nCtrl+C ile çıkılıyor...")
				// Burada ek temizleme işlemleri yapabilirsiniz.
				os.Exit(0)
			default:
				// Diğer işlemler burada çalışır.
				fmt.Println("Program devam ediyor...")
				time.Sleep(1 * time.Second)
			}
		}
	}()

	// Programın çalışmasını sürdür
	select {}

}
