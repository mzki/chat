package chat

import (
	"sync"

	"github.com/ipfans/echo-session"
	"github.com/labstack/echo"
	"github.com/mzki/chat/entity"
	"github.com/mzki/chat/entity/stub"
)

type UserForm struct {
	UserID   string `json:"user_id" form:"user_id" query:"user_id"`
	Password string `json:"password" form:"password" query:"password"`
}

// key for session value which is user loggedin information.
const KeyUserTableID = "user.id"

// LoginHandler handles login requests.
// it holds logged-in users, so that each request can reference
// any logged-in user.
type LoginHandler struct {
	userRepo entity.UserRepository

	mu            *sync.RWMutex
	loggedinUsers map[int64]entity.User
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{
		userRepo:      stub.UserRepository{},
		mu:            new(sync.RWMutex),
		loggedinUsers: make(map[int64]entity.User),
	}
}

func (lh *LoginHandler) Login(c echo.Context) error {
	u := new(UserForm)
	if err := c.Bind(u); err != nil {
		return err
	}

	user, err := lh.userRepo.Get(u.UserID, u.Password)
	if err != nil {
		// TODO show error message as html.
		return err
	}

	// login succeed, save it into session and redirect to next page.
	lh.mu.Lock()
	lh.loggedinUsers[user.ID] = user
	lh.mu.Unlock()

	sess := session.Default(c)
	sess.Set(KeyUserTableID, user.ID)
	sess.Save()

	// TODO c.Redirect(code, page)
	return nil
}

func (lh *LoginHandler) Logout(c echo.Context) error {
	sess := session.Default(c)
	id, ok := sess.Get(KeyUserTableID).(int64)
	if !ok {
		return nil
	}
	sess.Delete(KeyUserTableID)
	sess.Save()

	lh.mu.Lock()
	delete(lh.loggedinUsers, id)
	lh.mu.Unlock()
	return nil
}

func (lh *LoginHandler) LoginPage(c echo.Context) error {
	// TODO return html using template
	return nil
}

func (lh *LoginHandler) LogoutPage(c echo.Context) error {
	// TODO return html using template
	return nil
}

func (lh *LoginHandler) IsLoggedInRequest(c echo.Context) bool {
	id, ok := session.Default(c).Get(KeyUserTableID).(int64)
	if !ok {
		return false
	}
	// here, user exactlly logged in,
	// addionally we assert existance in loggedinUsers map.
	if !lh.IsLoggedInUser(id) {
		panic("session loggedin but loggedin user map is not.")
	}
	return true
}

func (lh *LoginHandler) IsLoggedInUser(id int64) bool {
	lh.mu.RLock()
	defer lh.mu.RUnlock()
	_, loggedin := lh.loggedinUsers[id]
	return loggedin
}

func (lh *LoginHandler) LoggedinFilter(c echo.Context) error {
	// TODO login filter using middleware
	return nil
}