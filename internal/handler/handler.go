package handler

import (
	"manga-library/internal/service"
	"manga-library/pkg/logger"
	"net/http"

	_ "manga-library/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
	logger   logger.Logger
}

func NewHandler(service *service.Service, logger logger.Logger) *Handler {
	return &Handler{services: service, logger: logger}
}

func (h *Handler) InitRoutes() *gin.Engine {
	h.logger.Infoln("routes initializated")
	router := gin.New()
	router.MaxMultipartMemory = 6

	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
	}

	api := router.Group("/api")
	{
		api.GET("/heartbeat", h.heartbeat)

		manga := api.Group("/manga")
		{
			manga.GET("/latest", h.getLatestManga)
			manga.GET("", h.getManga)
			manga.GET("/filter", h.getMangaByFilter)
			manga.POST("/", h.userIdentity, h.createManga)
			manga.DELETE("/:id", h.userIdentity, h.deleteManga)
			manga.PATCH("/:id", h.userIdentity, h.updateManga)
		}
		upload := api.Group("/upload", h.userIdentity)
		{
			upload.POST("/preview", h.uploadMangaPreview)
			upload.DELETE("/preview/:id", h.deleteMangaPreview)
		}
		chapter := api.Group("/chapter")
		{
			chapter.POST("/", h.userIdentity, h.uploadChapter)
			chapter.DELETE("/:mangaSlug/:volume/:number", h.userIdentity, h.deleteChapter)
		}
		user := api.Group("/user")
		{
			user.GET("/:username", h.getUserByUsername)
			user.PATCH("/:id", h.updateUser)
			user.DELETE("/:id", h.deleteUser)
		}
	}
	router.GET("/docs", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/swagger/index.html")
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
