package templates

const DocumentTemplate = `<html>
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
{{.Description}}
</section>

<section>
{{.Content}}
</section>

</article>

</body>
</html>`
