# EAAS

## GRPC

The GRPC server listens on port 50010.

`GetTracks`, `SearchTracks` and `EventStream` are server-streaming. A listing
arrives as one response per batch. `EventStream` stays open until the client
disconnects; a server that returns from it straight away makes the device
reconnect in a loop.

A unary server still answers a streaming client correctly as long as it sends
exactly one response, which is why the single-track demo library worked before.

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

### Performance data blobs

`TrackPerformanceData.beat_grid`, `overview_waveform` and
`overview_waveform_rgb` are framed the same way in the on-disk Engine library:
a `uint32` big-endian length of the uncompressed payload, then a zlib stream.

The uncompressed `beat_grid` payload is two big-endian `double`s (sample rate,
sample count), a `uint8` flag, then two beat grids. Each grid is a big-endian
`int64` count followed by 24-byte markers: little-endian `double` sample
offset, little-endian `int64` beat number, little-endian `int32` beats in the
run, little-endian `int32` unused. Everything after the second grid is copied
through verbatim.

An overview waveform payload is two big-endian `int64`s (both 1024), a
big-endian `double` of samples per point (sample count / 1024), then 1024
triples. In `overview_waveform` each triple holds low, mid and high band
levels; in `overview_waveform_rgb` the same three bands are stored as the
rendered color, which is why the two payloads for one track never match byte
for byte. engineOS keeps the rgb payload in a file rather than in its library
database, at
`Engine Library/Database2/OverviewData/<library uuid>/<track id>.rgb`.

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
