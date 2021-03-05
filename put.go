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

	// Create file
	if pointer == 0 {
		_ , err = database.Exec( "INSERT INTO files( pathname , ident , service ) VALUES ( $1 , $2 , $3 ) ;" , request.URL.Path , token.Ident , token.Service )
		if err != nil {
			response.WriteHeader( http.StatusInternalServerError )
			return
		}
		err = database.QueryRow( "SELECT pointer FROM files WHERE ident = $1 AND service = $2 AND pathname = $3 ;" , token.Ident , token.Service , request.URL.Path ).Scan( & pointer )
		if err != nil {
			response.WriteHeader( http.StatusInternalServerError )
			return
		}

	}

	// Add content type
	if request.Header.Get( "Content-Type" ) != "" {
		database.Exec( "UPDATE files SET mimetype = $1 WHERE pointer = $2 ;" , request.Header.Get( "Content-Type" ) , pointer ) }

	// Write to file
	file , err := os.Create( fmt.Sprintf( "%d" , pointer ) )
	defer file.Close( )
	_ , err = io.Copy( file , request.Body )
	if err != nil {
		response.WriteHeader( http.StatusInternalServerError )
		database.Exec( "DELETE FROM files WHERE pointer = $1 AND ident = $2 AND service = $3 AND pathname = $4 ;" , pointer , token.Ident , token.Service , request.URL.Path )	// Delete on error
		os.Remove( fmt.Sprintf( "%d" , pointer ) )
		return
	}

	// Do I have enough space?	FIXME Actually forbid uploads
	var usage uint64 = 0
	rows , err := database.Query( "SELECT pointer FROM files WHERE ident = $1 ;" , token.Ident )
	if err != nil {
		response.WriteHeader( http.StatusInternalServerError )
		return
	}
	for rows.Next( ) {
		err = rows.Scan( & pointer )
		if err != nil {
			response.WriteHeader( http.StatusInternalServerError )
			return
		}
		fileinfo , err := os.Stat( fmt.Sprintf( "%d" , pointer ) )
		if err == nil {
			usage += uint64( fileinfo.Size( ) ) }
	}
	fileinfo , err := file.Stat( )
	if err != nil {
		response.WriteHeader( http.StatusInternalServerError )
		return
	}
	usage -= uint64( fileinfo.Size( ) )
	if payment( token ) > 0 && usage > payment( token ) {
		response.WriteHeader( http.StatusPaymentRequired )
		return
	}

	// Done
	response.WriteHeader( http.StatusOK )
	return
}


