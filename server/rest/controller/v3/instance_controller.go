/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package v3

import (
	"net/http"

	"github.com/apache/servicecomb-service-center/pkg/rest"
	v4 "github.com/apache/servicecomb-service-center/server/resource/disco"
)

type MicroServiceInstanceService struct {
	v4.InstanceResource
}

func (this *MicroServiceInstanceService) URLPatterns() []rest.Route {
	return []rest.Route{
		{http.MethodGet, "/registry/v3/instances", this.FindInstances},
		{http.MethodGet, "/registry/v3/microservices/:serviceId/instances", this.ListInstance},
		{http.MethodGet, "/registry/v3/microservices/:serviceId/instances/:instanceId", this.GetInstance},
		{http.MethodPost, "/registry/v3/microservices/:serviceId/instances", this.RegisterInstance},
		{http.MethodDelete, "/registry/v3/microservices/:serviceId/instances/:instanceId", this.UnregisterInstance},
		{http.MethodPut, "/registry/v3/microservices/:serviceId/instances/:instanceId/properties", this.PutInstanceProperties},
		{http.MethodPut, "/registry/v3/microservices/:serviceId/instances/:instanceId/status", this.PutInstanceStatus},
		{http.MethodPut, "/registry/v3/microservices/:serviceId/instances/:instanceId/heartbeat", this.SendHeartbeat},
		{http.MethodPut, "/registry/v3/heartbeats", this.SendManyHeartbeat},
	}
}
