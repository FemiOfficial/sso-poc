package app

type AppController struct {
	appService *AppService
}

func CreateAppController(appService *AppService) *AppController {
	return &AppController{appService: appService}
}
