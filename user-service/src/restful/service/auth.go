package service

import (
	"time"
	"user-service/src/common/log"

	"user-service/src/cache"
	"user-service/src/common/dto/request"
	"user-service/src/common/dto/response"
	"user-service/src/common/errors"
	"user-service/src/common/model"
	v "user-service/src/common/pkg/validator"
	"user-service/src/common/utils"
	"user-service/src/factory"
	"user-service/src/publisher"
	"user-service/src/repository"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	Register(c fiber.Ctx, data *request.RegisterRequest) (email string, err error)
	VerifyRegister(c fiber.Ctx, data *request.VerifyRegister) error
	Login(c fiber.Ctx, data *request.LoginRequest) (*response.Login, error)
	Logout(c fiber.Ctx, refreshToken string) error
	RefreshToken(c fiber.Ctx, refreshToken string) (*response.Token, error)
}

type authImpl struct {
	UserRepository         repository.User
	CredentialRepository   repository.Credential
	RefreshTokenRepository repository.RefreshToken
	CacheRepository        cache.Cache
	RabbitMQPublisher      *publisher.RabitMQ
}

// function ini nanti di panggil di handler
func NewAuth(f *factory.Factory) Auth {
	return &authImpl{
		UserRepository:         f.UserRepository,
		CredentialRepository:   f.CredentialRepository,
		CacheRepository:        f.CacheRepository,
		RabbitMQPublisher:      f.RabbitMQPublisher,
		RefreshTokenRepository: f.RefreshTokenRepository,
	}
}
func (s *authImpl) Register(c fiber.Ctx, data *request.RegisterRequest) (email string, err error) {
	if err := v.Validate.Struct(data); err != nil {
		return "", err
	}

	user, err := s.UserRepository.Find(c.Context(), "email =?", []any{data.Email})
	if err != nil {
		return "", err
	}
	if user != nil {
		return "", &errors.Response{Message: "email already exist", HttpCode: 409}
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	data.Password = string(bcryptPassword)

	// Data user (dengan password hash) disimpan ke Redis dengan key "register:<email>" selama 20 menit
	if err := s.CacheRepository.Set(c.Context(), "register:"+data.Email, data, 20*time.Minute); err != nil {
		return "", err
	}
	//buat otp
	otp, err := utils.GenerateOTP()
	if err != nil {
		return "", err
	}
	sendOtpReq := request.SendOtpReq{
		Email: data.Email,
		Otp:   otp,
	}
	// Data OTP disimpan ke Redis dengan key "otp:<email>" selama 20 menit
	if err := s.CacheRepository.Set(c.Context(), "otp:"+data.Email, sendOtpReq, 20*time.Minute); err != nil {
		return "", err
	}
	// Kirim OTP melalui RabbitMQ
	if err := s.RabbitMQPublisher.Publish("user", "otp", &request.RabbitMqMessage{
		RequestID: c.Locals("request_id"),
		Message:   sendOtpReq,
	}); err != nil {
		return "", err
	}
	return data.Email, nil
}

func (s *authImpl) VerifyRegister(c fiber.Ctx, data *request.VerifyRegister) error {
	if err := v.Validate.Struct(data); err != nil {
		return err
	}
	//Mengambil data OTP dari Redis menggunakan key "otp:<email>"
	otp, err := s.CacheRepository.GetSendOtp(c.Context(), "otp:"+data.Email)
	if err != nil {
		return err
	}
	if otp == nil {
		return &errors.Response{Message: "otp not found", HttpCode: 404}
	}

	if otp.Otp != data.Otp {
		return &errors.Response{Message: "invalid otp", HttpCode: 400}
	}
	//Mengambil data registrasi dari Redis menggunakan key "register:<email>"
	register, err := s.CacheRepository.GetRegister(c.Context(), "register:"+data.Email)
	if err != nil {
		return err
	}
	if register == nil {
		return &errors.Response{Message: "register not found", HttpCode: 404}
	}

	// Generate UUID v4 baru
	userID := uuid.New()

	if err := s.UserRepository.Create(c.Context(), &model.User{
		ID:    userID,
		Name:  register.Name,
		Email: register.Email,
		Role:  "USER",
	}, register.Password); err != nil {
		return err
	}

	if err := s.CacheRepository.Delete(c.Context(), "register:"+data.Email); err != nil {
		log.Logger.Error(err.Error())
	}
	if err := s.CacheRepository.Delete(c.Context(), "otp:"+data.Email); err != nil {
		log.Logger.Error(err.Error())
	}
	return nil
}

func (s *authImpl) Login(c fiber.Ctx, data *request.LoginRequest) (*response.Login, error) {
	if err := v.Validate.Struct(data); err != nil {
		return nil, err
	}
	//cari emain di tabel users
	user, err := s.UserRepository.Find(c.Context(), "email = ?", []any{data.Email})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, &errors.Response{Message: "email tidak terdaftar", HttpCode: 404}
	}
	//cari password di tabel credential milik email tersebut
	credential, err := s.CredentialRepository.Find(c.Context(), "email = ?", []any{data.Email})
	if err != nil {
		return nil, err
	}
	if credential == nil {
		return nil, &errors.Response{Message: "credential not found", HttpCode: 404}
	}
	//mencocokkan password plain dari request dengan hash yang disimpan di DB
	if err := bcrypt.CompareHashAndPassword([]byte(credential.Password), []byte(data.Password)); err != nil {
		return nil, &errors.Response{Message: "password salah", HttpCode: 400}
	}
	//generate acces token
	accessToken, err := utils.GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}
	//generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, err
	}
	//simpan refresh token ke database
	if err := s.RefreshTokenRepository.Create(c.Context(), &model.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		return nil, err
	}
	//mengembalikan objece response.Login ke handler
	return &response.Login{
		Data:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authImpl) Logout(c fiber.Ctx, refreshToken string) error {
	if err := s.RefreshTokenRepository.Delete(c.Context(), "token = ?", []any{refreshToken}); err != nil {
		return err
	}
	return nil
}

func (s *authImpl) RefreshToken(c fiber.Ctx, refreshToken string) (*response.Token, error) {
	token, err := s.RefreshTokenRepository.Find(c.Context(), refreshToken)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, &errors.Response{Message: "refresh token not found", HttpCode: 404}
	}

	accessToken, err := utils.GenerateAccessToken(token.UserID.String(), token.Role)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := utils.GenerateRefreshToken(token.UserID.String())
	if err != nil {
		return nil, err
	}

	// Update token di database
	if err := s.RefreshTokenRepository.Update(c.Context(), &model.RefreshToken{
		Token:     newRefreshToken,
		UpdatedAt: time.Now(),
	}, refreshToken); err != nil {
		return nil, err
	}

	// Kembalikan token BARU ke response
	return &response.Token{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken, // <-- Pastikan mengembalikan token baru!
	}, nil
}
