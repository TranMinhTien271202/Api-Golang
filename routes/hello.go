package routes

import (
	"Test/controllers" // Sử dụng tên module của bạn
	"net/http"
)

// Hàm xử lý phương thức cho các route
func handleMethod(rMethod string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == rMethod {
			handlerFunc(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// Hàm giúp nhóm các route theo prefix
func groupRoutes(prefix string, routes map[string]http.HandlerFunc) {
	for route, handler := range routes {
		http.HandleFunc(prefix+route, handler)
	}
}

// Định nghĩa các route
func SetupRoutes() {
	apiRoutes := map[string]http.HandlerFunc{
		"/users":            handleMethod("GET", controllers.GetAllUsersHandler),
		"/post/users":    handleMethod("POST", controllers.InsertUserHandler),
		"/update/users": handleMethod("PUT", controllers.UpdateUserHandler),
		"/find/users":   handleMethod("GET", controllers.FindUserHandler),
		"/delete/users": handleMethod("DELETE", controllers.DeleteUserHandler),
	}
	groupRoutes("/api", apiRoutes) // Nhóm các route theo prefix /api

}
