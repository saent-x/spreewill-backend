package vendorx

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"spreewill-core/pkg/models"
	"spreewill-core/pkg/session"
	"spreewill-core/pkg/util"
	"strconv"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VendorService struct {
	session *session.ServiceSession
	db      *mgm.Collection
}

func NewVendorService(session *session.ServiceSession) *VendorService {
	return &VendorService{session: session, db: mgm.Coll(&models.Vendor{})}
}

func (v *VendorService) CreateVendor(w http.ResponseWriter, r *http.Request) {
	user_id, err := util.GetUserIdForFromAccessToken(w, r)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req models.Vendor

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := v.db.FindOne(mgm.Ctx(), bson.D{{"user_id", user_id}})

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			req.UserID = *user_id

			req.Verified = false
			err = v.db.Create(&req)

			if err != nil {
				util.SendError(w, http.StatusBadRequest, err.Error())
				return
			}

			util.SendSuccess(w, http.StatusOK, req)
			return
		}
	}

	util.SendError(w, http.StatusBadRequest, "invalid duplicate")
	return
}

func (v *VendorService) GetVendor(w http.ResponseWriter, r *http.Request) {
	user_id, err := util.GetUserIdForFromAccessToken(w, r)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	var vendor models.Vendor

	result := v.db.FindOne(mgm.Ctx(), bson.D{{"user_id", user_id}})
	if result.Err() != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = result.Decode(&vendor)
	if result.Err() != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.SendSuccess(w, http.StatusOK, vendor)
}

func (v *VendorService) GetVendors(w http.ResponseWriter, r *http.Request) {
	user_id, err := util.GetUserIdForFromAccessToken(w, r)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	paginated := true

	pageNos, err := strconv.ParseInt(r.URL.Query().Get("page_number"), 10, 64)
	if err != nil {
		fmt.Println(err)
		paginated = false
	}
	pageSize, err := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 64)
	if err != nil {
		fmt.Println(err)
		paginated = false
	}

	if paginated {
		fmt.Printf("page_size: %d | page_nos: %d \n", pageSize, pageNos)

		skips := pageSize * (pageNos - 1)
		fmt.Printf("skips: %d \n", skips)
		cur, err := v.db.Find(mgm.Ctx(), bson.D{{"user_id", user_id}}, &options.FindOptions{
			Limit: &pageSize,
			Skip:  &skips,
		})

		if err != nil {
			util.SendError(w, http.StatusBadRequest, err.Error())
			return
		}

		util.GetAllInCursor[models.Vendor](cur, w)
		return
	} else {
		cur, err := v.db.Find(mgm.Ctx(), bson.D{{"user_id", user_id}}, nil)
		if err != nil {
			util.SendError(w, http.StatusBadRequest, err.Error())
			return
		}
		util.GetAllInCursor[models.Vendor](cur, w)
		return
	}
}

func (v *VendorService) UpdateVendor(w http.ResponseWriter, r *http.Request) {
	user_id, err := util.GetUserIdForFromAccessToken(w, r)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req models.Vendor

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	update := bson.D{{"$set", req}}

	_, err = v.db.UpdateOne(mgm.Ctx(), bson.D{{"user_id", user_id}}, update)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.SendSuccess(w, http.StatusOK, "success")
}

func (v *VendorService) DeleteVendor(w http.ResponseWriter, r *http.Request) {
	user_id, err := util.GetUserIdForFromAccessToken(w, r)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := v.db.DeleteOne(mgm.Ctx(), bson.D{{"user_id", user_id}})
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.SendSuccess(w, http.StatusOK, result.DeletedCount)
}
