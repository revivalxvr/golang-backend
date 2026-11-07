package service

import (
	"context"
	"fmt"
	"time"
	"user-service/src/cache"
	"user-service/src/common/dto/request"
	"user-service/src/common/errors"
	"user-service/src/common/log"
	"user-service/src/common/model"
	"user-service/src/factory"
	"user-service/src/repository"

	"golang.org/x/crypto/bcrypt"
)

type ProfileService interface {
	GetProfile(ctx context.Context, userID string) (*model.User, error)
	UpdateProfileAdmin(ctx context.Context, userID string, req request.UpdateProfileRequest) (*model.User, error)
	UpdateProfileUser(ctx context.Context, userID string, req request.UpdateProfileRequest) (*model.User, error)
	UpdatePasswordAdmin(ctx context.Context, userID string, req request.ChangePasswordRequest) (*model.Credential, error)
	UpdatePasswordUser(ctx context.Context, userID string, req request.ChangePasswordRequest) (*model.Credential, error)
}

type profileImpl struct {
	CacheRepository      cache.Cache
	UserRepository       repository.User
	CredentialRepository repository.Credential
}

func NewProfileService(f *factory.Factory) ProfileService {
	return &profileImpl{
		CacheRepository:      f.CacheRepository,
		UserRepository:       f.UserRepository,
		CredentialRepository: f.CredentialRepository,
	}
}

func (s *profileImpl) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	// 1. Cek cache
	setUser, err := s.CacheRepository.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if setUser != nil {
		return setUser, nil
	}

	// 2. Query ke Database
	seUser, err := s.UserRepository.Find(ctx, "id = ?", []any{userID})
	if err != nil {
		return nil, err
	}

	// 3. Jika user tidak ditemukan di DB
	if seUser == nil {
		return nil, &errors.Response{Message: "user not found", HttpCode: 404}
	}

	// 4. Simpan ke Redis Cache (durasi 30 menit)
	if err := s.CacheRepository.Set(ctx, fmt.Sprintf("user:%s", userID), seUser, 30*time.Minute); err != nil {
		log.Logger.Error(err.Error())
	}

	// 5. Return data user dari DB
	return seUser, nil
}

// UpdateProfileAdmin
func (s *profileImpl) UpdateProfileAdmin(ctx context.Context, userID string, req request.UpdateProfileRequest) (*model.User, error) {
	// 1. Ambil data user yang ada saat ini dari repository
	user, err := s.UserRepository.FindByID(userID)
	if err != nil {
		return nil, err
	}
	// 2. Terapkan logika partial update (Opsi 1: Cek pointer nil atau tidak)
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		// Cek apakah email baru sudah digunakan oleh user lain
		existingUser, err := s.UserRepository.FindByEmailAndNotID(*req.Email, userID)
		if err != nil {
			return nil, err
		}
		if existingUser != nil {
			return nil, &errors.Response{Message: "email already exist", HttpCode: 409}
		}
		user.Email = *req.Email
	}
	// 3. Simpan perubahan melalui repository
	if err := s.UserRepository.Update(user); err != nil {
		return nil, &errors.Response{Message: "failed to update profile in database", HttpCode: 500}
	}

	return user, nil
}

func (s *profileImpl) UpdateProfileUser(ctx context.Context, userID string, req request.UpdateProfileRequest) (*model.User, error) {
	// 1. Ambil data user yang ada saat ini dari repository
	user, err := s.UserRepository.FindByID(userID)
	if err != nil {
		return nil, err
	}
	// 2. Terapkan logika partial update (Opsi 1: Cek pointer nil atau tidak)
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		// Cek apakah email baru sudah digunakan oleh user lain
		existingUser, err := s.UserRepository.FindByEmailAndNotID(*req.Email, userID)
		if err != nil {
			return nil, err
		}
		if existingUser != nil {
			return nil, &errors.Response{Message: "email already exist", HttpCode: 409}
		}
		user.Email = *req.Email
	}
	// 3. Simpan perubahan melalui repository
	if err := s.UserRepository.Update(user); err != nil {
		return nil, &errors.Response{Message: "failed to update profile in database", HttpCode: 500}
	}

	return user, nil
}

func (s *profileImpl) UpdatePasswordAdmin(ctx context.Context, userID string, req request.ChangePasswordRequest) (*model.Credential, error) {
	// 1. Ambil data credential yang ada saat ini dari repository
	credential, err := s.CredentialRepository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	// 2. Verifikasi password lama dengan hash yang ada di tabel credentials
	if err := bcrypt.CompareHashAndPassword([]byte(credential.Password), []byte(req.OldPassword)); err != nil {
		return nil, &errors.Response{Message: "old password is incorrect", HttpCode: 400}
	}
	// 3. Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

	if err != nil {
		return nil, &errors.Response{Message: "failed to hash new password", HttpCode: 500}
	}
	// 4. Update password di tabel credentials
	credential.Password = string(hashedPassword)
	// 5. Simpan perubahan melalui repository
	if err := s.CredentialRepository.Update(credential); err != nil {
		return nil, &errors.Response{Message: "failed to update password in database", HttpCode: 500}
	}

	return credential, nil
}

// h.service.UpdatePasswordUser(c.Context(), userID, req)
func (s *profileImpl) UpdatePasswordUser(ctx context.Context, userID string, req request.ChangePasswordRequest) (*model.Credential, error) {
	// 1. Ambil data credential yang ada saat ini dari repository
	credentialUser, err := s.CredentialRepository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	// 2. Verifikasi password lama dengan hash yang ada di tabel credentials
	if err := bcrypt.CompareHashAndPassword([]byte(credentialUser.Password), []byte(req.OldPassword)); err != nil {
		return nil, &errors.Response{Message: "old password is incorrect", HttpCode: 400}
	}
	// 3. Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

	if err != nil {
		return nil, &errors.Response{Message: "failed to hash new password", HttpCode: 500}
	}
	// 4. Update password di tabel credentials
	credentialUser.Password = string(hashedPassword)
	// 5. Simpan perubahan melalui repository
	if err := s.CredentialRepository.Update(credentialUser); err != nil {
		return nil, &errors.Response{Message: "failed to update password in database", HttpCode: 500}
	}

	return credentialUser, nil
}
