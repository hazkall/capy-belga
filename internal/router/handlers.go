package router

import (
	"net/http"

	middlewares "github.com/hazkall/capy-belga/internal/middleware"
)

func HandlersPipeline(deps *HandlerDeps) {
	http.Handle("/contrate/discount-club", middlewarePipeline(discountClubPostHandler(deps.ClubChannel)))
	http.Handle("/contrate/discount-club/signup", middlewarePipeline(discountClubSignupPostHandler(deps.ClubChannel)))
	http.Handle("/contrate/discount-club/user", middlewarePipeline(discountClubUserPostHandler(deps.ClubChannel)))
	http.Handle("/user/state", middlewarePipeline(userState(deps.UserService)))
	http.Handle("/user/cancel/club", middlewarePipeline(cancelUserClub(deps.SignupService)))
	http.Handle("/user/plan/status", middlewarePipeline(userPlanSignup(deps.SignupService)))

}

func middlewarePipeline(handler http.Handler) http.Handler {
	h := middlewares.RecoverMiddleware(handler)
	h = middlewares.OtelMiddleware(h)
	h = middlewares.RequestsCountMiddleware(h)
	return h
}
