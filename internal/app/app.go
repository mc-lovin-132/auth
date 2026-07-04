package app

import (
	"context"
	"net"
	"time"

	"github.com/mc-lovin-132/auth/config"
	userClient "github.com/mc-lovin-132/auth/internal/infrastructure/clients/grpc"
	detector "github.com/mc-lovin-132/auth/internal/infrastructure/delivery/grpc/detector"
	handlers "github.com/mc-lovin-132/auth/internal/infrastructure/delivery/grpc/handlers"
	interceptors "github.com/mc-lovin-132/auth/internal/infrastructure/delivery/grpc/interceptors"
	accessManager "github.com/mc-lovin-132/auth/internal/infrastructure/managers/access"
	refreshManager "github.com/mc-lovin-132/auth/internal/infrastructure/managers/refresh"
	sessionManager "github.com/mc-lovin-132/auth/internal/infrastructure/managers/sessions"
	repo "github.com/mc-lovin-132/auth/internal/infrastructure/repository/psql/common"
	refreshRepo "github.com/mc-lovin-132/auth/internal/infrastructure/repository/psql/refresh"
	accessRepo "github.com/mc-lovin-132/auth/internal/infrastructure/repository/redis/access"
	"github.com/mc-lovin-132/auth/internal/service"
	"github.com/mc-lovin-132/auth/pb"

	userspb "github.com/mc-lovin-132/users/pb"

	// "github.com/jmoiron/sqlx"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	cfg    *config.Config
	logger *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Start(ctx context.Context) error {
	// redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     a.cfg.RedisAddr(),
		Password: a.cfg.RedisPassword,
		DB:       a.cfg.RedisBDNumber,
	})

	// проверка подключения
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	a.logger.Info("successsfuly set connection with redis", zap.String("pong", pong))

	// postgres
	db, err := sqlx.Connect("postgres", a.cfg.DSN())
	if err != nil {
		return err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			a.logger.Error("err while closing db connection", zap.Error(err))
			return
		}
		a.logger.Info("db connection successfuly closed")
	}()
	a.logger.Info("successfuly init db connection")

	// migrations
	err = repo.RunMigrations(db.DB)
	if err != nil {
		return err
	}
	a.logger.Info("successfuly run migrations")

	// set grpc client connection
	conn, err := grpc.NewClient(
		a.cfg.UserAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	a.logger.Info("successfuly init user grpc connection")

	client := userspb.NewUserServiceClient(conn)

	// service
	srv := service.New(
		userClient.New(client),
		// пока что версии устанавливаются навсегда
		// TODO: подумать над удалением
		// можно либо просто установить ttl, либо добавить крон таймер который
		// будет чистить старые версии
		accessRepo.NewGlobalVersion(rdb, a.cfg.StartGlobalVersion, 0),
		accessRepo.NewDeviceVersion(rdb, a.cfg.StartDeviceVersion, 0),
		accessRepo.NewBlackList(rdb, time.Duration(a.cfg.AccessTTLMinutes)*time.Minute),
		accessRepo.NewBlackList(rdb, time.Duration(a.cfg.AccessTTLMinutes)*time.Minute),
		refreshRepo.New(db),
		detector.New(a.cfg.Salt),
		sessionManager.New(),
		refreshManager.New(time.Duration(a.cfg.RefreshTTLDays)*time.Hour*24),
		accessManager.New(a.cfg.SecretKey, time.Duration(a.cfg.AccessTTLMinutes)*time.Minute),
	)

	// handler
	handler := handlers.New(srv)

	// health handler
	healthService := service.NewHealthService(db, rdb)
	healthHandler := handlers.NewHealthHandler(healthService)

	// grpc server start
	lis, err := net.Listen("tcp", a.cfg.Addr())
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.NewLogging(a.logger)),
	)
	pb.RegisterAuthServiceServer(grpcServer, handler)
	pb.RegisterHealthServiceServer(grpcServer, healthHandler)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-gCtx.Done()
		grpcServer.GracefulStop()
		return gCtx.Err()
	})

	g.Go(func() error {
		a.logger.Info("gRPC server listening", zap.String("port", a.cfg.Port))
		if err := grpcServer.Serve(lis); err != nil {
			a.logger.Error("grpc server run failed", zap.Error(err))
			return err
		}
		return nil
	})

	return g.Wait()
}
