package main
import "crypto/rand"
import "database/sql"
import "encoding/base64"
import "github.com/ArcticKona/token"
import "net/http"


//
// SHARE file
func share( response http.ResponseWriter, request * http.Request , database * sql.DB ) {
	var err error
	var pointer int64
	var token token.Token
	var random [ 64 ]byte

	// Where is it?
	pointer , token = where( response , request , database )
	if pointer == 0 {
		response.WriteHeader( http.StatusNotFound )
		return
	}

	// Crypto random
	_ , err = rand.Read( random[ : ] )
	if err != nil {
		response.WriteHeader( http.StatusInternalServerError )
		return
	}

	// Add link	FIXME: Check if share link already exists
	_ , err = database.Exec( "UPDATE files SET sharelink = $1 WHERE pointer = $2 AND ident = $3 AND service = $4 AND pathname = $5 ;" , base64.RawURLEncoding.EncodeToString( random[ : ] ) , pointer , token.Ident , token.Service , request.URL.Path )
	if err != nil {
		response.WriteHeader( http.StatusInternalServerError )
		return
	}
	response.WriteHeader( http.StatusOK )
	response.Write( [ ]byte( base64.RawURLEncoding.EncodeToString( random[ : ] ) ) )

	return
}


