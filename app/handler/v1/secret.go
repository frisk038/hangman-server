package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/frisk038/hangman-server/business/usecase"
	"github.com/gin-gonic/gin"
)

type secretRsp struct {
	Number int64    `json:"number"`
	Secret []string `json:"secret"`
}

func GetSecret(c *gin.Context) {
	secret, err := usecase.GenerateDailySecret()
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
