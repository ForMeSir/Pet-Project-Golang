package repository

import (
	"encoding/base64"
	"fmt"

	"pet/internal/model"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	userTable    = "users"
	sessionTable = "sessions"
	itemTable    = "items"
	RefTimeLive  = 24 * time.Hour
	salt         = "nvkchbenwmksdicuhfdj" // Намеренно оставленно открытым
)

type AuthService struct {
	db *sqlx.DB
}

func NewAuthService(db *sqlx.DB) *AuthService {
	return &AuthService{db: db}
}

func (d *AuthService) CreateUser(name string, username string, password string, userRole string) (id uuid.UUID, err error) {
	data := []byte(salt + password)
	password = base64.StdEncoding.EncodeToString(data)
	id, err = uuid.NewUUID()
	if err != nil {
		fmt.Println(err)
		return
	}
	query := fmt.Sprintf("INSERT INTO %s (id, name, username, password_hash, user_role) values($1,$2,$3,$4,$5)", userTable)
	_, err = d.db.Exec(query, id, name, username, password, userRole)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func (d *AuthService) GetUser(username string, password string) (user model.User, err error) {
	data := []byte(salt + password)
	password = base64.StdEncoding.EncodeToString(data)
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=$1 AND password_hash=$2", userTable)
	row := d.db.QueryRow(query, username, password)
	if err = row.Scan(&user.ID, &user.Name, &user.Username, &user.PasswordHasn, &user.Role); err != nil {
		fmt.Println(err)
		return
	}
	return
}

func (d *AuthService) CreateSession(userId uuid.UUID, userRole string, id uuid.UUID, refreshtoken string) (err error) {
	query := fmt.Sprintf("INSERT INTO %s (id, user_id, refreshtoken, is_blocked, created_at, expirated_at) values($1, $2,$3,$4,$5,$6)", sessionTable)
	_, err = d.db.Exec(query, id, userId, refreshtoken, false, time.Now(), time.Now().Add(RefTimeLive))
	if err != nil {
		fmt.Println(err, 61)
		return
	}
	return
}

func (d *AuthService) FindSession(sessionid uuid.UUID) (session model.Session, err error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", sessionTable)
	row := d.db.QueryRow(query, sessionid)

	if err = row.Scan(&session.Token.SessionID, &session.Token.UserID, &session.Refreshtoken, &session.IsBlocked, &session.Token.CreatedAt, &session.Token.ExpiratedAt); err != nil {
		fmt.Println(err, 62)
		return
	}

	return
}