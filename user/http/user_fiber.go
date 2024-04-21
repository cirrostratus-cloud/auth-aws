package main

import (
	"fmt"

	"github.com/cirrostratus-cloud/auth/user"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type userRequest struct {
	Email            string `json:"email"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Password         string `json:"password"`
	PasswordRepeated string `json:"passwordRepeated"`
}

type userResponse struct {
	ID        string `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Enabled   bool   `json:"enabled,omitempty"`
}

type emailConfirmationRequest struct {
	ValidationToken string `json:"validationToken"`
}

type userAPI struct {
	api                    fiber.Router
	createUserUseCase      user.CreateUserUseCase
	getUserUseCase         user.GetUserUseCase
	updateProfileUseCase   user.UpdateUserProfileUseCase
	confirmateEmailUseCase user.ConfirmateEmailUseCase
}

func newUserAPI(createUserUseCase user.CreateUserUseCase, getUserUseCase user.GetUserUseCase, updateProfileUseCase user.UpdateUserProfileUseCase, confirmateEmailUseCase user.ConfirmateEmailUseCase) *userAPI {
	return &userAPI{
		createUserUseCase:      createUserUseCase,
		getUserUseCase:         getUserUseCase,
		updateProfileUseCase:   updateProfileUseCase,
		confirmateEmailUseCase: confirmateEmailUseCase,
	}
}

func (u *userAPI) setUp(app *fiber.App, stage string) (userAPI *userAPI) {
	log.
		WithField("Stage", stage).
		Info("Setting up stage.")
	u.api = app.Group(fmt.Sprintf("/%s", stage))
	u.api.Post("/users", u.createUser)
	u.api.Get("/users/:id", u.getUserByID)
	u.api.Put("/users/:id", u.updateUser)
	u.api.Post("/users/confirmate-email", u.confirmateEmail)
	return u
}

func (u *userAPI) createUser(c *fiber.Ctx) error {
	userRequest := new(userRequest)
	if err := c.BodyParser(userRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}
	user, err := u.createUserUseCase.NewUser(user.CreateUserRequest{
		Email:            userRequest.Email,
		FirstName:        userRequest.FirstName,
		LastName:         userRequest.LastName,
		Password:         userRequest.Password,
		PasswordRepeated: userRequest.PasswordRepeated,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	return c.Status(201).JSON(&fiber.Map{
		"user": userResponse{
			ID: user.UserID,
		},
		"message": "User created",
	})
}

func (u *userAPI) getUserByID(c *fiber.Ctx) error {
	UserID := c.Params("id")
	user, err := u.getUserUseCase.GetUserByID(user.UserByID{UserID: UserID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"user": userResponse{
			ID:        user.UserID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Enabled:   user.Enabled,
		},
		"message": fmt.Sprintf("User ID: %s", UserID),
	})
}

func (u *userAPI) updateUser(c *fiber.Ctx) error {
	UserID := c.Params("id")
	userRequest := new(userRequest)
	if err := c.BodyParser(userRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}
	user := user.UpdateUserProfileRequest{
		UserID:    UserID,
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
	}
	_, err := u.updateProfileUseCase.UpdateUserProfile(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"user": userResponse{
			ID:        UserID,
			Email:     userRequest.Email,
			FirstName: userRequest.FirstName,
			LastName:  userRequest.LastName,
		},
		"message": fmt.Sprintf("User ID: %s", UserID),
	})
}

func (u *userAPI) confirmateEmail(c *fiber.Ctx) error {
	emailConfirmationRequest := new(emailConfirmationRequest)
	if err := c.BodyParser(emailConfirmationRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}
	confirmateEmailResponse, err := u.confirmateEmailUseCase.ConfirmateEmail(user.ConfirmateEmailRequest{
		ValidationToken: emailConfirmationRequest.ValidationToken,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Email confirmed",
		"user": userResponse{
			ID: confirmateEmailResponse.UserID,
		},
	})
}
