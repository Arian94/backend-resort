package models

type AllUserInfo struct {
	Profile    Profile         `json:"profile" bson:"profile" mapstructure:"profile"`
	Hotel      []ClientRequest `json:"hotel" bson:"hotel" mapstructure:"hotel"`
	Restaurant []struct {
		OrdersAndPrice `mapstructure:",squash"`
		FullName       string `json:"fullname" bson:"receiver" mapstructure:"receiver" validate:"required"`
		OrderedDate    string `json:"orderedDate" mapstructure:"ordered_date"`
	} `json:"restaurant"`
}

type UserSignupInfo struct {
	GeneralUserInfo `mapstructure:",squash"`
	LoginInfo       `mapstructure:",squash"`
}

type GeneralUserInfo struct {
	Firstname   string `json:"firstname" bson:"firstname" mapstructure:"firstname" validate:"required"`
	Lastname    string `json:"lastname" bson:"lastname" mapstructure:"lastname" validate:"required"`
	PhoneNumber string `json:"phoneNumber" bson:"phone_number" mapstructure:"phone_number" validate:"required"`
	Address     string `json:"address" bson:"address" mapstructure:"address"`
}

type Profile struct {
	GeneralUserInfo `mapstructure:",squash"`
	LoginInfo       `mapstructure:",squash"` // email - password
}

type LoginInfo struct {
	Email    string `json:"email" bson:"email" mapstructure:"email" validate:"email,required"`
	Password string `json:"password" bson:"password" mapstructure:"password" validate:"required"`
}

type Fullname struct {
	Firstname string `bson:"firstname"`
	Lastname  string `bson:"lastname"`
}

type OrdersAndPrice struct {
	Orders []struct {
		Name          string `json:"name" bson:"name" validate:"required"`
		NumberOfMeals byte   `json:"numberOfMeals" bson:"number_of_meals" mapstructure:"number_of_meals" validate:"required"`
	} `json:"orders"`
	TotalPrice int64 `json:"totalPrice"  mapstructure:"total_price" validate:"required"`
}

type FoodOrder struct {
	OrdersAndPrice
	Customer struct {
		FullName    string `json:"fullName" bson:"receiver" validate:"required"`
		Address     string `json:"address" bson:"address" validate:"required"`
		PhoneNumber string `json:"phoneNumber" bson:"phone_number" validate:"required"`
	} `json:"customer"`
}

type NumberAndGenericSubtype struct {
	GenericSubtype string `json:"genericSubtype" bson:"room_subtype" mapstructure:"room_subtype"`      // 3 possibilites
	NumberOfRooms  byte   `json:"numberOfRooms" bson:"number_of_rooms" mapstructure:"number_of_rooms"` // 1 to 3
}
type HotelReservation struct {
	NumberAndGenericSubtype `mapstructure:",squash"`
	StartDate               string `json:"startDate" bson:"start_date" mapstructure:"start_date"`
	EndDate                 string `json:"endDate" bson:"end_date" mapstructure:"end_date"`
}
type ClientRequest struct {
	RoomType         string `json:"roomTypeSlug" bson:"room_type" mapstructure:"room_type"` // table name
	HotelReservation `mapstructure:",squash"`
}

type FoodList struct {
	Type        string `json:"type"`
	Origin      string `json:"origin"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Price       uint8  `json:"price"`
}
