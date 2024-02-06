-- sql/queries.sql
-- name: CreateUser :one
INSERT INTO users (name, phone_number, otp, otp_expiration_time) VALUES ($1, $2, $3, $4) RETURNING id;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT (1);

-- name: GetUserByPhoneNumber :one
SELECT * FROM users WHERE phone_number = $1;

-- name: IsOTPValid :one
SELECT COUNT(id) AS x FROM users WHERE phone_number = $1 AND otp=$2 AND otp_expiration_time < NOW();

-- name: UpdateOTPByPhoneNumber :one
UPDATE users SET otp = $3 ,otp_expiration_time = $2 WHERE phone_number = $1 RETURNING *;

-- name: CheckOTPExpire :one
SELECT COUNT(id) AS x FROM users WHERE phone_number = $1 AND otp=$2 AND otp_expiration_time >=NOW();

-- name: CheckOTPExist :one
SELECT COUNT(id) as x FROM users WHERE  otp = $1  ;

-- name: CHECKPHONEEXIST :one
SELECT COUNT(id) as x FROM users WHERE  phone_number = $1 ;


-- name: CHECKOTPPHONEEXIST :one
SELECT COUNT(id) as x FROM users WHERE phone_number = $1 AND otp = $2  ;

-- name: DeleteOTP :exec
DELETE FROM users WHERE id = $1;

