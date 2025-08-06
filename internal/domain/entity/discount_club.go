package entity

type Club struct {
	ID                 int64  `json:"id,omitempty"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	AquisitionChannel  string `json:"aquisition_channel"`
	AquisitionLocation string `json:"aquisition_location"`
	PlanType           string `json:"plan_type"`
}

type SignupPayload struct {
	Email    string `json:"email"`
	ClubName string `json:"club"`
}
