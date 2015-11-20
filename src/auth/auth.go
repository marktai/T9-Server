package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"db"
	"encoding/base64"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/big"
	mrand "math/rand"
	"net/http"
	"strconv"
	"time"
)

var secretMap = make(map[uint]*Secret)

type Secret struct {
	value      *big.Int
	expiration time.Time
}

func (s *Secret) Bytes() []byte {
	return s.value.Bytes()
}

var bitSize int64 = 512

var limit *big.Int

func newSecret() *Secret {
	if limit == nil {
		limit = big.NewInt(2)
		limit.Exp(big.NewInt(2), big.NewInt(bitSize), nil)
	}

	value, err := rand.Int(rand.Reader, limit)
	if err != nil {
		log.Panic(err)
	}
	retSecret := &Secret{value, time.Now()}
	return retSecret
}

func stringtoUint(s string) (uint, error) {
	i, err := strconv.Atoi(s)
	return uint(i), err
}

func getUniqueID() (uint, error) {

	mrand.Seed(time.Now().Unix())
	collision := 1
	times := 0
	var id uint

	for collision != 0 {

		id = uint(mrand.Int31n(65536))
		err := db.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=?)", id).Scan(&collision)
		if err != nil {
			return id, err
		}
		times++
		if times > 20 {
			return id, errors.New("Too many attempts to find a unique game ID")
		}
	}
	return id, nil
}

func MakeUser(user, pass string) (uint, error) {
	_, err := getUserID(user)
	if err == nil {
		return 0, errors.New("User already made")
	}

	id, err := getUniqueID()
	if err != nil {
		return 0, err
	}

	saltHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	saltHashString := base64.StdEncoding.EncodeToString(saltHash)

	err = db.Db.Ping()
	if err != nil {
		return 0, err
	}

	addUser, err := db.Db.Prepare("INSERT INTO users VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}

	_, err = addUser.Exec(id, user, saltHashString)

	if err != nil {
		return 0, err
	}

	return id, nil

}

func getUserID(user string) (uint, error) {

	var userID uint
	err := db.Db.QueryRow("SELECT id FROM users WHERE name=?", user).Scan(&userID)

	return userID, err
}

func getSaltHash(userID uint) ([]byte, error) {
	saltHashString := ""
	err := db.Db.QueryRow("SELECT salthash FROM users WHERE id=?", userID).Scan(&saltHashString)
	if err != nil {
		return nil, err
	}
	saltHash, err := base64.StdEncoding.DecodeString(saltHashString)
	return saltHash, err
}

func Login(user, pass string) (uint, *Secret, error) {

	userID, err := getUserID(user)
	hash, err := getSaltHash(userID)
	if err != nil {
		return 0, nil, err
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(pass))

	if err != nil {
		return 0, nil, err
	}

	if _, ok := secretMap[userID]; !ok {
		secret := newSecret()
		secretMap[userID] = secret
	}

	return userID, secretMap[userID], nil
}

// CheckMAC reports whether messageHMAC is a valid HMAC tag for message.
func checkMAC(message, messageHMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageHMAC, expectedMAC)
}

// returns userID, message used to generate HMAC, and HMAC from request
func parseRequest(r *http.Request) (uint, []byte, []byte, error) {
	userID, err := stringtoUint(r.FormValue("ID"))
	if err != nil {
		return 0, nil, nil, err
	}

	timeSlice, ok := r.Header["Time-Sent"]

	time := timeSlice[0]

	if err != nil {
		return 0, nil, nil, err
	}

	path := r.URL.Path

	message := append([]byte(time), []byte(path)...)

	messageHMACSlice, ok := r.Header["Hmac"]
	if !ok {
		return 0, nil, nil, errors.New("No HMAC header provided")
	}

	messageHMACString := messageHMACSlice[0]
	messageHMAC, err := base64.StdEncoding.DecodeString(messageHMACString)
	if err != nil {
		return 0, nil, nil, err
	}

	return userID, message, messageHMAC, nil
}

func AuthRequest(r *http.Request) (bool, error) {
	userID, message, messageHMAC, err := parseRequest(r)
	if err != nil {
		return false, err
	}
	secret, ok := secretMap[userID]
	if !ok {
		return false, errors.New("No secret for that user")
	}

	secretBytes := secret.Bytes()
	log.Println(secretBytes)
	return checkMAC(message, messageHMAC, secretBytes), nil
}
