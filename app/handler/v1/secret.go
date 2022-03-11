package v1

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
	Score     int       `json:"score" binding:"required"`
}

type BusinessSecret interface {
	GetSecret(ctx context.Context) (entity.Secret, error)
	ProcessScore(ctx context.Context, score entity.Score) error
}

type SecretHandler struct {
	businessSecret BusinessSecret
}

func NewSecretHandler(bs BusinessSecret) SecretHandler {
	return SecretHandler{businessSecret: bs}
}

func (sh SecretHandler) GetSecret(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
    c.Header("Access-Control-Allow-Methods", "GET")
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
	var score score
	c.BindJSON(&score)
	err := sh.businessSecret.ProcessScore(c.Request.Context(), entity.Score(score))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.AbortWithStatus(200)
}
