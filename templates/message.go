package templates

const messageTemplate = `<!DOCTYPE HTML>
<html lang="{{.LanguageTag}}">
<head>
	<title>{{.Title}}</title>
</head>

<body>

<article>

<header>
<h1>
{{.Title}}
</h1>
</header>

<section>
{{.Content}}
</section>

</article>

</body>
</html>`
