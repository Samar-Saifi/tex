#!/usr/bin/env bash

set -e

APP_NAME="tex"
SRC_DIR="./src"
BUILD_DIR="./build"

echo "📦 Building $APP_NAME..."

mkdir -p "$BUILD_DIR"

pushd "$SRC_DIR" >/dev/null
go mod tidy
go build -o "../build/$APP_NAME" .
popd >/dev/null

if [ "$EUID" -eq 0 ]; then
    INSTALL_PATH="/usr/local/bin"
else
    INSTALL_PATH="$HOME/.local/bin"
    mkdir -p "$INSTALL_PATH"
fi

cp "$BUILD_DIR/$APP_NAME" "$INSTALL_PATH/$APP_NAME"
chmod +x "$INSTALL_PATH/$APP_NAME"

echo "✅ Installed successfully to $INSTALL_PATH"

if [ "$EUID" -ne 0 ]; then
    case "$(basename "$SHELL")" in
        zsh)
            SHELL_RC="$HOME/.zshrc"
            ;;
        bash)
            SHELL_RC="$HOME/.bashrc"
            ;;
        fish)
            SHELL_RC="$HOME/.config/fish/config.fish"
            ;;
        *)
            SHELL_RC="$HOME/.profile"
            ;;
    esac

    echo "Using rc file: $SHELL_RC"


    if ! echo ":$PATH:" | grep -q ":$HOME/.local/bin:"; then
        export PATH="$HOME/.local/bin:$PATH"

        if [ ! -f "$SHELL_RC" ] || ! grep -qxF 'export PATH="$HOME/.local/bin:$PATH"' "$SHELL_RC"; then
            echo '' >> "$SHELL_RC"
            echo '# Added by tex installer' >> "$SHELL_RC"
            echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$SHELL_RC"

            echo ""
            echo "⚠️ Added ~/.local/bin to your PATH."
            echo "Run:"
            echo "    source $SHELL_RC"
        fi
    fi
fi

echo ""
echo "👉 Run it with:"
echo "   $APP_NAME"