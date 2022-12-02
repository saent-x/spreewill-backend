package comments

import (
	"encoding/json"
	"fmt"
	"net/http"
	"spreewill-core/pkg/models"
	"spreewill-core/pkg/session"
	"spreewill-core/pkg/util"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-chi/chi/v5"
)

type CommentService struct {
	session *session.ServiceSession
	db      *mgm.Collection
}

type CommentDTO struct {
	PostID  string    `json:"post_id,omitempty"`
	Content string    `json:"content,omitiempty"`
	Time    time.Time `json:"time,omitempty"`
	UserID  string    `json:"user_id,omitempty"`
}

func NewCommentService(session *session.ServiceSession) *CommentService {
	return &CommentService{session: session, db: mgm.Coll(&models.Comment{})}
}

func (v *CommentService) CreateComment(w http.ResponseWriter, r *http.Request) {
	accessToken := util.GetHeaderToken(w, r)

	if accessToken == "" {
		util.SendError(w, http.StatusBadRequest, "invalid authorization header")
		return
	}
	var req models.Comment
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = v.db.Create(&req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.SendSuccess(w, http.StatusOK, "success")
}

func (v *CommentService) GetComment(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")

	if len(param) == 0 {
		util.SendError(w, http.StatusBadRequest, "invalid url params")
		return
	}

	var comment models.Comment
	_id, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = v.db.FindByID(_id, &comment)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.SendSuccess(w, http.StatusOK, comment)
}

func (v *CommentService) GetComments(w http.ResponseWriter, r *http.Request) {
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

		util.GetAllInCursor[models.Comment](cur, w)
		return
	} else {
		cur, err := v.db.Find(mgm.Ctx(), bson.D{}, nil)
		if err != nil {
			util.SendError(w, http.StatusBadRequest, err.Error())
			return
		}
		util.GetAllInCursor[models.Comment](cur, w)
		return
	}
}

func (v *CommentService) UpdateComment(w http.ResponseWriter, r *http.Request) {
	var req models.Comment
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

func (v *CommentService) DeleteComment(w http.ResponseWriter, r *http.Request) {
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
