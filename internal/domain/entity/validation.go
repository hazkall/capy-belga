package entity

import (
	"fmt"
	"log/slog"
)

func (c *User) ValidateUser() error {

	slog.Info("Validating user entity", "user", c)

	if c == nil {
		return fmt.Errorf("User must be provided")
	}

	if c.Name == "" || len(c.Name) < 3 {
		return fmt.Errorf("User name must not be empty and must be at least 3 characters long")
	}

	if c.Email == "" {
		return fmt.Errorf("User email must not be empty")
	}

	return nil
}

func (c *Club) ValidateClub() error {
	if c == nil {
		return fmt.Errorf("Club must be provided")
	}

	if c.Name == "" || len(c.Name) < 3 {
		return fmt.Errorf("Club name must not be empty and must be at least 3 characters long")
	}

	if c.Description == "" {
		return fmt.Errorf("Club description must not be empty")
	}

	if c.AquisitionChannel == "" || len(c.AquisitionChannel) < 3 || c.AquisitionChannel != "online" && c.AquisitionChannel != "offline" {
		return fmt.Errorf("Club aquisition channel must not be empty and must be either 'online' or 'offline'")
	}

	if c.AquisitionLocation == "" || len(c.AquisitionLocation) < 3 || c.AquisitionLocation != "store" && c.AquisitionLocation != "website" {
		return fmt.Errorf("Club aquisition location must not be empty and must be either 'store' or 'website'")
	}

	if c.PlanType == "" || (c.PlanType != "basic" && c.PlanType != "premium") {
		return fmt.Errorf("Club plan type must not be empty and must be either 'basic' or 'premium'")
	}

	return nil
}
