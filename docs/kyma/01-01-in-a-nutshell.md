---
title: In a nutshell
type: Overview
---

Kyma allows you to extend applications with microservices and Functions. First, connect your application to a Kubernetes cluster and expose the application's API or events securely. Then, implement the business logic you require by creating microservices or Functions, and triggering them to react to particular events or calls to your application's API. To limit the time spent on coding, use the in-built cloud services from the Service Catalog, exposed by open service brokers from such cloud providers as GCP, Azure, and AWS.

Kyma comes equipped with these out-of-the-box functionalities:

- Service-to-service communication and proxying (Istio-based [Service Mesh](/components/service-mesh/#overview-overview))
- In-built [monitoring](/components/monitoring/#overview-overview), [tracing](/components/tracing/#overview-overview), and [logging](/components/logging/#overview-overview) (Grafana, Prometheus, Jaeger, Loki, Kiali)
- Secure [authentication and authorization](/components/security/#overview-overview) (Dex, Ory, Service Identity, TLS, Role Based Access Control)
- The catalog of services to choose from ([Service Catalog](/components/service-catalog/#overview-service-catalog), [Service Brokers](/components/service-catalog/#overview-service-brokers)
- The development platform to run lightweight Functions in a cost-efficient and scalable way ([Serverless](/components/serverless/#overview-overview))
- The endpoint to register Events and APIs of external applications ([Application Connector](/components/application-connector/#overview-overview))
- Secure API exposure ([API Gateway](/components/api-gateway/#overview-overview))
- The messaging channel to receive Events, enrich them, and trigger business flows using Functions or services ([Event Mesh](/components/event-mesh/#overview-overview), NATS)
- CLI supported by the intuitive UI ([Console](/components/console/#overview-overview))
- Asset management and storing tool ([Rafter](/components/rafter/#overview-overview), MinIO)
- Backup of Kyma clusters ([Kyma Backup](/root/kyma/#installation-back-up-kyma))
