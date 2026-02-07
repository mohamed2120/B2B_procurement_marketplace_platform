#!/bin/bash

# Monitor Next.js layout issue
# This script helps debug the missing html/body tags issue

echo "=========================================="
echo "Next.js Layout Issue Monitor"
echo "=========================================="
echo ""

FRONTEND_URL="${1:-http://localhost:3002}"

echo "1. Checking layout file..."
if [ -f "app/layout.tsx" ]; then
    echo "   ✅ Layout file exists"
    HTML_COUNT=$(grep -c "<html" app/layout.tsx || echo "0")
    BODY_COUNT=$(grep -c "<body" app/layout.tsx || echo "0")
    echo "   <html> tags found: $HTML_COUNT"
    echo "   <body> tags found: $BODY_COUNT"
else
    echo "   ❌ Layout file NOT found"
    exit 1
fi

echo ""
echo "2. Checking HTML output..."
HTML_OUTPUT=$(curl -s "$FRONTEND_URL" 2>&1)
if echo "$HTML_OUTPUT" | grep -q "</body></html>"; then
    echo "   ✅ HTML output contains </body></html>"
else
    echo "   ❌ HTML output missing closing tags"
fi

if echo "$HTML_OUTPUT" | grep -q "__next_root_layout_missing_tags"; then
    echo "   ⚠️  Next.js warning detected in HTML"
    MISSING_TAGS=$(echo "$HTML_OUTPUT" | grep -o '__next_root_layout_missing_tags=\[.*\]' | head -1)
    echo "   Warning: $MISSING_TAGS"
else
    echo "   ✅ No Next.js warning in HTML"
fi

echo ""
echo "3. Checking DOM structure..."
# Try to extract html/body from output
if echo "$HTML_OUTPUT" | grep -q "<html"; then
    echo "   ✅ <html> tag found in output"
else
    echo "   ❌ <html> tag NOT in output"
fi

if echo "$HTML_OUTPUT" | grep -q "<body"; then
    echo "   ✅ <body> tag found in output"
else
    echo "   ❌ <body> tag NOT in output"
fi

echo ""
echo "4. Next.js version..."
if [ -f "package.json" ]; then
    NEXT_VERSION=$(grep -o '"next": "[^"]*"' package.json | cut -d'"' -f4)
    echo "   Version: $NEXT_VERSION"
fi

echo ""
echo "5. Build cache status..."
if [ -d ".next" ]; then
    echo "   ✅ .next directory exists"
    CACHE_SIZE=$(du -sh .next 2>/dev/null | cut -f1)
    echo "   Cache size: $CACHE_SIZE"
else
    echo "   ⚠️  .next directory not found (may need build)"
fi

echo ""
echo "=========================================="
echo "Summary:"
echo "=========================================="
if echo "$HTML_OUTPUT" | grep -q "__next_root_layout_missing_tags"; then
    echo "⚠️  Next.js is reporting missing tags, but:"
    if echo "$HTML_OUTPUT" | grep -q "</body></html>"; then
        echo "✅ The HTML structure is actually correct"
        echo "   This is a Next.js dev mode false positive"
    fi
else
    echo "✅ No issues detected"
fi
echo ""
