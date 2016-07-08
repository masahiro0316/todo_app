package todo

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

type Task struct {
	Author   string
	Content  string
	Priority int
	Date     time.Time
}

/* make new key for todo task */
func todoKey(ctx appengine.Context) *datastore.Key {
	return datastore.NewKey(ctx, "Task", "defaulty_todo", 0, nil)
}

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/add", add)
	// http.HandleFunc("/delete", delete)
}

func root(w http.ResponseWriter, r *http.Request) {
	/* get Context */
	ctx := appengine.NewContext(r)
	/* get pointer to all existed queries */
	qr := datastore.NewQuery("Task").Ancestor(todoKey(ctx)).Order("Date")
	/* Get the number of Queries */
	cnt, err := qr.Count(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/* get Tasks */
	tasks := make([]Task, 0, cnt)
	if _, err := qr.GetAll(ctx, &tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/* Output html */
	if err := todoTemplate.Execute(w, tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/* define html template */
var todoTemplate = template.Must(template.ParseFiles(filepath.Join("templates", "todo.html")))

func add(w http.ResponseWriter, r *http.Request) {
	/* get Context */
	ctx := appengine.NewContext(r)
	/* set New Task object */
	tsk := Task{
		Content: r.FormValue("content"),
		Date:    time.Now(),
	}

	if usr := user.Current(ctx); usr != nil {
		tsk.Author = usr.String()
	}
	/* create new key */
	key := datastore.NewIncompleteKey(ctx, "Task", todoKey(ctx))
	_, err := datastore.Put(ctx, key, &tsk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

/*
func delete(w http.ResponseWriter, r *http.Request) {
}
*/
