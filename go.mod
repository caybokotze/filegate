module github.com/filegate/filegate

go 1.24.1

require (
	github.com/abema/go-mp4 v1.4.1 // indirect
	github.com/alecthomas/kong v1.13.0 // indirect
	github.com/anacrolix/dms v1.7.2 // indirect
	github.com/anacrolix/ffprobe v1.1.0 // indirect
	github.com/anacrolix/generics v0.0.1 // indirect
	github.com/anacrolix/log v0.15.2 // indirect
	github.com/dhowden/tag v0.0.0-20240417053706-3d75831295e8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	golang.org/x/exp v0.0.0-20240613232115-7f521ea00fb8 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
)

// Use local ffprobe fork to suppress "not found" warnings
replace github.com/anacrolix/ffprobe => ./internal/ffprobe
