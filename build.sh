npx tailwindcss -i ./main.css -o ./static/css/getsrc.css

echo '{{define "css"}}<link href="/css/getsrc.css?'`md5 -q static/css/getsrc.css`'" rel="stylesheet">{{end}}' > tmpl/gen.go.html