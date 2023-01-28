package main

import (
	"fmt"
	"github.com/danymarita/redis-otp-sample/dto"
	cache2 "github.com/danymarita/redis-otp-sample/pkg/cache"
	"github.com/danymarita/redis-otp-sample/pkg/config"
	"github.com/danymarita/redis-otp-sample/util"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"net/http"
	"strconv"
	"time"
)

const (
	otpCountKeyFormat  string = "user:%d:otp_request_count"
	otpValueKeyFormat  string = "user:%d:otp_request_value"
	otpForbidKeyFormat string = "user:%d:otp_forbid"
)

func main() {
	cfg := config.Config()
	cache := cache2.NewCache(cfg)
	e := echo.New()
	otpRoute := e.Group("/otp")
	otpRoute.POST("/request", func(c echo.Context) error {
		var err error
		req := new(dto.OtpRequest)
		if err = c.Bind(req); err != nil {
			return c.String(http.StatusBadRequest, "Invalid OTP request")
		}

		var (
			otpCountKey       string        = fmt.Sprintf(otpCountKeyFormat, req.UserID)
			otpValueKey       string        = fmt.Sprintf(otpValueKeyFormat, req.UserID)
			otpForbidKey      string        = fmt.Sprintf(otpForbidKeyFormat, req.UserID)
			otpDuration       time.Duration = cast.ToDuration("5m")
			otpForbidDuration time.Duration = cast.ToDuration("60m")
			otpCode           int
		)
		//Check OTP forbid cache
		otpForbidExist := cache.CheckCacheExists(otpForbidKey)
		if otpForbidExist {
			return c.String(http.StatusForbidden, "Your OTP request reach limit")
		}
		//Check OTP count cache
		otpCountExist := cache.CheckCacheExists(otpCountKey)
		if otpCountExist {
			countB, err := cache.ReadCache(otpCountKey)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Err get OTP count: "+err.Error())
			}
			count, err := strconv.Atoi(string(countB))
			if err != nil {
				return c.String(http.StatusInternalServerError, "Err convert OTP count: "+err.Error())
			}
			if count < 3 {
				otpCode, err = util.GenerateRandomNumber(6)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Err generate OTP code: "+err.Error())
				}
				_, err = cache.IncrementCache(otpCountKey)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Err increment OTP count: "+err.Error())
				}
				otpB := []byte(strconv.Itoa(otpCode))
				err = cache.WriteCache(otpValueKey, otpB, otpDuration)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Err write OTP value: "+err.Error())
				}
			} else {
				countB := []byte(strconv.Itoa(1))
				err = cache.WriteCache(otpForbidKey, countB, otpForbidDuration)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Err write OTP forbid: "+err.Error())
				}
				return c.String(http.StatusForbidden, "Your OTP request reach limit, wait for 60 minutes to make request again")
			}
		} else {
			otpCode, err = util.GenerateRandomNumber(6)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Err generate OTP code: "+err.Error())
			}
			countB := []byte(strconv.Itoa(1))
			err = cache.WriteCache(otpCountKey, countB, otpDuration)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Err write OTP count: "+err.Error())
			}
			otpB := []byte(strconv.Itoa(otpCode))
			err = cache.WriteCache(otpValueKey, otpB, otpDuration)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Err write OTP value: "+err.Error())
			}
		}
		resp := dto.OtpResponse{
			UserID:  req.UserID,
			OtpCode: otpCode,
		}
		return c.JSON(http.StatusOK, resp)
	})
	otpRoute.POST("/validate", func(c echo.Context) error {
		var err error
		req := new(dto.OtpValidateRequest)
		if err = c.Bind(req); err != nil {
			return c.String(http.StatusBadRequest, "Invalid OTP validation request")
		}
		var (
			otpValueKey string = fmt.Sprintf(otpValueKeyFormat, req.UserID)
			otpCountKey string = fmt.Sprintf(otpCountKeyFormat, req.UserID)
		)
		otpB, err := cache.ReadCache(otpValueKey)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Err get OTP value: "+err.Error())
		}
		otp, err := strconv.Atoi(string(otpB))
		if err != nil {
			return c.String(http.StatusInternalServerError, "Err convert OTP value: "+err.Error())
		}
		if otp != req.OtpCode {
			return c.String(http.StatusBadRequest, "Invalid OTP code")
		}
		//Delete OTP cache
		err = cache.DeleteCache(otpValueKey)
		if err != nil {
			e.Logger.Error("Err delete OTP value cache: " + err.Error())
		}
		err = cache.DeleteCache(otpCountKey)
		if err != nil {
			e.Logger.Error("Err delete OTP count cache: " + err.Error())
		}
		return c.String(http.StatusOK, "Validate OTP success")
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", cfg.AppHost, cast.ToInt(cfg.AppPort))))
}
