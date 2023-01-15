# environment

1. flag -a - SERVER_ADDRESS - server startup address (by default - 127.0.0.1:8080)
2. flag -b - BASE_URL - base address of the shortened URL (by default - http://127.0.0.1:8080 )
3. flag -d - DATABASE_DSN - DB connection address (Postgres) (by default - "", i.e. we are working without DB)
4. flag -f - FILE_STORAGE_PATH - file for storing shortened URLs (by default - "", i.e. we work without a file)

# storage
The following can be used as storage:
- in-memory storage (RAM),
- .txt file
- database server

# endpoints
The GET /api/user/urls endpoint reads the UserID from the request cookie and outputs all URLs saved by this user in the format of an array of JSON structures {"short_url":"<some_shorten_url>","original_url":"<some_original_url>"}

The GET /{id} endpoint takes the shortened URL identifier as a parameter and returns a response with the status 307 and the original URL in the HTTP Location header

The GET /ping endpoint checks the availability of the database, issues a response with the status 200 if available, and 500 if not available

The endpoint POST / accepts the URL string for shortening in the request body as text and returns a response with the code 201 and the shortened URL as a text string in the package body

The POST /api/shorten endpoint is similar to the previous one, but accepts a JSON object in the request body {"url":"<original_url>"} and returns a JSON object in the response body {"result":"<shorten_url>"}

The endpoint POST /api/shorten/batch, accepts in the request body a set of URLs for shortening in the format of an array of JSON structures {"correlation_id":"<some_id>","original_url":"<some_original_url>"} and returns shortened URLs in the format of an array of JSON structures {"correlation_id":"<some_id>","short_url":"<some_shorten_url>"}}

Endpoint DELETE /api/user/urls, accepts tasks to delete a list of previously generated URLs, issues a response with the status 202, and then asynchronously deletes records from the database
