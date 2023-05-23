// Code generated by counterfeiter. DO NOT EDIT.
package specsfakes

import (
	"net/http"
	"sync"

	"k8s.io/release/pkg/obs/specs"
	"k8s.io/release/pkg/release"
)

type FakeImpl struct {
	GetKubeVersionStub        func(release.VersionType) (string, error)
	getKubeVersionMutex       sync.RWMutex
	getKubeVersionArgsForCall []struct {
		arg1 release.VersionType
	}
	getKubeVersionReturns struct {
		result1 string
		result2 error
	}
	getKubeVersionReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	GetRequestStub        func(string) (*http.Response, error)
	getRequestMutex       sync.RWMutex
	getRequestArgsForCall []struct {
		arg1 string
	}
	getRequestReturns struct {
		result1 *http.Response
		result2 error
	}
	getRequestReturnsOnCall map[int]struct {
		result1 *http.Response
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeImpl) GetKubeVersion(arg1 release.VersionType) (string, error) {
	fake.getKubeVersionMutex.Lock()
	ret, specificReturn := fake.getKubeVersionReturnsOnCall[len(fake.getKubeVersionArgsForCall)]
	fake.getKubeVersionArgsForCall = append(fake.getKubeVersionArgsForCall, struct {
		arg1 release.VersionType
	}{arg1})
	stub := fake.GetKubeVersionStub
	fakeReturns := fake.getKubeVersionReturns
	fake.recordInvocation("GetKubeVersion", []interface{}{arg1})
	fake.getKubeVersionMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeImpl) GetKubeVersionCallCount() int {
	fake.getKubeVersionMutex.RLock()
	defer fake.getKubeVersionMutex.RUnlock()
	return len(fake.getKubeVersionArgsForCall)
}

func (fake *FakeImpl) GetKubeVersionCalls(stub func(release.VersionType) (string, error)) {
	fake.getKubeVersionMutex.Lock()
	defer fake.getKubeVersionMutex.Unlock()
	fake.GetKubeVersionStub = stub
}

func (fake *FakeImpl) GetKubeVersionArgsForCall(i int) release.VersionType {
	fake.getKubeVersionMutex.RLock()
	defer fake.getKubeVersionMutex.RUnlock()
	argsForCall := fake.getKubeVersionArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeImpl) GetKubeVersionReturns(result1 string, result2 error) {
	fake.getKubeVersionMutex.Lock()
	defer fake.getKubeVersionMutex.Unlock()
	fake.GetKubeVersionStub = nil
	fake.getKubeVersionReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeImpl) GetKubeVersionReturnsOnCall(i int, result1 string, result2 error) {
	fake.getKubeVersionMutex.Lock()
	defer fake.getKubeVersionMutex.Unlock()
	fake.GetKubeVersionStub = nil
	if fake.getKubeVersionReturnsOnCall == nil {
		fake.getKubeVersionReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getKubeVersionReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeImpl) GetRequest(arg1 string) (*http.Response, error) {
	fake.getRequestMutex.Lock()
	ret, specificReturn := fake.getRequestReturnsOnCall[len(fake.getRequestArgsForCall)]
	fake.getRequestArgsForCall = append(fake.getRequestArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetRequestStub
	fakeReturns := fake.getRequestReturns
	fake.recordInvocation("GetRequest", []interface{}{arg1})
	fake.getRequestMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeImpl) GetRequestCallCount() int {
	fake.getRequestMutex.RLock()
	defer fake.getRequestMutex.RUnlock()
	return len(fake.getRequestArgsForCall)
}

func (fake *FakeImpl) GetRequestCalls(stub func(string) (*http.Response, error)) {
	fake.getRequestMutex.Lock()
	defer fake.getRequestMutex.Unlock()
	fake.GetRequestStub = stub
}

func (fake *FakeImpl) GetRequestArgsForCall(i int) string {
	fake.getRequestMutex.RLock()
	defer fake.getRequestMutex.RUnlock()
	argsForCall := fake.getRequestArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeImpl) GetRequestReturns(result1 *http.Response, result2 error) {
	fake.getRequestMutex.Lock()
	defer fake.getRequestMutex.Unlock()
	fake.GetRequestStub = nil
	fake.getRequestReturns = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeImpl) GetRequestReturnsOnCall(i int, result1 *http.Response, result2 error) {
	fake.getRequestMutex.Lock()
	defer fake.getRequestMutex.Unlock()
	fake.GetRequestStub = nil
	if fake.getRequestReturnsOnCall == nil {
		fake.getRequestReturnsOnCall = make(map[int]struct {
			result1 *http.Response
			result2 error
		})
	}
	fake.getRequestReturnsOnCall[i] = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeImpl) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getKubeVersionMutex.RLock()
	defer fake.getKubeVersionMutex.RUnlock()
	fake.getRequestMutex.RLock()
	defer fake.getRequestMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeImpl) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ specs.Impl = new(FakeImpl)
