APP_NAME="tex"
SRC_DIR="./src"
BUILD_DIR="./build"

echo "📦 Building $APP_NAME..."

mkdir -p "$BUILD_DIR"

cd "$SRC_DIR"
go mod tidy

go build -o "../build/$APP_NAME" .

cd ..

INSTALL_PATH="/usr/local/bin"

if [ -w "$INSTALL_PATH" ]; then

    cp "$BUILD_DIR/$APP_NAME" "$INSTALL_PATH/$APP_NAME"

    chmod +x "$INSTALL_PATH/$APP_NAME"

    echo "✅ Installed successfully!"

    echo ""
    echo "👉 Run it with:"
    echo "   tex"
else
    echo "Installation Failed: Permission Denied"
fi