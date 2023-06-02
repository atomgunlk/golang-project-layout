package handler

import (
	"errors"

	"github.com/atomgunlk/YOUR-REPO-NAME/cmd/_your_app_/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/iancoleman/strcase"
)

type Handler struct {
	Config   *Config
	Deps     *Dependency
	Validate *validator.Validate
}

type Config struct {
	JWTSecret string
}

type Dependency struct {
	Service *service.Service
}

// New handelr
func New(config *Config, deps *Dependency) (*Handler, error) {
	return &Handler{
		Config:   config,
		Deps:     deps,
		Validate: validator.New(),
	}, nil
}

func (h *Handler) InitRoutes(f *fiber.App) error {

	v1NoAuth := f.Group("/v1")
	{
		v1NoAuth.Post("/login", h.Login)
	}

	// V1 With JWT Middleware
	v1 := v1NoAuth.Use(jwtware.New(
		jwtware.Config{
			TokenLookup:    "header:Authorization",
			AuthScheme:     "Bearer",
			SigningKey:     []byte(h.Config.JWTSecret),
			SigningMethod:  jwt.SigningMethodHS256.Name,
			ErrorHandler:   AuthError,
			SuccessHandler: AuthSuccess,
		},
	))
	{
		v1.Post("/user", h.User)
	}

	return nil
}

func AuthError(c *fiber.Ctx, e error) error {
	c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Unauthorized",
		"msg":   e.Error(),
	})
	return nil
}

func AuthSuccess(c *fiber.Ctx) error {
	c.Next()
	return nil
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := ApiError{
				ErrorMsg: make([]ErrorMsg, len(ve)),
			}

			for idx, fe := range ve {
				out.ErrorMsg[idx] = ErrorMsg{
					Field: strcase.ToLowerCamel(fe.Field()), Msg: msgForTag(fe.Tag()),
				}

			}

			return out
		}

		return err
	}

	return nil
}
