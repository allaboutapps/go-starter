package api

import "allaboutapps.at/aw/go-mranftl-sample/pgtestpool"

type Server struct {
	M      *pgtestpool.Manager
	Config ServerConfig
}
