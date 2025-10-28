# Makefile for building the go-check TUI and its Windows installer

# --- Variables ---
# Define the path to your Go executable inside the Installer dir
GO_EXECUTABLE=Installer/go-check.exe

# Define the path to your WiX project file
WIX_PROJECT=Installer/go-check.wixproj

# --- Targets ---

# Default target: build everything.
# By typing just 'make', this will run the 'msi' target.
.PHONY: all
all: msi

# Target to build the Windows executable.
# This builds the .exe directly into the 'Installer' folder,
# where the .wxs file expects to find it.
$(GO_EXECUTABLE):
	@echo "Building Go executable for Windows..."
	@go build -o $(GO_EXECUTABLE) .

# Target to build the MSI installer.
# This 'depends' on the $(GO_EXECUTABLE) target, so 'make' will
# automatically run that target first.
.PHONY: msi
msi: $(GO_EXECUTABLE)
	@echo "Building Windows MSI installer..."
	@dotnet build $(WIX_PROJECT)

build: $(GO_EXECUTABLE)
	@cpy $(GO_EXECUTABLE) .
# Target to clean up build artifacts.
.PHONY: clean
clean:
	@echo "Cleaning up build artifacts..."
	@rm -f $(GO_EXECUTABLE)
	@rm -rf Installer/bin
	@rm -rf Installer/obj

