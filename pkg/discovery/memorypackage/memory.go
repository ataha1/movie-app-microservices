package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ataha1/movie-app/pkg/discovery"
)

type serviceName string
type instanceID string

type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[instanceID]*serviceInstance
}

type serviceInstance struct{
	hostPort string
	lastActive time.Time
}

func NewRegistry()*Registry{
	return &Registry{
		serviceAddrs: map[serviceName]map[instanceID]*serviceInstance{},		
	}
}

func (r *Registry)Register(ctx context.Context, instID string, svcName string, hostPort string) error{
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName(svcName)]; !ok {
		r.serviceAddrs[serviceName(svcName)] = map[instanceID]*serviceInstance{}
	}
	r.serviceAddrs[serviceName(svcName)][instanceID(instID)] = &serviceInstance{
		hostPort: hostPort,
		lastActive: time.Now(),
	}
	return nil
}

func (r *Registry)Deregister(ctx context.Context, instID string, svcName string) error{
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName(svcName)]; !ok{
		return nil
	}
	delete(r.serviceAddrs[serviceName(svcName)], instanceID(instID))
	return nil
}

func (r *Registry)ServiceAddresses(ctx context.Context, svcName string)([]string, error){
	r.RLock()
	defer r.RLock()
	if len(r.serviceAddrs[serviceName(svcName)]) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, i := range r.serviceAddrs[serviceName(svcName)]{
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)){
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}

func(r *Registry) ReportHealthyState(instID, svcName string) error{
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName(svcName)]; !ok {
		return errors.New("service is not registered yet")
	}
	if _, ok := r.serviceAddrs[serviceName(svcName)][instanceID(instID)]; !ok{
		return errors.New("service instance is not registered yet")
	}
	r.serviceAddrs[serviceName(svcName)][instanceID(instID)].lastActive = time.Now()
	return nil
}