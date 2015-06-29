Here is an example show how to use the gomodel package and cmd tool.

This example is about `User`-`Follow`, `User` can `Follow` and be `Follow`ed by other `User`s.

* Main files:
    - `db.go`: `DB` and tables configurations
    - `user.go`: `User` structure
    - `follow.go`: `Follow` relationship structure

* Other files:
    - `model_gen.go`: generated code for `User` and `Follow`
    - `model.tmpl`: template file for `gmodel` cmd's output
    - `Makefile`: automation
    - `example_test.go`: test cases

