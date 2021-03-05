// Basic tests
// tests require functioning database
package main
import "bytes"
import "database/sql"
import "github.com/ArcticKona/token"
import _ "github.com/lib/pq"
import "io/ioutil"
import "net/http"
import "net/http/httptest"
import "net/url"
import "testing"

func TestMain( testing * testing.T ) {
	var err error
	var response httptest.ResponseRecorder
	var request http.Request
	var database * sql.DB

	// Prepare database
	databaseStr := "host=localhost port=5432 user=store.akona.me password=3dbaace1e81f9ac69ef4d86c5c030c5b dbname=store.akona.me sslmode=disable"	// Change this
	database , _ = sql.Open( "postgres" , databaseStr )
	_ , err = database.Exec( "SELECT * FROM files;" )
	if err != nil {
		testing.Fatalf( "Database prepare failed: %v\r\n" , err ) }

	// Prepare request
	var location url.URL
	var token token.Token
	var bodybuf bytes.Buffer
	location.Scheme = "https:"
	location.Host = "store.akona.me:443"
	location.Path = "/test0"
	token.Expire = 1024 ^ 3
	request.URL = & location
	request.Proto = "HTTP/1.1"
	request.ProtoMajor = 1
	request.ProtoMinor = 1
	request.Header = map[ string ][ ]string{ "Authorization" : { token.String( ) } }
	request.ContentLength = -1
	request.Host = "store.akona.me"
	response.Body = & bodybuf

	// Upload file
	request.Body = ioutil.NopCloser( bytes.NewReader( [ ]byte( "body" ) ) )
	put( & response , & request , database )
	if response.Code != 200 {
		testing.Fatalf( "upload returned: %v\r\n" , response.Code ) }

	// Make file link
	request.Method = "LINK"
	share( & response , & request , database )
	if response.Code != 200 {
		testing.Fatalf( "share returned: %v\r\n" , response.Code ) }

	// Try to get file
	var buffer [ ]byte
	buffer , err = ioutil.ReadAll( response.Body )
	request.Method = "GET"
	request.URL.Path = "/" + string( buffer )
	get( & response , & request , database )
	if response.Code != 200 {
		testing.Fatalf( "get returned: %v\r\n" , response.Code ) }
	buffer , err = ioutil.ReadAll( response.Body )
	if string( buffer ) != "body" {
		testing.Fatalf( "get body mismatch: %v\r\n" , string( buffer ) ) }

	// Delete file
	request.Method = "DELETE"
	request.URL.Path = "/test0"
	remove( & response , & request , database )
	if response.Code != 200 {
		testing.Fatalf( "delete returned: %v\r\n" , response.Code ) }

	// Done
}


