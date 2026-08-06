get_version() {
    sed -n 's/.*Version = "\(.*\)"/\1/p' \
        "$PROJECT_ROOT/backend/internal/version/version.go"
}