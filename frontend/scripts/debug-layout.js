#!/usr/bin/env node

/**
 * Debug script to check Next.js root layout configuration
 * Run: node scripts/debug-layout.js
 */

const fs = require('fs');
const path = require('path');

console.log('========================================');
console.log('Next.js Root Layout Debugger');
console.log('========================================\n');

const appDir = path.join(__dirname, '..', 'app');
const layoutPath = path.join(appDir, 'layout.tsx');

console.log('1. Checking layout file existence...');
if (fs.existsSync(layoutPath)) {
  console.log('   ✅ Layout file exists:', layoutPath);
} else {
  console.log('   ❌ Layout file NOT found:', layoutPath);
  process.exit(1);
}

console.log('\n2. Reading layout file content...');
const layoutContent = fs.readFileSync(layoutPath, 'utf-8');
console.log('   ✅ File read successfully');
console.log('   File size:', layoutContent.length, 'bytes');
console.log('   Line count:', layoutContent.split('\n').length);

console.log('\n3. Checking for required tags...');
const hasHtml = /<html[^>]*>/i.test(layoutContent);
const hasBody = /<body[^>]*>/i.test(layoutContent);
const hasClosingHtml = /<\/html>/i.test(layoutContent);
const hasClosingBody = /<\/body>/i.test(layoutContent);

console.log('   <html> tag:', hasHtml ? '✅ Found' : '❌ Missing');
console.log('   </html> tag:', hasClosingHtml ? '✅ Found' : '❌ Missing');
console.log('   <body> tag:', hasBody ? '✅ Found' : '❌ Missing');
console.log('   </body> tag:', hasClosingBody ? '✅ Found' : '❌ Missing');

console.log('\n4. Checking layout structure...');
const hasDefaultExport = /export\s+default\s+function\s+RootLayout/i.test(layoutContent);
const hasChildren = /children:\s*React\.ReactNode/i.test(layoutContent);
const hasReturn = /return\s*\(/i.test(layoutContent);

console.log('   Default export RootLayout:', hasDefaultExport ? '✅ Found' : '❌ Missing');
console.log('   children prop type:', hasChildren ? '✅ Found' : '❌ Missing');
console.log('   return statement:', hasReturn ? '✅ Found' : '❌ Missing');

console.log('\n5. Checking for common issues...');
const issues = [];

if (layoutContent.includes("'use client'")) {
  issues.push('⚠️  Layout is marked as client component (should be server component)');
}

if (!hasHtml || !hasBody) {
  issues.push('❌ Missing required html/body tags');
}

if (!hasDefaultExport) {
  issues.push('❌ Missing default export RootLayout function');
}

const htmlTagMatch = layoutContent.match(/<html[^>]*>/i);
const bodyTagMatch = layoutContent.match(/<body[^>]*>/i);

if (htmlTagMatch) {
  console.log('\n6. HTML tag details:');
  console.log('   Tag:', htmlTagMatch[0]);
}

if (bodyTagMatch) {
  console.log('\n7. Body tag details:');
  console.log('   Tag:', bodyTagMatch[0]);
}

console.log('\n8. File structure:');
const lines = layoutContent.split('\n');
lines.forEach((line, index) => {
  if (line.includes('<html') || line.includes('<body') || line.includes('</html') || line.includes('</body')) {
    console.log(`   Line ${index + 1}: ${line.trim()}`);
  }
});

if (issues.length > 0) {
  console.log('\n⚠️  Issues found:');
  issues.forEach(issue => console.log('   ', issue));
} else {
  console.log('\n✅ No issues detected in layout file structure');
}

console.log('\n9. Checking for other layout files...');
const allLayouts = [];
function findLayouts(dir) {
  const files = fs.readdirSync(dir);
  files.forEach(file => {
    const filePath = path.join(dir, file);
    const stat = fs.statSync(filePath);
    if (stat.isDirectory()) {
      findLayouts(filePath);
    } else if (file === 'layout.tsx' || file === 'layout.js') {
      allLayouts.push(filePath);
    }
  });
}

findLayouts(appDir);
console.log(`   Found ${allLayouts.length} layout file(s):`);
allLayouts.forEach(layout => {
  const relPath = path.relative(path.join(__dirname, '..'), layout);
  console.log(`   - ${relPath}`);
});

console.log('\n========================================');
console.log('Debug complete');
console.log('========================================\n');
