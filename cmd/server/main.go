// sessions.go
package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/sessions"
    "text/template"
    "time"
    "github.com/computermouth/cookies/pkg/dynamic"
    "encoding/json"
)

var (
    // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
    key = []byte("super-secret-kei")
    store = sessions.NewCookieStore(key)
)

var (
	db = []dynamic.Entry{
		{
			Username: "ben",
			Password: "benpw",
			Projects: []dynamic.Project{
				{
					Id: 0,
					Name: "robocar",
					Percent: 0,
					Status: dynamic.StatCode(dynamic.SUCCEEDED).String(),
				},
				{
					Id: 1,
					Name: "kiosk",
					Percent: 0,
					Status: dynamic.StatCode(dynamic.FAILED).String(),
				},
				{
					Id: 2,
					Name: "smart tv",
					Percent: 0,
					Status: dynamic.StatCode(dynamic.PENDING).String(),
				},
				{
					Id: 3,
					Name: "gif camera",
					Percent: 10,
					Status: dynamic.StatCode(dynamic.BUILDING).String(),
				},
				{
					Id: 4,
					Name: "smart speaker",
					Percent: 47,
					Status: dynamic.StatCode(dynamic.BUILDING).String(),
				},
			},
		},
		{
			Username: "lucas",
			Password: "lucaspw",
			Projects: []dynamic.Project{
				{
					Id: 0,
					Name: "hoverboard",
					Percent: 10,
					Status: dynamic.StatCode(dynamic.BUILDING).String(),
				},
				{
					Id: 1,
					Name: "quadcopter",
					Percent: 0,
					Status: dynamic.StatCode(dynamic.PENDING).String(),
				},
			},
		},
	}
)

var (
	homestaticheader = `
<html>
<head/>
<body>
<h1>{{.Username}}<h1>
<hr>
<div class="dynamic">
`

	homestaticfooter = `
</div>
<hr>
<a href=/logout>Log out</a>
<script src="/home.js" type="text/javascript"></script>
</body>
</html>
`
)

func home(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")
    
    errormsg := "E: either you need to log in, or something went wrong"

    // Check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    // Print table
    for i := 0; i < len(db); i++ {
		if db[i].Username == session.Values["username"] {
			t := template.Must(template.New("homestaticheader").Parse(homestaticheader))
			err := t.Execute(w, db[i])
			if err != nil {
				fmt.Println(err)
				fmt.Fprintln(w, errormsg)
				break
			}
			t = template.Must(template.New("homedynamic").Parse(dynamic.HomeBody))
			err = t.Execute(w, db[i].Projects)
			if err != nil {
				fmt.Println(err)
				fmt.Fprintln(w, errormsg)
				break
			}
			fmt.Fprintln(w, homestaticfooter)
			return
		}
	}
	
    fmt.Fprintln(w, errormsg)
}

func projects(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")
    
    errormsg := "E: either you need to log in, or something went wrong"

    // Check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    // Print table
    for i := 0; i < len(db); i++ {
		if db[i].Username == session.Values["username"] {
			p, _ := json.Marshal(db[i].Projects)
			fmt.Fprintln(w, string(p))
			return
		}
	}
	
    fmt.Fprintln(w, errormsg)
}

func login(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Authentication goes here
    for i := 0; i < len(db); i++ {
		if ( db[i].Username == r.FormValue("username") &&
			 db[i].Password == r.FormValue("password")) {
				// Set user as authenticated
				session.Values["authenticated"] = true
				session.Values["username"] = db[i].Username 
				session.Save(r, w)
		}
	}
	
	if session.Values["authenticated"] == true {
		//~ fmt.Fprintln(w, "<html>Success!  <a href=\"/secret\">See your secret</a></html>")
		http.Redirect(w, r, "/home", 301)
	} else {
		fmt.Fprintln(w, "<html>Failure! Username or password is incorrect <a href=\"/\">Try logging in</a></html>")
	}
	
}

func logout(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Revoke users authentication
    session.Values["authenticated"] = false
    session.Values["username"] = ""
    session.Save(r, w)
    http.Redirect(w, r, "/", 301)
}

var (
	roothtml = `
<html>
<h1>Login</h1>
<form action="/login" method="POST">
    <label>Username:</label><br />
    <input type="text" name="username"><br />
    <label>Password:</label><br />
    <input type="password" name="password"><br />
    <input type="submit">
</form>
</html>	
`
)

func root(w http.ResponseWriter, r *http.Request) {
	
	fmt.Fprintln(w, roothtml)
	
}

func main() {
    http.HandleFunc("/home", home)
    http.HandleFunc("/login", login)
    http.HandleFunc("/logout", logout)
    http.HandleFunc("/projects", projects)
    http.HandleFunc("/home.js", func(w http.ResponseWriter, r *http.Request){ http.ServeFile (w, r, "./home.js")})
    http.HandleFunc("/", root)
    
    go func(){
		for {
			for i := 0; i < len(db); i++ {
				for j := 0; j < len(db[i].Projects); j++ {
					if db[i].Projects[j].Status == dynamic.StatCode(dynamic.BUILDING).String() {
							// increment and roll progress
							db[i].Projects[j].Percent++
							db[i].Projects[j].Percent %= 100
					}
				}
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

    http.ListenAndServe(":8080", nil)
}
