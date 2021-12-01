package routing

import (
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// GET '/'
func (RM *RoutingMemory) GET_0_(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")

	w.Write([]byte("Old Digibox"))
	return
}

// GET '/generate_pwd/'
func (RM *RoutingMemory) GET_0_generate_pwd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if !strings.HasPrefix(r.Host, "127.0.0.1") {
		w.WriteHeader(403)
		w.Write([]byte("Forbidden."))
		return
	}

	if r.FormValue("pwd") == "" {
		w.Write([]byte("NO password."))
		return
	}

	pw, err := HashPassword(r.FormValue("pwd"))

	if err == nil {
		w.Write([]byte(pw))
	} else {
		w.Write([]byte(err.Error()))
	}

	return

}

// GET '/auth'
func (RM *RoutingMemory) GET_0_auth(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.FormValue("name") == "" {
		// no username.
		w.WriteHeader(404)

		w.Write([]byte("No NAME provided."))
		return
	}

	if r.FormValue("psk") == "" {
		// no psk?
		w.WriteHeader(404)

		w.Write([]byte("No PSK provided."))
		return
	}

	user, e := RM.database.Getuser(r.FormValue("name"))

	if !e {
		// user not exists
		w.WriteHeader(404)

		w.Write([]byte("User does not exist."))
		return
	}

	if CheckPasswordHash(r.FormValue("psk"), user.PassHash) {
		// password correct.

		//w.Write([]byte("Streaming permission granted."))

		if r.FormValue("call") == "play" {
			if user.PermLevel == 1 || user.PermLevel == 3 {
				w.WriteHeader(201)
				w.Write([]byte("Play permission granted."))
			} else {
				w.WriteHeader(404)
				w.Write([]byte("Play permission denied."))
			}
		} else if r.FormValue("call") == "publish" {
			if user.PermLevel == 2 || user.PermLevel == 3 {
				w.WriteHeader(201)
				w.Write([]byte("Publish permission granted."))
			} else {
				w.WriteHeader(404)
				w.Write([]byte("Publish permission denied."))
			}
		}
	} else {
		// pwd incorrect
		w.WriteHeader(404)

		w.Write([]byte("Streaming permission denied (psk)."))
	}

	return
}
