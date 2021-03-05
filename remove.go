package main
import "database/sql"
import "fmt"
import "github.com/ArcticKona/token"
import "net/http"
import "os"

//
// DELETE file
func remove( response http.ResponseWriter, request * http.Request , database * sql.DB ) {
	var err error
	var pointer int64
	var token token.Token

	// Where is it?
	pointer , token = where( response , request , database )
	if pointer == 0 {
		response.WriteHeader( http.StatusNotFound )
		return
	}

	// Delete
	_ , err = database.Exec( "DELETE FROM files WHERE pointer = $1 AND ident = $2 AND service = $3 AND pathname = $4 ;" , pointer , token.Ident , token.Service , request.URL.Path )
	if err != nil {
		response.WriteHeader( http.StatusInternalServerError )
		return
	}
	response.WriteHeader( http.StatusOK )
	os.Remove( fmt.Sprintf( "%d" , pointer ) )

	return
}


