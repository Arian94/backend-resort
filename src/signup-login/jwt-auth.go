package signupLogin

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//jwt service
type JWTService interface {
	GenerateToken(email string) string
	ValidateToken(token string) (*jwt.Token, error)
}
type authCustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type jwtServices struct {
	secretKey string
	issuer    string
}

//auth-jwt
func JWTAuthService() JWTService {
	return &jwtServices{
		secretKey: "secret",
		issuer:    "Arian",
	}
}

// func getSecretKey() string {
// 	secret := os.Getenv("JWT_SECRET")
// 	if secret == "" {
// 		secret = "secret"
// 	}
// 	return secret
// }

func (service *jwtServices) GenerateToken(email string) string {
	claims := &authCustomClaims{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 102).Unix(),
			Issuer:    service.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (service *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			var err error
			log.Println("invalid token:", token.Header["alg"])
			return nil, err
		}
		return []byte(service.secretKey), nil
	})

}
