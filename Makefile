examples:
	cd examples && go build && \
	./examples -dsn "postgres://postgres:@localhost/migrate_test?sslmode=disable" \
		-root-dsn "postgres://postgres:@localhost/postgres?sslmode=disable" \
		-dbname migrate_test

.PHONY: examples
