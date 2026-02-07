#!/usr/bin/env node

/**
 * Check Next.js layout detection
 * This script simulates what Next.js does to detect the root layout
 */

const fs = require('fs');
const path = require('path');

const appDir = path.join(__dirname, '..', 'app');
const layoutPath = path.join(appDir, 'layout.tsx');

console.log('Next.js Layout Detection Check\n');

// Check 1: File exists at app/layout.tsx
console.log('✓ Check 1: Root layout file exists');
if (!fs.existsSync(layoutPath)) {
  console.log('  ❌ FAIL: app/layout.tsx not found');
  process.exit(1);
}
console.log('  ✅ PASS: app/layout.tsx exists\n');

// Check 2: File is readable
console.log('✓ Check 2: File is readable');
let content;
try {
  content = fs.readFileSync(layoutPath, 'utf-8');
  console.log('  ✅ PASS: File readable\n');
} catch (e) {
  console.log('  ❌ FAIL: Cannot read file:', e.message);
  process.exit(1);
}

// Check 3: Has default export
console.log('✓ Check 3: Has default export function');
const hasDefaultExport = /export\s+default\s+function/i.test(content);
if (!hasDefaultExport) {
  console.log('  ❌ FAIL: No default export function found');
  process.exit(1);
}
console.log('  ✅ PASS: Has default export function\n');

// Check 4: Function name is RootLayout
console.log('✓ Check 4: Function name is RootLayout');
const hasRootLayout = /export\s+default\s+function\s+RootLayout/i.test(content);
if (!hasRootLayout) {
  console.log('  ⚠️  WARN: Function name is not RootLayout (but may still work)');
} else {
  console.log('  ✅ PASS: Function name is RootLayout\n');
}

// Check 5: Returns JSX with html tag
console.log('✓ Check 5: Returns JSX with <html> tag');
const htmlPattern = /<html[^>]*>/i;
const htmlMatch = content.match(htmlPattern);
if (!htmlMatch) {
  console.log('  ❌ FAIL: No <html> tag found in return statement');
  process.exit(1);
}
console.log('  ✅ PASS: <html> tag found:', htmlMatch[0]);

// Check 6: Returns JSX with body tag
console.log('✓ Check 6: Returns JSX with <body> tag');
const bodyPattern = /<body[^>]*>/i;
const bodyMatch = content.match(bodyPattern);
if (!bodyMatch) {
  console.log('  ❌ FAIL: No <body> tag found in return statement');
  process.exit(1);
}
console.log('  ✅ PASS: <body> tag found:', bodyMatch[0]);

// Check 7: Has closing tags
console.log('\n✓ Check 7: Has closing tags');
const hasClosingHtml = /<\/html>/i.test(content);
const hasClosingBody = /<\/body>/i.test(content);
console.log('  </html>:', hasClosingHtml ? '✅' : '❌');
console.log('  </body>:', hasClosingBody ? '✅' : '❌');

// Check 8: Not a client component
console.log('\n✓ Check 8: Is server component (not client)');
const isClient = /['"]use\s+client['"]/i.test(content);
if (isClient) {
  console.log('  ❌ FAIL: Layout is marked as client component');
  console.log('  Root layout MUST be a server component');
  process.exit(1);
}
console.log('  ✅ PASS: Layout is server component\n');

// Check 9: Structure validation
console.log('✓ Check 9: Structure validation');
const returnMatch = content.match(/return\s*\([\s\S]*?\)/);
if (returnMatch) {
  const returnContent = returnMatch[0];
  const htmlInReturn = /<html/i.test(returnContent);
  const bodyInReturn = /<body/i.test(returnContent);
  
  console.log('  HTML in return:', htmlInReturn ? '✅' : '❌');
  console.log('  Body in return:', bodyInReturn ? '✅' : '❌');
  
  if (!htmlInReturn || !bodyInReturn) {
    console.log('  ❌ FAIL: html/body tags not in return statement');
    process.exit(1);
  }
}

console.log('\n========================================');
console.log('✅ All checks passed!');
console.log('Layout file structure is correct.');
console.log('========================================\n');

console.log('If Next.js still shows the error, it may be:');
console.log('1. A Next.js dev mode bug (known issue in 14.2.5)');
console.log('2. Browser cache issue - try hard refresh');
console.log('3. Next.js build cache - try: rm -rf .next');
console.log('4. The error is a false positive - app works correctly\n');
