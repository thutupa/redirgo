package app

import (
	"crypto/md5"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/crhym3/go-endpoints/endpoints"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

const (
	accountKind = "Account"
	actionKind  = "Action"
)

var (
	scopes    = []string{endpoints.EmailScope}
	clientIDs = []string{clientID, endpoints.ApiExplorerClientId}
	// in case we'll want to use TicTacToe API from an Android app
	audiences = []string{clientID}
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
type ActionsListReq struct {
	Phrase string `json:"phrase"`
}

func lookupActionsForPhrase(c appengine.Context, email, phrase string) (actions []*Action, err error) {
	userKey := makeUserKey(c, email)
	q := datastore.NewQuery(actionKind).Ancestor(userKey).Order("date")
	if len(phrase) > 0 {
		for _, w := range strings.Split(phrase, " ") {
			q = q.Filter("actionwords =", w)
		}
	}
	keys, err := q.GetAll(c, &actions)
	if err != nil {
		return actions, err
	}

	for i, k := range keys {
		actions[i].Key = k
	}
	return actions, nil
}

// List returns a list of matching actions
func (as *ActionsService) List(r *http.Request, req *ActionsListReq, resp *ActionsListResp) error {
	c := endpoints.NewContext(r)
	u, err := getUser(c)
	if err != nil {
		return err
	}
	if resp.Items, err = lookupActionsForPhrase(c, u.Email, req.Phrase); err != nil {
		return err
	}
	return nil
}

// ActionAddResp is a response type of ActionsService.List method
type ActionAddResp struct{}

//Request type for ActionsService.List
type ActionAddReq struct {
	Words    string `json:"actionWords"`
	Redirect string `json:"redirectLink"`
}

func makeUserKey(c appengine.Context, userEmail string) *datastore.Key {
	userID := fmt.Sprintf("%x", md5.Sum([]byte(userEmail)))
	return datastore.NewKey(c, accountKind, userID, 0, nil)
}

// Add adds an action.
func (as *ActionsService) Add(r *http.Request, req *ActionAddReq, resp *ActionAddResp) error {
	c := endpoints.NewContext(r)
	u, err := getUser(c)
	if err != nil {
		return err
	}
	if err = ValidateRedirect(req.Redirect); err != nil {
		return err
	}
	act := &Action{
		Key:          nil,
		ActionWords:  strings.Split(req.Words, " "),
		RedirectLink: req.Redirect,
		Date:         time.Now(),
		UserID:       u.ID,
	}
	putKey := datastore.NewIncompleteKey(c, actionKind, makeUserKey(c, u.Email)) // no id, let it auto generate.
	_, err = datastore.Put(c, putKey, act)
	if err != nil {
		return err
	}
	return nil
}

// ActionEditResp is a response type of ActionsService.List method
type ActionEditResp struct{}

//Request type for ActionsService.List
type ActionEditReq struct {
	KeyString string `json:"id"`
	Words     string `json:"actionWords"`
	Redirect  string `json:"redirectLink"`
}

// Edit adds an action.
func (as *ActionsService) Edit(r *http.Request, req *ActionEditReq, resp *ActionEditResp) error {
	c := endpoints.NewContext(r)
	u, err := getUser(c)
	if err != nil {
		return err
	}
	if err = ValidateRedirect(req.Redirect); err != nil {
		return err
	}
	key, err := datastore.DecodeKey(req.KeyString)
	if err != nil {
		return err
	}
	act := &Action{
		Key:          key,
		ActionWords:  strings.Split(req.Words, " "),
		RedirectLink: req.Redirect,
		Date:         time.Now(),
		UserID:       u.ID,
	}
	_, err = datastore.Put(c, key, act)
	if err != nil {
		return err
	}
	return nil
}

func ValidateRedirect(redirect string) error {
	if u, err := url.Parse(redirect); len(u.Scheme) == 0 || err != nil {
		if err != nil {
			return fmt.Errorf("Url not well formed %v", err)
		} else {
			return fmt.Errorf("Url not well formed")
		}
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
	info.Scopes, info.ClientIds, info.Audiences = scopes, clientIDs, audiences

	add := api.MethodByName("Add").Info()
	add.Name, add.HttpMethod, add.Path, add.Desc =
		"add", "PUT", "add", "Add an action."
	add.Scopes, add.ClientIds, add.Audiences = scopes, clientIDs, audiences

	edit := api.MethodByName("Edit").Info()
	edit.Name, edit.HttpMethod, edit.Path, edit.Desc =
		"edit", "PUT", "edit", "Edit an action."
	edit.Scopes, edit.ClientIds, edit.Audiences = scopes, clientIDs, audiences

	endpoints.HandleHttp()
	http.HandleFunc("/breathe", breatheHandler)
	http.HandleFunc("/redirect", redirectHandler)
	http.HandleFunc("/", mainHandler)
}

func getUser(c endpoints.Context) (*user.User, error) {
	u, err := endpoints.CurrentUser(c, scopes, audiences, clientIDs)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("Unauthorized: Please, sign in.")
	}
	c.Debugf("Current user: %#v", u)
	return u, nil
}

func templatePath(fname string) string {
	return "templates/" + fname
}

type TemplateParams struct {
	ClientID string
}

func handler(w http.ResponseWriter, r *http.Request, templateFile string) {
	basePageTemplate, err := template.New("basePagetemplate").Delims("<<<", ">>>").ParseFiles(templatePath(templateFile))
	if err != nil {
		http.Error(w, "Yeah!"+err.Error(), http.StatusInternalServerError)
		return
	}
	err = basePageTemplate.ExecuteTemplate(w, templateFile, TemplateParams{ClientID: clientID})
	if err != nil {
		http.Error(w, "Eooh!"+err.Error(), http.StatusInternalServerError)
		return
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	handler(w, r, "base.html")
}

func breatheHandler(w http.ResponseWriter, r *http.Request) {
	handler(w, r, "breathe.html")
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, _ := user.LoginURL(c, r.URL.RequestURI())
		http.Redirect(w, r, url, http.StatusFound)
		return
	}
	phrase := r.FormValue("p")
	c.Infof("Phrase = %v", phrase)
	c.Infof("Email = %v", u.Email)
	actions, err := lookupActionsForPhrase(c, u.Email, phrase)
	if err != nil {
		http.Error(w, "Eooh!"+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(actions) == 1 {
		http.Redirect(w, r, actions[0].RedirectLink, http.StatusFound)
		return
	}
	if len(actions) == 0 {
		fmt.Fprintf(w, "Found none.")
	} else {
		// TODO(syam): Print a nice form.
		fmt.Fprintf(w, "Found too many.")
	}
}
