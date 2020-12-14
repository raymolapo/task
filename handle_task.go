package task

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// IndexHandler responds to requests with our greeting.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello, World!")
}

// Handler processes task requests.
func Handler(w http.ResponseWriter, r *http.Request) {
	taskName := r.Header.Get("X-Appengine-Taskname")
	if taskName == "" {
		// You may use the presence of the X-Appengine-Taskname header to validate
		// the request comes from Cloud Tasks.
		log.Println("Invalid Task: No X-Appengine-Taskname request header found")
		http.Error(w, "Bad Request - Invalid Task", http.StatusBadRequest)
		return
	}

	// Pull useful headers from Task request.
	queueName := r.Header.Get("X-Appengine-Queuename")

	// Extract the request body for further task details.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ReadAll: %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	// Log & output details of the task.
	output := fmt.Sprintf("Completed task: task queue(%s), task name(%s), payload(%s)",
		queueName,
		taskName,
		string(body),
	)
	log.Println(output)

	// Set a non-2xx status code to indicate a failure in task processing that should be retried.
	// For example, http.Error(w, "Internal Server Error: Task Processing", http.StatusInternalServerError)
	fmt.Fprintln(w, output)
}
