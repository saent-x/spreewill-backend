package models

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type (
	ResponseObject struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
		Data    any    `json:"data"`
	}

	Vendor struct {
		mgm.DefaultModel `bson:",inline"`
		UserID           string  `json:"user_id,omitempty" bson:"user_id,omitempty"`
		BusinessName     string  `json:"business_name,omitempty" bson:"business_name,omitempty"`
		BusinessAddress  string  `json:"business_address,omitempty" bson:"business_address,omitempty"`
		Firstname        string  `json:"firstname,omitempty" bson:"firstname,omitempty"`
		Lastname         string  `json:"lastname,omitempty" bson:"lastname,omitempty"`
		Gender           string  `json:"gender,omitempty" bson:"gender,omitempty"`
		RCNumber         string  `json:"rc_number,omitempty" bson:"rc_number,omitempty"`
		Phone            string  `json:"phone,omitempty" bson:"phone,omitempty"`
		Email            string  `json:"email,omitempty" bson:"email,omitempty"`
		Industry         string  `json:"industry,omitempty" bson:"industry,omitempty"`
		IGAccount        string  `json:"ig,omitempty" bson:"ig,omitempty"`
		FBAccount        string  `json:"fb,omitempty" bson:"fb,omitempty"`
		Verified         bool    `json:"verified,omitempty" bson:"verified,omitempty"`
		Posts            []Post  `json:"posts,omitempty" bson:"posts,omitempty"`
		Stories          []Story `json:"stories,omitempty" bson:"stories,omitempty"`
	}

	Story struct {
		mgm.DefaultModel `bson:",inline"`
		VendorID         uint    `json:"vendor_id,omitempty" bson:"vendor_id,omitempty"`
		Images           []Image `json:"images,omitempty" bson:"images,omitempty"`
		Views            uint    `json:"views,omitempty" bson:"views,omitempty"`
	}

	Post struct {
		mgm.DefaultModel `bson:",inline"`
		Images           []Image   `json:"images,omitempty" bson:"images,omitempty"`
		VendorID         uint      `json:"vendor_id,omitempty" bson:"vendor_id,omitempty"`
		Caption          string    `json:"caption,omitempty" bson:"caption,omitempty"`
		Likes            uint      `json:"like,omitempty" bson:"like,omitempty"`
		Dislikes         uint      `json:"dislikes,omitempty" bson:"dislikes,omitempty"`
		Views            uint      `json:"views,omitempty" bson:"views,omitempty"`
		Comments         []Comment `json:"comments,omitempty" bson:"comments,omitempty"`
	}

	Comment struct {
		mgm.DefaultModel `bson:",inline"`
		PostID           uint      `json:"post_id,omitempty" bson:"post_id,omitempty"`
		Content          string    `json:"content,omitempty" bson:"content,omitempty"`
		DateTime         time.Time `json:"time,omitempty" bson:"time,omitempty"`
		UserID           string    `json:"user_id,omitempty" bson:"user_id,omitempty"`
	}

	Image struct {
		mgm.DefaultModel `bson:",inline"`
		PostID           uint
		StoryID          uint
		Link             string `json:"link,omitempty" bson:"link,omitempty"`
	}

	Customer struct {
		mgm.DefaultModel `bson:",inline"`
		UserID           string `json:"user_id,omitempty" bson:"user_id,omitempty"`
		Username         string `json:"username,omitempty" bson:"username,omitempty"`
		Firstname        string `json:"firstname,omitempty" bson:"firstname,omitempty"`
		Lastname         string `json:"lastname,omitempty" bson:"lastname,omitempty"`
		Address          string `json:"address,omitempty" bson:"address,omitempty"`
		Phone            string `json:"phone,omitempty" bson:"phone,omitempty"`
		Email            string `json:"email,omitempty" bson:"email,omitempty"`
		Interests        string `json:"interests,omitempty" bson:"interests,omitempty"`
		Verified         bool   `json:"-" gorm:"type:bool;default:false" bson:"verified"`
	}

	Entity interface {
		Customer | Image | Post | Comment | Story | Vendor
	}
)
