package goodidea

import (
    "fmt"
	"time"
    "strings"
	"context"
    "net/http"
    "crypto/rand"
	"crypto/sha512"
	"encoding/hex"

	//"github.com/rs/xid"
)

// users within the system
//
//	id - unique ID for the user, an xid prefixed with 'u'
//	name - the username for the given user
//	salt - the salt used befored hashing the passwd
//	passwd - the user's password used for logging in
//	createdAt - the date only, day the user signed up
//	admin - if the user is an admin user or not
type user struct {
	ID        string    `json:"id" validate:"required,u20"`
	Name      string    `json:"name" validate:"required,max=48"`
	Salt      string    `json:"-" validate:"omitempty,len=8"`
	Passwd    string    `json:"passwd" validate:"required,sha256"`
	CreatedAt time.Time `json:"createdAt"`
	Admin     bool      `json:"admin"`
}

// sessions are for user logins
type session struct {
    //session IDs should be cryptographically secure
	SessionID string    `json:"sessionId"`
    //the user in question
	UserID    string    `json:"userId" validate:"required,u20"`
    //When the token was 
	CreatedAt time.Time `json:"createdAt"`
}

func handleSignUp(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
	if err != nil {
		Logr.Error("Error parsing sign-up form", "error", err)
		fmt.Fprintf(w, "<p>Oops! Looks like the form wasn't submitted correctly</p>")
		return
	}

    pw := strings.TrimSpace(r.FormValue("password"))
	if pw == "" {
		fmt.Fprintf(w, "<p>Oops! Looks like the no password was submitted</p>")
		return
	}
    if  len(pw) < 12 {
		fmt.Fprintf(w, "<p>Oops! Looks like the password is not long enough (<12)</p>")
		return
    }
    u := strings.TrimSpace(r.FormValue("username"))
	if u != "" {
		fmt.Fprintf(w, "<p>Oops! Looks like the no username was submitted</p>")
        return
	}
    exists, err := checkNameExistence(u) 
    if err != nil {
		Logr.Error("Error checking for username existence", "error", err)
		fmt.Fprintf(w, "<p>Oops! Server error attempting to sign you up</p>")
		return
    }
    if exists {
		fmt.Fprintf(w, "<p>Sorry, that username is already taken.</p>")
		return
    }

    if err := createUser(u, pw); err != nil {
		Logr.Error("Error checking for username existence", "error", err)
		fmt.Fprintf(w, "<p>Oops! Server error attempting to sign you up</p>")
		return
    }

    http.Redirect(w, r, "/login", http.StatusFound)
}

func createUser(name, passwd string) error {
    s := generateRandomSalt()
    p := hashPassword(passwd, s)
    return insertUser(name, p, string(s))
}

func checkNameExistence(u string) (bool, error) {
    ctx := context.TODO()
    var exists bool
    q := "SELECT EXISTS(SELECT 1 FROM users WHERE name= $1)"
    if err := DB.QueryRow(ctx, q).Scan(&exists); err != nil {
        Logr.Error("Unable to query the database for username", "name", u, "error", err)
        return exists, err
    }
    return exists, nil
}

//persistSession in-case of server reload
func persistSession(s *session) error {
    ctx := context.TODO()
    q := "INSERT INTO sessions(session_id, user_id, created_at) VALUES ($1, $2, $3)"
	_, err := DB.Exec(ctx, q, s.SessionID, s.UserID, s.CreatedAt)
    return err
}

//Create a new user
func insertUser(name, passwd, salt string) error {
    ctx := context.TODO()
    q := "INSERT INTO users(name, passwd, salt) VALUES ($1, $2, $3)"
	_, err := DB.Exec(ctx, q, name, passwd, salt)
    return err
}

func createUserSession(u *user) (*session, error) {
    b := make([]byte, 32)
	_, err := rand.Read(b[:])
    if err != nil {
        return nil, err
    }

	var sessionId = hex.EncodeToString(b)
    fmt.Println(sessionId)
    fmt.Println(len(sessionId))

    return &session{
        SessionID: sessionId,
        UserID: u.ID,
        CreatedAt: time.Now(),
    }, nil
}

//getUser by username to check password in a loging
func getUser(name string) (u *user, err error) {
	ctx := context.Background()
	q := `SELECT id, salt, passwd FROM users WHERE name = $1`
	if err := DB.QueryRow(ctx, q).Scan(&u.ID, &u.Salt, &u.Passwd); err != nil {
        return u, err
	}
    u.Name = name
    return u, nil
}


// Check if two passwords match
func doPasswordsMatch(hashedPassword, currPassword string, salt []byte) bool {
	var currPasswordHash = hashPassword(currPassword, salt)

	return hashedPassword == currPasswordHash
}

func generateRandomSalt() []byte {
	var salt = make([]byte, 8)
	_, err := rand.Read(salt[:])

	if err != nil {
		panic(err)
	}

	return salt
}

// Combine password and salt then hash them using the SHA-512
// hashing algorithm and then return the hashed password
// as a hex string
func hashPassword(password string, salt []byte) string {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Append salt to password
	passwordBytes = append(passwordBytes, salt...)

	var sha512Hasher = sha512.New()
	sha512Hasher.Write(passwordBytes)
	var hashedPasswordBytes = sha512Hasher.Sum(nil)
	// Convert the hashed password to a hex string
	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHex
}

