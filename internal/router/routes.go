package router

import (
	"net/http"

	"github.com/hazkall/capy-belga/internal/controller"
	"github.com/hazkall/capy-belga/internal/domain/service"
)

func discountClubPostHandler(ch chan *controller.Message) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.ControllerCreateDiscountClub(w, r, ch)
	}
}

func discountClubUserPostHandler(ch chan *controller.Message) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.ControllerCreateUser(w, r, ch)
	}
}

func discountClubSignupPostHandler(ch chan *controller.Message) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.ControllerCreateClubSignup(w, r, ch)
	}
}

func userState(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.ControllerUserState(w, r, userService)
	}
}

func userPlanSignup(signupService *service.SignupService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.ControllerGetClubStatus(w, r, signupService)
	}
}

func cancelUserClub(signupService *service.SignupService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.ControllerCancelClubSignup(w, r, signupService)
	}
}
