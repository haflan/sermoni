#!/bin/sh

# index.html for development
cat <<EOF > index.html
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
cp index.html html.go
sed -i 's/`/`\+"`"\+`/g' html.go
cat <<EOF > html.go
package http

var websiteHTML = []byte(\`
$(cat html.go)
\`)
EOF
