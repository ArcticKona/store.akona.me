package main
import "database/sql"
import "flag"
import "fmt"
import _ "github.com/lib/pq"
import "net/http"
import "os"

func main( ) {

	//
	// Parse arguments
	flag.Usage = func( ) {
		fmt.Fprintf( flag.CommandLine.Output( ) , "For help and usage please contact software provider.\r\nCopyright © 2020 Kona Arctic. All rights reserved. NO WARRANTY! https://akona.me mailto:arcticjieer@gmail.com\r\n" )
	}
	databaseStr := flag.String( "database" , "host=localhost port=5432 user=store.akona.me password=3dbaace1e81f9ac69ef4d86c5c030c5b dbname=store.akona.me sslmode=disable" , "database connection e.g. user:password@tcp(127.0.0.1:3306)/hello" )
	listenStr := flag.String( "listen" , ":8080" , "listen here e.g. [::]:8080" )
	flag.Parse( )

	//
	// Prepare database
	database , err := sql.Open( "postgres" , * databaseStr )
	_ , err = database.Exec( "SELECT * FROM files;" )
	if ( err != nil ) {
		panic( err ) }

	//
	// Consistancy test
	// TODO


	//
	// Handle request
	http.HandleFunc( "/" , func( response http.ResponseWriter, request * http.Request ) {

		if request.Header.Get( "Origin" ) != "" {
			response.Header( ).Set( "Access-Control-Allow-Origin" , request.Header.Get( "Origin" ) )
			response.Header( ).Set( "Vary" , "Origin" )
		}

		if request.Method == http.MethodGet {
			get( response , request , database )
		} else if request.Method == http.MethodPut || request.Method == http.MethodPost {
			put( response , request , database )

		} else if request.Method == "LINK" || request.Method == "SHARE" {
			share( response , request , database )

		} else if request.Method == http.MethodDelete {
			remove( response , request , database )

		} else {
			response.WriteHeader( http.StatusNotImplemented ) }

		return
	} )

	//
	// Serve
	fmt.Printf( "listening on %s\r\n" , * listenStr ) 
	err = http.ListenAndServe( * listenStr , nil )
	if err != nil {
		panic( err ) }

	os.Exit( 1 )
}


