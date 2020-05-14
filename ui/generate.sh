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
INDEX_HTML_SANITIZED=${INDEX_HTML}_sane
cp $INDEX_HTML $INDEX_HTML_SANITIZED
sed -i 's/`/`\+"`"\+`/g' $INDEX_HTML_SANITIZED
cat <<EOF > $HTML_GO
// +build PRODUCTION

package http

const PRODUCTION = true;

func getWebsite() []byte {
	return []byte(\`
$(cat $INDEX_HTML_SANITIZED)
	\`)
}
EOF
rm $INDEX_HTML_SANITIZED
