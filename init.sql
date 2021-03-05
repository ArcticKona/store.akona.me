BEGIN;
--DROP TABLE files;
CREATE TABLE files( pointer BIGSERIAL PRIMARY KEY , pathname VARCHAR NOT NULL , ident BIGINT NOT NULL , service INT NOT NULL , sharelink VARCHAR UNIQUE DEFAULT NULL , mimetype VARCHAR DEFAULT NULL , modifytime TIMESTAMP DEFAULT CURRENT_TIMESTAMP , filesize BIGINT , alive BOOL DEFAULT TRUE );
INSERT INTO files( pathname , ident , service , sharelink ) VALUES ( '/' , 0 , 0 , '/test0' );
CREATE INDEX ON files( pathname );
CREATE INDEX ON files( ident );
CREATE INDEX ON files( sharelink );
COMMIT;
