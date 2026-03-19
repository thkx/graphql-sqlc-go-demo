package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/thkx/graphql-sqlc-go-dome/internal/auth"
	"github.com/thkx/graphql-sqlc-go-dome/internal/auth/jwt"
	"github.com/thkx/graphql-sqlc-go-dome/internal/config"
	"github.com/thkx/graphql-sqlc-go-dome/internal/graph"
	"github.com/thkx/graphql-sqlc-go-dome/internal/repository/postgres"
	"github.com/thkx/graphql-sqlc-go-dome/internal/usecase"
	"github.com/thkx/graphql-sqlc-go-dome/internal/usecase/comment"
	"github.com/thkx/graphql-sqlc-go-dome/internal/usecase/integration"
	"github.com/thkx/graphql-sqlc-go-dome/internal/usecase/post"
	"github.com/thkx/graphql-sqlc-go-dome/internal/usecase/user"
	"github.com/thkx/graphql-sqlc-go-dome/pkg/app"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	cfg := config.Load()
	// 注入配置到 JWT
	jwt.SetupJWTConfig(cfg.JWTSecret, cfg.JWTExpiry)

	// repo := memory.NewMemoryRepository()
	// repo, err := sqlite.NewSQLiteRepository(cfg)
	repo, err := postgres.NewPostgresRepository(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	userUsecase := user.NewUsecase(repo)
	postUsecase := post.NewUsecase(repo)
	commentUsecase := comment.NewUsecase(repo)
	integrationUsecase := integration.NewUsecase(repo)

	u := usecase.NewUsecase(userUsecase, postUsecase, commentUsecase, integrationUsecase)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: graph.NewResolver(u),
		Directives: graph.DirectiveRoot{
			RequiresAuth: auth.AuthFieldMiddleware,
		},
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	if cfg.PlaygroundEnabled {
		log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.Port)
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}

	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", auth.Middleware()(auth.WrapGraphQLHandler(srv)))

	s := app.New(app.WithHandler(http.DefaultServeMux), app.WithAddr(":"+cfg.Port))

	// log.Fatal(http.ListenAndServe(":"+port, nil))

	log.Printf("Server starting on port %s (env: %s)", cfg.Port, cfg.Env)
	// log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
	s.Start()
	s.Stop()
}
