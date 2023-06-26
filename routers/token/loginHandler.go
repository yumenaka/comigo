package token

//
//import (
//	jwt "github.com/appleboy/gin-jwt/v2"
//)
//
//// LoginHandler can be used by clients to get a jwt token.
//// Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
//// Reply will be of the form {"token": "TOKEN"}.
//func (mw *jwt.GinJWTMiddleware) LoginHandler(c *gin.Context) {
//
//	"github.com/yumenaka/comi/common"
//	if mw.Authenticator == nil {
//		mw.unauthorized(c, http.StatusInternalServerError, mw.HTTPStatusMessageFunc(ErrMissingAuthenticatorFunc, c))
//		return
//	}
//
//	data, err := mw.Authenticator(c)
//	if err != nil {
//		mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, c))
//		return
//	}
//
//	// Create the token
//	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
//	claims := token.Claims.(jwt.MapClaims)
//
//	if mw.PayloadFunc != nil {
//		for key, value := range mw.PayloadFunc(data) {
//			claims[key] = value
//		}
//	}
//
//	expire := mw.TimeFunc().Add(mw.Timeout)
//	claims["exp"] = expire.Unix()
//	claims["orig_iat"] = mw.TimeFunc().Unix()
//	tokenString, err := mw.signedString(token)
//	if err != nil {
//		mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrFailedTokenCreation, c))
//		return
//	}
//
//	// set cookie
//	if mw.SendCookie {
//		expireCookie := mw.TimeFunc().Add(mw.CookieMaxAge)
//		maxage := int(expireCookie.Unix() - mw.TimeFunc().Unix())
//
//		if mw.CookieSameSite != 0 {
//			c.SetSameSite(mw.CookieSameSite)
//		}
//
//		c.SetCookie(
//			mw.CookieName,
//			tokenString,
//			maxage,
//			"/",
//			mw.CookieDomain,
//			mw.SecureCookie,
//			mw.CookieHTTPOnly,
//		)
//	}
//
//	mw.LoginResponse(c, http.StatusOK, tokenString, expire)
//}
