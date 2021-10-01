package models

import "database/sql"

type AllUserInfo struct {
	Profile    Profile         `json:"profile" bson:"profile" mapstructure:"profile"`
	Hotel      []ClientRequest `json:"hotel" bson:"hotel" mapstructure:"hotel"`
	Restaurant []struct {
		OrdersAndPrice `mapstructure:",squash"`
		FullName       string `json:"fullName" bson:"receiver" mapstructure:"receiver" validate:"required"`
		OrderedDate    string `json:"orderedDate" mapstructure:"orderedDate"`
	} `json:"restaurant"`
}

type UserSignupInfo struct {
	GeneralUserInfo `mapstructure:",squash"`
	LoginInfo       `mapstructure:",squash"`
}

type GeneralUserInfo struct {
	FirstName   string `json:"firstName" bson:"firstName" mapstructure:"firstName" validate:"required"`
	LastName    string `json:"lastName" bson:"lastName" mapstructure:"lastName" validate:"required"`
	PhoneNumber string `json:"phoneNumber" bson:"phoneNumber" mapstructure:"phoneNumber" validate:"required"`
	Address     string `json:"address" bson:"address" mapstructure:"address"`
}

type Profile struct {
	GeneralUserInfo `mapstructure:",squash"`
	LoginInfo       `mapstructure:",squash"`
}

type LoginInfo struct {
	Email    string `json:"email" bson:"email" mapstructure:"email" validate:"email,required"`
	Password string `json:"password,omitempty" bson:"password" mapstructure:"password" validate:"required"`
}

type FullName struct {
	FirstName string `bson:"firstName" mapstructure:"firstName"`
	LastName  string `bson:"lastName" mapstructure:"lastName"`
}

type OrdersAndPrice struct {
	Orders []struct {
		Name          string `json:"name" bson:"name" validate:"required"`
		NumberOfMeals byte   `json:"numberOfMeals" bson:"numberOfMeals" mapstructure:"numberOfMeals" validate:"required"`
	} `json:"orders"`
	TotalPrice int64 `json:"totalPrice"  mapstructure:"totalPrice" validate:"required"`
}

type FoodOrder struct {
	OrdersAndPrice
	CustomerForm struct {
		FullName    string `json:"fullName" bson:"receiver" validate:"required"`
		Address     string `json:"address" bson:"address" validate:"required"`
		PhoneNumber string `json:"phoneNumber" bson:"phoneNumber" validate:"required"`
	} `json:"customerForm"`
}

type NumberAndGenericSubtype struct {
	GenericSubtype string `json:"genericSubtype" bson:"roomSubtype" mapstructure:"roomSubtype"`    // 3 possibilites
	NumberOfRooms  byte   `json:"numberOfRooms" bson:"numberOfRooms" mapstructure:"numberOfRooms"` // 1 to 3
}
type HotelReservation struct {
	NumberAndGenericSubtype `mapstructure:",squash"`
	StartDate               string `json:"startDate" bson:"startDate" mapstructure:"startDate"`
	EndDate                 string `json:"endDate" bson:"endDate" mapstructure:"endDate"`
}
type ClientRequest struct {
	RoomType         string `json:"roomTypeSlug" bson:"roomType" mapstructure:"roomType"` // table name
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

type Bookings struct {
	Id       int16  `json:"id" bson:"id" mapstructure:"id"`
	FullName string `json:"fullName" bson:"fullName" mapstructure:"fullName"`
	// RoomMark uint8  `json:"roomMark" bson:"roomMark" mapstructure:"roomMark"`
	RoomMark         sql.NullInt32 `json:"roomMark" bson:"roomMark" mapstructure:"roomMark"`
	Email            string        `json:"email" bson:"email" mapstructure:"email"`
	HotelReservation `mapstructure:",squash"`
}
