files := user.go follow.go

.PHONY: model test

model: model_gen.go

model_gen.go: $(files)
	@if [[ -f $@ ]]; then\
		rm $@;\
	fi
	gomodel -model -sql -t model.tmpl $^
	@gofmt -w $@

test:
	@go test ./...
