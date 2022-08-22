package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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
	UserAgent string    `json:"user_agent"`
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

type topUser struct {
	Username string `json:"user_name" binding:"required"`
	Score    int    `json:"score"`
}

type topUserRsp struct {
	TopUser []topUser `json:"leaderboard"`
	Status  string    `json:"status"`
}

type gifWinRsp struct {
	Url string `json:"url"`
	Msg string `json:"msg"`
}

type BusinessSecret interface {
	GetSecret(ctx context.Context) (entity.Secret, error)
	ProcessScore(ctx context.Context, score entity.Score) error
	UpdateUserName(ctx context.Context, score entity.Score) error
	GetTopPlayer(ctx context.Context) ([]entity.Score, error)
	GetWeeklyTopPlayer(ctx context.Context, secretNum int) (entity.Score, error)
}

type businessGif interface {
	GetGif() (string, error)
}

type SecretHandler struct {
	businessSecret BusinessSecret
	bgif           businessGif
}

func NewSecretHandler(bs BusinessSecret, bg businessGif) SecretHandler {
	return SecretHandler{businessSecret: bs, bgif: bg}
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
	c.JSON(http.StatusOK, secretRsp{
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
		UserName:  strings.ToUpper(score.Username),
		UserAgent: fmt.Sprintf("%s|%s",
			c.GetHeader("X-forwarded-for"),
			score.UserAgent,
		),
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
		UserName:  strings.ToUpper(user.Username),
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

func (sh SecretHandler) SelectTopUser(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	topPlayers, err := sh.businessSecret.GetTopPlayer(c.Request.Context())
	switch err {
	case business.SecretNumNotValid:
		c.JSON(http.StatusBadRequest, topUserRsp{Status: err.Error()})
	case nil:
		jsTopPlayers := []topUser{}
		for _, v := range topPlayers {
			jsTopPlayers = append(jsTopPlayers, topUser{Username: strings.ToUpper(v.UserName), Score: v.Score})
		}
		c.JSON(http.StatusOK, topUserRsp{TopUser: jsTopPlayers, Status: "Ok"})
	default:
		c.JSON(http.StatusInternalServerError, topUserRsp{Status: err.Error()})
	}
}

func (sh SecretHandler) GetSuccessGif(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	url, err := sh.bgif.GetGif()
	switch err {
	case nil:
		c.JSON(http.StatusOK, gifWinRsp{Url: url})
	default:
		c.JSON(http.StatusInternalServerError, gifWinRsp{Msg: err.Error()})

	}
}

func (sh SecretHandler) GetWeeklyTopPlayer(c *gin.Context) {
	var err error
	c.Header("Access-Control-Allow-Origin", "*")
	secretNumStr, ok := c.GetQuery("secretnum")
	if !ok || secretNumStr == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("secretnum is required"))
		return
	}
	secretNum, err := strconv.Atoi(secretNumStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("secretnum is not valid"))
		return
	}

	topPlayer, err := sh.businessSecret.GetWeeklyTopPlayer(c.Request.Context(), secretNum)
	switch err {
	case business.SecretNumNotValid:
		c.JSON(http.StatusBadRequest, topUserRsp{Status: err.Error()})
	case nil:
		c.JSON(http.StatusOK, topUser{Username: topPlayer.UserName, Score: topPlayer.Score})
	default:
		c.JSON(http.StatusInternalServerError, topUserRsp{Status: err.Error()})
	}
}
