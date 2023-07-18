v1.2 (18.07.2023)
- Splitted services: eventstore, projection, projector
- Add kubernetes-deployment.yml

v1.1 (01.03.2023)
- Use testutil/TestCase for http tests
- Splitted services: eventstore, projection
- Folder refactoring:
    - main.go -> eventstore/cmd/main.go
    - controller/controller.go -> eventstore/handler/http/http.go
    - controller/router.go -> eventstore/handler/http/http.go
    - controller/controller_test.go -> eventstore/handler/http/http_test.go
    - controller/test_router.go -> eventstore/handler/http/http_test.go
    - db/answer_utility.go -> eventstore/service/eventstore/utility.go
    - db/event_store/event_store_db.go -> eventstore/repository/esdb/esdb.go
    - db/read_repository/read_repository_mongo.go -> projection/repository/mongodb/mongodb.go
    - service/projector.go -> projection/service/projector/service.go
    - service/qa_projection.go -> projection/service/projection/service.go
    - service/qa_service.go -> eventstore/service/eventstore/service.go
    - service/qa_service_integration_test.go -> test/integration/db/integration_test.go
    - service/qa_service_test.go -> eventstore/service/eventstore/service_test.go
