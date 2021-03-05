package main
import "database/sql"
import "github.com/ArcticKona/token"
import "net/http"

func where( response http.ResponseWriter, request * http.Request , database * sql.DB ) ( int64 , token.Token ) {
	var err error
	var pointer int64 = 0
	var tok token.Token

	for _ , tok = range token.Authorization( response , request ) {

		err = database.QueryRow( "SELECT pointer FROM files WHERE ident = $1 AND service = $2 AND pathname = $3 ;" , tok.Ident , tok.Service , request.URL.Path ).Scan( & pointer )

		if err != nil && err != sql.ErrNoRows {
			response.WriteHeader( http.StatusInternalServerError )
			response.( http.Flusher ).Flush( )
			return 0 , tok
		}

		break

	}

	return pointer , tok
}


