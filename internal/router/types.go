package router

import (
	"github.com/hazkall/capy-belga/internal/controller"
	"github.com/hazkall/capy-belga/internal/domain/service"
)

type HandlerDeps struct {
	ClubChannel   chan *controller.Message
	UserService   *service.UserService
	ClubService   *service.ClubService
	SignupService *service.SignupService
}
