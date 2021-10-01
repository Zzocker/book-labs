package domain

type UserProfile struct {
	ID           string `bson:"_id"`
	Email        string `bson:"email"`
	Name         string `bson:"name"`
	Password     string `bson:"password"`
	ProfilePicID string `bson:"profile_pic_id"`
}
