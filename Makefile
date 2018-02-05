BIN := vault-spec

all: $(BIN)

$(BIN): 
	bazel build //:vault-spec

clean:
	rm -f $(BIN)

test:
	bazel test $(shell bazel query 'kind("go_test", //... except filter(//vendor, //...))') --test_output=all

.PHONY: clean
.PHONY: all
.PHONY: $(BIN)
