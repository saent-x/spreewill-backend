package customer

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"spreewill-core/pkg/models"
	"spreewill-core/pkg/services/auth"
	"spreewill-core/pkg/session"
	"spreewill-core/pkg/util"
	"strconv"

	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-chi/chi/v5"
)

type CustomerService struct {
	session *session.ServiceSession
	db      *mgm.Collection
}

func NewCustomerService(session *session.ServiceSession) *CustomerService {
	return &CustomerService{session: session, db: mgm.Coll(&models.Customer{})}
}

func (v *CustomerService) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	accessToken := util.GetHeaderToken(w, r)

	if accessToken == "" {
		util.SendError(w, http.StatusBadRequest, "invalid authorization header")
		return
	}
	var req models.Customer

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: verify that the userID exists in cognito
	cognitoClient, ok := r.Context().Value("CognitoClient").(*auth.CognitoClient)
	if !ok {
		util.SendError(w, http.StatusInternalServerError, "could not retrieve cognitoClient from context")
		return
	}

	getUserInput := &cip.GetUserInput{
		AccessToken: &accessToken,
	}

	output, err := cognitoClient.GetUser(r.Context(), getUserInput)

	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := v.db.FindOne(mgm.Ctx(), bson.D{{"user_id", output.Username}})

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			req.UserID = *output.Username

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

func (v *CustomerService) GetCustomer(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")

	if len(param) == 0 {
		util.SendError(w, http.StatusBadRequest, "invalid url params")
		return
	}

	var customer models.Customer

	_id, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = v.db.FindByID(_id, &customer)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.SendSuccess(w, http.StatusOK, customer)
}

func (v *CustomerService) GetCustomers(w http.ResponseWriter, r *http.Request) {
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
		cur, err := v.db.Find(mgm.Ctx(), bson.D{}, &options.FindOptions{
			Limit: &pageSize,
			Skip:  &skips,
		})

		if err != nil {
			util.SendError(w, http.StatusBadRequest, err.Error())
			return
		}

		util.GetAllInCursor[models.Customer](cur, w)
		return
	} else {
		cur, err := v.db.Find(mgm.Ctx(), bson.D{}, nil)
		if err != nil {
			util.SendError(w, http.StatusBadRequest, err.Error())
			return
		}
		util.GetAllInCursor[models.Customer](cur, w)
		return
	}
}

func (v *CustomerService) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	var req models.Customer

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	update := bson.D{{"$set", req}}

	result, err := v.db.UpdateOne(mgm.Ctx(), bson.D{{"_id", req.ID}}, update)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if result.ModifiedCount == 0 {
		util.SendError(w, http.StatusBadRequest, "invalid id")
		return
	}

	util.SendSuccess(w, http.StatusOK, "success")
}

func (v *CustomerService) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")

	if len(param) == 0 {
		util.SendError(w, http.StatusBadRequest, "invalid url params")
		return
	}

	_id, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, "invalid id")
		return
	}

	result, err := v.db.DeleteOne(mgm.Ctx(), bson.D{{"_id", _id}})
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.SendSuccess(w, http.StatusOK, result.DeletedCount)
}
