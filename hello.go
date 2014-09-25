package hello

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)
import "github.com/crhym3/go-endpoints/endpoints"

// Action is a
// It also serves as (a part of) a response of ActionsService.
type Action struct {
	Key          *datastore.Key `json:"id" datastore:"-"`
	ActionWords  []string       `json:"actionWords" datastore:"actionwords"`
	RedirectLink string         `json:"redirectLink" datastore:"redirect_link,noindex"`
	Date         time.Time      `json:"date" datastore:"date"`
}

// ActionsList is a response type of ActionsService.List method
type ActionsList struct {
	Actions []*Action `json:"items"`
}

// Request type for ActionsService.List
type ActionsListReq struct{}

// ActionsService can sign the guesbook, list all actions and delete
// a action from the guestbook.
type ActionsService struct {
}

// List responds with a list of all actions ordered by Date field.
// Most recent greets come first.
func (as *ActionsService) List(r *http.Request, req *ActionsListReq, resp *ActionsList) error {
	c := endpoints.NewContext(r)
	u, err := getUser(c)
	if err != nil {
		return err
	}
	userKey := datastore.NewKey(c, "User", u.ID, 0, nil)
	q := datastore.NewQuery("Action").Ancestor(userKey)
	var actions []*Action
	keys, err := q.GetAll(c, &actions)
	if err != nil {
		return err
	}

	for i, k := range keys {
		actions[i].Key = k
	}
	resp.Actions = actions
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
		"actions.list", "GET", "actions", "List most recent actions."

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
		http.Error(w, "Yeah!" + err.Error(), http.StatusInternalServerError)
	}
	err = basePageTemplate.ExecuteTemplate(w, "base.html", "")
	if err != nil {
		http.Error(w, "Eooh!" + err.Error(), http.StatusInternalServerError)
	}
}
