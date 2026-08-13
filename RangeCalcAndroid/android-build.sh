#!/bin/zsh
#
# android-build.sh - Helper for building RangeCalcAndroid from CLI
# (Android project on JohnDovey external drive)
#
# Configuration cache is enabled (see gradle.properties) for much faster
# repeated builds after the first configuration.
#   Override for one run: ./android-build.sh --no-configuration-cache ...
#
# Usage examples:
#   ./android-build.sh                  # assembleDebug (default)
#   ./android-build.sh assembleRelease
#   ./android-build.sh installDebug
#   ./android-build.sh clean
#   ./android-build.sh tasks            # list available tasks
#   ./android-build.sh assembleDebug installDebug
#   ./android-build.sh --no-configuration-cache assembleDebug   # bypass CC for one run
#   ./android-build.sh --warning-mode all :app:assembleDebug    # see all deprecation details
#
# Make sure the JohnDovey drive is mounted.

set -e

# Activate JohnDovey environment (quietly)
if [[ -f ~/source-john-dovey.sh ]]; then
    source ~/source-john-dovey.sh --quiet
fi

# Move to project root (directory containing this script)
cd "$(dirname "$0")"

# Ensure wrapper is executable
chmod +x gradlew 2>/dev/null || true

# Default to assembleDebug if no arguments. All args (tasks + gradle flags like
# --no-configuration-cache, --warning-mode all, -x test, etc.) are passed through.
if [ $# -eq 0 ]; then
    set -- "assembleDebug"
fi

echo "🚀 Building RangeCalcAndroid on JohnDovey..."
echo "   Args: $@"
echo "   ANDROID_HOME=$ANDROID_HOME"
echo ""

./gradlew "$@"

echo ""
echo "✅ Build finished"
