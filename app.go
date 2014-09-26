package hello

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)
import "github.com/crhym3/go-endpoints/endpoints"

const (
	accountKind = "Account"
	actionKind   = "Action"
)

// Action is a
// It also serves as (a part of) a response of ActionsService.
type Action struct {
	Key          *datastore.Key `json:"id" datastore:"-"`
	ActionWords  []string       `json:"actionWords" datastore:"actionwords"`
	RedirectLink string         `json:"redirectLink" datastore:"redirect_link,noindex"`
	Date         time.Time      `json:"date" datastore:"date"`
	UserID       string         `json:"-" datastore:"user_id"`
}

// ActionsService can sign the guesbook, list all actions and delete
// a action from the guestbook.
type ActionsService struct {
}

// ActionsListResp is a response type of ActionsService.List method
type ActionsListResp struct {
	Items []*Action `json:"items"`
}

// Request type for ActionsService.List
type ActionsListReq struct{}

// List returns a list of matching actions
func (as *ActionsService) List(r *http.Request, req *ActionsListReq, resp *ActionsListResp) error {
	c := endpoints.NewContext(r)
	u, err := getUser(c)
	if err != nil {
		return err
	}
	userKey := makeUserKey(c, u.ID)
	q := datastore.NewQuery(actionKind).Ancestor(userKey)
	var actions []*Action
	keys, err := q.GetAll(c, &actions)
	if err != nil {
		return err
	}

	for i, k := range keys {
		actions[i].Key = k
	}
	resp.Items = actions
	return nil
}

// ActionAddResp is a response type of ActionsService.List method
type ActionAddResp struct{}

//Request type for ActionsService.List
type ActionAddReq struct {
	Words    string
	Redirect string
}

func makeUserKey(c appengine.Context, userID string) *datastore.Key {
	return datastore.NewKey(c, accountKind, userID, 0, nil)
}

// Add adds an action.
func (as *ActionsService) Add(r *http.Request, req *ActionAddReq, resp *ActionAddResp) error {
	c := endpoints.NewContext(r)
	u, err := getUser(c)
	if err != nil {
		return err
	}
	act := &Action{
		Key:          nil,
		ActionWords:  strings.Split(req.Words, " "),
		RedirectLink: req.Redirect,
		Date:         time.Now(),
		UserID:       u.ID,
	}
	putKey := datastore.NewIncompleteKey(c, actionKind, makeUserKey(c, u.ID)) // no id, let it auto generate.
	_, err = datastore.Put(c, putKey, act)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	actionsService := &ActionsService{}
	api, err := endpoints.RegisterService(actionsService,
		"action", "v1", "Actions API", true)
	if err != nil {
		panic(err.Error())
	}

	info := api.MethodByName("List").Info()
	info.Name, info.HttpMethod, info.Path, info.Desc =
		"list", "GET", "list", "List most recent actions."

	add := api.MethodByName("Add").Info()
	add.Name, add.HttpMethod, add.Path, add.Desc =
		"add", "PUT", "add", "Add an action."

	endpoints.HandleHttp()
	http.HandleFunc("/", handler)
}

func getUser(ctx appengine.Context) (*user.User, error) {
	u := user.Current(ctx)
	if u == nil {
		return nil, fmt.Errorf("Not Logged in")
	}
	return u, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
	basePageTemplate, err := template.New("basePagetemplate").ParseFiles("templates/base.html")
	if err != nil {
		http.Error(w, "Yeah!"+err.Error(), http.StatusInternalServerError)
	}
	err = basePageTemplate.ExecuteTemplate(w, "base.html", "")
	if err != nil {
		http.Error(w, "Eooh!"+err.Error(), http.StatusInternalServerError)
	}
}
