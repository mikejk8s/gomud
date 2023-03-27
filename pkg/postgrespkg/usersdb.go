import (
"time"

"gorm.io/gorm"

"github.com/mikejk8s/gmud/pkg/models"
)

type SqlConn struct {
DB *gorm.DB
}

func (s *SqlConn) CreateUsersTable() {
s.DB.Exec(`CREATE TABLE IF NOT EXISTS users (
        id              SERIAL PRIMARY KEY,
        created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at      TIMESTAMP,
        deleted_at      TIMESTAMP,
        name            VARCHAR(255),
        password_hash   VARCHAR(255),
        remember_hash   VARCHAR(255)
    );`)
}

func (s *SqlConn) CreateNewUser(userInfo LoginReq) error {
user := &models.User{
CreatedAt:    time.Now(),
Name:         userInfo.Name,
PasswordHash: userInfo.Password,
RememberHash: "",
}

if err := s.DB.Create(user).Error; err != nil {
return err
}
return nil
}
