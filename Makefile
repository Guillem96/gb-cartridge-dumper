BUILDDIR=build
TARGET=gb-dumper

directories:
	mkdir -p $(BUILDDIR)

build: directories
	go build -o $(BUILDDIR)/$(TARGET) main.go

run: build
	$(BUILDDIR)/$(TARGET)