package models

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"gopkg.in/guregu/null.v4"
	"learning.com/learning/utils/helper"
)

type User struct {
	ID   int64  `json:"id" db:"id"`
	Type string `json:"type" db:"type"`
}

// 匿名用户
var AnnoymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnnoymousUser
}

type UserBrief struct {
	User                            // 继承 User
	Name              string        `json:"name" db:"name"`
	Gender            null.String   `json:"gender,omitempty" db:"gender"`
	Avatar            null.String   `json:"avatar,omitempty" db:"avatar"`
	PasswdHash        []byte        `json:"-" db:"password_hash"`
	CreatedAt         time.Time     `json:"created_at" db:"created_at"`
	PurchasingCourses pq.Int64Array `json:"purchasing_courses" db:"purchasing_courses"`
	PurchasedCourses  pq.Int64Array `json:"purchased_courses" db:"purchased_courses"`
}

type UserModel struct {
	DB     *sqlx.DB
	helper *helper.Helper
}

// 通过账号、邮箱、手机号、微信号或 ID 获取 UserBrief
func (m *UserModel) GetUserBriefByAccountOrID(account string, id int64) (*UserBrief, error) {
	query := `
		SELECT
			id, name, type, gender, avatar, password_hash, created_at,
			(
				SELECT ARRAY_AGG(course_id) FROM business.shopping_carts
				WHERE user_id = users.id
			) AS purchasing_courses,
			(
				SELECT ARRAY_AGG(course_id) FROM business.orders
				WHERE user_id = users.id
			) AS purchased_courses
		FROM appuser.users
		WHERE email = $1 OR phone_number = $1 OR id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user UserBrief
	err := m.DB.QueryRowxContext(ctx, query, account, id).StructScan(&user)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, helper.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

// 通过 token 获取 User
//
// 不是 UserBrief。主要用于验证 token 的有效性
func (m *UserModel) GetUserForToken(tokenScope, tokenPlaintext string) (*User, error) {
	// Calculate the SHA-256 hash of the plaintext token provided by the client.
	// Remember that this returns a byte *array* with length 32, not a slice.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	// Set up the SQL query to retrieve the user data for the given token.
	query := `
		SELECT
			u.id as id,
			u.type as type
		FROM appuser.users u
		INNER JOIN appuser.tokens t
		ON u.id = t.user_id
		WHERE t.hash = $1
		AND t.scope = $2
		AND t.expiry > $3`

	// Create a slice containing the query arguments. Notice how we use the [:] operator
	// to get a slice containing the token hash, rather than passing in the array (which
	// is not supported by the pq driver), and that we pass the current time as the
	// value to check against the token expiry.
	args := []any{tokenHash[:], tokenScope, time.Now()}

	// Create a context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query and scan the result into the User struct. If no matching record
	// is found, we return an ErrRecordNotFound error.
	var user User
	err := m.DB.QueryRowxContext(ctx, query, args...).StructScan(&user)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, helper.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

// 通过 ID 获取头像
func (m *UserModel) GetUserAvatarById(id int64) (*string, error) {
	query := `
		SELECT avatar FROM appuser.users
		WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var avatar null.String
	err := m.DB.QueryRowxContext(ctx, query, id).Scan(&avatar)
	if err != nil {
		return nil, err
	}

	if avatar.Valid {
		return &avatar.String, nil
	} else {
		return nil, nil
	}
}
