files := user.go follow.go

.PHONY: model test

model: model_gen.go

model_gen.go: $(files)
	@if [[ -f $@ ]]; then\
		rm $@;\
	fi
	gomodel -model -sql $^
	@sed -i 's/package\ gomodel/package\ $(shell basename `pwd`)/' $@ # replace package name with name of current directory
	@gofmt -w $@

test:
	@go test ./...
