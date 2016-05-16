.PHONY: tests

tests:
	APP_ENV=test go test $$(glide novendor)
