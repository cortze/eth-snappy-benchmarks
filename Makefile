
GOCC=go

BUILD_PATH="./build/"
BIN="$(BUILD_PATH)/snappy-benchmark"

.PHONY: clean install build dependencies

build:
	mkdir -p $(BUILD_PATH)
	$(GOCC) build -o $(BIN)

install:
	$(GOCC) install

clean:
	rm -r $(BUILD_PATH)
