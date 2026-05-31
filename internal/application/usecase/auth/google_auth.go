package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"budget-book-go/internal/application/dto"
	"budget-book-go/internal/domain/entity"
	"budget-book-go/internal/domain/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type GoogleAuthUseCase struct {
	userRepo    repository.UserRepository
	oauthConfig *oauth2.Config
	jwtSecret   string
}

func NewGoogleAuthUseCase(
	userRepo repository.UserRepository,
	oauthConfig *oauth2.Config,
	jwtSecret string,
) *GoogleAuthUseCase {
	return &GoogleAuthUseCase{
		userRepo:    userRepo,
		oauthConfig: oauthConfig,
		jwtSecret:   jwtSecret,
	}
}

type GoogleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (uc *GoogleAuthUseCase) GetAuthURL(state string) string {
	return uc.oauthConfig.AuthCodeURL(state)
}

func (uc *GoogleAuthUseCase) Authenticate(ctx context.Context, code string) (*dto.AuthResult, error) {
	// 1. 認証コードでアクセストークン取得
	token, err := uc.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("認証コードの交換に失敗しました: %w", err)
	}

	// 2. Googleユーザー情報取得
	client := uc.oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("ユーザー情報の取得に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンスの読み取りに失敗しました: %w", err)
	}

	var googleUser GoogleUserInfo
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, fmt.Errorf("ユーザー情報のパースに失敗しました: %w", err)
	}

	// 3. ユーザー検索 or 作成
	user, err := uc.findOrCreateUser(ctx, googleUser)
	if err != nil {
		return nil, err
	}

	// 4. JWT発行
	jwtToken, err := uc.generateJWT(user)
	if err != nil {
		return nil, err
	}

	name := ""
	if user.Name != nil {
		name = *user.Name
	}

	return &dto.AuthResult{
		Token:  jwtToken,
		UserID: user.ID.String(),
		Email:  user.Email,
		Name:   name,
	}, nil
}

func (uc *GoogleAuthUseCase) findOrCreateUser(ctx context.Context, googleUser GoogleUserInfo) (*entity.User, error) {
	// まずprovider_idで検索
	user, err := uc.userRepo.FindByProviderID(ctx, "GOOGLE", googleUser.ID)
	if err == nil {
		return user, nil
	}

	// 次にemailで検索
	user, err = uc.userRepo.FindByEmail(ctx, googleUser.Email)
	if err == nil {
		return user, nil
	}

	// なければ新規作成
	newUser := &entity.User{
		Email:      googleUser.Email,
		Provider:   "GOOGLE",
		ProviderID: &googleUser.ID,
		Name:       &googleUser.Name,
	}

	return uc.userRepo.Save(ctx, newUser)
}

func (uc *GoogleAuthUseCase) generateJWT(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"userId": user.ID.String(),
		"email":  user.Email,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}