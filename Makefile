.PHONY: tests

tests:
	APP_ENV=test go test -v $$(glide novendor)
