module demo-embeddings-from-code

go 1.22.1

require github.com/parakeet-nest/parakeet v0.0.9

require (
	go.etcd.io/bbolt v1.3.10 // indirect
	golang.org/x/sys v0.20.0 // indirect
)

replace github.com/parakeet-nest/parakeet => ../..
