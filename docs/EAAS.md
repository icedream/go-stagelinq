# EAAS

## GRPC

The GRPC server listens on port 50010.

### `networktrust.v1.NetworkTrustService/CreateTrust`

#### Request parameters

- `string` - UUID of the device requesting trust
- `string` - Name of the device requesting trust

#### Response

GRPC status is always 0. Status can be one of the values for Granted, Denied or Busy.

### `enginelibrary.v1.EngineLibraryService/GetLibrary`
### `enginelibrary.v1.EngineLibraryService/EventStream`
### `enginelibrary.v1.EngineLibraryService/GetTracks`
### `enginelibrary.v1.EngineLibraryService/GetTrack`
### `enginelibrary.v1.EngineLibraryService/PutEvents`
### `enginelibrary.v1.EngineLibraryService/SearchTracks`
### `enginelibrary.v1.EngineLibraryService/GetSearchFilters`
### `enginelibrary.v1.EngineLibraryService/GetLibraries`
### `enginelibrary.v1.EngineLibraryService/GetHistorySessions`
### `enginelibrary.v1.EngineLibraryService/GetHistoryPlayedTracks`
### `enginelibrary.v1.EngineLibraryService/GetCredentials`



## HTTP

The HTTP server listens on port 50020.

### `GET /ping`

Returns an empty `200 OK` response to ensure the client the server is still
reachable.

### `GET /download/:path`

Provide a given audio file.

#### Request parameters

- `:path` - Full path to audio file to deliver. Normally URL-encoded.

#### Response

The response is the audio file itself.

#### Remarks

Server should return proper `Content-length`, however `Content-type` does not
seem to matter as it can be set to `application/octet-stream`.
