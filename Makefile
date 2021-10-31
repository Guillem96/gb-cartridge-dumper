BUILDDIR=build
TARGET=gb-dumper

RPI_USER=pi
RPI_HOST=192.168.1.54
RPI_DESTINATION=/home/pi/$(TARGET)


directories:
	mkdir -p $(BUILDDIR)

build: directories
	go build -o $(BUILDDIR)/$(TARGET) main.go

rpi-build:
	GOOS=linux GOARCH=arm GOARM=5 go build -o $(BUILDDIR)/rpi-$(TARGET)
	chmod +x $(BUILDDIR)/rpi-$(TARGET)

rpi-upload: rpi-build
	scp $(BUILDDIR)/rpi-$(TARGET) $(RPI_USER)@$(RPI_HOST):$(RPI_DESTINATION)

run: build
	$(BUILDDIR)/$(TARGET)