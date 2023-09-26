package user

import (
	"backend/internal/user/controllers"
	"backend/internal/user/integrations"
	"backend/internal/user/repositories"
	"backend/internal/user/services"
	"backend/pkg/database"
	"backend/pkg/eventbus/rabbitmq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/mongo"
)

func Run() {
	_, db := database.Initialize(database.DatabaseConfig{Url: "mongodb://localhost:27017", DatabaseName: "tests"})

	service_user, err := initializeUserService(db)
	if err != nil {
		panic(err)
	}
	service_account, err := initializeAccountService(db)
	if err != nil {
		panic(err)
	}
	service_session, err := initializeSessionService(db)
	if err != nil {
		panic(err)
	}
	broker, err := initializeBroker(service_user)
	if err != nil {
		panic(err)
	}
	defer broker.Close()

	controller, err := controllers.NewUserController(service_user, service_session, service_account)
	if err != nil {
		panic(err)
	}

	_, err = initializeApp(controller)
	if err != nil {
		panic(err)
	}

}

func initializeBroker(user *services.UserService) (*rabbitmq.RabbitMQ, error) {
	broker, err := rabbitmq.NewRabbitMQ(rabbitmq.RabbitMQConfig{URL: "amqp://localhost:5672", ExchangeName: "Replica.User", AcK: true, ExchangePrefix: "IntegrationEvent", ExchangeSuffix: "User"})
	if err != nil {
		return nil, err

	}
	// Define Base Integrations with needed Services.
	handlers, err := integrations.NewIntegrations(user)
	// Define Subscriptions
	if err != nil {
		return nil, err
	}
	_, err = broker.AddSubscription(integrations.GetUserIntegrationEvent{}, handlers.GetUserIntegrationEventHandler)
	if err != nil {
		return nil, err
	}
	// Defined

	return broker, nil
}

func initializeUserService(db *mongo.Database) (*services.UserService, error) {

	repository_user, err := repositories.NewUserRepository(db, "users")
	if err != nil {
		return nil, err
	}

	service, err := services.NewUserService(*repository_user)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func initializeAccountService(db *mongo.Database) (*services.AccountService, error) {
	repository_account, err := repositories.NewAccountRepository(db, "accounts")
	if err != nil {
		return nil, err
	}

	service, err := services.NewAccountService(*repository_account)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func initializeSessionService(db *mongo.Database) (*services.SessionService, error) {
	repository_session, err := repositories.NewSessionRepository(db, "sessions")
	if err != nil {
		return nil, err
	}

	service, err := services.NewSessionService(repository_session)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func initializeApp(controller *controllers.UserController) (any, error) {
	app := fiber.New()

	app.Use(cors.New())

	app.Listen("0.0.0.0:8080")

	//sigCh := make(chan os.Signal, 1)
	//signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	return nil, nil
}
