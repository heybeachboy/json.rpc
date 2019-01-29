package json_rpc

import (
	"fmt"
	"reflect"
	"strings"
)

/**
 * json rpc service implements
 */

type JsonRpcService struct {
	 services  ServiceMap
	 methods   MethodMap
	 callbacks CallbackMap
}

/**
 * register service
 */
func (j *JsonRpcService) RegisterService(serviceName string, service interface{}) {

	if j.services == nil {
		j.services = make(ServiceMap)
	}

	if j.judgeServiceIsExist(serviceName) {
		fmt.Errorf("service name already register : %s", serviceName)
		return
	}

	var svr Service
	svr.ServerName = serviceName
	svr.ServiceTyp = reflect.TypeOf(service)
	svr.ServiceVal = reflect.ValueOf(service)
	resp := j.reflectCallback(svr)

	if !resp {
		fmt.Errorf("register service  %s : method is null", svr.ServerName)
		return
	}
	j.services[svr.ServerName] = &svr

}

//func(j *JsonRpcService)

/**
 *reflect service callback function
 */

func (j *JsonRpcService) reflectCallback(service Service) (bool) {

	count := 0
	for i := 0; i < service.ServiceTyp.NumMethod(); i++ {
		var c Callback
		method := service.ServiceTyp.Method(i)
		mTyp := method.Type
		name := FormatName(method.Name)
		c.ServiceName = service.ServerName
		c.Method = method
		c.MethodTyp = mTyp
		c.MethodName = service.ServerName + `.` + name
		if j.judgeMethodIsExist(c.MethodName) {
			fmt.Errorf("service of method already register : %s", c.MethodName)
			continue
		}
		count++

		if j.callbacks == nil {
			j.callbacks = make(CallbackMap)
		}
		j.callbacks[c.MethodName] = &c
	}
	if count > 0 {
		return true
	}
	return false

}

/**
 *handle json rpc request
 */

func (j *JsonRpcService) ServerHandleRequest(json JsonRpcIf) {
	defer json.Destroy()
	req, err := json.ReadJsonRpcRequestHeaders()

	if err != nil || len(req) < 1 {
		json.WriteJsonRpcResponse(json.CreateExceptionResponse(123, -32700, err))
		return
	}

	json.WriteJsonRpcResponse(req[0])

}

/**
 *


func (j *JsonRpcService) checkSingleRequest(reqList []JsonRpcRequest) (bool, string) {

	for _, req := range reqList {

		if !j.checkServiceAndMethod(req.Method) {

		}

	}

	return true, ""

}* /

func (j *JsonRpcService) parseServiceAndMethod(name string)(status bool,server string,method string) {
	list :=strings.Split(name,`_`)

	if len(list)  < 2 {
		return false,name,list[0]
	}

	if len(list) == 2 {
		return  true,list[0],list[1]
	}

	return  false,name,list[0]


}

func (j *JsonRpcService) checkService(serviceName string) (bool) {

	return true

}

func (j *JsonRpcService) checkMethod(methodName string) (bool) {
	return true
}

func (j *JsonRpcService) checkInParameter() (bool) {
	return true

}

/**
 *
 */
func (j *JsonRpcService) judgeServiceIsExist(serviceName string) (bool) {
	if _, ok := j.services[serviceName]; ok {
		return true
	}
	return false
}

func (j *JsonRpcService) judgeMethodIsExist(callbackKey string) (bool) {

	if _, ok := j.callbacks[callbackKey]; ok {
		return true
	}
	return false
}

/**
 * format string name handle
 */
func FormatName(name string) (string) {

	return strings.ToLower(strings.Trim(name, " "))

}
