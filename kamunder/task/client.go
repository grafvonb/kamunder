package task

import (
	pdsvc "github.com/grafvonb/camunder/internal/services/processdefinition"
	pisvc "github.com/grafvonb/camunder/internal/services/processinstance"
)

type API interface{}

type client struct {
	pdApi pdsvc.API
	piApi pisvc.API
}

func New(pdApi pdsvc.API, piApi pisvc.API) API {
	return &client{
		pdApi: pdApi,
		piApi: piApi,
	}
}
