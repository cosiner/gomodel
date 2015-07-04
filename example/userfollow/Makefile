files := user.go follow.go

.PHONY: model test

model: model_gen.go

model_gen.go: $(files)
	@if [[ -f $@ ]]; then\
		rm $@;\
	fi
	gomodel -model -sql -t ../model.tmpl $^
	@sed -i 's/package\ example/package\ $(shell basename `pwd`)/' $@ # replace package name with name of current directory
	@gofmt -w $@

test:
	@go test ./...
