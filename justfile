# Show available recipes
default:
    @{{just_executable()}} --list

build:
    go build -o gitrpt ./cmd/gitrpt

run *args:
    go run ./cmd/gitrpt {{args}}

test:
    go test ./...

# =============================================================================
# i18n Commands - go-i18n Workflow
# =============================================================================
#
# Standard workflow (3 steps):
#   1. just i18n-extract    - Extract English messages from Go code
#   2. just i18n-translate  - Generate translate.*.json files for translation
#   3. just i18n-merge      - Merge translations into active.*.json
#
# Warning: i18n-translate overwrites existing translate.*.json files!
# Run i18n-merge first if you have unfinished translations.

# Install goi18n CLI tool globally
i18n-install:
    go install github.com/nicksnyder/go-i18n/v2/goi18n@latest
    @echo "goi18n installed. Run 'goi18n -help' for usage."

# Step 1: Extract English messages from Go code.
# Creates/updates active.en.json from &i18n.Message{...} literals in source files.
i18n-extract:
    goi18n extract -sourceLanguage en -format json -outdir ./internal/i18n/messages ./internal/i18n
    @echo ""
    @echo "✓ Extracted English messages to active.en.json"
    @echo "  Note: English text is embedded in code via DefaultMessage"
    @echo "  Next: just i18n-translate"

# Step 2: Generate translate.*.json files for all supported languages.
# WARNING: Overwrites existing translate.*.json files!
# Run 'just i18n-merge' first if you have unfinished translations.
i18n-translate:
    #!/usr/bin/env bash
    set -e
    echo "⚠️  Warning: This will overwrite existing translate.*.json files"
    echo "   Ensure you have merged previous translations (just i18n-merge)"
    echo ""
    goi18n merge \
        -sourceLanguage en \
        -format json \
        -outdir ./internal/i18n/messages \
        ./internal/i18n/messages/active.*.json
    echo ""
    echo "✓ Generated translate.*.json files"
    echo "  Translate all 'other' fields in these files"
    echo "  Then run: just i18n-merge"

# Step 3: Merge translated messages into active.*.json files.
# Merges translate.*.json into active.*.json and removes translate files.
i18n-merge:
    #!/usr/bin/env bash
    set -e
    goi18n merge \
        -sourceLanguage en \
        -format json \
        -outdir ./internal/i18n/messages \
        ./internal/i18n/messages/active.*.json \
        ./internal/i18n/messages/translate.*.json
    rm -f ./internal/i18n/messages/translate.*.json
    echo ""
    echo "✓ Merged translations and cleaned up translate files"
    echo "  Rebuild: just build"

# Create a new language translation file.
# Usage: just i18n-new-lang de
# Note: For multiple new languages, use i18n-translate after extracting.
i18n-new-lang lang:
    #!/usr/bin/env bash
    set -e
    if [ -f "./internal/i18n/messages/active.{{lang}}.json" ]; then
        echo "Error: active.{{lang}}.json already exists"
        exit 1
    fi
    # Create empty file, then merge to populate with English source
    echo '{}' > ./internal/i18n/messages/active.{{lang}}.json
    goi18n merge \
        -sourceLanguage en \
        -format json \
        -outdir ./internal/i18n/messages \
        ./internal/i18n/messages/active.en.json \
        ./internal/i18n/messages/active.{{lang}}.json
    echo ""
    echo "✓ Created active.{{lang}}.json"
    echo "  Translate the 'other' fields in this file"
    echo "  Rebuild: just build"
