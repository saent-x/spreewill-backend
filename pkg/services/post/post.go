package post

import (
	"encoding/json"
	"fmt"
	"net/http"
	"spreewill-core/pkg/models"
	"spreewill-core/pkg/session"
	"spreewill-core/pkg/util"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-chi/chi/v5"
)

type PostService struct {
	session *session.ServiceSession
	db      *mgm.Collection
}

func NewPostService(session *session.ServiceSession) *PostService {
	return &PostService{session: session, db: mgm.Coll(&models.Post{})}
}

func (v *PostService) CreatePost(w http.ResponseWriter, r *http.Request) {
	accessToken := util.GetHeaderToken(w, r)

	if accessToken == "" {
		util.SendError(w, http.StatusBadRequest, "invalid authorization header")
		return
	}
	var req models.Post
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	req.Likes = 0
	req.Dislikes = 0

	err = v.db.Create(&req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.SendSuccess(w, http.StatusOK, "success")
}

func (v *PostService) Like(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	post_id := req["id"].(string)
	var post models.Post

	_id, err := primitive.ObjectIDFromHex(post_id)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = v.db.FindByID(_id, &post)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, "id not found")
		return
	}

	like := req["like"].(float64)

	var newCount uint
	if like == 0 {
		if post.Likes == 0 {
			newCount = 0
		} else {
			newCount = post.Likes - 1
		}
	} else if like == 1 {
		newCount = post.Likes + 1
	}

	update := bson.D{{"$set", bson.D{{"likes", newCount}}}}

	result, err := v.db.UpdateOne(mgm.Ctx(), bson.D{{"_id", _id}}, update)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if result.ModifiedCount == 0 {
		util.SendError(w, http.StatusNotModified, "no action")
		return
	}

	util.SendSuccess(w, http.StatusOK, "success")

}

func (v *PostService) Dislike(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	post_id := req["id"].(string)
	var post models.Post

	_id, err := primitive.ObjectIDFromHex(post_id)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = v.db.FindByID(_id, &post)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, "id not found")
		return
	}

	dislike := req["dislike"].(float64)

	var newCount uint
	if dislike == 0 {
		if post.Dislikes == 0 {
			newCount = 0
		} else {
			newCount = post.Dislikes - 1
		}
	} else if dislike == 1 {
		newCount = post.Dislikes + 1
	} else {
		util.SendError(w, http.StatusBadRequest, "invalid parameters")
		return
	}

	fmt.Println(post.Dislikes)
	update := bson.D{{"$set", bson.D{{"dislikes", newCount}}}}

	result, err := v.db.UpdateOne(mgm.Ctx(), bson.D{{"_id", _id}}, update)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if result.ModifiedCount == 0 {
		util.SendError(w, http.StatusNotModified, "no action")
		return
	}

	util.SendSuccess(w, http.StatusOK, "success")
}

func (v *PostService) GetPost(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")

	if len(param) == 0 {
		util.SendError(w, http.StatusBadRequest, "invalid url params")
		return
	}

	var post models.Post
	_id, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = v.db.FindByID(_id, &post)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.SendSuccess(w, http.StatusOK, post)
}

func (v *PostService) GetPosts(w http.ResponseWriter, r *http.Request) {
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

		util.GetAllInCursor[models.Post](cur, w)
		return
	} else {
		cur, err := v.db.Find(mgm.Ctx(), bson.D{}, nil)
		if err != nil {
			util.SendError(w, http.StatusBadRequest, err.Error())
			return
		}
		util.GetAllInCursor[models.Post](cur, w)
		return
	}
}

func (v *PostService) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var req models.Post
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

func (v *PostService) DeletePost(w http.ResponseWriter, r *http.Request) {
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
