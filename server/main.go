package main

import (
	"majoo-backend/pkg/env"
	"majoo-backend/pkg/interfacepkg"
	"majoo-backend/pkg/jwe"
	"majoo-backend/pkg/jwt"
	"majoo-backend/pkg/logruslogger"
	"majoo-backend/pkg/pg"
	"majoo-backend/pkg/str"
	boot "majoo-backend/server/bootstrap"
	"majoo-backend/usecase"
	"os"
	"strings"

	"github.com/rs/xid"

	"github.com/rs/cors"

	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/go-chi/chi"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v7"
	validator "gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
	idTranslations "gopkg.in/go-playground/validator.v9/translations/id"
)

var (
	_, b, _, _      = runtime.Caller(0)
	basepath        = filepath.Dir(b)
	debug           = false
	host            string
	validatorDriver *validator.Validate
	uni             *ut.UniversalTranslator
	translator      ut.Translator
	envConfig       map[string]string
	corsDomainList  []string
)

// Init first time running function
func init() {
	// Load env variable from .env file
	envConfig = env.NewEnvConfig("../.env")

	host = envConfig["APP_HOST"]
	if str.StringToBool(envConfig["APP_DEBUG"]) {
		debug = true
		log.Printf("Running on Debug Mode: On at host [%v]", host)
	}
}

func main() {
	ctx := "main"

	// Connect to redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     envConfig["REDIS_HOST"],
		Password: envConfig["REDIS_PASSWORD"],
		DB:       0,
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}

	// Postgre DB connection
	dbInfo := pg.Connection{
		Host:    envConfig["DATABASE_HOST"],
		DB:      envConfig["DATABASE_DB"],
		User:    envConfig["DATABASE_USER"],
		Pass:    envConfig["DATABASE_PASSWORD"],
		Port:    str.StringToInt(envConfig["DATABASE_PORT"]),
		SslMode: "disable",
	}
	db, err := dbInfo.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// JWT credential
	jwtCredential := jwt.Credential{
		Secret:           envConfig["TOKEN_SECRET"],
		ExpSecret:        str.StringToInt(envConfig["TOKEN_EXP_SECRET"]),
		RefreshSecret:    envConfig["TOKEN_REFRESH_SECRET"],
		RefreshExpSecret: str.StringToInt(envConfig["TOKEN_EXP_REFRESH_SECRET"]),
	}

	// JWE credential
	jweCredential := jwe.Credential{
		KeyLocation: envConfig["APP_PRIVATE_KEY_LOCATION"],
		Passphrase:  envConfig["APP_PRIVATE_KEY_PASSPHRASE"],
	}

	// Validator initialize
	validatorInit()

	// Load contract struct
	contractUC := usecase.ContractUC{
		ReqID:     xid.New().String(),
		DB:        db,
		Redis:     redisClient,
		EnvConfig: envConfig,
		Jwt:       jwtCredential,
		Jwe:       jweCredential,
	}

	r := chi.NewRouter()
	// Cors setup
	r.Use(cors.New(cors.Options{
		AllowedOrigins: corsDomainList,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}).Handler)

	// load application bootstrap
	bootApp := boot.Bootup{
		R:          r,
		CorsDomain: corsDomainList,
		EnvConfig:  envConfig,
		DB:         db,
		Redis:      redisClient,
		Validator:  validatorDriver,
		Translator: translator,
		ContractUC: contractUC,
		Jwt:        jwtCredential,
		Jwe:        jweCredential,
	}

	// register middleware
	bootApp.RegisterMiddleware()

	// register routes
	bootApp.RegisterRoutes()

	// Log start server
	startBody := map[string]interface{}{
		"Host":     host,
		"Location": str.DefaultData(envConfig["APP_DEFAULT_LOCATION"], "Asia/Jakarta"),
	}
	logruslogger.Log(logruslogger.InfoLevel, interfacepkg.Marshall(startBody), ctx, "server_start", "")

	// Create static folder for file uploading
	filePath := envConfig["FILE_STATIC_FILE"]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Mkdir(filePath, os.ModePerm)
	}

	// Register folder for a go static folder
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, filePath)
	FileServer(r, envConfig["FILE_PATH"], http.Dir(filesDir))

	// Cors
	handler := cors.AllowAll().Handler(r)

	// Run the app
	http.ListenAndServe(host, handler)
}

func validatorInit() {
	en := en.New()
	id := id.New()
	uni = ut.New(en, id)

	transEN, _ := uni.GetTranslator("en")
	transID, _ := uni.GetTranslator("id")

	validatorDriver = validator.New()

	enTranslations.RegisterDefaultTranslations(validatorDriver, transEN)
	idTranslations.RegisterDefaultTranslations(validatorDriver, transID)

	switch envConfig["APP_LOCALE"] {
	case "en":
		translator = transEN
	case "id":
		translator = transID
	}
}

// FileServer ...
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
