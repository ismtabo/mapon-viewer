package dao

// User represents an application user DAO.
type User struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}
