#!/bin/sh

INDEX_HTML=dist/index.html
HTML_GO=dist/html.go

# index.html for development
cat <<EOF > $INDEX_HTML
<html>
<head>
<meta name="viewport" content="width=device-width,initial-scale=1">
</head>
<body>
<div id="app"></div>
<script>
$(cat ./dist/sermoni.js)
</script>
<body>
</html>
EOF

# html.go for production, ` must be replaced by `+"`"+`
cp $INDEX_HTML $HTML_GO
sed -i 's/`/`\+"`"\+`/g' $HTML_GO
cat <<EOF > $HTML_GO
// +build PRODUCTION
package http

var websiteHTML = []byte(\`
$(cat $HTML_GO)
\`)
EOF
