package main

import (
	"context"

	"github.com/kingstonduy/mcs-money-transfer/bootstrap"
	"github.com/kingstonduy/mcs-money-transfer/infras/data"
	"github.com/kingstonduy/mcs-money-transfer/presentation/http"
	"github.com/kingstonduy/mcs-money-transfer/presentation/http/controller"
	"github.com/kingstonduy/mcs-money-transfer/usecase"
	"github.com/lengocson131002/go-clean-core/logger"
	"go.uber.org/fx"
)

var configModule = fx.Module("config",
	fx.Provide(bootstrap.GetLogger),
	fx.Provide(bootstrap.GetConfigure),
	fx.Provide(bootstrap.GetServerConfig),
	fx.Provide(bootstrap.GetValidator),
	fx.Provide(bootstrap.GetTracer),
	fx.Provide(bootstrap.GetKafkaBroker),
	fx.Provide(bootstrap.NewHealthChecker),
	fx.Provide(bootstrap.NewElasticSearchClient),
	fx.Provide(bootstrap.GetPrometheusMetricer),
	fx.Provide(bootstrap.GetConfig),
	fx.Provide(bootstrap.GetTemporalClient),
)

var controllerModule = fx.Module("controller",
	fx.Provide(controller.NewAccountController),
	fx.Provide(controller.NewMoneyTransferController),
	// add new controller in here

)

var pipelineModule = fx.Module("pipeline",
	fx.Provide(bootstrap.NewMetricBehavior),
	fx.Provide(bootstrap.NewTracingBehavior),
	fx.Provide(bootstrap.NewRequestLoggingBehavior),
	fx.Provide(bootstrap.NewErrorHandlingBehavior),
	fx.Provide(usecase.NewMoneyTransferHandler),
	fx.Invoke(bootstrap.RegisterPipelineBehaviors),
	fx.Invoke(bootstrap.RegisterPipelineHandlers),
)

var serverModule = fx.Module("server",
	fx.Provide(http.NewHttpServer),
	// fx.Provide(broker.NewBrokerServer),
	// fx.Provide(grpc.NewGrpcServer),
	// fx.Provide(cron.NewCronServer),
)

var infrasModule = fx.Module("infras",
	fx.Provide(data.NewAccountBalanceData),
)

func main() {
	fx.New(
		configModule,
		pipelineModule,
		serverModule,
		infrasModule,
		controllerModule,
		fx.Invoke(run),
	).Run()
}

func run(
	lc fx.Lifecycle,
	httpServer *http.HttpServer,
	// brokerServer *broker.BrokerServer,
	// grpcServer *grpc.GrpcServer,
	// cronServer *cron.CronServer,
	log logger.Logger,
	shutdowner fx.Shutdowner) {
	gCtx, cancel := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := httpServer.Start(gCtx); err != nil {
					log.Fatal(ctx, "Failed to start HTTP server: %s", err)
					cancel()
					shutdowner.Shutdown()
				}
			}()

			// go func() {
			// 	if err := brokerServer.Start(gCtx); err != nil {
			// 		log.Fatalf(ctx, "Failed to start Broker server: %s", err)
			// 		cancel()
			// 		shutdowner.Shutdown()
			// 	}
			// }()

			// go func() {
			// 	if err := grpcServer.Start(gCtx); err != nil {
			// 		log.Fatalf(ctx, "Failed to start GRPC server: %s", err)
			// 		cancel()
			// 		shutdowner.Shutdown()
			// 	}
			// }()

			// go func() {
			// 	if err := cronServer.Start(gCtx); err != nil {
			// 		log.Fatalf(ctx, "Failed to start Cron server: %s", err)
			// 		cancel()
			// 		shutdowner.Shutdown()
			// 	}
			// }()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			cancel()
			return nil
		},
	})
}
