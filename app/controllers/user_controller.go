package controllers

import (
	"fmt"
	"time" // <-- tambah

	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/app/repository"
	"github.com/kooroshh/fiber-boostrap/pkg/jwt_token"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *fiber.Ctx) error {
	user := new(models.User)
	err := ctx.BodyParser(user)
	if err != nil {
		errResponse := fmt.Errorf("Failed to parse request: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	err = user.Validate()
	if err != nil {
		errResponse := fmt.Errorf("Failed to validate request: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		errResponse := fmt.Errorf("Failed to encrypt the password: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}
	user.Password = string(hashPassword)

	err = repository.InsertNewUser(ctx.Context(), user)
	if err != nil {
		errResponse := fmt.Errorf("failed to insert new user: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	resp := user
	resp.Password = ""
	return response.SendSuccessResponse(ctx, resp)
}

func Login(ctx *fiber.Ctx) error {
	loginReq := new(models.LoginRequest)
	var resp models.LoginResponse // <-- bukan fungsi

	err := ctx.BodyParser(loginReq)
	if err != nil {
		errResponse := fmt.Errorf("Failed to parse request: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	err = loginReq.Validate()
	if err != nil {
		errResponse := fmt.Errorf("Failed to validate request: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	user, err := repository.GetUserByUsername(ctx.Context(), loginReq.Username)
	if err != nil {
		errResponse := fmt.Errorf("Failed to get username: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "Username Failed", nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		errResponse := fmt.Errorf("Failed to check password: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "Password Failed", nil)
	}

	// MINIMAL FIX: pakai ctx.UserContext() + time.Now() sesuai signature lib kamu
	now := time.Now()
	token, err := jwt_token.GenerateToken(ctx.UserContext(), user.Username, user.FullName, "token", now)
	if err != nil {
		errResponse := fmt.Errorf("Failed to generate token: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "Terjadi kesalahan pada system", nil)
	}

	refreshToken, err := jwt_token.GenerateToken(ctx.UserContext(), user.Username, user.FullName, "refresh_token", now)
	if err != nil {
		errResponse := fmt.Errorf("Failed to generate token: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "Terjadi kesalahan pada system", nil)
	}

	resp.Username = user.Username
	resp.FullName = user.FullName
	resp.Token = token // <-- pakai variabel, bukan user.token
	resp.RefreshToken = refreshToken

	return response.SendSuccessResponse(ctx, resp)
}
