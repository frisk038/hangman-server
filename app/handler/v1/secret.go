package v1

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/frisk038/hangman-server/business/entity"
	"github.com/gin-gonic/gin"
)

type secretRsp struct {
	Number int      `json:"number"`
	Secret []string `json:"secret"`
}

type BusinessSecret interface {
	GetSecret(ctx context.Context) (entity.Secret, error)
}

type SecretHandler struct {
	businessSecret BusinessSecret
}

func NewSecretHandler(bs BusinessSecret) SecretHandler {
	return SecretHandler{businessSecret: bs}
}

func (sh SecretHandler) GetSecret(c *gin.Context) {
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
