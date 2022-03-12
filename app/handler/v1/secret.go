package v1

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/frisk038/hangman-server/business"
	"github.com/frisk038/hangman-server/business/entity"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type secretRsp struct {
	Number int      `json:"number"`
	Secret []string `json:"secret"`
}

type score struct {
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	SecretNum int       `json:"secret_num" binding:"required"`
	Score     int       `json:"score"`
	Username  string    `json:"user_name"`
}

type username struct {
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	SecretNum int       `json:"secret_num" binding:"required"`
	Username  string    `json:"user_name" binding:"required"`
}

type userRsp struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type BusinessSecret interface {
	GetSecret(ctx context.Context) (entity.Secret, error)
	ProcessScore(ctx context.Context, score entity.Score) error
	UpdateUserName(ctx context.Context, score entity.Score) error
}

type SecretHandler struct {
	businessSecret BusinessSecret
}

func NewSecretHandler(bs BusinessSecret) SecretHandler {
	return SecretHandler{businessSecret: bs}
}

func (sh SecretHandler) GetSecret(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	secret, err := sh.businessSecret.GetSecret(c.Request.Context())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	secretArr := strings.Split(secret.SecretWord, "")
	if len(secretArr) == 0 {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("split failed"))
		return
	}
	c.JSON(200, secretRsp{
		Number: secret.Number,
		Secret: secretArr,
	})
}

func (sh SecretHandler) PostScore(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST")
	c.Header("Access-Control-Allow-Headers", "Content-Type")

	var score score
	c.BindJSON(&score)
	err := sh.businessSecret.ProcessScore(c.Request.Context(), entity.Score{
		UserID:    score.UserID,
		SecretNum: score.SecretNum,
		Score:     score.Score,
		UserName:  score.Username,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.AbortWithStatus(200)
}

func (sh SecretHandler) UpdateUserName(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST")
	c.Header("Access-Control-Allow-Headers", "Content-Type")

	var user username
	c.BindJSON(&user)
	err := sh.businessSecret.UpdateUserName(c.Request.Context(), entity.Score{
		UserID:    user.UserID,
		SecretNum: user.SecretNum,
		UserName:  user.Username,
	})
	switch err {
	case business.ScoreNotValid, business.SecretNumNotValid,
		business.UsernameNotValid:
		c.JSON(http.StatusBadRequest, userRsp{Status: "KO", Reason: err.Error()})
	case nil:
		c.JSON(http.StatusOK, userRsp{Status: "Ok"})
	default:
		c.JSON(http.StatusInternalServerError, userRsp{Status: "KO", Reason: err.Error()})
	}
}
