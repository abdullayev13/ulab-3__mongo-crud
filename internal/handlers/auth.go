package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"my_app/internal/config"
	"my_app/internal/models"
	"my_app/internal/pkg/mongodb"
	"net/http"
	"strings"
	"time"
)

func Register(c *gin.Context) {
	data := new(models.User)
	err := c.ShouldBind(data)
	if err != nil {
		fail(c, err)
		return
	}

	if data.Username == "" {
		failMsg(c, "username is required")
		return
	}

	coll := mongodb.GetColl("users")

	_, err = coll.Find(c, Obj{"username": data.Username})
	if !errors.Is(err, mongo.ErrNoDocuments) {
		failMsg(c, "username already exists")
		return
	}

	data.Id = primitive.ObjectID{}
	res, err := coll.InsertOne(c, data)
	if err != nil {
		return
	}

	data.Id = res.InsertedID.(primitive.ObjectID)

	token, err := createToken(data)
	if err != nil {
		fail(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  data,
	})

}

func Login(c *gin.Context) {
	data := new(models.User)
	err := c.ShouldBind(data)
	if err != nil {
		fail(c, err)
		return
	}

	if data.Username == "" {
		failMsg(c, "username is required")
		return
	}

	coll := mongodb.GetColl("users")

	cur, err := coll.Find(c, Obj{"username": data.Username})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			failMsg(c, "user not found")
			return
		}
		fail(c, err)
		return
	}

	model := new(models.User)
	err = cur.Decode(model)
	if err != nil {
		fail(c, err)
		return
	}

	token, err := createToken(model)
	if err != nil {
		fail(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  model,
	})

}

func Auth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		failMsg(c, "error: Missing token")
		c.Abort()
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		failMsg(c, "invalid auth header")
		c.Abort()
		return
	}

	if len(headerParts[1]) == 0 {
		failMsg(c, "token is empty")
		c.Abort()
		return
	}

	claims, err := parseToken(headerParts[1])
	if err != nil {
		fail(c, err)
		c.Abort()
		return
	}

	c.Set(userCtx, claims.UserId)
	c.Next()

}

func getUserInfo(c *gin.Context) userInfo {
	value, ok := c.Get(userCtx)
	if !ok {
		panic("user not found in ctx")
	}

	id := value.(primitive.ObjectID)
	return userInfo{id}
}

const userCtx = "user_id"

type userInfo struct {
	UserId primitive.ObjectID
}
type jwtCustomClaim struct {
	jwt.StandardClaims
	UserId primitive.ObjectID `json:"user_id"`
}

func createToken(user *models.User) (string, error) {
	claims := &jwtCustomClaim{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(config.JwtDuration).Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte(config.JwtSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func parseToken(token string) (*jwtCustomClaim, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwtCustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(config.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*jwtCustomClaim)
	if !ok {
		return nil, errors.New("token claims are not of type *jwtCustomClaim")
	}

	return claims, nil
}
