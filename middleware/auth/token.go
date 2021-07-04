package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

type TokenMetadata struct {
	Uuid     string
	UserId   uint
	UserName string
}

type Tokens struct {
	AccessToken  string
	AccessUuid   string
	AtExpires    int64
	WsToken      string
	WsUuid       string
	WsExpires    int64
	RefreshToken string
	RefreshUuid  string
	RtExpires    int64
}

func CreateTokens(userId uint, userName string) (*Tokens, error) {
	var err error
	tokens := &Tokens{}

	//Creating Access Token
	tokens.AtExpires = time.Now().Add(time.Minute * 60).Unix() //expires after 60 min
	tokens.AccessUuid = uuid.NewV4().String()

	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = tokens.AccessUuid
	atClaims["user_id"] = userId
	atClaims["user_name"] = userName
	atClaims["exp"] = tokens.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)

	tokens.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	//Creating Websocket Token
	tokens.WsExpires = tokens.AtExpires
	tokens.WsUuid = fmt.Sprintf("%s++WS", tokens.AccessUuid)

	wstClaims := jwt.MapClaims{}
	wstClaims["ws_uuid"] = tokens.WsUuid
	wstClaims["user_id"] = userId
	wstClaims["user_name"] = userName
	wstClaims["exp"] = tokens.WsExpires
	wst := jwt.NewWithClaims(jwt.SigningMethodHS512, wstClaims)

	tokens.WsToken, err = wst.SignedString([]byte(os.Getenv("WS_SECRET")))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	tokens.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	tokens.RefreshUuid = fmt.Sprintf("%s++%d", tokens.AccessUuid, userId)

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = tokens.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["user_name"] = userName
	rtClaims["exp"] = tokens.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS512, rtClaims)

	tokens.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func extractToken(r *http.Request, isCookie bool) string {
	if isCookie {
		token, _ := r.Cookie("auth_token")
		return token.Value
	} else {
		token := r.Header.Get("Authorization")
		strs := strings.Split(token, " ")
		if len(strs) == 2 {
			return strs[1]
		}
		return ""
	}
}

func verifyToken(r *http.Request, secret string, isRefresh bool) (*jwt.Token, error) {
	tokenString := extractToken(r, isRefresh)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, errors.New("unauthorized")
	}
	return token, nil
}

func extract(token *jwt.Token, uuid string) (*TokenMetadata, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uuid, ok := claims[uuid].(string)
		userId := uint(claims["user_id"].(float64))
		userName, userNameOk := claims["user_name"].(string)
		if !ok || !userNameOk {
			return nil, errors.New("unauthorized")
		} else {
			return &TokenMetadata{
				Uuid:     uuid,
				UserId:   userId,
				UserName: userName,
			}, nil
		}
	}
	return nil, errors.New("something went wrong")
}

func VerifyAccessToken(r *http.Request) (*jwt.Token, error) {
	return verifyToken(r, os.Getenv("ACCESS_SECRET"), false)
}

func ExtractFromAccess(token *jwt.Token) (*TokenMetadata, error) {
	return extract(token, "access_uuid")
}

func ExtractMetadataFromAccess(r *http.Request) (*TokenMetadata, error) {
	token, err := VerifyAccessToken(r)
	if err != nil {
		return nil, err
	}
	metadata, err := ExtractFromAccess(token)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}

func VerifyWsToken(r *http.Request) (*jwt.Token, error) {
	return verifyToken(r, os.Getenv("WS_SECRET"), true)
}

func ExtractFromWs(token *jwt.Token) (*TokenMetadata, error) {
	return extract(token, "ws_uuid")
}

func ExtractMetadataFromWs(r *http.Request) (*TokenMetadata, error) {
	token, err := VerifyWsToken(r)
	if err != nil {
		return nil, err
	}
	metadata, err := ExtractFromWs(token)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}

func VerifyRefreshToken(r *http.Request) (*jwt.Token, error) {
	return verifyToken(r, os.Getenv("REFRESH_SECRET"), true)
}

func ExtractFromRefresh(token *jwt.Token) (*TokenMetadata, error) {
	return extract(token, "refresh_uuid")
}

func ExtractMetadataFromRefresh(r *http.Request) (*TokenMetadata, error) {
	token, err := VerifyRefreshToken(r)
	if err != nil {
		return nil, err
	}
	metadata, err := ExtractFromRefresh(token)
	if err != nil {
		return nil, err
	}
	userId, err := strconv.Atoi(strings.SplitAfter(metadata.Uuid, "++")[1])
	if err != nil {
		return nil, err
	}
	if metadata.UserId != uint(userId) {
		// Try to build a refresh token for someone, tss
		return nil, errors.New("unauthorized")
	}

	return metadata, nil
}
