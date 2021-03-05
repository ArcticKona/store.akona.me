package main
import "database/sql"
import "fmt"
import "github.com/ArcticKona/token"
import "io"
import "net/http"
import "os"

//
// PUT or POST file
func put( response http.ResponseWriter, request * http.Request , database * sql.DB ) {
	var err error
	var pointer int64
	var token token.Token

	// Does it already exist?
	pointer , token = where( response , request , database )

	// Currently, POST does not create files
	if pointer == 0 && request.Method == http.MethodPost {
		response.WriteHeader( http.StatusNotFound )
		return
	}

	// Do I have enough space? TODO
//	var diskusage uint64 = 0
//	

	// Create file
	if pointer == 0 {
		err = database.QueryRow( "INSERT INTO files( pathname , ident , service ) VALUES ( $1 , $2 , $3 ) RETURNING pointer;" , request.URL.Path , token.Ident , token.Service ).Scan( & pointer )
		if err != nil {
			response.WriteHeader( http.StatusInternalServerError )
			return
		}

	}

	// Write to file
	file , err := os.Create( fmt.Sprintf( "%d" , pointer ) )
	_ , err = io.Copy( file , request.Body )
	if err != nil {	// Delete on error
		response.WriteHeader( http.StatusInternalServerError )
		database.Exec( "DELETE FROM files WHERE pointer = $1 AND ident = $2 AND service = $3 AND pathname = $4 ;" , pointer , token.Ident , token.Service , request.URL.Path )
		os.Remove( fmt.Sprintf( "%d" , pointer ) )
		return
	}
	file.Close( )
	response.WriteHeader( http.StatusOK )

	return
}


