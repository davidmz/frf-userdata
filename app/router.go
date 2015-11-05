package app

func (app *App) InitRouter() {
	app.Router.HandleFunc("/public/{site}/{username}/{paramName}", app.ApiCall(app.checkVars(app.GetPublicParam))).Methods("GET")
	app.Router.HandleFunc("/public/{site}/{username}/{paramName}", app.ApiCall(app.checkVars(app.PostPublicParam))).Methods("POST")
	app.Router.NotFoundHandler = app.ApiCall(app.NotFound)
}
