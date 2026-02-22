#!/bin/bash

# setup.sh - Helper script to initialize a new project from a template
# Usage: ./setup.sh

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Go Project Setup ===${NC}"

# 1. List available templates
echo -e "\nAvailable templates:"
templates=($(find . -maxdepth 1 -type d -name "template-*" | sed 's|./||' | sort))

if [ ${#templates[@]} -eq 0 ]; then
    echo -e "${RED}No templates found in the current directory.${NC}"
    exit 1
fi

for i in "${!templates[@]}"; do
    echo "$((i+1))) ${templates[$i]}"
done

# 2. Select a template
read -p "Select a template (enter number): " template_idx
template_idx=$((template_idx-1))

if [ -z "${templates[$template_idx]}" ]; then
    echo -e "${RED}Invalid selection.${NC}"
    exit 1
fi

SOURCE_TEMPLATE="${templates[$template_idx]}"
echo -e "${GREEN}Selected: ${SOURCE_TEMPLATE}${NC}"

# 3. Get Project Name
read -p "Enter new project name (directory name): " PROJECT_NAME
if [ -z "$PROJECT_NAME" ]; then
    echo -e "${RED}Project name cannot be empty.${NC}"
    exit 1
fi

# 4. Get Module Name
read -p "Enter Go module name (e.g., github.com/user/my-project): " MODULE_NAME
if [ -z "$MODULE_NAME" ]; then
    echo -e "${RED}Module name cannot be empty.${NC}"
    exit 1
fi

# 5. Copy Template
TARGET_DIR="../$PROJECT_NAME"

if [ -d "$TARGET_DIR" ]; then
    echo -e "${RED}Directory $TARGET_DIR already exists!${NC}"
    read -p "Do you want to overwrite it? (y/N): " overwrite
    if [[ "$overwrite" != "y" && "$overwrite" != "Y" ]]; then
        echo "Aborting."
        exit 1
    fi
    rm -rf "$TARGET_DIR"
fi

echo -e "\n${BLUE}Creating project in $TARGET_DIR...${NC}"
cp -R "$SOURCE_TEMPLATE" "$TARGET_DIR"

# 6. Process the new project
cd "$TARGET_DIR"

# Get old module name
OLD_MODULE=$(grep "^module" go.mod | awk '{print $2}')
echo "Old module: $OLD_MODULE"
echo "New module: $MODULE_NAME"

# Rename module in go.mod
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS sed
    sed -i '' "s|^module .*|module ${MODULE_NAME}|g" go.mod
    find . -type f -name "*.go" -exec sed -i '' "s|${OLD_MODULE}|${MODULE_NAME}|g" {} +
else
    # Linux/GNU sed
    sed -i "s|^module .*|module ${MODULE_NAME}|g" go.mod
    find . -type f -name "*.go" -exec sed -i "s|${OLD_MODULE}|${MODULE_NAME}|g" {} +
fi

# Update config.yaml app name if present
if [ -f "config/config.yaml" ]; then
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "s|name: \"${SOURCE_TEMPLATE}\"|name: \"${PROJECT_NAME}\"|g" config/config.yaml
        sed -i '' "s|name: \"go-${SOURCE_TEMPLATE}\"|name: \"${PROJECT_NAME}\"|g" config/config.yaml
    else
        sed -i "s|name: \"${SOURCE_TEMPLATE}\"|name: \"${PROJECT_NAME}\"|g" config/config.yaml
        sed -i "s|name: \"go-${SOURCE_TEMPLATE}\"|name: \"${PROJECT_NAME}\"|g" config/config.yaml
    fi
fi

# Initialize Git
echo -e "\n${BLUE}Initializing Git...${NC}"
rm -rf .git
git init

# Run Go Mod Tidy
echo -e "\n${BLUE}Running go mod tidy...${NC}"
go mod tidy

echo -e "\n${GREEN}Success! Project created at: ${TARGET_DIR}${NC}"
echo -e "To get started:"
echo -e "  cd ${TARGET_DIR}"
echo -e "  make run"
