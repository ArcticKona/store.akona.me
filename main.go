package main
import "database/sql"
import "flag"
import "fmt"
import _ "github.com/lib/pq"
import "io/ioutil"
import "net/http"
import "os"

func main( ) {

	//
	// Parse arguments
	flag.Usage = func( ) {
		fmt.Fprintf( flag.CommandLine.Output( ) , "For help and usage please contact software provider.\r\nCopyright © 2020 Kona Arctic. All rights reserved. NO WARRANTY! https://akona.me mailto:arcticjieer@gmail.com\r\n" )
	}
	databaseStr := flag.String( "database" , os.Getenv( "AKONA_POSTGRESQL_GOPQ" ) , "database connection e.g. user:password@tcp(127.0.0.1:3306)/hello" )
	listenStr := flag.String( "listen" , ":" + os.Getenv( "PORT" ) , "listen here e.g. [::]:8080" )
	flag.Parse( )

	//
	// Prepare database
	database , err := sql.Open( "postgres" , * databaseStr )

	// Create table if not exist
	file , err := os.Open( "init.sql" )
	if err != nil {
		panic( err ) }
	defer file.Close( )
	buffer , err := ioutil.ReadAll( file )
	database.Exec( string( buffer ) )

	// Consistency test
	rows , err := database.Query( "SELECT pointer FROM files;" )
	if err != nil {
		panic( err ) }
	for rows.Next( ) {
		var pointer int64
		err = rows.Scan( & pointer )
		if err != nil {
			panic( err ) }
		_ , err := os.Stat( fmt.Sprintf( "%d" , pointer ) )
		if err != nil && pointer != 0 {
			fmt.Fprintf( os.Stderr , "Warning: cannot access %d: %v\r\n" , pointer , err ) }
	}

	//
	// Handle request
	http.HandleFunc( "/" , func( response http.ResponseWriter, request * http.Request ) {
		if request.URL.Path == "/" {
			response.Header( ).Set( "Location" , "https://akona.me/" )
			response.WriteHeader( http.StatusTemporaryRedirect )
			return
		}
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


