package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/hazkall/capy-belga/internal/domain/entity"
	"github.com/hazkall/capy-belga/internal/domain/service"
	"github.com/hazkall/capy-belga/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func ControllerCreateDiscountClub(w http.ResponseWriter, r *http.Request, ch chan *Message) {

	club := new(entity.Club)

	m := new(Message)

	if err := json.NewDecoder(r.Body).Decode(&club); err != nil {
		http.Error(w, "Erro ao decodificar o JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	m.Type = "create_discount_club"
	m.Data, _ = json.Marshal(club)

	if err := club.ValidateClub(); err != nil {
		http.Error(w, "Erro de validação: "+err.Error(), http.StatusBadRequest)
		return
	}

	ch <- m

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Club de Desconto sendo Processado!"))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "close")
}

func ControllerCreateUser(w http.ResponseWriter, r *http.Request, ch chan *Message) {

	user := new(entity.User)

	m := new(Message)

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Erro ao decodificar o JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	m.Type = "users"
	m.Data, _ = json.Marshal(user)

	if err := user.ValidateUser(); err != nil {
		http.Error(w, "Erro de validação: "+err.Error(), http.StatusBadRequest)
		return
	}

	slog.Info("Publicando mensagem na fila para usuário", "email", user.Email)

	ch <- m

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Usuário criado!"))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "close")
}

func ControllerHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "close")
}

func ControllerCreateClubSignup(w http.ResponseWriter, r *http.Request, ch chan *Message) {
	signup := new(entity.SignupPayload)

	m := new(Message)

	if err := json.NewDecoder(r.Body).Decode(&signup); err != nil {
		http.Error(w, "Erro ao decodificar o JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	m.Type = "discount_club_signup"
	m.Data, _ = json.Marshal(signup)

	ch <- m

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Inscrição no Clube de Desconto sendo Processada!"))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "close")
}

func ControllerUserState(w http.ResponseWriter, r *http.Request, userService *service.UserService) {

	user := new(entity.User)

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := user.ValidateUser(); err != nil {
		http.Error(w, "Erro de validação: "+err.Error(), http.StatusBadRequest)
		return
	}

	slog.Info("Verificando estado do usuário", "email", user.Email)

	state, err := userService.UserState(r.Context(), user.Email)
	if err != nil {
		http.Error(w, "Erro ao obter estado do usuário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"email": user.Email,
		"name":  user.Name,
		"state": state,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Erro ao criar resposta JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "close")
}

func ControllerGetClubStatus(w http.ResponseWriter, r *http.Request, signupService *service.SignupService) {

	signup := new(entity.SignupPayload)

	if err := json.NewDecoder(r.Body).Decode(&signup); err != nil {
		http.Error(w, "Erro ao decodificar o JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	status, plan, err := signupService.UserClubStatus(r.Context(), signup)
	if err != nil {
		http.Error(w, "Erro ao obter status do clube: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"email":  signup.Email,
		"status": status,
		"plan":   plan,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Erro ao criar resposta JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "close")
}

func ControllerCancelClubSignup(w http.ResponseWriter, r *http.Request, signupService *service.SignupService) {

	plan := new(entity.SignupPayload)

	if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
		http.Error(w, "Erro ao decodificar o JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := signupService.CancelSignup(r.Context(), plan); err != nil {
		http.Error(w, "Erro ao cancelar a inscrição: "+err.Error(), http.StatusInternalServerError)
		return
	}

	telemetry.PlanGauge.Record(
		r.Context(),
		-1,
		metric.WithAttributes(
			attribute.String("email", plan.Email),
			attribute.String("club_name", plan.ClubName),
		),
	)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Inscrição no Clube de Desconto cancelada!"))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "close")
}
