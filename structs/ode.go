package structs

type Signup struct {
	Id       string `bson:"_id"`
	Email    string
	Login    string
	Password string
}

type Login struct {
	Id       string `bson:"_id"`
	Login    string
	Password string
}

type AddProduct struct {
	Id               string `bson:"_id"`
	Image            string
	ProductName      string
	ProductPrice     int
	ShortDescription string
	OwnerID          string
}
type GetList struct {
	Id       string `bson:"id"`
	Login    string
	Password string
}
type Delete struct {
	ItemID string
}
