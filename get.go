package main
import "database/sql"
import "fmt"
import "io"
import "net/http"
import "os"

//
// GET file
func get( response http.ResponseWriter, request * http.Request , database * sql.DB ) {
	var err error
	var pointer int64

	// Is it shared?
	err = database.QueryRow( "SELECT pointer FROM files WHERE sharelink = $1 ;" , request.URL.Path[ 1 : ] ).Scan( & pointer )
	if err != nil && err != sql.ErrNoRows {
		response.WriteHeader( http.StatusInternalServerError )
		return

	// Is it in my storage?
	} else if err != nil {
		pointer , _ = where( response , request , database )
		if pointer == 0 {
			response.WriteHeader( http.StatusNotFound )
			return
		}
	}

	// Check headers
	var mimetype sql.NullString
	err = database.QueryRow( "SELECT mimetype FROM files WHERE pointer = $1 ;" , pointer ).Scan( & mimetype )
	if err != nil {
		response.WriteHeader( http.StatusInternalServerError )
		return
	}
	if mimetype.Valid {
		response.Header( ).Add( "Content-Type" , mimetype.String ) }
	response.Header( ).Add( "Content-Security-Policy" , "default-src 'none', sandbox;" )	// This is not a hosting provider

	// Get file if possible
	file , err := os.Open( fmt.Sprintf( "%d" , pointer ) )
	if err != nil {
		response.WriteHeader( http.StatusInternalServerError )
		return
	}
	response.WriteHeader( http.StatusOK )
	io.Copy( response , file )
	file.Close( )

	return
}


