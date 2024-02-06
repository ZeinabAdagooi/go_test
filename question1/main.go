package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	sqlc "github.com/zeinab/question1/usersqlc"
)

var dbconn *pgx.Conn
var queries *sqlc.Queries

func initdb() (*pgx.Conn, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("SSL_MODE"))
	print(connString)
	var errConn error
	//connString = "postgresql://aquser:123@localhost:5433/aqrdb?sslmode=disable"

	dbconn, errConn = pgx.Connect(context.Background(), connString)
	if errConn != nil {
		log.Fatal(errConn)
	}

	return dbconn, nil
}
func main() {

	r := gin.Default()
	// Initialize database

	dbconn, err := initdb()
	if err != nil {
		log.Fatal(err)
	}

	queries := sqlc.New(dbconn)
	initRoutes(r, queries)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}

func initRoutes(r *gin.Engine, db *sqlc.Queries) {
	// Initialize the repository
	qdb := db
	// User routes
	userGroup := r.Group("api/")
	{
		userGroup.POST("/createUser", createUser(qdb))
		userGroup.POST("/generateotp", generateOTP(qdb))
		userGroup.POST("/verifyOTP", verifyOTP(qdb))
	}
}
func createUser(repo *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse input

		var newUser sqlc.User

		//var err error
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		X, err := repo.CHECKPHONEEXIST(context.Background(), newUser.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if X > 0 {
			c.JSON(400, "phone number exists")
			return
		}

		result, err := repo.CreateUser(context.Background(), sqlc.CreateUserParams{
			Name:              newUser.Name,
			PhoneNumber:       newUser.PhoneNumber,
			Otp:               newUser.Otp,
			OtpExpirationTime: newUser.OtpExpirationTime},
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Println(result)

		// Return user data

		c.JSON(http.StatusCreated, newUser)
	}
}

func generateOTP(repo *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse input

		var request sqlc.UpdateOTPByPhoneNumberParams
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		X, err := repo.CHECKPHONEEXIST(context.Background(), request.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if X == 0 {
			c.JSON(401, "phone number not exists")
			return
		}
		// Generate random 4-digit OTP
		otp := fmt.Sprintf("%04d", rand.Intn(10000))
		otptime := time.Now().Add(1 * time.Minute)
		var pgtime pgtype.Timestamp
		pgtime.Scan(otptime)

		// Generate OTP
		result, err := repo.UpdateOTPByPhoneNumber(context.Background(), sqlc.UpdateOTPByPhoneNumberParams{PhoneNumber: request.PhoneNumber,
			OtpExpirationTime: pgtime,
			Otp:               otp})

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Return OTP data
		c.JSON(200, result)
	}
}

func verifyOTP(repo *sqlc.Queries) gin.HandlerFunc {

	return func(c *gin.Context) {
		// Parse input
		var request sqlc.CHECKOTPPHONEEXISTParams
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		Z, err := repo.CheckOTPExist(context.Background(), request.Otp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if Z == 0 {
			c.JSON(401, " OTP not exists")
			return
		}

		Y, err := repo.CHECKPHONEEXIST(context.Background(), request.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if Y == 0 {
			c.JSON(401, "Phone Number not exists")
			return
		}

		X, err := repo.CHECKOTPPHONEEXIST(context.Background(), request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if X == 0 {
			c.JSON(401, "phone number with otp not exists")
			return
		}

		// Verify OTP
		A, err := repo.CheckOTPExpire(context.Background(),
			sqlc.CheckOTPExpireParams{PhoneNumber: request.PhoneNumber, Otp: request.Otp})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if A == 0 {
			c.JSON(400, " OTP Expired")
			return
		}

		result, err := repo.UpdateOTPByPhoneNumber(context.Background(), sqlc.UpdateOTPByPhoneNumberParams{PhoneNumber: request.PhoneNumber,
			OtpExpirationTime: pgtype.Timestamp{},
			Otp:               ""})

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		log.Println(result)

		// Return success
		c.JSON(200, gin.H{"message": "OTP verified successfully"})

	}
}
